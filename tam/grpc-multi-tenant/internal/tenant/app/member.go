package app

import (
	"context"
	"time"

	"github.com/google/uuid"
	"github.com/tuannguyenandpadcojp/go-training/tam/grpc-multi-tenant/internal/tenant/domain"
)

type CreateMemberInput struct {
	TenantID string
	UserID   string
}

type MemberOutput struct {
	ID         string
	TenantID   string
	UserID     string
	Status     string
	InvitedAt  *time.Time
	AcceptedAt *time.Time
}
type UpdateMemberInput struct {
	MemberID string
}

type UpdateMemberOutput struct {
	Status string
}

func (a *application) CreateMember(ctx context.Context, input CreateMemberInput) (*MemberOutput, error) {
	memberId := uuid.New().String()
	invitedAt := time.Now()

	member := domain.Member{
		ID:        memberId,
		TenantID:  input.TenantID,
		UserID:    input.UserID,
		Status:    "pending",
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
	status := "accepted"
	acceptedAt := time.Now()

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
	if member.TenantID == tenantID && member.Status == "accepted" {
		return true
	}
	return false
}
