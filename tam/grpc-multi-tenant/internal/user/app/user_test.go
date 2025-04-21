package app

import (
	"context"
	"errors"
	"testing"
	"time"

	"github.com/google/uuid"
	domain "github.com/tuannguyenandpadcojp/go-training/tam/grpc-multi-tenant/internal/user/domain/user"
)

// Mock repository for testing
type mockUserRepository struct {
	createFunc      func(ctx context.Context, user domain.User) (*domain.User, error)
	getUserByIDFunc func(ctx context.Context, userID string) (*domain.User, error)
}

func (m *mockUserRepository) Create(ctx context.Context, user domain.User) (*domain.User, error) {
	return m.createFunc(ctx, user)
}

func (m *mockUserRepository) GetUserByID(ctx context.Context, userID string) (*domain.User, error) {
	return m.getUserByIDFunc(ctx, userID)
}

// Mock tenant adapter for testing
type mockTenantAdapter struct {
	getTenantByIDFunc func(ctx context.Context, id string) (*GetTenantByIDOutput, error)
}

func (m *mockTenantAdapter) GetTenantByID(ctx context.Context, id string) (*GetTenantByIDOutput, error) {
	return m.getTenantByIDFunc(ctx, id)
}

func TestApplication_CreateUser_Success(t *testing.T) {
	// Define test time
	now := time.Now()

	// Setup mocks
	mockRepo := &mockUserRepository{
		createFunc: func(ctx context.Context, user domain.User) (*domain.User, error) {
			// Verify input
			if user.Name != "Test User" {
				t.Errorf("Expected user name to be 'Test User', got %s", user.Name)
			}
			if user.Email != "user@example.com" {
				t.Errorf("Expected user email to be 'user@example.com', got %s", user.Email)
			}
			if user.TenantID != "tenant-123" {
				t.Errorf("Expected tenant ID to be 'tenant-123', got %s", user.TenantID)
			}

			// Return success response
			return &domain.User{
				ID:        user.ID,
				Name:      user.Name,
				Email:     user.Email,
				TenantID:  user.TenantID,
				CreatedAt: now,
			}, nil
		},
	}

	mockTenantAdapter := &mockTenantAdapter{
		getTenantByIDFunc: func(ctx context.Context, id string) (*GetTenantByIDOutput, error) {
			// Return a tenant to simulate it exists
			return &GetTenantByIDOutput{
				ID:   id,
				Name: "Test Tenant",
			}, nil
		},
	}

	// Create application with mock dependencies
	app := NewApplication(mockRepo, mockTenantAdapter)

	// Test the CreateUser function
	input := Input{
		Name:     "Test User",
		Email:    "user@example.com",
		TenantID: "tenant-123",
	}

	result, err := app.CreateUser(context.Background(), input)

	// Verify results
	if err != nil {
		t.Errorf("Expected no error, got %v", err)
	}
	if result == nil {
		t.Fatal("Expected result to not be nil")
	}
	if result.Name != "Test User" {
		t.Errorf("Expected result name to be 'Test User', got %s", result.Name)
	}
	if result.Email != "user@example.com" {
		t.Errorf("Expected result email to be 'user@example.com', got %s", result.Email)
	}
	if result.TenantID != "tenant-123" {
		t.Errorf("Expected result tenant ID to be 'tenant-123', got %s", result.TenantID)
	}
}

func TestApplication_CreateUser_TenantNotFound(t *testing.T) {
	// Setup mocks
	mockRepo := &mockUserRepository{
		createFunc: func(ctx context.Context, user domain.User) (*domain.User, error) {
			// This should not be called
			t.Error("Create should not be called when tenant is not found")
			return nil, nil
		},
	}

	mockTenantAdapter := &mockTenantAdapter{
		getTenantByIDFunc: func(ctx context.Context, id string) (*GetTenantByIDOutput, error) {
			// Return error to simulate tenant not found
			return nil, errors.New("tenant not found")
		},
	}

	// Create application with mock dependencies
	app := NewApplication(mockRepo, mockTenantAdapter)

	// Test the CreateUser function with nonexistent tenant
	input := Input{
		Name:     "Test User",
		Email:    "user@example.com",
		TenantID: "nonexistent-tenant",
	}

	result, err := app.CreateUser(context.Background(), input)

	// Verify results
	if err == nil {
		t.Error("Expected error for nonexistent tenant, got nil")
	}
	if result != nil {
		t.Errorf("Expected nil result, got %+v", result)
	}
}

