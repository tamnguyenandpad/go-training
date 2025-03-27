package main

import (
	"database/sql"
	"fmt"
	"log"
	"net"
	"os"

	_ "github.com/go-sql-driver/mysql"

	tenant_pb "github.com/tuannguyenandpadcojp/go-training/tam/grpc-multi-tenant/gen/go/tenant/v1"
	user_pb "github.com/tuannguyenandpadcojp/go-training/tam/grpc-multi-tenant/gen/go/user/v1"
	tenant_app "github.com/tuannguyenandpadcojp/go-training/tam/grpc-multi-tenant/internal/tenant/app"
	tenant_datastore "github.com/tuannguyenandpadcojp/go-training/tam/grpc-multi-tenant/internal/tenant/infra/datastore"
	tenant_ports "github.com/tuannguyenandpadcojp/go-training/tam/grpc-multi-tenant/internal/tenant/ports"
	user_app "github.com/tuannguyenandpadcojp/go-training/tam/grpc-multi-tenant/internal/user/app"
	user_datastore "github.com/tuannguyenandpadcojp/go-training/tam/grpc-multi-tenant/internal/user/infra/datastore"
	user_ports "github.com/tuannguyenandpadcojp/go-training/tam/grpc-multi-tenant/internal/user/ports"
	"google.golang.org/grpc"
	"google.golang.org/grpc/reflection"
)

func main() {
	//DB connection
	dbUser := os.Getenv("DB_USER")
	dbPassword := os.Getenv("DB_PASSWORD")
	dbHost := os.Getenv("DB_HOST")
	dbPort := os.Getenv("DB_PORT")
	dbName := os.Getenv("DB_NAME")

	dsn := fmt.Sprintf("%s:%s@tcp(%s:%s)/%s?parseTime=true", dbUser, dbPassword, dbHost, dbPort, dbName)
	db, err := sql.Open("mysql", dsn)
	if err != nil {
		log.Fatalf("Failed to connect to database: %v", err)
	}
	if err := db.Ping(); err != nil {
		log.Fatalf("Failed to ping to database: %v", err)
	}
	defer db.Close()

	// Create gRPC server
	listener, err := net.Listen("tcp", ":50051")
	if err != nil {
		log.Fatalf("Failed to listen: %v", err)
	}
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
	tenantGrpcService := tenant_ports.NewGrpcServer(tenantApplication)
	userGrpcService := user_ports.NewGrpcServer(userApplication)

	// Register service
	tenant_pb.RegisterTenantServiceServer(grpcServer, tenantGrpcService)
	user_pb.RegisterUserServiceServer(grpcServer, userGrpcService)

	if err := grpcServer.Serve(listener); err != nil {
		log.Fatalf("Failed to serve: %v", err)
	}
}
