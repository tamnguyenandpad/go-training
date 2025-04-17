package adapter

import (
	"context"

	"github.com/tuannguyenandpadcojp/go-training/tam/grpc-multi-tenant/internal/tenant/app"
	"github.com/tuannguyenandpadcojp/go-training/tam/grpc-multi-tenant/internal/tenant/domain/user"
)

type UserAdapter struct {
	userRepo user.Repository
}

func NewUserAdapter(userRepo user.Repository) app.UserAdapter {
	return &UserAdapter{userRepo: userRepo}
}

func (a *UserAdapter) GetUserByID(ctx context.Context, userID string) (*app.Output, error) {
	user, err := a.userRepo.GetUserByID(ctx, userID)
	if err != nil {
		return nil, err
	}
	return &app.Output{
		ID:        user.ID,
		TenantID:  user.TenantID,
		Email:     user.Email,
		Name:      user.Name,
		CreatedAt: user.CreatedAt,
	}, nil
}
