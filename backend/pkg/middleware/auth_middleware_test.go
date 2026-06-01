package middleware_test

import (
	"context"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/gofiber/fiber/v2"
	"github.com/google/uuid"
	"github.com/stretchr/testify/assert"
	"golang.org/x/crypto/bcrypt"

	"github.com/velora-teknologi/velora/internal/domain/models"
	"github.com/velora-teknologi/velora/pkg/middleware"
)

type mockUserRepo struct {
	users map[string]*models.User
}

func newMockUserRepo() *mockUserRepo {
	return &mockUserRepo{users: make(map[string]*models.User)}
}

func (m *mockUserRepo) Create(ctx context.Context, user *models.User) error {
	m.users[user.ID] = user
	return nil
}
func (m *mockUserRepo) GetByID(ctx context.Context, id string) (*models.User, error) {
	return m.users[id], nil
}
func (m *mockUserRepo) GetByEmail(ctx context.Context, email string) (*models.User, error) {
	for _, user := range m.users {
		if user.Email == email {
			return user, nil
		}
	}
	return nil, nil
}
func (m *mockUserRepo) Update(ctx context.Context, user *models.User) error {
	m.users[user.ID] = user
	return nil
}
func (m *mockUserRepo) Delete(ctx context.Context, id string) error {
	delete(m.users, id)
	return nil
}
func (m *mockUserRepo) List(ctx context.Context, skip, limit int) ([]*models.User, error) {
	return nil, nil
}
func (m *mockUserRepo) ListByTenantID(ctx context.Context, tenantID string, skip, limit int) ([]*models.User, error) {
	return nil, nil
}

type mockAPIKeyRepo struct {
	apiKeys map[string]*models.APIKey
}

func newMockAPIKeyRepo() *mockAPIKeyRepo {
	return &mockAPIKeyRepo{apiKeys: make(map[string]*models.APIKey)}
}

func (m *mockAPIKeyRepo) Create(ctx context.Context, apiKey *models.APIKey) error {
	m.apiKeys[apiKey.ID] = apiKey
	return nil
}
func (m *mockAPIKeyRepo) GetByID(ctx context.Context, id string) (*models.APIKey, error) {
	return m.apiKeys[id], nil
}
func (m *mockAPIKeyRepo) GetByKeyID(ctx context.Context, keyID string) (*models.APIKey, error) {
	for _, apiKey := range m.apiKeys {
		if apiKey.KeyID == keyID {
			return apiKey, nil
		}
	}
	return nil, nil
}
func (m *mockAPIKeyRepo) Update(ctx context.Context, apiKey *models.APIKey) error {
	m.apiKeys[apiKey.ID] = apiKey
	return nil
}
func (m *mockAPIKeyRepo) Delete(ctx context.Context, id string) error {
	delete(m.apiKeys, id)
	return nil
}
func (m *mockAPIKeyRepo) List(ctx context.Context, userID string, skip, limit int) ([]*models.APIKey, error) {
	return nil, nil
}

func TestAuthMiddleware_ApiKey(t *testing.T) {
	app := fiber.New()

	userRepo := newMockUserRepo()
	apiKeyRepo := newMockAPIKeyRepo()

	user := &models.User{ID: uuid.NewString(), Email: "apikey@example.com", FullName: "API Key User", Role: "user"}
	userRepo.Create(context.Background(), user)

	secret := "super-secret-value"
	hash, err := bcrypt.GenerateFromPassword([]byte(secret), bcrypt.DefaultCost)
	assert.NoError(t, err)

	apiKey := &models.APIKey{
		ID:       uuid.NewString(),
		KeyID:    uuid.NewString(),
		UserID:   user.ID,
		KeyHash:  string(hash),
		IsActive: true,
	}
	apiKeyRepo.Create(context.Background(), apiKey)

	app.Use(middleware.AuthMiddleware("ignored", apiKeyRepo, userRepo))
	app.Get("/protected", func(c *fiber.Ctx) error {
		claims := c.Locals("user").(map[string]interface{})
		return c.JSON(fiber.Map{"sub": claims["sub"], "role": claims["role"]})
	})

	req := httptest.NewRequest(http.MethodGet, "/protected", nil)
	req.Header.Set("Authorization", "ApiKey "+apiKey.KeyID+"."+secret)

	resp, err := app.Test(req)
	assert.NoError(t, err)
	assert.Equal(t, http.StatusOK, resp.StatusCode)
}

func TestAuthMiddleware_InvalidApiKey(t *testing.T) {
	app := fiber.New()

	userRepo := newMockUserRepo()
	apiKeyRepo := newMockAPIKeyRepo()

	app.Use(middleware.AuthMiddleware("ignored", apiKeyRepo, userRepo))
	app.Get("/protected", func(c *fiber.Ctx) error {
		return c.SendStatus(http.StatusOK)
	})

	req := httptest.NewRequest(http.MethodGet, "/protected", nil)
	req.Header.Set("Authorization", "ApiKey invalid.key")

	resp, err := app.Test(req)
	assert.NoError(t, err)
	assert.Equal(t, http.StatusUnauthorized, resp.StatusCode)
}
