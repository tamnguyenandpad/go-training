package domain

import (
	"context"
	"time"
)

type MemberUpdateData struct {
	MemberID   string
	AcceptedAt *time.Time
	Status     string
}

type Repository interface {
	Create(ctx context.Context, tenant Tenant) (*Tenant, error)
	GetTenantByID(ctx context.Context, tenantID string) (*Tenant, error)
	CreateMember(ctx context.Context, member Member) (*Member, error)
	UpdateMember(ctx context.Context, memberUpdateData MemberUpdateData) (*Member, error)
	GetMemberById(ctx context.Context, memberID string) (*Member, error)
	GetMemberByUserID(ctx context.Context, userID string) (*Member, error)
}
