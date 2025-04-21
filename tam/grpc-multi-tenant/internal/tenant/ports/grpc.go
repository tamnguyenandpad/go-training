package ports

import (
	"context"

	pb "github.com/tuannguyenandpadcojp/go-training/tam/grpc-multi-tenant/gen/go/tenant/v1"
	"github.com/tuannguyenandpadcojp/go-training/tam/grpc-multi-tenant/internal/pkg"
	pkgerrors "github.com/tuannguyenandpadcojp/go-training/tam/grpc-multi-tenant/internal/pkg/errors"
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
		return nil, pkgerrors.FromError(err)
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
		return pkgerrors.InvalidArgument("name and owner_email are required")
	}
	if !pkg.IsValidEmail(req.OwnerEmail) {
		return pkgerrors.InvalidArgument("invalid email format")
	}
	return nil
}

func (s *GrpcServer) InviteMember(ctx context.Context, req *pb.InviteMemberRequest) (*pb.InviteMemberResponse, error) {
	if req.TenantId == "" || req.UserId == "" {
		return nil, pkgerrors.InvalidArgument("tenant_id and user_id are required")
	}

	member, err := s.app.CreateMember(ctx, app.CreateMemberInput{
		TenantID: req.TenantId,
		UserID:   req.UserId,
	})
	if err != nil {
		// Handle specific error cases
		errMsg := err.Error()
		switch {
		case errMsg == "tenant not found":
			return nil, pkgerrors.NotFound("tenant", req.TenantId)
		case errMsg == "user not found":
			return nil, pkgerrors.NotFound("user", req.UserId)
		case errMsg == "user already joined":
			return nil, pkgerrors.AlreadyExists("user is already a member of this tenant")
		case errMsg == "member is still pending":
			return nil, pkgerrors.FailedPrecondition("member invitation is still pending")
		default:
			return nil, pkgerrors.InternalError(errMsg)
		}
	}

	return &pb.InviteMemberResponse{
		MemberId: member.ID,
	}, nil
}

func (s *GrpcServer) AcceptInvitation(ctx context.Context, req *pb.AcceptInvitationRequest) (*pb.AcceptInvitationResponse, error) {
	if req.MemberId == "" {
		return nil, pkgerrors.InvalidArgument("member_id is required")
	}

	member, err := s.app.UpdateMember(ctx, app.UpdateMemberInput{
		MemberID: req.MemberId,
	})
	if err != nil {
		// Handle specific error cases
		errMsg := err.Error()
		switch {
		case errMsg == "member not found":
			return nil, pkgerrors.NotFound("member", req.MemberId)
		case errMsg == "member already accepted":
			return nil, pkgerrors.AlreadyExists("invitation has already been accepted")
		default:
			return nil, pkgerrors.InternalError(errMsg)
		}
	}

	return &pb.AcceptInvitationResponse{
		Status: string(member.Status),
	}, nil
}
