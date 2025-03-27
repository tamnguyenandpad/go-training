package app

import (
	"context"

	"github.com/tuannguyenandpadcojp/go-training/tam/grpc-multi-tenant/internal/tenant/domain"
)

type Application interface {
	CreateTenant(ctx context.Context, input CreateTenantInput) (*CreateTenantOutput, error)
	CreateMember(ctx context.Context, input CreateMemberInput) (*CreateMemberOutput, error)
	UpdateMember(ctx context.Context, input UpdateMemberInput) (*UpdateMemberOutput, error)
}

type appplication struct {
	tenantRepo domain.Repository
}

func NewApplication(tenantRepo domain.Repository) Application {
	return &appplication{tenantRepo: tenantRepo}
}
