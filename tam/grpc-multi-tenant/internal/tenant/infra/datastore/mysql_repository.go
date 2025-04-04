package datastore

import (
	"context"
	"database/sql"
	"fmt"

	"github.com/tuannguyenandpadcojp/go-training/tam/grpc-multi-tenant/internal/tenant/domain"
)

type TenantMysqlRepository struct {
	db *sql.DB
}

func NewTenantMysqlRepository(db *sql.DB) *TenantMysqlRepository {
	return &TenantMysqlRepository{
		db: db,
	}
}

func (t *TenantMysqlRepository) Create(ctx context.Context, tenant domain.Tenant) (*domain.Tenant, error) {
	query := "INSERT INTO tenants (id, name, owner_email, created_at) VALUES (?, ?, ?, ?)"

	_, err := t.db.ExecContext(ctx, query, tenant.ID, tenant.Name, tenant.OwnerEmail, tenant.CreatedAt)
	if err != nil {
		return nil, fmt.Errorf("failed to insert tenant: %v", err)
	}

	return &tenant, nil
}

func (t *TenantMysqlRepository) CreateMember(ctx context.Context, member domain.Member) (*domain.Member, error) {
	query := "INSERT INTO members (id, tenant_id, user_id, status, invited_at) VALUES (?, ?, ?, ?, ?)"

	_, err := t.db.ExecContext(ctx, query, member.ID, member.TenantID, member.UserID, member.Status, member.InvitedAt)
	if err != nil {
		return nil, fmt.Errorf("failed to insert member: %v", err)
	}

	return &member, nil
}

func (t *TenantMysqlRepository) GetMemberByUserID(ctx context.Context, userID string) (*domain.Member, error) {
	query := "SELECT id, tenant_id, user_id, status, invited_at, accepted_at FROM members WHERE user_id = ?"
	row := t.db.QueryRowContext(ctx, query, userID)
	var member domain.Member
	if err := row.Scan(&member.ID, &member.TenantID, &member.UserID, &member.Status, &member.InvitedAt, &member.AcceptedAt); err != nil {
		if err == sql.ErrNoRows {
			return nil, fmt.Errorf("member not found: %v", err)
		}
		return nil, fmt.Errorf("failed to fetch member: %v", err)
	}
	return &member, nil
}

func (t *TenantMysqlRepository) UpdateMember(ctx context.Context, memberUpdateData domain.MemberUpdateData) (*domain.Member, error) {
	updateQuery := "UPDATE members SET status = ?, accepted_at = ? WHERE id = ?"
	_, err := t.db.ExecContext(ctx, updateQuery, memberUpdateData.Status, memberUpdateData.AcceptedAt, memberUpdateData.MemberID)
	if err != nil {
		return nil, fmt.Errorf("failed to update member: %v", err)
	}

	selectQuery := "SELECT id, tenant_id, user_id, status, invited_at, accepted_at FROM members WHERE id = ?"
	row := t.db.QueryRowContext(ctx, selectQuery, memberUpdateData.MemberID)

	var updatedMember domain.Member
	if err := row.Scan(&updatedMember.ID, &updatedMember.TenantID, &updatedMember.UserID, &updatedMember.Status, &updatedMember.InvitedAt, &updatedMember.AcceptedAt); err != nil {
		return nil, fmt.Errorf("failed to fetch updated member: %v", err)
	}

	return &updatedMember, nil
}

func (t *TenantMysqlRepository) GetTenantByID(ctx context.Context, tenantID string) (*domain.Tenant, error) {
	query := "SELECT id, name, owner_email, created_at FROM tenants WHERE id = ?"
	row := t.db.QueryRowContext(ctx, query, tenantID)

	var tenant domain.Tenant
	if err := row.Scan(&tenant.ID, &tenant.Name, &tenant.OwnerEmail, &tenant.CreatedAt); err != nil {
		if err == sql.ErrNoRows {
			return nil, fmt.Errorf("tenant not found from repository: %v", err)
		}
		return nil, fmt.Errorf("failed to fetch tenant: %v", err)
	}

	return &tenant, nil
}

func (t *TenantMysqlRepository) GetMemberById(ctx context.Context, memberID string) (*domain.Member, error) {
	query := "SELECT id, tenant_id, user_id, status, invited_at, accepted_at FROM members WHERE id = ?"
	row := t.db.QueryRowContext(ctx, query, memberID)

	var member domain.Member
	if err := row.Scan(&member.ID, &member.TenantID, &member.UserID, &member.Status, &member.InvitedAt, &member.AcceptedAt); err != nil {
		if err == sql.ErrNoRows {
			return nil, fmt.Errorf("member not found: %v", err)
		}
		return nil, fmt.Errorf("failed to fetch member: %v", err)
	}

	return &member, nil
}
