package services

import (
	"context"
	"encoding/json"

	"github.com/velora-teknologi/velora/internal/application/dtos"
	"github.com/velora-teknologi/velora/internal/domain/models"
	"github.com/velora-teknologi/velora/internal/domain/repositories"
	customErrors "github.com/velora-teknologi/velora/pkg/errors"
	"go.uber.org/zap"
)

type AgentService interface {
	CreateAgent(ctx context.Context, userID string, req *dtos.CreateAgentRequest) (*dtos.AgentResponse, error)
	GetAgent(ctx context.Context, userID, id string) (*dtos.AgentResponse, error)
	UpdateAgent(ctx context.Context, userID, id string, req *dtos.UpdateAgentRequest) (*dtos.AgentResponse, error)
	DeleteAgent(ctx context.Context, userID, id string) error
	ListAgents(ctx context.Context, userID string, skip, limit int) ([]*dtos.AgentResponse, error)
	ExecuteAgent(ctx context.Context, userID, agentID string, input interface{}) (*dtos.AgentExecutionResponse, error)
}

type agentService struct {
	repo      repositories.AgentRepository
	publisher Publisher
	logger    *zap.SugaredLogger
}

func NewAgentService(repo repositories.AgentRepository, publisher Publisher, logger *zap.SugaredLogger) AgentService {
	return &agentService{repo: repo, publisher: publisher, logger: logger}
}

func (s *agentService) CreateAgent(ctx context.Context, userID string, req *dtos.CreateAgentRequest) (*dtos.AgentResponse, error) {
	if req.Name == "" {
		return nil, customErrors.NewValidationError("Agent name is required")
	}

	configJSON, err := json.Marshal(req.Config)
	if err != nil {
		s.logger.Errorf("Error marshaling agent config: %v", err)
		return nil, customErrors.NewInternalError("Failed to create agent")
	}

	agent := &models.Agent{
		Name:        req.Name,
		Description: req.Description,
		UserID:      userID,
		Config:      string(configJSON),
		Status:      "active",
	}

	if err := s.repo.Create(ctx, agent); err != nil {
		s.logger.Errorf("Error creating agent: %v", err)
		return nil, customErrors.NewInternalError("Failed to create agent")
	}

	return agentToResponse(agent), nil
}

func (s *agentService) GetAgent(ctx context.Context, userID, id string) (*dtos.AgentResponse, error) {
	agent, err := s.repo.GetByID(ctx, id)
	if err != nil {
		s.logger.Errorf("Error retrieving agent: %v", err)
		return nil, customErrors.NewInternalError("Failed to get agent")
	}
	if agent == nil || agent.UserID != userID {
		return nil, customErrors.NewNotFoundError("Agent")
	}
	return agentToResponse(agent), nil
}

func (s *agentService) UpdateAgent(ctx context.Context, userID, id string, req *dtos.UpdateAgentRequest) (*dtos.AgentResponse, error) {
	agent, err := s.repo.GetByID(ctx, id)
	if err != nil {
		s.logger.Errorf("Error retrieving agent: %v", err)
		return nil, customErrors.NewInternalError("Failed to update agent")
	}
	if agent == nil || agent.UserID != userID {
		return nil, customErrors.NewNotFoundError("Agent")
	}

	if req.Name != "" {
		agent.Name = req.Name
	}
	if req.Description != "" {
		agent.Description = req.Description
	}
	if req.Config != nil {
		configJSON, err := json.Marshal(req.Config)
		if err != nil {
			s.logger.Errorf("Error marshaling agent config: %v", err)
			return nil, customErrors.NewInternalError("Failed to update agent")
		}
		agent.Config = string(configJSON)
	}

	if err := s.repo.Update(ctx, agent); err != nil {
		s.logger.Errorf("Error updating agent: %v", err)
		return nil, customErrors.NewInternalError("Failed to update agent")
	}

	return agentToResponse(agent), nil
}

func (s *agentService) DeleteAgent(ctx context.Context, userID, id string) error {
	agent, err := s.repo.GetByID(ctx, id)
	if err != nil {
		s.logger.Errorf("Error retrieving agent: %v", err)
		return customErrors.NewInternalError("Failed to delete agent")
	}
	if agent == nil || agent.UserID != userID {
		return customErrors.NewNotFoundError("Agent")
	}
	if err := s.repo.Delete(ctx, id); err != nil {
		s.logger.Errorf("Error deleting agent: %v", err)
		return customErrors.NewInternalError("Failed to delete agent")
	}
	return nil
}

func (s *agentService) ListAgents(ctx context.Context, userID string, skip, limit int) ([]*dtos.AgentResponse, error) {
	agents, err := s.repo.List(ctx, userID, skip, limit)
	if err != nil {
		s.logger.Errorf("Error listing agents: %v", err)
		return nil, customErrors.NewInternalError("Failed to list agents")
	}

	var responses []*dtos.AgentResponse
	for _, agent := range agents {
		responses = append(responses, agentToResponse(agent))
	}
	return responses, nil
}

func (s *agentService) ExecuteAgent(ctx context.Context, userID, agentID string, input interface{}) (*dtos.AgentExecutionResponse, error) {
	agent, err := s.repo.GetByID(ctx, agentID)
	if err != nil {
		s.logger.Errorf("Error retrieving agent: %v", err)
		return nil, customErrors.NewInternalError("Failed to execute agent")
	}
	if agent == nil || agent.UserID != userID {
		return nil, customErrors.NewNotFoundError("Agent")
	}

	if s.publisher == nil {
		s.logger.Error("Publisher not configured")
		return nil, customErrors.NewInternalError("Agent runtime unavailable")
	}

	// build event
	event := map[string]interface{}{
		"agent_id": agent.ID,
		"user_id":  userID,
		"config":   agent.Config,
		"input":    input,
	}
	payload, err := json.Marshal(event)
	if err != nil {
		s.logger.Errorf("Error marshaling agent execution event: %v", err)
		return nil, customErrors.NewInternalError("Failed to execute agent")
	}

	if err := s.publisher.Publish("agent.execute", payload); err != nil {
		s.logger.Errorf("Error publishing agent execution event: %v", err)
		return nil, customErrors.NewInternalError("Failed to queue agent execution")
	}

	return &dtos.AgentExecutionResponse{AgentID: agent.ID, Status: "queued", Subject: "agent.execute"}, nil
}

func agentToResponse(agent *models.Agent) *dtos.AgentResponse {
	var config interface{}
	if agent.Config != "" {
		_ = json.Unmarshal([]byte(agent.Config), &config)
	}
	return &dtos.AgentResponse{
		ID:          agent.ID,
		Name:        agent.Name,
		Description: agent.Description,
		Status:      agent.Status,
		Config:      config,
	}
}
