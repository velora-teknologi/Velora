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

type WorkflowService interface {
	CreateWorkflow(ctx context.Context, userID string, req *dtos.CreateWorkflowRequest) (*dtos.WorkflowResponse, error)
	GetWorkflow(ctx context.Context, userID, id string) (*dtos.WorkflowResponse, error)
	UpdateWorkflow(ctx context.Context, userID, id string, req *dtos.UpdateWorkflowRequest) (*dtos.WorkflowResponse, error)
	DeleteWorkflow(ctx context.Context, userID, id string) error
	ListWorkflows(ctx context.Context, userID, agentID string, skip, limit int) ([]*dtos.WorkflowResponse, error)
	ExecuteWorkflow(ctx context.Context, userID, workflowID string) (*dtos.WorkflowExecutionResponse, error)
}

type workflowService struct {
	repo      repositories.WorkflowRepository
	agentRepo repositories.AgentRepository
	publisher Publisher
	logger    *zap.SugaredLogger
}

func NewWorkflowService(repo repositories.WorkflowRepository, agentRepo repositories.AgentRepository, publisher Publisher, logger *zap.SugaredLogger) WorkflowService {
	return &workflowService{repo: repo, agentRepo: agentRepo, publisher: publisher, logger: logger}
}

func (s *workflowService) CreateWorkflow(ctx context.Context, userID string, req *dtos.CreateWorkflowRequest) (*dtos.WorkflowResponse, error) {
	if req.Name == "" {
		return nil, customErrors.NewValidationError("Workflow name is required")
	}
	if req.AgentID == "" {
		return nil, customErrors.NewValidationError("Agent ID is required")
	}

	agent, err := s.agentRepo.GetByID(ctx, req.AgentID)
	if err != nil {
		s.logger.Errorf("Error retrieving agent: %v", err)
		return nil, customErrors.NewInternalError("Failed to create workflow")
	}
	if agent == nil || agent.UserID != userID {
		return nil, customErrors.NewForbiddenError("Agent not found or unauthorized")
	}

	definitionJSON, err := json.Marshal(req.Definition)
	if err != nil {
		s.logger.Errorf("Error marshaling workflow definition: %v", err)
		return nil, customErrors.NewInternalError("Failed to create workflow")
	}

	workflow := &models.Workflow{
		TenantID:    agent.TenantID,
		Name:        req.Name,
		Description: req.Description,
		AgentID:     req.AgentID,
		Definition:  string(definitionJSON),
		IsActive:    true,
	}
	if err := s.repo.Create(ctx, workflow); err != nil {
		s.logger.Errorf("Error creating workflow: %v", err)
		return nil, customErrors.NewInternalError("Failed to create workflow")
	}

	return workflowToResponse(workflow), nil
}

func (s *workflowService) GetWorkflow(ctx context.Context, userID, id string) (*dtos.WorkflowResponse, error) {
	workflow, err := s.repo.GetByID(ctx, id)
	if err != nil {
		s.logger.Errorf("Error retrieving workflow: %v", err)
		return nil, customErrors.NewInternalError("Failed to get workflow")
	}
	if workflow == nil {
		return nil, customErrors.NewNotFoundError("Workflow")
	}

	agent, err := s.agentRepo.GetByID(ctx, workflow.AgentID)
	if err != nil {
		s.logger.Errorf("Error retrieving agent: %v", err)
		return nil, customErrors.NewInternalError("Failed to get workflow")
	}
	if agent == nil || agent.UserID != userID {
		return nil, customErrors.NewNotFoundError("Workflow")
	}

	return workflowToResponse(workflow), nil
}

func (s *workflowService) UpdateWorkflow(ctx context.Context, userID, id string, req *dtos.UpdateWorkflowRequest) (*dtos.WorkflowResponse, error) {
	workflow, err := s.repo.GetByID(ctx, id)
	if err != nil {
		s.logger.Errorf("Error retrieving workflow: %v", err)
		return nil, customErrors.NewInternalError("Failed to update workflow")
	}
	if workflow == nil {
		return nil, customErrors.NewNotFoundError("Workflow")
	}

	agent, err := s.agentRepo.GetByID(ctx, workflow.AgentID)
	if err != nil {
		s.logger.Errorf("Error retrieving agent: %v", err)
		return nil, customErrors.NewInternalError("Failed to update workflow")
	}
	if agent == nil || agent.UserID != userID {
		return nil, customErrors.NewNotFoundError("Workflow")
	}

	if req.Name != "" {
		workflow.Name = req.Name
	}
	if req.Description != "" {
		workflow.Description = req.Description
	}
	if req.Definition != nil {
		definitionJSON, err := json.Marshal(req.Definition)
		if err != nil {
			s.logger.Errorf("Error marshaling workflow definition: %v", err)
			return nil, customErrors.NewInternalError("Failed to update workflow")
		}
		workflow.Definition = string(definitionJSON)
	}
	workflow.IsActive = req.IsActive

	if err := s.repo.Update(ctx, workflow); err != nil {
		s.logger.Errorf("Error updating workflow: %v", err)
		return nil, customErrors.NewInternalError("Failed to update workflow")
	}

	return workflowToResponse(workflow), nil
}

