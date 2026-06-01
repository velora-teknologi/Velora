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
	"go.uber.org/zap"
)

type mockWorkflowService struct{}

func (m *mockWorkflowService) CreateWorkflow(_ context.Context, userID string, req *dtos.CreateWorkflowRequest) (*dtos.WorkflowResponse, error) {
	return &dtos.WorkflowResponse{ID: "workflow-1", Name: req.Name, Description: req.Description, Definition: req.Definition, IsActive: true}, nil
}

func (m *mockWorkflowService) GetWorkflow(_ context.Context, userID, id string) (*dtos.WorkflowResponse, error) {
	return &dtos.WorkflowResponse{ID: id, Name: "Workflow One", Description: "Test workflow", Definition: map[string]interface{}{"steps": []string{"start"}}, IsActive: true}, nil
}

func (m *mockWorkflowService) ExecuteWorkflow(_ context.Context, userID, id string) (*dtos.WorkflowExecutionResponse, error) {
	return &dtos.WorkflowExecutionResponse{WorkflowID: id, AgentID: "agent-1", Status: "queued", Subject: "workflow.execute"}, nil
}

func (m *mockWorkflowService) UpdateWorkflow(_ context.Context, userID, id string, req *dtos.UpdateWorkflowRequest) (*dtos.WorkflowResponse, error) {
	return nil, nil
}

func (m *mockWorkflowService) DeleteWorkflow(_ context.Context, userID, id string) error {
	return nil
}

func (m *mockWorkflowService) ListWorkflows(_ context.Context, userID, agentID string, skip, limit int) ([]*dtos.WorkflowResponse, error) {
	return []*dtos.WorkflowResponse{{ID: "workflow-1", Name: "Workflow One", Description: "Test workflow", Definition: map[string]interface{}{"steps": []string{"start"}}, IsActive: true}}, nil
}

func TestWorkflowHandlerCreateWorkflow(t *testing.T) {
	logger := zap.NewNop().Sugar()
	service := &mockWorkflowService{}
	handler := handlers.NewWorkflowHandler(service, logger)

	app := fiber.New()
	app.Post("/workflows", func(c *fiber.Ctx) error {
		c.Locals("user", map[string]interface{}{"sub": "test-user-id"})
		return handler.CreateWorkflow(c)
	})

	body, _ := json.Marshal(dtos.CreateWorkflowRequest{Name: "Workflow One", AgentID: "agent-1", Description: "Test workflow"})
	req := httptest.NewRequest(http.MethodPost, "/workflows", bytes.NewReader(body))
	req.Header.Set("Content-Type", "application/json")

	resp, err := app.Test(req)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if resp.StatusCode != http.StatusCreated {
		t.Fatalf("expected status %d but got %d", http.StatusCreated, resp.StatusCode)
	}
}

func TestWorkflowHandlerListWorkflows(t *testing.T) {
	logger := zap.NewNop().Sugar()
	service := &mockWorkflowService{}
	handler := handlers.NewWorkflowHandler(service, logger)

	app := fiber.New()
	app.Get("/workflows", func(c *fiber.Ctx) error {
		c.Locals("user", map[string]interface{}{"sub": "test-user-id"})
		return handler.ListWorkflows(c)
	})

	req := httptest.NewRequest(http.MethodGet, "/workflows", nil)
	resp, err := app.Test(req)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if resp.StatusCode != http.StatusOK {
		t.Fatalf("expected status %d but got %d", http.StatusOK, resp.StatusCode)
	}
}

func TestWorkflowHandlerGetWorkflow(t *testing.T) {
	logger := zap.NewNop().Sugar()
	service := &mockWorkflowService{}
	handler := handlers.NewWorkflowHandler(service, logger)

	app := fiber.New()
	app.Get("/workflows/:id", func(c *fiber.Ctx) error {
		c.Locals("user", map[string]interface{}{"sub": "test-user-id"})
		return handler.GetWorkflow(c)
	})

	req := httptest.NewRequest(http.MethodGet, "/workflows/workflow-1", nil)
	resp, err := app.Test(req)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if resp.StatusCode != http.StatusOK {
		t.Fatalf("expected status %d but got %d", http.StatusOK, resp.StatusCode)
	}
}
