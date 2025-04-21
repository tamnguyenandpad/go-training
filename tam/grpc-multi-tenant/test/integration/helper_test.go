package integration

import (
	"context"
	"database/sql"
	"fmt"
	"net"
	"testing"
	"time"

	"github.com/google/go-cmp/cmp"
	"github.com/google/go-cmp/cmp/cmpopts"
	tenant_pb "github.com/tuannguyenandpadcojp/go-training/tam/grpc-multi-tenant/gen/go/tenant/v1"
	user_pb "github.com/tuannguyenandpadcojp/go-training/tam/grpc-multi-tenant/gen/go/user/v1"
	"github.com/tuannguyenandpadcojp/go-training/tam/grpc-multi-tenant/internal/pkg/logger"
	user_adapter "github.com/tuannguyenandpadcojp/go-training/tam/grpc-multi-tenant/internal/tenant/adapter"
	tenant_app "github.com/tuannguyenandpadcojp/go-training/tam/grpc-multi-tenant/internal/tenant/app"
	tenant_datastore "github.com/tuannguyenandpadcojp/go-training/tam/grpc-multi-tenant/internal/tenant/infra/datastore"
	tenant_ports "github.com/tuannguyenandpadcojp/go-training/tam/grpc-multi-tenant/internal/tenant/ports"
	tenant_adapter "github.com/tuannguyenandpadcojp/go-training/tam/grpc-multi-tenant/internal/user/adapter"
	user_app "github.com/tuannguyenandpadcojp/go-training/tam/grpc-multi-tenant/internal/user/app"
	user_datastore "github.com/tuannguyenandpadcojp/go-training/tam/grpc-multi-tenant/internal/user/infra/datastore"
	user_ports "github.com/tuannguyenandpadcojp/go-training/tam/grpc-multi-tenant/internal/user/ports"
	"github.com/tuannguyenandpadcojp/go-training/tam/grpc-multi-tenant/test/integration/testutil"

	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
	"google.golang.org/grpc/reflection"
	"google.golang.org/grpc/test/bufconn"
)

const (
	host = "localhost"
	port = "50051"
)

var (
	helperWithoutCtx = &ServiceTestHelper{}
	log              = logger.InitLogger(logger.Debug, "integration-tests")
)

type ServiceTestHelper struct {
	cli     tenant_pb.TenantServiceClient
	userCli user_pb.UserServiceClient
	ctx     context.Context
	DB      *sql.DB
}

func (r *ServiceTestHelper) CreateServiceTestHelper(t *testing.T) *ServiceTestHelper {
	log.Info("Creating service test helper", map[string]interface{}{"test": t.Name()})
	db, dbName := testutil.InitDB(t)
	log.Debug("Test database initialized", map[string]interface{}{"dbName": dbName})

	// Create gRPC server
	grpcServer := grpc.NewServer()
	reflection.Register(grpcServer)
	log.Info("gRPC test server initialized", nil)

	// Initialize repo
	tenantRepository := tenant_datastore.NewTenantMysqlRepository(db)
	userRepository := user_datastore.NewUserMysqlRepository(db)
	log.Debug("Test repositories initialized", nil)

	// Initialize adapter
	tenantAdapter := tenant_adapter.NewTenantAdapter(tenantRepository)
	userAdapter := user_adapter.NewUserAdapter(userRepository)
	log.Debug("Test adapters initialized", nil)

	// Initialize application
	tenantApplication := tenant_app.NewApplication(tenantRepository, userAdapter)
	userApplication := user_app.NewApplication(userRepository, tenantAdapter)
	log.Debug("Test applications initialized", nil)

	// Initialize service
	tenantGrpcService := tenant_ports.NewGrpcServer(tenantApplication)
	userGrpcService := user_ports.NewGrpcServer(userApplication)
	log.Debug("Test gRPC services initialized", nil)

	// Register service
	tenant_pb.RegisterTenantServiceServer(grpcServer, tenantGrpcService)
	user_pb.RegisterUserServiceServer(grpcServer, userGrpcService)
	log.Debug("Test services registered", nil)

	listener := bufconn.Listen(1024 * 1024)

	go func() {
		defer grpcServer.GracefulStop()
		if err := grpcServer.Serve(listener); err != nil {
			log.Error("Failed to serve test gRPC server", map[string]interface{}{"error": err.Error()})
		}
	}()

	bufDialer := func(ctx context.Context, address string) (net.Conn, error) {
		// notlint: wrapcheck
		return listener.Dial()
	}

	var ctx context.Context
	if r.ctx == nil {
		deadlineCtx, cancelFunc := context.WithTimeout(context.Background(), 5*time.Second)
		t.Cleanup(cancelFunc)
		ctx = deadlineCtx
	} else {
		ctx = r.ctx
	}

	conn, err := grpc.NewClient(
		fmt.Sprintf("%s:%s", host, port),
		grpc.WithContextDialer(bufDialer),
		grpc.WithTransportCredentials(insecure.NewCredentials()))
	if err != nil {
		log.Error("Failed to create gRPC client", map[string]interface{}{"error": err.Error()})
		t.Fatal(err)
	}

	client := tenant_pb.NewTenantServiceClient(conn)
	userClient := user_pb.NewUserServiceClient(conn)
	log.Info("Test clients initialized", nil)

	t.Cleanup(func() {
		log.Debug("Cleaning up test resources", map[string]interface{}{"test": t.Name()})
		conn.Close()
		grpcServer.GracefulStop()
	})

	return &ServiceTestHelper{
		cli:     client,
		userCli: userClient,
		ctx:     ctx,
		DB:      db,
	}
}

func CreateServiceTestHelper(t *testing.T) *ServiceTestHelper {
	return helperWithoutCtx.CreateServiceTestHelper(t)
}

func batchIgnoreProtoUnexportedFields(typs ...interface{}) []cmp.Option {
	opts := make([]cmp.Option, len(typs))
	for i, typ := range typs {
		opts[i] = cmpopts.IgnoreFields(typ, "state", "sizeCache", "unknownFields")
	}
	return opts
}
