package app

import (
	"context"

	"github.com/tuannguyenandpadcojp/go-training/tam/grpc-multi-tenant/internal/tenant/domain"
)

type Application interface {
	CreateTenant(ctx context.Context, input CreateTenantInput) (*CreateTenantOutput, error)
	CreateMember(ctx context.Context, input CreateMemberInput) (*MemberOutput, error)
	UpdateMember(ctx context.Context, input UpdateMemberInput) (*UpdateMemberOutput, error)
	GetTenantByID(ctx context.Context, id string) (*GetTenantByIDOutput, error)
	GetMemberByID(ctx context.Context, memberID string) (*MemberOutput, error)
	GetMemberByUserID(ctx context.Context, userID string) (*MemberOutput, error)
	CheckUserAlreadyAMember(ctx context.Context, tenantID string, userID string) bool
}

type application struct {
	tenantRepo domain.Repository
}

func NewApplication(tenantRepo domain.Repository) Application {
	return &application{tenantRepo: tenantRepo}
}
