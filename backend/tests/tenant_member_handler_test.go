package services_test

import (
	"context"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/gofiber/fiber/v2"
	"github.com/velora-teknologi/velora/internal/application/dtos"
	"github.com/velora-teknologi/velora/internal/infrastructure/handlers"
	customErrors "github.com/velora-teknologi/velora/pkg/errors"
	"go.uber.org/zap"
)

type mockTenantMemberServiceHandler struct{}

func (m *mockTenantMemberServiceHandler) ListTenantMembers(ctx context.Context, userID, tenantID string, skip, limit int) ([]*dtos.TenantMemberResponse, error) {
	if tenantID == "tenant-1" && userID == "owner-1" {
		return []*dtos.TenantMemberResponse{{ID: "user-1", Email: "user1@example.com", FullName: "User One", Role: "member", IsActive: true}}, nil
	}
	return nil, customErrors.NewForbiddenError("unauthorized")
}

func TestTenantMemberHandlerListSuccess(t *testing.T) {
	logger := zap.NewNop().Sugar()
	service := &mockTenantMemberServiceHandler{}
	handler := handlers.NewTenantMemberHandler(service, logger)

	app := fiber.New()
	app.Get("/tenants/:id/members", func(c *fiber.Ctx) error {
		c.Locals("user", map[string]interface{}{"sub": "owner-1"})
		return handler.ListTenantMembers(c)
	})

	req := httptest.NewRequest(http.MethodGet, "/tenants/tenant-1/members", nil)
	resp, err := app.Test(req)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if resp.StatusCode != http.StatusOK {
		t.Fatalf("expected status %d but got %d", http.StatusOK, resp.StatusCode)
	}
}

func TestTenantMemberHandlerListUnauthorized(t *testing.T) {
	logger := zap.NewNop().Sugar()
	service := &mockTenantMemberServiceHandler{}
	handler := handlers.NewTenantMemberHandler(service, logger)

	app := fiber.New()
	app.Get("/tenants/:id/members", func(c *fiber.Ctx) error {
		c.Locals("user", map[string]interface{}{"sub": "other-owner"})
		return handler.ListTenantMembers(c)
	})

	req := httptest.NewRequest(http.MethodGet, "/tenants/tenant-1/members", nil)
	resp, err := app.Test(req)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if resp.StatusCode == http.StatusOK {
		t.Fatalf("expected non-OK status for unauthorized access, got %d", resp.StatusCode)
	}
}
