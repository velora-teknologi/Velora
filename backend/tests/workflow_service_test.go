package services_test

import (
	"context"
	"testing"

	"github.com/google/uuid"
	"github.com/velora-teknologi/velora/internal/application/dtos"
	"github.com/velora-teknologi/velora/internal/application/services"
	"github.com/velora-teknologi/velora/internal/domain/models"
	"go.uber.org/zap"
)

type MockWorkflowRepository struct {
	workflows map[string]*models.Workflow
}

func NewMockWorkflowRepository() *MockWorkflowRepository {
	return &MockWorkflowRepository{workflows: make(map[string]*models.Workflow)}
}

func (m *MockWorkflowRepository) Create(ctx context.Context, workflow *models.Workflow) error {
	if workflow.ID == "" {
		workflow.ID = uuid.NewString()
	}
	m.workflows[workflow.ID] = workflow
	return nil
}

func (m *MockWorkflowRepository) GetByID(ctx context.Context, id string) (*models.Workflow, error) {
	return m.workflows[id], nil
}

func (m *MockWorkflowRepository) GetByAgentID(ctx context.Context, agentID string) ([]*models.Workflow, error) {
	var workflows []*models.Workflow
	for _, workflow := range m.workflows {
		if workflow.AgentID == agentID {
			workflows = append(workflows, workflow)
		}
	}
	return workflows, nil
}

func (m *MockWorkflowRepository) Update(ctx context.Context, workflow *models.Workflow) error {
	m.workflows[workflow.ID] = workflow
	return nil
}

func (m *MockWorkflowRepository) Delete(ctx context.Context, id string) error {
	delete(m.workflows, id)
	return nil
}

func (m *MockWorkflowRepository) List(ctx context.Context, agentID string, skip, limit int) ([]*models.Workflow, error) {
	var workflows []*models.Workflow
	for _, workflow := range m.workflows {
		if workflow.AgentID == agentID {
			workflows = append(workflows, workflow)
		}
	}
	if skip >= len(workflows) {
		return []*models.Workflow{}, nil
	}
	end := skip + limit
	if end > len(workflows) {
		end = len(workflows)
	}
	return workflows[skip:end], nil
}

type MockAgentRepositoryForWorkflow struct {
	agents map[string]*models.Agent
}

func NewMockAgentRepositoryForWorkflow() *MockAgentRepositoryForWorkflow {
	return &MockAgentRepositoryForWorkflow{agents: make(map[string]*models.Agent)}
}

func (m *MockAgentRepositoryForWorkflow) Create(ctx context.Context, agent *models.Agent) error {
	if agent.ID == "" {
		agent.ID = uuid.NewString()
	}
	m.agents[agent.ID] = agent
	return nil
}

func (m *MockAgentRepositoryForWorkflow) GetByID(ctx context.Context, id string) (*models.Agent, error) {
	return m.agents[id], nil
}

func (m *MockAgentRepositoryForWorkflow) GetByUserID(ctx context.Context, userID string) ([]*models.Agent, error) {
	var agents []*models.Agent
	for _, agent := range m.agents {
		if agent.UserID == userID {
			agents = append(agents, agent)
		}
	}
	return agents, nil
}

func (m *MockAgentRepositoryForWorkflow) List(ctx context.Context, userID string, skip, limit int) ([]*models.Agent, error) {
	var agents []*models.Agent
	for _, agent := range m.agents {
		if agent.UserID == userID {
			agents = append(agents, agent)
		}
	}
	if skip >= len(agents) {
		return []*models.Agent{}, nil
	}
	end := skip + limit
	if end > len(agents) {
		end = len(agents)
	}
	return agents[skip:end], nil
}

func (m *MockAgentRepositoryForWorkflow) Update(ctx context.Context, agent *models.Agent) error {
	m.agents[agent.ID] = agent
	return nil
}

func (m *MockAgentRepositoryForWorkflow) Delete(ctx context.Context, id string) error {
	delete(m.agents, id)
	return nil
}

func TestCreateWorkflow(t *testing.T) {
	agentRepo := NewMockAgentRepositoryForWorkflow()
	workflowRepo := NewMockWorkflowRepository()
	logger := zap.NewNop().Sugar()
	service := services.NewWorkflowService(workflowRepo, agentRepo, logger)

	userID := uuid.NewString()
	agent := &models.Agent{ID: uuid.NewString(), UserID: userID, Name: "My Agent"}
	agentRepo.Create(context.Background(), agent)

	req := &dtos.CreateWorkflowRequest{
		Name:       "Data Pipeline",
		AgentID:    agent.ID,
		Definition: map[string]interface{}{"steps": []interface{}{"fetch", "process"}},
	}

	resp, err := service.CreateWorkflow(context.Background(), userID, req)
	if err != nil {
		t.Fatalf("expected no error, got %v", err)
	}
	if resp.Name != req.Name {
		t.Fatalf("expected workflow name %s, got %s", req.Name, resp.Name)
	}
}

func TestListWorkflows(t *testing.T) {
	agentRepo := NewMockAgentRepositoryForWorkflow()
	workflowRepo := NewMockWorkflowRepository()
	logger := zap.NewNop().Sugar()
	service := services.NewWorkflowService(workflowRepo, agentRepo, logger)

	userID := uuid.NewString()
	agent := &models.Agent{ID: uuid.NewString(), UserID: userID, Name: "My Agent"}
	agentRepo.Create(context.Background(), agent)

	workflowRepo.Create(context.Background(), &models.Workflow{ID: uuid.NewString(), Name: "Workflow A", AgentID: agent.ID})
	workflowRepo.Create(context.Background(), &models.Workflow{ID: uuid.NewString(), Name: "Workflow B", AgentID: agent.ID})

	resp, err := service.ListWorkflows(context.Background(), userID, agent.ID, 0, 10)
	if err != nil {
		t.Fatalf("expected no error, got %v", err)
	}
	if len(resp) != 2 {
		t.Fatalf("expected 2 workflows, got %d", len(resp))
	}
}
