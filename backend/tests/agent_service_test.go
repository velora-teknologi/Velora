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

type MockAgentRepository struct {
	agents map[string]*models.Agent
}

func NewMockAgentRepository() *MockAgentRepository {
	return &MockAgentRepository{agents: make(map[string]*models.Agent)}
}

func (m *MockAgentRepository) Create(ctx context.Context, agent *models.Agent) error {
	if agent.ID == "" {
		agent.ID = uuid.NewString()
	}
	m.agents[agent.ID] = agent
	return nil
}

func (m *MockAgentRepository) GetByID(ctx context.Context, id string) (*models.Agent, error) {
	return m.agents[id], nil
}

func (m *MockAgentRepository) GetByUserID(ctx context.Context, userID string) ([]*models.Agent, error) {
	var agents []*models.Agent
	for _, agent := range m.agents {
		if agent.UserID == userID {
			agents = append(agents, agent)
		}
	}
	return agents, nil
}

func (m *MockAgentRepository) Update(ctx context.Context, agent *models.Agent) error {
	m.agents[agent.ID] = agent
	return nil
}

func (m *MockAgentRepository) Delete(ctx context.Context, id string) error {
	delete(m.agents, id)
	return nil
}

func (m *MockAgentRepository) List(ctx context.Context, userID string, skip, limit int) ([]*models.Agent, error) {
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

func TestCreateAgent(t *testing.T) {
	repo := NewMockAgentRepository()
	logger := zap.NewNop().Sugar()
	service := services.NewAgentService(repo, nil, logger)

	userID := uuid.NewString()
	req := &dtos.CreateAgentRequest{
		Name:        "Assistant",
		Description: "Personal AI assistant",
		Config:      map[string]interface{}{"model": "gpt-4"},
	}

	resp, err := service.CreateAgent(context.Background(), userID, req)
	if err != nil {
		t.Fatalf("expected no error, got %v", err)
	}
	if resp.Name != req.Name {
		t.Fatalf("expected name %s, got %s", req.Name, resp.Name)
	}
	if resp.Description != req.Description {
		t.Fatalf("expected description %s, got %s", req.Description, resp.Description)
	}
}

func TestListAgents(t *testing.T) {
	repo := NewMockAgentRepository()
	logger := zap.NewNop().Sugar()
	service := services.NewAgentService(repo, nil, logger)

	userID := uuid.NewString()
	repo.Create(context.Background(), &models.Agent{ID: uuid.NewString(), Name: "Agent A", UserID: userID})
	repo.Create(context.Background(), &models.Agent{ID: uuid.NewString(), Name: "Agent B", UserID: userID})

	resp, err := service.ListAgents(context.Background(), userID, 0, 10)
	if err != nil {
		t.Fatalf("expected no error, got %v", err)
	}
	if len(resp) != 2 {
		t.Fatalf("expected 2 agents, got %d", len(resp))
	}
}
