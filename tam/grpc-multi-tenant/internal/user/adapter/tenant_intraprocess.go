package adapter

import (
	"context"

	"github.com/tuannguyenandpadcojp/go-training/tam/grpc-multi-tenant/internal/user/app"
	"github.com/tuannguyenandpadcojp/go-training/tam/grpc-multi-tenant/internal/user/domain/tenant"
)

type TenantAdapter struct {
	tenantRepo tenant.Repository
}

func NewTenantAdapter(tenantRepo tenant.Repository) app.TenantAdapter {
	return &TenantAdapter{tenantRepo: tenantRepo}
}

func (a *TenantAdapter) GetTenantByID(ctx context.Context, id string) (*app.GetTenantByIDOutput, error) {
	tenant, err := a.tenantRepo.GetTenantByID(ctx, id)
	if err != nil {
		return nil, err
	}
	return &app.GetTenantByIDOutput{
		ID:   tenant.ID,
		Name: tenant.Name,
	}, nil
}
