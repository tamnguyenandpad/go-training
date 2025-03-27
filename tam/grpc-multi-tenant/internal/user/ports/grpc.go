package ports

import (
	"context"
	"fmt"

	pb "github.com/tuannguyenandpadcojp/go-training/tam/grpc-multi-tenant/gen/go/user/v1"
	"github.com/tuannguyenandpadcojp/go-training/tam/grpc-multi-tenant/internal/user/app"
)

type GrpcServer struct {
	pb.UnimplementedUserServiceServer
	app app.Application
}

func NewGrpcServer(app app.Application) pb.UserServiceServer {
	return &GrpcServer{app: app}
}

func (s *GrpcServer) CreateUser(ctx context.Context, req *pb.CreateUserRequest) (*pb.CreateUserResponse, error) {
	if req.TenantId == "" || req.Email == "" || req.Name == "" {
		return nil, fmt.Errorf("name tenantId, email, name are required")
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
