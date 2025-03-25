package server

import (
	"context"
	"database/sql"
	"fmt"
	"time"

	"github.com/google/uuid"
	pb "github.com/tuannguyenandpadcojp/go-training/tam/grpc-multi-tenant/gen/go/tenant/v1"
)

type TenantServiceServer struct {
	pb.UnimplementedTenantServiceServer
	DB *sql.DB
}

func (s *TenantServiceServer) CreateTenant(ctx context.Context, req *pb.CreateTenantRequest) (*pb.CreateTenantResponse, error) {
	// Validate input
	if req.Name == "" || req.OwnerEmail == "" {
		return nil, fmt.Errorf("name and owner_email are required")
	}

	// Prepare query
	query := "INSERT INTO tenants (id, name, owner_email, created_at) VALUES (?, ?, ?, ?)"

	tenantID := uuid.New().String()
	createdAt := time.Now()

	// Execute insert query
	_, err := s.DB.ExecContext(ctx, query, tenantID, req.Name, req.OwnerEmail, createdAt)
	if err != nil {
		return nil, fmt.Errorf("failed to insert tenant: %v", err)
	}

	return &pb.CreateTenantResponse{Id: tenantID}, nil
}
