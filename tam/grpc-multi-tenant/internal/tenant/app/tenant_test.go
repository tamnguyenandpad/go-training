package app

import (
	"context"
	"errors"
	"reflect"
	"testing"
	"time"

	domain "github.com/tuannguyenandpadcojp/go-training/tam/grpc-multi-tenant/internal/tenant/domain/tenant"
)

// Mock repository for testing
type mockTenantRepository struct {
	createFunc        func(ctx context.Context, tenant domain.Tenant) (*domain.Tenant, error)
	getTenantByIDFunc func(ctx context.Context, tenantID string) (*domain.Tenant, error)
}

func (m *mockTenantRepository) Create(ctx context.Context, tenant domain.Tenant) (*domain.Tenant, error) {
	return m.createFunc(ctx, tenant)
}

func (m *mockTenantRepository) GetTenantByID(ctx context.Context, tenantID string) (*domain.Tenant, error) {
	return m.getTenantByIDFunc(ctx, tenantID)
}

// Mock adapter required by application constructor
type mockUserAdapter struct{}

func (m *mockUserAdapter) GetUserByID(ctx context.Context, userID string) (*Output, error) {
	return nil, nil
}

// Additional methods required by the Repository interface
func (m *mockTenantRepository) CreateMember(ctx context.Context, member domain.Member) (*domain.Member, error) {
	return nil, nil
}

func (m *mockTenantRepository) GetMemberByUserID(ctx context.Context, userID string) (*domain.Member, error) {
	return nil, nil
}

func (m *mockTenantRepository) GetMemberById(ctx context.Context, memberID string) (*domain.Member, error) {
	return nil, nil
}

func (m *mockTenantRepository) UpdateMember(ctx context.Context, memberUpdateData domain.MemberUpdateData) (*domain.Member, error) {
	return nil, nil
}

func TestApplication_CreateTenant(t *testing.T) {
	// Setup mock repository
	mockRepo := &mockTenantRepository{
		createFunc: func(ctx context.Context, tenant domain.Tenant) (*domain.Tenant, error) {
			// Verify input
			if tenant.Name != "Test Tenant" {
				t.Errorf("Expected tenant name to be 'Test Tenant', got %s", tenant.Name)
			}
			if tenant.OwnerEmail != "owner@example.com" {
				t.Errorf("Expected owner email to be 'owner@example.com', got %s", tenant.OwnerEmail)
			}

			// Return success response
			return &domain.Tenant{
				ID:         tenant.ID,
				Name:       tenant.Name,
				OwnerEmail: tenant.OwnerEmail,
				CreatedAt:  tenant.CreatedAt,
			}, nil
		},
	}

	// Create application with mock dependencies
	app := NewApplication(mockRepo, &mockUserAdapter{})

	// Test the CreateTenant function
	input := CreateTenantInput{
		Name:       "Test Tenant",
		OwnerEmail: "owner@example.com",
	}

	result, err := app.CreateTenant(context.Background(), input)

	// Verify results
	if err != nil {
		t.Errorf("Expected no error, got %v", err)
	}
	if result == nil {
		t.Fatal("Expected result to not be nil")
	}
	if result.Name != "Test Tenant" {
		t.Errorf("Expected result name to be 'Test Tenant', got %s", result.Name)
	}
	if result.OwnerEmail != "owner@example.com" {
		t.Errorf("Expected result owner email to be 'owner@example.com', got %s", result.OwnerEmail)
	}
}

func TestApplication_CreateTenant_Error(t *testing.T) {
	// Setup mock repository that returns an error
	mockRepo := &mockTenantRepository{
		createFunc: func(ctx context.Context, tenant domain.Tenant) (*domain.Tenant, error) {
			return nil, errors.New("database error")
		},
	}

	// Create application with mock dependencies
	app := NewApplication(mockRepo, &mockUserAdapter{})

	// Test the CreateTenant function
	input := CreateTenantInput{
		Name:       "Test Tenant",
		OwnerEmail: "owner@example.com",
	}

	result, err := app.CreateTenant(context.Background(), input)

	// Verify results
	if err == nil {
		t.Error("Expected error, got nil")
	}
	if result != nil {
		t.Errorf("Expected nil result, got %+v", result)
	}
}

func TestApplication_GetTenantByID(t *testing.T) {
	createdAt := time.Now()
	// Setup mock repository
	mockRepo := &mockTenantRepository{
		getTenantByIDFunc: func(ctx context.Context, tenantID string) (*domain.Tenant, error) {
			// Verify input
			if tenantID != "tenant-123" {
				t.Errorf("Expected tenant ID to be 'tenant-123', got %s", tenantID)
			}

			// Return success response
			return &domain.Tenant{
				ID:         "tenant-123",
				Name:       "Test Tenant",
				OwnerEmail: "owner@example.com",
				CreatedAt:  createdAt,
			}, nil
		},
	}

	// Create application with mock dependencies
	app := NewApplication(mockRepo, &mockUserAdapter{})

	// Test the GetTenantByID function
	result, err := app.GetTenantByID(context.Background(), "tenant-123")

	// Verify results
	if err != nil {
		t.Errorf("Expected no error, got %v", err)
	}
	if result == nil {
		t.Fatal("Expected result to not be nil")
	}
	if result.ID != "tenant-123" {
		t.Errorf("Expected result ID to be 'tenant-123', got %s", result.ID)
	}
	if result.Name != "Test Tenant" {
		t.Errorf("Expected result name to be 'Test Tenant', got %s", result.Name)
	}
	if result.OwnerEmail != "owner@example.com" {
		t.Errorf("Expected result owner email to be 'owner@example.com', got %s", result.OwnerEmail)
	}
	if !reflect.DeepEqual(result.CreatedAt, createdAt) {
		t.Errorf("Expected result created at to be %v, got %v", createdAt, result.CreatedAt)
	}
}

func TestApplication_GetTenantByID_Error(t *testing.T) {
	// Setup mock repository that returns an error
	mockRepo := &mockTenantRepository{
		getTenantByIDFunc: func(ctx context.Context, tenantID string) (*domain.Tenant, error) {
			return nil, errors.New("tenant not found")
		},
	}

	// Create application with mock dependencies
	app := NewApplication(mockRepo, &mockUserAdapter{})

	// Test the GetTenantByID function
	result, err := app.GetTenantByID(context.Background(), "tenant-123")

	// Verify results
	if err == nil {
		t.Error("Expected error, got nil")
	}
	if result != nil {
		t.Errorf("Expected nil result, got %+v", result)
	}
}