func (s *workflowService) DeleteWorkflow(ctx context.Context, userID, id string) error {
	workflow, err := s.repo.GetByID(ctx, id)
	if err != nil {
		s.logger.Errorf("Error retrieving workflow: %v", err)
		return customErrors.NewInternalError("Failed to delete workflow")
	}
	if workflow == nil {
		return customErrors.NewNotFoundError("Workflow")
	}

	agent, err := s.agentRepo.GetByID(ctx, workflow.AgentID)
	if err != nil {
		s.logger.Errorf("Error retrieving agent: %v", err)
		return customErrors.NewInternalError("Failed to delete workflow")
	}
	if agent == nil || agent.UserID != userID {
		return customErrors.NewNotFoundError("Workflow")
	}

	if err := s.repo.Delete(ctx, id); err != nil {
		s.logger.Errorf("Error deleting workflow: %v", err)
		return customErrors.NewInternalError("Failed to delete workflow")
	}
	return nil
}

func (s *workflowService) ListWorkflows(ctx context.Context, userID, agentID string, skip, limit int) ([]*dtos.WorkflowResponse, error) {
	if agentID == "" {
		return nil, customErrors.NewValidationError("agent_id is required")
	}

	agent, err := s.agentRepo.GetByID(ctx, agentID)
	if err != nil {
		s.logger.Errorf("Error retrieving agent: %v", err)
		return nil, customErrors.NewInternalError("Failed to list workflows")
	}
	if agent == nil || agent.UserID != userID {
		return nil, customErrors.NewForbiddenError("Agent not found or unauthorized")
	}

	workflows, err := s.repo.List(ctx, agentID, skip, limit)
	if err != nil {
		s.logger.Errorf("Error listing workflows: %v", err)
		return nil, customErrors.NewInternalError("Failed to list workflows")
	}

	var responses []*dtos.WorkflowResponse
	for _, workflow := range workflows {
		responses = append(responses, workflowToResponse(workflow))
	}
	return responses, nil
}

func (s *workflowService) ExecuteWorkflow(ctx context.Context, userID, workflowID string) (*dtos.WorkflowExecutionResponse, error) {
	workflow, err := s.repo.GetByID(ctx, workflowID)
	if err != nil {
		s.logger.Errorf("Error retrieving workflow: %v", err)
		return nil, customErrors.NewInternalError("Failed to execute workflow")
	}
	if workflow == nil {
		return nil, customErrors.NewNotFoundError("Workflow")
	}

	agent, err := s.agentRepo.GetByID(ctx, workflow.AgentID)
	if err != nil {
		s.logger.Errorf("Error retrieving agent: %v", err)
		return nil, customErrors.NewInternalError("Failed to execute workflow")
	}
	if agent == nil || agent.UserID != userID {
		return nil, customErrors.NewNotFoundError("Workflow")
	}

	if s.publisher == nil {
		s.logger.Error("NATS publisher is not configured")
		return nil, customErrors.NewInternalError("Workflow runtime is unavailable")
	}

	var definition interface{}
	if workflow.Definition != "" {
		if err := json.Unmarshal([]byte(workflow.Definition), &definition); err != nil {
			s.logger.Errorf("Error unmarshalling workflow definition: %v", err)
			return nil, customErrors.NewInternalError("Failed to execute workflow")
		}
	}

	event := map[string]interface{}{
		"workflow_id": workflow.ID,
		"agent_id":    workflow.AgentID,
		"user_id":     userID,
		"definition":  definition,
	}

	payload, err := json.Marshal(event)
	if err != nil {
		s.logger.Errorf("Error marshaling workflow execution event: %v", err)
		return nil, customErrors.NewInternalError("Failed to execute workflow")
	}

	if err := s.publisher.Publish("workflow.execute", payload); err != nil {
		s.logger.Errorf("Error publishing workflow execution event: %v", err)
		return nil, customErrors.NewInternalError("Failed to queue workflow execution")
	}

	return &dtos.WorkflowExecutionResponse{
		WorkflowID: workflow.ID,
		AgentID:    workflow.AgentID,
		Status:     "queued",
		Subject:    "workflow.execute",
	}, nil
}

func workflowToResponse(workflow *models.Workflow) *dtos.WorkflowResponse {
	var definition interface{}
	if workflow.Definition != "" {
		_ = json.Unmarshal([]byte(workflow.Definition), &definition)
	}
	return &dtos.WorkflowResponse{
		ID:          workflow.ID,
		Name:        workflow.Name,
		Description: workflow.Description,
		Definition:  definition,
		IsActive:    workflow.IsActive,
	}
}
