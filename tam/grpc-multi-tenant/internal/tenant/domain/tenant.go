package domain

import "time"

type Tenant struct {
	ID         string
	Name       string
	OwnerEmail string
	CreatedAt  time.Time
}
