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
	"github.com/velora-teknologi/velora/pkg/errors"
	"go.uber.org/zap"
)

type mockAgentService struct {
	createCalled bool
	agents       []*dtos.AgentResponse
}

func (m *mockAgentService) CreateAgent(_ context.Context, userID string, req *dtos.CreateAgentRequest) (*dtos.AgentResponse, error) {
	m.createCalled = true
	return &dtos.AgentResponse{ID: "agent-1", Name: req.Name, Description: req.Description, Status: "active", Config: req.Config}, nil
}

func (m *mockAgentService) GetAgent(_ context.Context, userID, id string) (*dtos.AgentResponse, error) {
	return nil, errors.NewNotFoundError("Agent")
}

func (m *mockAgentService) UpdateAgent(_ context.Context, userID, id string, req *dtos.UpdateAgentRequest) (*dtos.AgentResponse, error) {
	return nil, errors.NewNotFoundError("Agent")
}

func (m *mockAgentService) DeleteAgent(_ context.Context, userID, id string) error {
	return nil
}

func (m *mockAgentService) ListAgents(_ context.Context, userID string, skip, limit int) ([]*dtos.AgentResponse, error) {
	return []*dtos.AgentResponse{{ID: "agent-1", Name: "Agent One", Description: "Test agent", Status: "active"}}, nil
}

func (m *mockAgentService) ExecuteAgent(_ context.Context, userID, agentID string, input interface{}) (*dtos.AgentExecutionResponse, error) {
	return nil, errors.NewNotFoundError("Agent")
}

func TestAgentHandlerCreateAgent(t *testing.T) {
	logger := zap.NewNop().Sugar()
	service := &mockAgentService{}
	handler := handlers.NewAgentHandler(service, nil, logger)

	app := fiber.New()
	app.Post("/agents", func(c *fiber.Ctx) error {
		c.Locals("user", map[string]interface{}{"sub": "test-user-id"})
		return handler.CreateAgent(c)
	})

	body, _ := json.Marshal(dtos.CreateAgentRequest{Name: "Agent One", Description: "Test agent"})
	req := httptest.NewRequest(http.MethodPost, "/agents", bytes.NewReader(body))
	req.Header.Set("Content-Type", "application/json")

	resp, err := app.Test(req)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if resp.StatusCode != http.StatusCreated {
		t.Fatalf("expected status %d but got %d", http.StatusCreated, resp.StatusCode)
	}
}

func TestAgentHandlerListAgents(t *testing.T) {
	logger := zap.NewNop().Sugar()
	service := &mockAgentService{}
	handler := handlers.NewAgentHandler(service, nil, logger)

	app := fiber.New()
	app.Get("/agents", func(c *fiber.Ctx) error {
		c.Locals("user", map[string]interface{}{"sub": "test-user-id"})
		return handler.ListAgents(c)
	})

	req := httptest.NewRequest(http.MethodGet, "/agents", nil)
	resp, err := app.Test(req)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if resp.StatusCode != http.StatusOK {
		t.Fatalf("expected status %d but got %d", http.StatusOK, resp.StatusCode)
	}
}

func TestAgentHandlerListAgentWorkflows(t *testing.T) {
	logger := zap.NewNop().Sugar()
	service := &mockAgentService{}
	workflowService := &mockWorkflowService{}
	handler := handlers.NewAgentHandler(service, workflowService, logger)

	app := fiber.New()
	app.Get("/agents/:id/workflows", func(c *fiber.Ctx) error {
		c.Locals("user", map[string]interface{}{"sub": "test-user-id"})
		return handler.ListAgentWorkflows(c)
	})

	req := httptest.NewRequest(http.MethodGet, "/agents/agent-1/workflows", nil)
	resp, err := app.Test(req)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if resp.StatusCode != http.StatusOK {
		t.Fatalf("expected status %d but got %d", http.StatusOK, resp.StatusCode)
	}
}
