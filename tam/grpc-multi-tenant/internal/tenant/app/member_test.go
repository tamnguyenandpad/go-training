package app

import (
	"context"
	"errors"
	"fmt"
	"testing"
	"time"

	domain_pkg "github.com/tuannguyenandpadcojp/go-training/tam/grpc-multi-tenant/internal/tenant/domain/pkg"
	domain "github.com/tuannguyenandpadcojp/go-training/tam/grpc-multi-tenant/internal/tenant/domain/tenant"
)

// Mock repository for member-related tests
type mockMemberRepository struct {
	mockTenantRepository
	createMemberFunc      func(ctx context.Context, member domain.Member) (*domain.Member, error)
	getMemberByUserIDFunc func(ctx context.Context, userID string) (*domain.Member, error)
	getMemberByIdFunc     func(ctx context.Context, memberID string) (*domain.Member, error)
	updateMemberFunc      func(ctx context.Context, memberUpdateData domain.MemberUpdateData) (*domain.Member, error)
}

func (m *mockMemberRepository) CreateMember(ctx context.Context, member domain.Member) (*domain.Member, error) {
	return m.createMemberFunc(ctx, member)
}

func (m *mockMemberRepository) GetMemberByUserID(ctx context.Context, userID string) (*domain.Member, error) {
	return m.getMemberByUserIDFunc(ctx, userID)
}

func (m *mockMemberRepository) GetMemberById(ctx context.Context, memberID string) (*domain.Member, error) {
	return m.getMemberByIdFunc(ctx, memberID)
}

func (m *mockMemberRepository) UpdateMember(ctx context.Context, memberUpdateData domain.MemberUpdateData) (*domain.Member, error) {
	return m.updateMemberFunc(ctx, memberUpdateData)
}

// Mock user adapter for member-related tests
type mockUserAdapterForMember struct {
	getUserByIDFunc func(ctx context.Context, userID string) (*Output, error)
}

func (m *mockUserAdapterForMember) GetUserByID(ctx context.Context, userID string) (*Output, error) {
	return m.getUserByIDFunc(ctx, userID)
}

func TestApplication_CreateMember_Success(t *testing.T) {
	// Define test data
	invitedAt := time.Now()
	testMember := domain.Member{
		ID:        "member-123",
		TenantID:  "tenant-123",
		UserID:    "user-123",
		Status:    domain_pkg.MemberStatusPending,
		InvitedAt: &invitedAt,
	}

	// Setup mocks
	mockRepo := &mockMemberRepository{
		mockTenantRepository: mockTenantRepository{
			getTenantByIDFunc: func(ctx context.Context, tenantID string) (*domain.Tenant, error) {
				// Return a tenant to simulate it exists
				return &domain.Tenant{ID: tenantID}, nil
			},
		},
		getMemberByUserIDFunc: func(ctx context.Context, userID string) (*domain.Member, error) {
			return &domain.Member{
				ID:        "member-123",
				TenantID:  "tenant-123",
				UserID:    userID,
				InvitedAt: &invitedAt,
			}, nil
		},
		createMemberFunc: func(ctx context.Context, member domain.Member) (*domain.Member, error) {
			// Return the created member
			return &testMember, nil
		},
	}

	mockUserAdapter := &mockUserAdapterForMember{
		getUserByIDFunc: func(ctx context.Context, userID string) (*Output, error) {
			// Return a user to simulate it exists
			return &Output{ID: userID}, nil
		},
	}

	// Create application with mock dependencies
	app := NewApplication(mockRepo, mockUserAdapter)

	// Test the CreateMember function
	input := CreateMemberInput{
		TenantID: "tenant-123",
		UserID:   "user-123",
	}

	result, err := app.CreateMember(context.Background(), input)

	fmt.Println(err)

	// Verify results
	if err != nil {
		t.Errorf("Expected no error, got %v", err)
	}
	if result == nil {
		t.Fatal("Expected result to not be nil")
	}
	if result.ID != "member-123" {
		t.Errorf("Expected member ID to be 'member-123', got %s", result.ID)
	}
	if result.Status != domain_pkg.MemberStatusPending {
		t.Errorf("Expected status to be 'pending', got %s", result.Status)
	}
}

func TestApplication_CreateMember_TenantNotFound(t *testing.T) {
	// Setup mocks
	mockRepo := &mockMemberRepository{
		mockTenantRepository: mockTenantRepository{
			getTenantByIDFunc: func(ctx context.Context, tenantID string) (*domain.Tenant, error) {
				// Return error to simulate tenant not found
				return nil, errors.New("tenant not found")
			},
		},
	}

	mockUserAdapter := &mockUserAdapterForMember{
		getUserByIDFunc: func(ctx context.Context, userID string) (*Output, error) {
			// Return a user to simulate it exists
			return &Output{ID: userID}, nil
		},
	}

	// Create application with mock dependencies
	app := NewApplication(mockRepo, mockUserAdapter)

	// Test the CreateMember function
	input := CreateMemberInput{
		TenantID: "nonexistent-tenant",
		UserID:   "user-123",
	}

	result, err := app.CreateMember(context.Background(), input)

	// Verify results
	if err == nil {
		t.Error("Expected error for nonexistent tenant, got nil")
	}
	if result != nil {
		t.Errorf("Expected nil result, got %+v", result)
	}
}

