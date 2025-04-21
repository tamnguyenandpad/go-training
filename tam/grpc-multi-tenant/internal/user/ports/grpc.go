package ports

import (
	"context"

	pb "github.com/tuannguyenandpadcojp/go-training/tam/grpc-multi-tenant/gen/go/user/v1"
	"github.com/tuannguyenandpadcojp/go-training/tam/grpc-multi-tenant/internal/pkg"
	pkgerrors "github.com/tuannguyenandpadcojp/go-training/tam/grpc-multi-tenant/internal/pkg/errors"
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
	if err := s.validateCreateUserRequest(req); err != nil {
		return nil, err
	}

	user, err := s.app.CreateUser(ctx, app.Input{
		Name:     req.Name,
		Email:    req.Email,
		TenantID: req.TenantId,
	})
	if err != nil {
		errMsg := err.Error()
		if errMsg == "tenant not found" {
			return nil, pkgerrors.NotFound("tenant", req.TenantId)
		}
		return nil, pkgerrors.FromError(err)
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
		return pkgerrors.InvalidArgument("tenant_id, email, and name are required")
	}
	if !pkg.IsValidEmail(req.Email) {
		return pkgerrors.InvalidArgument("invalid email format")
	}
	return nil
}
