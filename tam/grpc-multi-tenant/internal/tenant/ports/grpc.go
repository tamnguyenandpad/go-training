package ports

import (
	"context"
	"fmt"

	pb "github.com/tuannguyenandpadcojp/go-training/tam/grpc-multi-tenant/gen/go/tenant/v1"
	"github.com/tuannguyenandpadcojp/go-training/tam/grpc-multi-tenant/internal/pkg"
	"github.com/tuannguyenandpadcojp/go-training/tam/grpc-multi-tenant/internal/tenant/app"
)

type GrpcServer struct {
	pb.UnimplementedTenantServiceServer
	app app.Application
}

func NewGrpcServer(app app.Application) pb.TenantServiceServer {
	return &GrpcServer{app: app}
}

func (s *GrpcServer) CreateTenant(ctx context.Context, req *pb.CreateTenantRequest) (*pb.CreateTenantResponse, error) {
	if err := s.validateCreateTenantRequest(req); err != nil {
		return nil, err
	}

	tenant, err := s.app.CreateTenant(ctx, app.CreateTenantInput{
		Name:       req.Name,
		OwnerEmail: req.OwnerEmail,
	})
	if err != nil {
		return nil, err
	}

	return &pb.CreateTenantResponse{
		Id:         tenant.ID,
		Name:       tenant.Name,
		OwnerEmail: tenant.OwnerEmail,
		CreatedAt:  tenant.CreatedAt.Format("2006-01-02T15:04:05Z"),
	}, nil
}

func (s *GrpcServer) validateCreateTenantRequest(req *pb.CreateTenantRequest) error {
	if req.Name == "" || req.OwnerEmail == "" {
		return fmt.Errorf("name and owner_email are required")
	}
	if !pkg.IsValidEmail(req.OwnerEmail) {
		return fmt.Errorf("invalid email format")
	}
	return nil
}

func (s *GrpcServer) InviteMember(ctx context.Context, req *pb.InviteMemberRequest) (*pb.InviteMemberResponse, error) {
	if req.TenantId == "" || req.UserId == "" {
		return nil, fmt.Errorf("tenant_id and user_id are required")
	}

	member, err := s.app.CreateMember(ctx, app.CreateMemberInput{
		TenantID: req.TenantId,
		UserID:   req.UserId,
	})
	if err != nil {
		return nil, err
	}

	return &pb.InviteMemberResponse{
		MemberId: member.ID,
	}, nil
}

func (s *GrpcServer) AcceptInvitation(ctx context.Context, req *pb.AcceptInvitationRequest) (*pb.AcceptInvitationResponse, error) {
	if req.MemberId == "" {
		return nil, fmt.Errorf("member_id are required")
	}

	member, err := s.app.UpdateMember(ctx, app.UpdateMemberInput{
		MemberID: req.MemberId,
	})
	if err != nil {
		return nil, err
	}

	return &pb.AcceptInvitationResponse{
		Status: string(member.Status),
	}, nil
}