func TestApplication_CreateMember_UserNotFound(t *testing.T) {
	// Setup mocks
	mockRepo := &mockMemberRepository{
		mockTenantRepository: mockTenantRepository{
			getTenantByIDFunc: func(ctx context.Context, tenantID string) (*domain.Tenant, error) {
				// Return a tenant to simulate it exists
				return &domain.Tenant{ID: tenantID}, nil
			},
		},
	}

	mockUserAdapter := &mockUserAdapterForMember{
		getUserByIDFunc: func(ctx context.Context, userID string) (*Output, error) {
			// Return error to simulate user not found
			return nil, errors.New("user not found")
		},
	}

	// Create application with mock dependencies
	app := NewApplication(mockRepo, mockUserAdapter)

	// Test the CreateMember function
	input := CreateMemberInput{
		TenantID: "tenant-123",
		UserID:   "nonexistent-user",
	}

	result, err := app.CreateMember(context.Background(), input)

	// Verify results
	if err == nil {
		t.Error("Expected error for nonexistent user, got nil")
	}
	if result != nil {
		t.Errorf("Expected nil result, got %+v", result)
	}
}

func TestApplication_UpdateMember_Success(t *testing.T) {
	// Define test time
	acceptedAt := time.Now()
	invitedAt := acceptedAt.Add(-time.Hour)

	// Setup mocks
	mockRepo := &mockMemberRepository{
		getMemberByIdFunc: func(ctx context.Context, memberID string) (*domain.Member, error) {
			// Return a pending member to simulate it exists
			return &domain.Member{
				ID:        memberID,
				Status:    domain_pkg.MemberStatusPending,
				InvitedAt: &invitedAt,
			}, nil
		},
		updateMemberFunc: func(ctx context.Context, memberUpdateData domain.MemberUpdateData) (*domain.Member, error) {
			// Return updated member
			return &domain.Member{
				ID:         memberUpdateData.MemberID,
				Status:     domain_pkg.MemberStatusAccepted,
				InvitedAt:  &invitedAt,
				AcceptedAt: memberUpdateData.AcceptedAt,
			}, nil
		},
	}

	// Create application with mock dependencies
	app := NewApplication(mockRepo, &mockUserAdapter{})

	// Test the UpdateMember function
	input := UpdateMemberInput{
		MemberID: "member-123",
	}

	result, err := app.UpdateMember(context.Background(), input)

	// Verify results
	if err != nil {
		t.Errorf("Expected no error, got %v", err)
	}
	if result == nil {
		t.Fatal("Expected result to not be nil")
	}
	if result.Status != domain_pkg.MemberStatusAccepted {
		t.Errorf("Expected status to be 'accepted', got %s", result.Status)
	}
}

func TestApplication_UpdateMember_MemberNotFound(t *testing.T) {
	// Setup mocks
	mockRepo := &mockMemberRepository{
		getMemberByIdFunc: func(ctx context.Context, memberID string) (*domain.Member, error) {
			// Return error to simulate member not found
			return nil, errors.New("member not found")
		},
	}

	// Create application with mock dependencies
	app := NewApplication(mockRepo, &mockUserAdapter{})

	// Test the UpdateMember function
	input := UpdateMemberInput{
		MemberID: "nonexistent-member",
	}

	result, err := app.UpdateMember(context.Background(), input)

	// Verify results
	if err == nil {
		t.Error("Expected error for nonexistent member, got nil")
	}
	if result != nil {
		t.Errorf("Expected nil result, got %+v", result)
	}
}

func TestApplication_UpdateMember_AlreadyAccepted(t *testing.T) {
	// Define test time
	acceptedAt := time.Now()
	invitedAt := acceptedAt.Add(-time.Hour)

	// Setup mocks
	mockRepo := &mockMemberRepository{
		getMemberByIdFunc: func(ctx context.Context, memberID string) (*domain.Member, error) {
			// Return an already accepted member
			return &domain.Member{
				ID:         memberID,
				Status:     domain_pkg.MemberStatusAccepted,
				InvitedAt:  &invitedAt,
				AcceptedAt: &acceptedAt,
			}, nil
		},
	}

	// Create application with mock dependencies
	app := NewApplication(mockRepo, &mockUserAdapter{})

	// Test the UpdateMember function
	input := UpdateMemberInput{
		MemberID: "member-123",
	}

	result, err := app.UpdateMember(context.Background(), input)

	// Verify results
	if err == nil {
		t.Error("Expected error for already accepted member, got nil")
	}
	if result != nil {
		t.Errorf("Expected nil result, got %+v", result)
	}
}

