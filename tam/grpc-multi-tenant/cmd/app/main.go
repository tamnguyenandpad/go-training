package main

import (
	"database/sql"
	"net"

	_ "github.com/go-sql-driver/mysql"

	tenant_pb "github.com/tuannguyenandpadcojp/go-training/tam/grpc-multi-tenant/gen/go/tenant/v1"
	user_pb "github.com/tuannguyenandpadcojp/go-training/tam/grpc-multi-tenant/gen/go/user/v1"
	"github.com/tuannguyenandpadcojp/go-training/tam/grpc-multi-tenant/internal/pkg/config"
	"github.com/tuannguyenandpadcojp/go-training/tam/grpc-multi-tenant/internal/pkg/logger"
	user_adapter "github.com/tuannguyenandpadcojp/go-training/tam/grpc-multi-tenant/internal/tenant/adapter"
	tenant_app "github.com/tuannguyenandpadcojp/go-training/tam/grpc-multi-tenant/internal/tenant/app"
	tenant_datastore "github.com/tuannguyenandpadcojp/go-training/tam/grpc-multi-tenant/internal/tenant/infra/datastore"
	tenant_ports "github.com/tuannguyenandpadcojp/go-training/tam/grpc-multi-tenant/internal/tenant/ports"
	tenant_adapter "github.com/tuannguyenandpadcojp/go-training/tam/grpc-multi-tenant/internal/user/adapter"
	user_app "github.com/tuannguyenandpadcojp/go-training/tam/grpc-multi-tenant/internal/user/app"
	user_datastore "github.com/tuannguyenandpadcojp/go-training/tam/grpc-multi-tenant/internal/user/infra/datastore"
	user_ports "github.com/tuannguyenandpadcojp/go-training/tam/grpc-multi-tenant/internal/user/ports"
	"google.golang.org/grpc"
	"google.golang.org/grpc/reflection"
)

func main() {
	// Load application configuration from environment variables
	cfg, err := config.LoadConfig()
	if err != nil {
		// Use standard log here as our logger is not initialized yet
		logger.Get().Fatal("Failed to load configuration", map[string]interface{}{"error": err.Error()})
	}

	// Initialize logger with config settings
	log := logger.InitLogger(cfg.LogLevel, "grpc-multi-tenant")

	// Log successful configuration loading
	log.Info("Configuration loaded successfully", map[string]interface{}{
		"server_port": cfg.ServerPort,
		"db_host":     cfg.DBHost,
		"db_port":     cfg.DBPort,
		"db_name":     cfg.DBName,
		"log_level":   cfg.LogLevel,
	})

	// Connect to database using config
	db, err := sql.Open("mysql", cfg.DBDSN)
	if err != nil {
		log.Fatal("Failed to connect to database", map[string]interface{}{"error": err.Error()})
	}
	if err := db.Ping(); err != nil {
		log.Fatal("Failed to ping to database", map[string]interface{}{"error": err.Error()})
	}
	defer db.Close()
	log.Info("Database connection established", nil)

	// Create server address from config
	serverAddr := ":" + cfg.ServerPort

	// Create gRPC server
	listener, err := net.Listen("tcp", serverAddr)
	if err != nil {
		log.Fatal("Failed to listen", map[string]interface{}{"error": err.Error(), "address": serverAddr})
	}
	grpcServer := grpc.NewServer()
	reflection.Register(grpcServer)
	log.Info("gRPC server initialized", map[string]interface{}{"port": cfg.ServerPort})

	// Initialize repo
	tenantRepository := tenant_datastore.NewTenantMysqlRepository(db)
	userRepository := user_datastore.NewUserMysqlRepository(db)
	log.Debug("Repositories initialized", nil)

	// Initialize adapter
	tenantAdapter := tenant_adapter.NewTenantAdapter(tenantRepository)
	userAdapter := user_adapter.NewUserAdapter(userRepository)
	log.Debug("Adapters initialized", nil)

	// Initialize application
	tenantApplication := tenant_app.NewApplication(tenantRepository, userAdapter)
	userApplication := user_app.NewApplication(userRepository, tenantAdapter)
	log.Debug("Applications initialized", nil)

	// Initialize service
	tenantGrpcService := tenant_ports.NewGrpcServer(tenantApplication)
	userGrpcService := user_ports.NewGrpcServer(userApplication)
	log.Debug("gRPC services initialized", nil)

	// Register service
	tenant_pb.RegisterTenantServiceServer(grpcServer, tenantGrpcService)
	user_pb.RegisterUserServiceServer(grpcServer, userGrpcService)
	log.Info("gRPC services registered", map[string]interface{}{
		"services": []string{"TenantService", "UserService"},
	})

	log.Info("Starting gRPC server", map[string]interface{}{"address": serverAddr})
	if err := grpcServer.Serve(listener); err != nil {
		log.Fatal("Failed to serve", map[string]interface{}{"error": err.Error()})
	}
}
