package tenant

import (
	"context"
	"time"

	domain "github.com/tuannguyenandpadcojp/go-training/tam/grpc-multi-tenant/internal/tenant/domain/tenant"
)

type MemberUpdateData struct {
	MemberID   string
	AcceptedAt *time.Time
	Status     string
}

type Repository interface {
	GetTenantByID(ctx context.Context, tenantID string) (*domain.Tenant, error)
}
