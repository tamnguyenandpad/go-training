package app

import (
	"context"
	"fmt"
	"time"

	"github.com/google/uuid"
	domain "github.com/tuannguyenandpadcojp/go-training/tam/grpc-multi-tenant/internal/user/domain/user"
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

type GetTenantByIDOutput struct {
	ID         string
	Name       string
	OwnerEmail string
	CreatedAt  time.Time
}

type TenantAdapter interface {
	GetTenantByID(ctx context.Context, id string) (*GetTenantByIDOutput, error)
}

type Application interface {
	CreateUser(ctx context.Context, input Input) (*Output, error)
	GetUserByID(ctx context.Context, userID string) (*Output, error)
}

type application struct {
	userRepo      domain.Repository
	tenantAdapter TenantAdapter
}

func NewApplication(userRepo domain.Repository, tenantAdapter TenantAdapter) Application {
	return &application{userRepo: userRepo, tenantAdapter: tenantAdapter}
}

func (a *application) CreateUser(ctx context.Context, input Input) (*Output, error) {
	createdAt := time.Now()
	userId := uuid.New().String()

	if _, err := a.tenantAdapter.GetTenantByID(context.Background(), input.TenantID); err != nil {
		return nil, fmt.Errorf("tenant not found")
	}

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

func (a *application) GetUserByID(ctx context.Context, userID string) (*Output, error) {
	user, err := a.userRepo.GetUserByID(ctx, userID)
	if err != nil {
		return nil, err
	}
	return &Output{
		ID:        user.ID,
		TenantID:  user.TenantID,
		Email:     user.Email,
		Name:      user.Name,
		CreatedAt: user.CreatedAt,
	}, nil
}
