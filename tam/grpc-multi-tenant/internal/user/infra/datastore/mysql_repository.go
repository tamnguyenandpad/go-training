package datastore

import (
	"context"
	"database/sql"
	"fmt"

	"github.com/tuannguyenandpadcojp/go-training/tam/grpc-multi-tenant/internal/user/domain"
)

type UserMysqlRepository struct {
	db *sql.DB
}

func NewUserMysqlRepository(db *sql.DB) *UserMysqlRepository {
	return &UserMysqlRepository{
		db: db,
	}
}

func (t *UserMysqlRepository) Create(ctx context.Context, user domain.User) (*domain.User, error) {
	query := "INSERT INTO users (id, tenant_id, email, name, created_at) VALUES (?, ?, ?, ?, ?)"

	_, err := t.db.ExecContext(ctx, query, user.ID, user.TenantID, user.Email, user.Name, user.CreatedAt)
	if err != nil {
		return nil, fmt.Errorf("failed to insert user: %v", err)
	}

	return &user, nil
}

func (t *UserMysqlRepository) GetUserByID(ctx context.Context, userID string) (*domain.User, error) {
	query := "SELECT id, tenant_id, email, name, created_at FROM users WHERE id = ?"
	row := t.db.QueryRowContext(ctx, query, userID)

	var user domain.User
	if err := row.Scan(&user.ID, &user.TenantID, &user.Email, &user.Name, &user.CreatedAt); err != nil {
		if err == sql.ErrNoRows {
			return nil, fmt.Errorf("user not found: %v", err)
		}
		return nil, fmt.Errorf("failed to fetch user: %v", err)
	}

	return &user, nil
}
