package domain

import "time"

type Member struct {
	ID         string
	TenantID   string
	UserID     string
	Status     string
	InvitedAt  *time.Time
	AcceptedAt *time.Time
}
