package services_test

import (
	"context"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
	"github.com/velora-teknologi/velora/internal/application/services"
	"github.com/velora-teknologi/velora/internal/domain/models"
	customErrors "github.com/velora-teknologi/velora/pkg/errors"
	"go.uber.org/zap/zaptest"
)

type mockUserRepoForTenantMembers struct {
	mock.Mock
}

func (m *mockUserRepoForTenantMembers) Create(ctx context.Context, user *models.User) error {
	return nil
}
func (m *mockUserRepoForTenantMembers) GetByID(ctx context.Context, id string) (*models.User, error) {
	args := m.Called(ctx, id)
	if user, ok := args.Get(0).(*models.User); ok {
		return user, args.Error(1)
	}
	return nil, args.Error(1)
}
func (m *mockUserRepoForTenantMembers) GetByEmail(ctx context.Context, email string) (*models.User, error) {
	return nil, nil
}
func (m *mockUserRepoForTenantMembers) Update(ctx context.Context, user *models.User) error {
	return nil
}
func (m *mockUserRepoForTenantMembers) Delete(ctx context.Context, id string) error { return nil }
func (m *mockUserRepoForTenantMembers) List(ctx context.Context, skip, limit int) ([]*models.User, error) {
	return nil, nil
}
func (m *mockUserRepoForTenantMembers) ListByTenantID(ctx context.Context, tenantID string, skip, limit int) ([]*models.User, error) {
	args := m.Called(ctx, tenantID, skip, limit)
	if users, ok := args.Get(0).([]*models.User); ok {
		return users, args.Error(1)
	}
	return nil, args.Error(1)
}

type mockTenantRepoForTenantMembers struct {
	mock.Mock
}

func (m *mockTenantRepoForTenantMembers) Create(ctx context.Context, tenant *models.Tenant) error {
	return nil
}
func (m *mockTenantRepoForTenantMembers) GetByID(ctx context.Context, id string) (*models.Tenant, error) {
	args := m.Called(ctx, id)
	if tenant, ok := args.Get(0).(*models.Tenant); ok {
		return tenant, args.Error(1)
	}
	return nil, args.Error(1)
}
func (m *mockTenantRepoForTenantMembers) Update(ctx context.Context, tenant *models.Tenant) error {
	return nil
}
func (m *mockTenantRepoForTenantMembers) Delete(ctx context.Context, id string) error { return nil }
func (m *mockTenantRepoForTenantMembers) List(ctx context.Context, ownerID string, skip, limit int) ([]*models.Tenant, error) {
	return nil, nil
}

func TestListTenantMembersSuccess(t *testing.T) {
	logger := zaptest.NewLogger(t).Sugar()
	userRepo := new(mockUserRepoForTenantMembers)
	tenantRepo := new(mockTenantRepoForTenantMembers)

	tenant := &models.Tenant{ID: "tenant-1", OwnerID: "owner-1"}
	users := []*models.User{{ID: "user-1", Email: "user1@example.com", FullName: "User One", Role: "member", IsActive: true}}

	tenantRepo.On("GetByID", mock.Anything, "tenant-1").Return(tenant, nil)
	userRepo.On("ListByTenantID", mock.Anything, "tenant-1", 0, 20).Return(users, nil)

	service := services.NewTenantMemberService(userRepo, tenantRepo, logger)

	result, err := service.ListTenantMembers(context.Background(), "owner-1", "tenant-1", 0, 20)
	assert.NoError(t, err)
	assert.Len(t, result, 1)
	assert.Equal(t, "user-1", result[0].ID)
	assert.Equal(t, "user1@example.com", result[0].Email)
}

func TestListTenantMembersForbidden(t *testing.T) {
	logger := zaptest.NewLogger(t).Sugar()
	userRepo := new(mockUserRepoForTenantMembers)
	tenantRepo := new(mockTenantRepoForTenantMembers)

	tenant := &models.Tenant{ID: "tenant-1", OwnerID: "owner-2"}
	tenantRepo.On("GetByID", mock.Anything, "tenant-1").Return(tenant, nil)

	service := services.NewTenantMemberService(userRepo, tenantRepo, logger)

	result, err := service.ListTenantMembers(context.Background(), "owner-1", "tenant-1", 0, 20)
	assert.Nil(t, result)
	assert.Error(t, err)
	assert.Equal(t, customErrors.ForbiddenError, err.(*customErrors.CustomError).Type)
}
