package ports

import (
	"context"
	"fmt"

	pb "github.com/tuannguyenandpadcojp/go-training/tam/grpc-multi-tenant/gen/go/user/v1"
	"github.com/tuannguyenandpadcojp/go-training/tam/grpc-multi-tenant/internal/pkg"
	tenant_app "github.com/tuannguyenandpadcojp/go-training/tam/grpc-multi-tenant/internal/tenant/app"
	"github.com/tuannguyenandpadcojp/go-training/tam/grpc-multi-tenant/internal/user/app"
)

type GrpcServer struct {
	pb.UnimplementedUserServiceServer
	app       app.Application
	tenantApp tenant_app.Application
}

func NewGrpcServer(app app.Application, tenantApp tenant_app.Application) pb.UserServiceServer {
	return &GrpcServer{app: app, tenantApp: tenantApp}
}

func (s *GrpcServer) CreateUser(ctx context.Context, req *pb.CreateUserRequest) (*pb.CreateUserResponse, error) {
	if err := s.validateCreateUserRequest(req); err != nil {
		return nil, err
	}

	user, err := s.app.CreateUser(ctx, app.Input{
		Name:     req.Name,
		Email:    req.Email,
		TenantID: req.TenantId,
	})
	if err != nil {
		return nil, err
	}

	return &pb.CreateUserResponse{
		Id:        user.ID,
		TenantId:  user.TenantID,
		Email:     user.Email,
		Name:      user.Name,
		CreatedAt: user.CreatedAt.Format("2006-01-02T15:04:05Z"),
	}, nil
}

func (s *GrpcServer) validateCreateUserRequest(req *pb.CreateUserRequest) error {
	if req.TenantId == "" || req.Email == "" || req.Name == "" {
		return fmt.Errorf("name tenantId, email, name are required")
	}
	if !pkg.EmailRegex.MatchString(req.Email) {
		return fmt.Errorf("invalid email format")
	}
	if _, err := s.tenantApp.GetTenantByID(context.Background(), req.TenantId); err != nil {
		return fmt.Errorf("tenant not found")
	}
	return nil
}
