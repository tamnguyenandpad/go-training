package domain

import (
	"time"

	domain_pkg "github.com/tuannguyenandpadcojp/go-training/tam/grpc-multi-tenant/internal/tenant/domain/pkg"
)

type Member struct {
	ID         string
	TenantID   string
	UserID     string
	Status     domain_pkg.MemberStatus
	InvitedAt  *time.Time
	AcceptedAt *time.Time
}
