package app

import (
	"context"
	"time"

	"github.com/google/uuid"
	"github.com/tuannguyenandpadcojp/go-training/tam/grpc-multi-tenant/internal/tenant/domain"
)

type CreateTenantInput struct {
	Name       string
	OwnerEmail string
}

type CreateTenantOutput struct {
	ID         string
	Name       string
	OwnerEmail string
	CreatedAt  time.Time
}

type GetTenantByIDOutput struct {
	ID         string
	Name       string
	OwnerEmail string
	CreatedAt  time.Time
}

func (a *application) CreateTenant(ctx context.Context, input CreateTenantInput) (*CreateTenantOutput, error) {
	tenantID := uuid.New().String()
	createdAt := time.Now()

	tenant := domain.Tenant{ID: tenantID, Name: input.Name, OwnerEmail: input.OwnerEmail, CreatedAt: createdAt}
	res, err := a.tenantRepo.Create(ctx, tenant)
	if err != nil {
		return nil, err
	}
	return &CreateTenantOutput{
		ID:         res.ID,
		Name:       res.Name,
		OwnerEmail: res.OwnerEmail,
		CreatedAt:  res.CreatedAt,
	}, nil
}

func (a *application) GetTenantByID(ctx context.Context, tenantID string) (*GetTenantByIDOutput, error) {
	tenant, err := a.tenantRepo.GetTenantByID(ctx, tenantID)
	if err != nil {
		return nil, err
	}
	return &GetTenantByIDOutput{
		ID:         tenant.ID,
		Name:       tenant.Name,
		OwnerEmail: tenant.OwnerEmail,
		CreatedAt:  tenant.CreatedAt,
	}, nil
}
