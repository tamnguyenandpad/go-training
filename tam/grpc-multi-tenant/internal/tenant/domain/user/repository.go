package user

import (
	"context"

	domain "github.com/tuannguyenandpadcojp/go-training/tam/grpc-multi-tenant/internal/user/domain/user"
)

type Repository interface {
	GetUserByID(ctx context.Context, userID string) (*domain.User, error)
}
