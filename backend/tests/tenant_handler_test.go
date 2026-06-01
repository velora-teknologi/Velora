package services_test

import (
	"bytes"
	"context"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/gofiber/fiber/v2"
	"github.com/velora-teknologi/velora/internal/application/dtos"
	"github.com/velora-teknologi/velora/internal/infrastructure/handlers"
	customErrors "github.com/velora-teknologi/velora/pkg/errors"
	"go.uber.org/zap"
)

type mockTenantService struct {
	createCalled bool
}

func (m *mockTenantService) CreateTenant(_ context.Context, userID string, req *dtos.CreateTenantRequest) (*dtos.TenantResponse, error) {
	m.createCalled = true
	return &dtos.TenantResponse{ID: "tenant-1", Name: req.Name, Description: req.Description, OwnerID: userID, IsActive: true}, nil
}

func (m *mockTenantService) GetTenant(_ context.Context, userID, id string) (*dtos.TenantResponse, error) {
	if id != "tenant-1" {
		return nil, customErrors.NewNotFoundError("Tenant")
	}
	return &dtos.TenantResponse{ID: id, Name: "Acme", Description: "Tenant", OwnerID: userID, IsActive: true}, nil
}

func (m *mockTenantService) UpdateTenant(_ context.Context, userID, id string, req *dtos.UpdateTenantRequest) (*dtos.TenantResponse, error) {
	return &dtos.TenantResponse{ID: id, Name: req.Name, Description: req.Description, OwnerID: userID, IsActive: true}, nil
}

func (m *mockTenantService) DeleteTenant(_ context.Context, userID, id string) error {
	return nil
}

func (m *mockTenantService) ListTenants(_ context.Context, userID string, skip, limit int) ([]*dtos.TenantResponse, error) {
	return []*dtos.TenantResponse{{ID: "tenant-1", Name: "Acme", Description: "Tenant", OwnerID: userID, IsActive: true}}, nil
}

func TestTenantHandlerCreateTenant(t *testing.T) {
	logger := zap.NewNop().Sugar()
	service := &mockTenantService{}
	handler := handlers.NewTenantHandler(service, logger)

	app := fiber.New()
	app.Post("/tenants", func(c *fiber.Ctx) error {
		c.Locals("user", map[string]interface{}{"sub": "owner-1"})
		return handler.CreateTenant(c)
	})

	body, _ := json.Marshal(dtos.CreateTenantRequest{Name: "Acme Corp", Description: "Tenant for Acme"})
	req := httptest.NewRequest(http.MethodPost, "/tenants", bytes.NewReader(body))
	req.Header.Set("Content-Type", "application/json")

	resp, err := app.Test(req)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if resp.StatusCode != http.StatusCreated {
		t.Fatalf("expected status %d but got %d", http.StatusCreated, resp.StatusCode)
	}
}

func TestTenantHandlerListTenants(t *testing.T) {
	logger := zap.NewNop().Sugar()
	service := &mockTenantService{}
	handler := handlers.NewTenantHandler(service, logger)

	app := fiber.New()
	app.Get("/tenants", func(c *fiber.Ctx) error {
		c.Locals("user", map[string]interface{}{"sub": "owner-1"})
		return handler.ListTenants(c)
	})

	req := httptest.NewRequest(http.MethodGet, "/tenants", nil)
	resp, err := app.Test(req)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if resp.StatusCode != http.StatusOK {
		t.Fatalf("expected status %d but got %d", http.StatusOK, resp.StatusCode)
	}
}
