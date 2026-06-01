package services_test

import (
	"context"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/gofiber/fiber/v2"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
	"github.com/velora-teknologi/velora/internal/application/services"
	"github.com/velora-teknologi/velora/internal/domain/models"
	"github.com/velora-teknologi/velora/internal/infrastructure/handlers"
	"go.uber.org/zap"
)

// mock implementations for repos used by tenant member service
type mockUserRepoE2E struct {
	mock.Mock
}

func (m *mockUserRepoE2E) Create(ctx context.Context, user *models.User) error { return nil }
func (m *mockUserRepoE2E) GetByID(ctx context.Context, id string) (*models.User, error) {
	args := m.Called(ctx, id)
	if u, ok := args.Get(0).(*models.User); ok {
		return u, args.Error(1)
	}
	return nil, args.Error(1)
}
func (m *mockUserRepoE2E) GetByEmail(ctx context.Context, email string) (*models.User, error) {
	return nil, nil
}
func (m *mockUserRepoE2E) Update(ctx context.Context, user *models.User) error { return nil }
func (m *mockUserRepoE2E) Delete(ctx context.Context, id string) error         { return nil }
func (m *mockUserRepoE2E) List(ctx context.Context, skip, limit int) ([]*models.User, error) {
	return nil, nil
}
func (m *mockUserRepoE2E) ListByTenantID(ctx context.Context, tenantID string, skip, limit int) ([]*models.User, error) {
	args := m.Called(ctx, tenantID, skip, limit)
	if us, ok := args.Get(0).([]*models.User); ok {
		return us, args.Error(1)
	}
	return nil, args.Error(1)
}

type mockTenantRepoE2E struct{ mock.Mock }

func (m *mockTenantRepoE2E) Create(ctx context.Context, tenant *models.Tenant) error { return nil }
func (m *mockTenantRepoE2E) GetByID(ctx context.Context, id string) (*models.Tenant, error) {
	args := m.Called(ctx, id)
	if t, ok := args.Get(0).(*models.Tenant); ok {
		return t, args.Error(1)
	}
	return nil, args.Error(1)
}
func (m *mockTenantRepoE2E) Update(ctx context.Context, tenant *models.Tenant) error { return nil }
func (m *mockTenantRepoE2E) Delete(ctx context.Context, id string) error             { return nil }
func (m *mockTenantRepoE2E) List(ctx context.Context, ownerID string, skip, limit int) ([]*models.Tenant, error) {
	return nil, nil
}

// test middleware: read X-User header and set c.Locals("user", map[string]interface{"sub": userID})
func testAuthMiddleware(c *fiber.Ctx) error {
	user := c.Get("X-User")
	if user != "" {
		c.Locals("user", map[string]interface{}{"sub": user})
	}
	return c.Next()
}

func TestE2ETenantMembers_Route(t *testing.T) {
	logger := zap.NewNop().Sugar()
	userRepo := new(mockUserRepoE2E)
	tenantRepo := new(mockTenantRepoE2E)

	service := services.NewTenantMemberService(userRepo, tenantRepo, logger)
	handler := handlers.NewTenantMemberHandler(service, logger)

	app := fiber.New()
	app.Use(testAuthMiddleware)
	app.Get("/api/v1/tenants/:id/members", handler.ListTenantMembers)

	// happy path: owner requests members
	tenant := &models.Tenant{ID: "tenant-1", OwnerID: "owner-1"}
	users := []*models.User{{ID: "user-1", Email: "user1@example.com", FullName: "User One", Role: "member", IsActive: true}}
	tenantRepo.On("GetByID", mock.Anything, "tenant-1").Return(tenant, nil)
	userRepo.On("ListByTenantID", mock.Anything, "tenant-1", 0, 20).Return(users, nil)

	req := httptest.NewRequest(http.MethodGet, "/api/v1/tenants/tenant-1/members", nil)
	req.Header.Set("X-User", "owner-1")
	resp, err := app.Test(req)
	assert.NoError(t, err)
	assert.Equal(t, http.StatusOK, resp.StatusCode)

	// unauthorized: different user
	tenantRepo.ExpectedCalls = nil
	userRepo.ExpectedCalls = nil
	tenantRepo.On("GetByID", mock.Anything, "tenant-1").Return(tenant, nil)

	req2 := httptest.NewRequest(http.MethodGet, "/api/v1/tenants/tenant-1/members", nil)
	req2.Header.Set("X-User", "other-user")
	resp2, err := app.Test(req2)
	assert.NoError(t, err)
	assert.NotEqual(t, http.StatusOK, resp2.StatusCode)
}
