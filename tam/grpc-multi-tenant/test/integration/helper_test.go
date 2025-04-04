package integration

import (
	"context"
	"database/sql"
	"log"
	"net"
	"testing"
	"time"

	"github.com/google/go-cmp/cmp"
	"github.com/google/go-cmp/cmp/cmpopts"
	tenant_pb "github.com/tuannguyenandpadcojp/go-training/tam/grpc-multi-tenant/gen/go/tenant/v1"
	user_pb "github.com/tuannguyenandpadcojp/go-training/tam/grpc-multi-tenant/gen/go/user/v1"
	tenant_app "github.com/tuannguyenandpadcojp/go-training/tam/grpc-multi-tenant/internal/tenant/app"
	tenant_datastore "github.com/tuannguyenandpadcojp/go-training/tam/grpc-multi-tenant/internal/tenant/infra/datastore"
	tenant_ports "github.com/tuannguyenandpadcojp/go-training/tam/grpc-multi-tenant/internal/tenant/ports"
	user_app "github.com/tuannguyenandpadcojp/go-training/tam/grpc-multi-tenant/internal/user/app"
	user_datastore "github.com/tuannguyenandpadcojp/go-training/tam/grpc-multi-tenant/internal/user/infra/datastore"
	user_ports "github.com/tuannguyenandpadcojp/go-training/tam/grpc-multi-tenant/internal/user/ports"
	"github.com/tuannguyenandpadcojp/go-training/tam/grpc-multi-tenant/test/integration/testutil"
	"google.golang.org/grpc"
	"google.golang.org/grpc/reflection"
	"google.golang.org/grpc/test/bufconn"
)

var (
	helperWithoutCtx = &ServiceTestHelper{}
)

type ServiceTestHelper struct {
	cli     tenant_pb.TenantServiceClient
	userCli user_pb.UserServiceClient
	ctx     context.Context
	DB      *sql.DB
}

func (r *ServiceTestHelper) CreateServiceTestHelper(t *testing.T) *ServiceTestHelper {
	db, _ := testutil.InitDB(t)

	// Create gRPC server
	grpcServer := grpc.NewServer()
	reflection.Register(grpcServer)
	log.Println("gRPC server is running on port 50051...")

	// Initialize repo
	tenantRepository := tenant_datastore.NewTenantMysqlRepository(db)
	userRepository := user_datastore.NewUserMysqlRepository(db)

	// Initialize application
	tenantApplication := tenant_app.NewApplication(tenantRepository)
	userApplication := user_app.NewApplication(userRepository)

	// Initialize service
	tenantGrpcService := tenant_ports.NewGrpcServer(tenantApplication, userApplication)
	userGrpcService := user_ports.NewGrpcServer(userApplication, tenantApplication)

	// Register service
	tenant_pb.RegisterTenantServiceServer(grpcServer, tenantGrpcService)
	user_pb.RegisterUserServiceServer(grpcServer, userGrpcService)

	// if err := grpcServer.Serve(listener); err != nil {
	// 	log.Fatalf("Failed to serve: %v", err)
	// }

	listener := bufconn.Listen(1024 * 1024)

	go func() {
		defer grpcServer.GracefulStop() // ????
		if err := grpcServer.Serve(listener); err != nil {
			log.Fatalf("Failed to serve: %v", err)
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

	conn, err := grpc.DialContext(
		ctx,
		"bufnet",
		grpc.WithContextDialer(bufDialer),
		grpc.WithInsecure())
	if err != nil {
		t.Fatal(err)
	}

	client := tenant_pb.NewTenantServiceClient(conn)
	userClient := user_pb.NewUserServiceClient(conn)

	t.Cleanup(func() {
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