func TestApplication_GetMemberByID(t *testing.T) {
	// Define test time
	now := time.Now()

	// Setup mocks
	mockRepo := &mockMemberRepository{
		getMemberByIdFunc: func(ctx context.Context, memberID string) (*domain.Member, error) {
			// Return a member to simulate it exists
			return &domain.Member{
				ID:         memberID,
				TenantID:   "tenant-123",
				UserID:     "user-123",
				Status:     domain_pkg.MemberStatusAccepted,
				InvitedAt:  &now,
				AcceptedAt: &now,
			}, nil
		},
	}

	// Create application with mock dependencies
	app := NewApplication(mockRepo, &mockUserAdapter{})

	// Test the GetMemberByID function
	result, err := app.GetMemberByID(context.Background(), "member-123")

	// Verify results
	if err != nil {
		t.Errorf("Expected no error, got %v", err)
	}
	if result == nil {
		t.Fatal("Expected result to not be nil")
	}
	if result.ID != "member-123" {
		t.Errorf("Expected member ID to be 'member-123', got %s", result.ID)
	}
	if result.TenantID != "tenant-123" {
		t.Errorf("Expected tenant ID to be 'tenant-123', got %s", result.TenantID)
	}
	if result.UserID != "user-123" {
		t.Errorf("Expected user ID to be 'user-123', got %s", result.UserID)
	}
	if result.Status != domain_pkg.MemberStatusAccepted {
		t.Errorf("Expected status to be 'accepted', got %s", result.Status)
	}
}

func TestApplication_CheckUserAlreadyAMember_True(t *testing.T) {
	// Define test time
	now := time.Now()

	// Setup mocks
	mockRepo := &mockMemberRepository{
		getMemberByUserIDFunc: func(ctx context.Context, userID string) (*domain.Member, error) {
			// Return a member with matching tenant ID
			return &domain.Member{
				ID:        "member-123",
				TenantID:  "tenant-123",
				UserID:    userID,
				Status:    domain_pkg.MemberStatusAccepted,
				InvitedAt: &now,
			}, nil
		},
	}

	// Create application with mock dependencies
	app := NewApplication(mockRepo, &mockUserAdapter{})

	// Test the CheckUserAlreadyAMember function
	result := app.CheckUserAlreadyAMember(context.Background(), "tenant-123", "user-123")

	// Verify results
	if !result {
		t.Error("Expected CheckUserAlreadyAMember to return true, got false")
	}
}

func TestApplication_CheckUserAlreadyAMember_False_DifferentTenant(t *testing.T) {
	// Define test time
	now := time.Now()

	// Setup mocks
	mockRepo := &mockMemberRepository{
		getMemberByUserIDFunc: func(ctx context.Context, userID string) (*domain.Member, error) {
			// Return a member with different tenant ID
			return &domain.Member{
				ID:        "member-123",
				TenantID:  "different-tenant",
				UserID:    userID,
				Status:    domain_pkg.MemberStatusAccepted,
				InvitedAt: &now,
			}, nil
		},
	}

	// Create application with mock dependencies
	app := NewApplication(mockRepo, &mockUserAdapter{})

	// Test the CheckUserAlreadyAMember function
	result := app.CheckUserAlreadyAMember(context.Background(), "tenant-123", "user-123")

	// Verify results
	if result {
		t.Error("Expected CheckUserAlreadyAMember to return false, got true")
	}
}

func TestApplication_CheckUserAlreadyAMember_False_NotAccepted(t *testing.T) {
	// Define test time
	now := time.Now()

	// Setup mocks
	mockRepo := &mockMemberRepository{
		getMemberByUserIDFunc: func(ctx context.Context, userID string) (*domain.Member, error) {
			// Return a member that's pending (not accepted)
			return &domain.Member{
				ID:        "member-123",
				TenantID:  "tenant-123",
				UserID:    userID,
				Status:    domain_pkg.MemberStatusPending,
				InvitedAt: &now,
			}, nil
		},
	}

	// Create application with mock dependencies
	app := NewApplication(mockRepo, &mockUserAdapter{})

	// Test the CheckUserAlreadyAMember function
	result := app.CheckUserAlreadyAMember(context.Background(), "tenant-123", "user-123")

	// Verify results
	if result {
		t.Error("Expected CheckUserAlreadyAMember to return false, got true")
	}
}

func TestApplication_CheckUserAlreadyAMember_False_Error(t *testing.T) {
	// Setup mocks
	mockRepo := &mockMemberRepository{
		getMemberByUserIDFunc: func(ctx context.Context, userID string) (*domain.Member, error) {
			// Return error to simulate member not found
			return nil, errors.New("member not found")
		},
	}

	// Create application with mock dependencies
	app := NewApplication(mockRepo, &mockUserAdapter{})

	// Test the CheckUserAlreadyAMember function
	result := app.CheckUserAlreadyAMember(context.Background(), "tenant-123", "user-123")

	// Verify results
	if result {
		t.Error("Expected CheckUserAlreadyAMember to return false when error occurs, got true")
	}
}
