package app

import (
	"context"
	"fmt"
	"time"

	"github.com/google/uuid"
	domain_pkg "github.com/tuannguyenandpadcojp/go-training/tam/grpc-multi-tenant/internal/tenant/domain/pkg"
	domain "github.com/tuannguyenandpadcojp/go-training/tam/grpc-multi-tenant/internal/tenant/domain/tenant"
)

type CreateMemberInput struct {
	TenantID string
	UserID   string
}

type MemberOutput struct {
	ID         string
	TenantID   string
	UserID     string
	Status     domain_pkg.MemberStatus
	InvitedAt  *time.Time
	AcceptedAt *time.Time
}
type UpdateMemberInput struct {
	MemberID string
}

type UpdateMemberOutput struct {
	Status domain_pkg.MemberStatus
}

func (a *application) CreateMember(ctx context.Context, input CreateMemberInput) (*MemberOutput, error) {
	memberId := uuid.New().String()
	invitedAt := time.Now()

	if _, err := a.GetTenantByID(ctx, input.TenantID); err != nil {
		return nil, fmt.Errorf("tenant not found: %v", err)
	}
	if _, err := a.userAdapter.GetUserByID(ctx, input.UserID); err != nil {
		return nil, fmt.Errorf("user not found: %v", err)
	}
	if alreadyJoined := a.CheckUserAlreadyAMember(ctx, input.TenantID, input.UserID); alreadyJoined {
		return nil, fmt.Errorf("user already joined")
	}
	mem, err := a.GetMemberByUserID(ctx, input.UserID)
	if err != nil {
		return nil, nil
	}
	if mem.Status == domain_pkg.MemberStatusPending {
		return nil, fmt.Errorf("member is still pending")
	}

	member := domain.Member{
		ID:        memberId,
		TenantID:  input.TenantID,
		UserID:    input.UserID,
		Status:    domain_pkg.MemberStatusPending,
		InvitedAt: &invitedAt,
	}
	res, err := a.tenantRepo.CreateMember(ctx, member)
	if err != nil {
		return nil, err
	}
	return &MemberOutput{
		ID:         res.ID,
		TenantID:   res.TenantID,
		UserID:     res.UserID,
		Status:     res.Status,
		InvitedAt:  res.InvitedAt,
		AcceptedAt: res.AcceptedAt,
	}, nil
}

func (a *application) UpdateMember(ctx context.Context, input UpdateMemberInput) (*UpdateMemberOutput, error) {
	status := domain_pkg.MemberStatusAccepted
	acceptedAt := time.Now()

	member, err := a.GetMemberByID(ctx, input.MemberID)
	if err != nil {
		return nil, fmt.Errorf("member not found: %v", err)
	}
	if member.Status == domain_pkg.MemberStatusAccepted {
		return nil, fmt.Errorf("member already accepted")
	}

	memberUpdateData := domain.MemberUpdateData{
		MemberID:   input.MemberID,
		AcceptedAt: &acceptedAt,
		Status:     status,
	}
	res, err := a.tenantRepo.UpdateMember(ctx, memberUpdateData)
	if err != nil {
		return nil, err
	}
	return &UpdateMemberOutput{
		Status: res.Status,
	}, nil
}

func (a *application) GetMemberByUserID(ctx context.Context, userID string) (*MemberOutput, error) {
	member, err := a.tenantRepo.GetMemberByUserID(ctx, userID)
	if err != nil {
		return nil, err
	}
	return &MemberOutput{
		ID:         member.ID,
		TenantID:   member.TenantID,
		UserID:     member.UserID,
		Status:     member.Status,
		InvitedAt:  member.InvitedAt,
		AcceptedAt: member.AcceptedAt,
	}, nil
}

func (a *application) GetMemberByID(ctx context.Context, memberID string) (*MemberOutput, error) {
	member, err := a.tenantRepo.GetMemberById(ctx, memberID)
	if err != nil {
		return nil, err
	}
	return &MemberOutput{
		ID:         member.ID,
		TenantID:   member.TenantID,
		UserID:     member.UserID,
		Status:     member.Status,
		InvitedAt:  member.InvitedAt,
		AcceptedAt: member.AcceptedAt,
	}, nil
}

func (a *application) CheckUserAlreadyAMember(ctx context.Context, tenantID string, userID string) bool {
	member, err := a.tenantRepo.GetMemberByUserID(ctx, userID)
	if err != nil {
		return false
	}
	if member.TenantID == tenantID && member.Status == domain_pkg.MemberStatusAccepted {
		return true
	}
	return false
}