func TestApplication_CreateUser_RepositoryError(t *testing.T) {
	// Setup mocks
	mockRepo := &mockUserRepository{
		createFunc: func(ctx context.Context, user domain.User) (*domain.User, error) {
			// Return error to simulate database error
			return nil, errors.New("database error")
		},
	}

	mockTenantAdapter := &mockTenantAdapter{
		getTenantByIDFunc: func(ctx context.Context, id string) (*GetTenantByIDOutput, error) {
			// Return a tenant to simulate it exists
			return &GetTenantByIDOutput{
				ID:   id,
				Name: "Test Tenant",
			}, nil
		},
	}

	// Create application with mock dependencies
	app := NewApplication(mockRepo, mockTenantAdapter)

	// Test the CreateUser function with repository that returns error
	input := Input{
		Name:     "Test User",
		Email:    "user@example.com",
		TenantID: "tenant-123",
	}

	result, err := app.CreateUser(context.Background(), input)

	// Verify results
	if err == nil {
		t.Error("Expected error from repository, got nil")
	}
	if result != nil {
		t.Errorf("Expected nil result, got %+v", result)
	}
}

func TestApplication_GetUserByID_Success(t *testing.T) {
	// Define test time
	now := time.Now()
	userId := uuid.New().String()

	// Setup mock repository
	mockRepo := &mockUserRepository{
		getUserByIDFunc: func(ctx context.Context, userID string) (*domain.User, error) {
			// Verify input
			if userID != userId {
				t.Errorf("Expected user ID to be '%s', got %s", userId, userID)
			}

			// Return success response
			return &domain.User{
				ID:        userId,
				Name:      "Test User",
				Email:     "user@example.com",
				TenantID:  "tenant-123",
				CreatedAt: now,
			}, nil
		},
	}

	// Create application with mock dependencies
	app := NewApplication(mockRepo, &mockTenantAdapter{})

	// Test the GetUserByID function
	result, err := app.GetUserByID(context.Background(), userId)

	// Verify results
	if err != nil {
		t.Errorf("Expected no error, got %v", err)
	}
	if result == nil {
		t.Fatal("Expected result to not be nil")
	}
	if result.ID != userId {
		t.Errorf("Expected result ID to be '%s', got %s", userId, result.ID)
	}
	if result.Name != "Test User" {
		t.Errorf("Expected result name to be 'Test User', got %s", result.Name)
	}
	if result.Email != "user@example.com" {
		t.Errorf("Expected result email to be 'user@example.com', got %s", result.Email)
	}
	if result.TenantID != "tenant-123" {
		t.Errorf("Expected result tenant ID to be 'tenant-123', got %s", result.TenantID)
	}
}

func TestApplication_GetUserByID_NotFound(t *testing.T) {
	// Setup mock repository
	mockRepo := &mockUserRepository{
		getUserByIDFunc: func(ctx context.Context, userID string) (*domain.User, error) {
			// Return error to simulate user not found
			return nil, errors.New("user not found")
		},
	}

	// Create application with mock dependencies
	app := NewApplication(mockRepo, &mockTenantAdapter{})

	// Test the GetUserByID function with nonexistent user
	result, err := app.GetUserByID(context.Background(), "nonexistent-user")

	// Verify results
	if err == nil {
		t.Error("Expected error for nonexistent user, got nil")
	}
	if result != nil {
		t.Errorf("Expected nil result, got %+v", result)
	}
}
