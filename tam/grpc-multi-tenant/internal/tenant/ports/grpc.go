package ports

import (
	"context"
	"fmt"

	pb "github.com/tuannguyenandpadcojp/go-training/tam/grpc-multi-tenant/gen/go/tenant/v1"
	"github.com/tuannguyenandpadcojp/go-training/tam/grpc-multi-tenant/internal/pkg"
	"github.com/tuannguyenandpadcojp/go-training/tam/grpc-multi-tenant/internal/tenant/app"
	user_app "github.com/tuannguyenandpadcojp/go-training/tam/grpc-multi-tenant/internal/user/app"
)

type GrpcServer struct {
	pb.UnimplementedTenantServiceServer
	app     app.Application
	userApp user_app.Application
}

func NewGrpcServer(app app.Application, userApp user_app.Application) pb.TenantServiceServer {
	return &GrpcServer{app: app, userApp: userApp}
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
	if !pkg.EmailRegex.MatchString(req.OwnerEmail) {
		return fmt.Errorf("invalid email format")
	}
	return nil
}

func (s *GrpcServer) InviteMember(ctx context.Context, req *pb.InviteMemberRequest) (*pb.InviteMemberResponse, error) {
	if err := s.validateInviteMemberRequest(ctx, req); err != nil {
		return nil, err
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

func (s *GrpcServer) validateInviteMemberRequest(ctx context.Context, req *pb.InviteMemberRequest) error {
	if req.TenantId == "" || req.UserId == "" {
		return fmt.Errorf("tenant_id and user_id are required")
	}
	if _, err := s.app.GetTenantByID(ctx, req.TenantId); err != nil {
		return fmt.Errorf("tenant not found: %v", err)
	}
	if _, err := s.userApp.GetUserByID(ctx, req.UserId); err != nil {
		return fmt.Errorf("user not found: %v", err)
	}
	if alreadyJoined := s.app.CheckUserAlreadyAMember(ctx, req.TenantId, req.UserId); alreadyJoined {
		return fmt.Errorf("user already joined")
	}
	member, err := s.app.GetMemberByUserID(ctx, req.UserId)
	if err != nil {
		return nil
	}
	if member.Status == "pending" {
		return fmt.Errorf("member is still pending")
	}

	return nil
}

func (s *GrpcServer) AcceptInvitation(ctx context.Context, req *pb.AcceptInvitationRequest) (*pb.AcceptInvitationResponse, error) {
	if err := s.validateAcceptInvitationRequest(ctx, req); err != nil {
		return nil, err
	}

	member, err := s.app.UpdateMember(ctx, app.UpdateMemberInput{
		MemberID: req.MemberId,
	})
	if err != nil {
		return nil, err
	}

	return &pb.AcceptInvitationResponse{
		Status: member.Status,
	}, nil
}

func (s *GrpcServer) validateAcceptInvitationRequest(ctx context.Context, req *pb.AcceptInvitationRequest) error {
	if req.MemberId == "" {
		return fmt.Errorf("member_id are required")
	}
	member, err := s.app.GetMemberByID(ctx, req.MemberId)
	if err != nil {
		return fmt.Errorf("member not found: %v", err)
	}
	if member.Status == "accepted" {
		return fmt.Errorf("member already accepted")
	}
	return nil
}
