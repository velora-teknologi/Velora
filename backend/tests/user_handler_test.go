package services_test

import (
	"context"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"

	"github.com/gofiber/fiber/v2"
	"github.com/google/uuid"
	"github.com/velora-teknologi/velora/internal/application/dtos"
	"github.com/velora-teknologi/velora/internal/infrastructure/handlers"
	"go.uber.org/zap"
)

type mockUserService struct {
	createdUser *dtos.UserResponse
	err         error
	users       []*dtos.UserResponse
}

func (m *mockUserService) CreateUser(ctx context.Context, req *dtos.CreateUserRequest) (*dtos.UserResponse, error) {
	if m.err != nil {
		return nil, m.err
	}
	return &dtos.UserResponse{
		ID:       uuid.NewString(),
		Email:    req.Email,
		FullName: req.FullName,
		Role:     "user",
		IsActive: true,
	}, nil
}

func (m *mockUserService) GetUser(ctx context.Context, id string) (*dtos.UserResponse, error) {
	return nil, nil
}

func (m *mockUserService) UpdateUser(ctx context.Context, id string, req *dtos.UpdateUserRequest) (*dtos.UserResponse, error) {
	return nil, nil
}

func (m *mockUserService) DeleteUser(ctx context.Context, id string) error {
	return nil
}

func (m *mockUserService) ListUsers(ctx context.Context, skip, limit int) ([]*dtos.UserResponse, error) {
	return m.users, nil
}

func (m *mockUserService) LoginUser(ctx context.Context, req *dtos.LoginRequest) (*dtos.LoginResponse, error) {
	return nil, nil
}

func TestRegisterUser_Success(t *testing.T) {
	app := fiber.New()
	service := &mockUserService{}
	handler := handlers.NewUserHandler(service, zap.NewNop().Sugar())

	app.Post("/api/v1/auth/register", handler.RegisterUser)

	body := `{"email":"register@example.com","password":"secure123","full_name":"Register User"}`
	req := httptest.NewRequest(http.MethodPost, "/api/v1/auth/register", strings.NewReader(body))
	req.Header.Set("Content-Type", "application/json")

	resp, err := app.Test(req)
	if err != nil {
		t.Fatalf("unexpected error from request: %v", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusCreated {
		t.Fatalf("expected status %d, got %d", http.StatusCreated, resp.StatusCode)
	}

	var response struct {
		Data dtos.UserResponse `json:"data"`
	}
	if err := json.NewDecoder(resp.Body).Decode(&response); err != nil {
		t.Fatalf("failed to decode response: %v", err)
	}
	if response.Data.Email != "register@example.com" {
		t.Fatalf("expected email register@example.com, got %s", response.Data.Email)
	}
}

func TestRegisterUser_InvalidPayload(t *testing.T) {
	app := fiber.New()
	service := &mockUserService{}
	handler := handlers.NewUserHandler(service, zap.NewNop().Sugar())

	app.Post("/api/v1/auth/register", handler.RegisterUser)

	body := `{"email":"invalid@example.com"}`
	req := httptest.NewRequest(http.MethodPost, "/api/v1/auth/register", strings.NewReader(body))
	req.Header.Set("Content-Type", "application/json")

	resp, err := app.Test(req)
	if err != nil {
		t.Fatalf("unexpected error from request: %v", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusBadRequest {
		t.Fatalf("expected status %d, got %d", http.StatusBadRequest, resp.StatusCode)
	}
}

func TestListUsers_AdminAllowed(t *testing.T) {
	app := fiber.New()
	service := &mockUserService{
		users: []*dtos.UserResponse{{ID: uuid.NewString(), Email: "admin@example.com", FullName: "Admin User", Role: "admin", IsActive: true}},
	}
	handler := handlers.NewUserHandler(service, zap.NewNop().Sugar())

	app.Use(func(c *fiber.Ctx) error {
		c.Locals("user", map[string]interface{}{"sub": "admin-id", "role": "admin"})
		return c.Next()
	})

	app.Get("/api/v1/users", handler.ListUsers)

	req := httptest.NewRequest(http.MethodGet, "/api/v1/users", nil)
	resp, err := app.Test(req)
	if err != nil {
		t.Fatalf("unexpected error from request: %v", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		t.Fatalf("expected status %d, got %d", http.StatusOK, resp.StatusCode)
	}
}

func TestListUsers_NonAdminForbidden(t *testing.T) {
	app := fiber.New()
	service := &mockUserService{
		users: []*dtos.UserResponse{{ID: uuid.NewString(), Email: "user@example.com", FullName: "Regular User", Role: "user", IsActive: true}},
	}
	handler := handlers.NewUserHandler(service, zap.NewNop().Sugar())

	app.Use(func(c *fiber.Ctx) error {
		c.Locals("user", map[string]interface{}{"sub": "user-id", "role": "user"})
		return c.Next()
	})

	app.Get("/api/v1/users", handler.ListUsers)

	req := httptest.NewRequest(http.MethodGet, "/api/v1/users", nil)
	resp, err := app.Test(req)
	if err != nil {
		t.Fatalf("unexpected error from request: %v", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusForbidden {
		t.Fatalf("expected status %d, got %d", http.StatusForbidden, resp.StatusCode)
	}
}
