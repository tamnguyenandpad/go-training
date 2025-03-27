package domain

import "context"

type Repository interface {
	Create(ctx context.Context, user User) (*User, error)
}
