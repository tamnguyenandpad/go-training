package app

import (
	"context"
	"time"

	"github.com/google/uuid"
	"github.com/tuannguyenandpadcojp/go-training/tam/grpc-multi-tenant/internal/user/domain"
)

type Input struct {
	Name     string
	Email    string
	TenantID string
}

type Output struct {
	ID        string
	TenantID  string
	Email     string
	Name      string
	CreatedAt time.Time
}

type Application interface {
	CreateUser(ctx context.Context, input Input) (*Output, error)
}

type application struct {
	userRepo domain.Repository
}

func NewApplication(userRepo domain.Repository) Application {
	return &application{userRepo: userRepo}
}

func (a *application) CreateUser(ctx context.Context, input Input) (*Output, error) {
	createdAt := time.Now()
	userId := uuid.New().String()

	user := domain.User{
		ID:        userId,
		Name:      input.Name,
		Email:     input.Email,
		TenantID:  input.TenantID,
		CreatedAt: createdAt,
	}

	res, err := a.userRepo.Create(ctx, user)
	if err != nil {
		return nil, err
	}
	return &Output{
		ID:        res.ID,
		TenantID:  res.TenantID,
		Email:     res.Email,
		Name:      res.Name,
		CreatedAt: res.CreatedAt,
	}, nil

}
