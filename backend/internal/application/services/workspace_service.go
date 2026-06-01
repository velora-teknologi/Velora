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

type WorkspaceService interface {
	CreateWorkspace(ctx context.Context, userID string, req *dtos.CreateWorkspaceRequest) (*dtos.WorkspaceResponse, error)
	GetWorkspace(ctx context.Context, userID, id string) (*dtos.WorkspaceResponse, error)
	UpdateWorkspace(ctx context.Context, userID, id string, req *dtos.UpdateWorkspaceRequest) (*dtos.WorkspaceResponse, error)
	DeleteWorkspace(ctx context.Context, userID, id string) error
	ListWorkspaces(ctx context.Context, userID, tenantID string, skip, limit int) ([]*dtos.WorkspaceResponse, error)
}

type workspaceService struct {
	repo       repositories.WorkspaceRepository
	tenantRepo repositories.TenantRepository
	publisher  Publisher
	logger     *zap.SugaredLogger
}

func NewWorkspaceService(repo repositories.WorkspaceRepository, tenantRepo repositories.TenantRepository, publisher Publisher, logger *zap.SugaredLogger) WorkspaceService {
	return &workspaceService{repo: repo, tenantRepo: tenantRepo, publisher: publisher, logger: logger}
}

func (s *workspaceService) CreateWorkspace(ctx context.Context, userID string, req *dtos.CreateWorkspaceRequest) (*dtos.WorkspaceResponse, error) {
	if req.Name == "" {
		return nil, customErrors.NewValidationError("Workspace name is required")
	}

	tenant, err := s.tenantRepo.GetByID(ctx, req.TenantID)
	if err != nil {
		s.logger.Errorf("Error retrieving tenant: %v", err)
		return nil, customErrors.NewInternalError("Failed to create workspace")
	}
	if tenant == nil || tenant.OwnerID != userID {
		return nil, customErrors.NewForbiddenError("Tenant not found or unauthorized")
	}

	workspace := &models.Workspace{
		TenantID:    req.TenantID,
		Name:        req.Name,
		Description: req.Description,
		IsActive:    true,
	}

	if err := s.repo.Create(ctx, workspace); err != nil {
		s.logger.Errorf("Error creating workspace: %v", err)
		return nil, customErrors.NewInternalError("Failed to create workspace")
	}

	if s.publisher != nil {
		eventData, err := json.Marshal(map[string]interface{}{
			"workspace_id": workspace.ID,
			"tenant_id":    workspace.TenantID,
			"name":         workspace.Name,
		})
		if err == nil {
			if err := s.publisher.Publish("workspace.created", eventData); err != nil {
				s.logger.Warnf("Failed to publish workspace.created event: %v", err)
			}
		} else {
			s.logger.Warnf("Failed to marshal workspace.created event: %v", err)
		}
	}

	return workspaceToResponse(workspace), nil
}

func (s *workspaceService) GetWorkspace(ctx context.Context, userID, id string) (*dtos.WorkspaceResponse, error) {
	workspace, err := s.repo.GetByID(ctx, id)
	if err != nil {
		s.logger.Errorf("Error retrieving workspace: %v", err)
		return nil, customErrors.NewInternalError("Failed to get workspace")
	}
	if workspace == nil {
		return nil, customErrors.NewNotFoundError("Workspace")
	}

	tenant, err := s.tenantRepo.GetByID(ctx, workspace.TenantID)
	if err != nil {
		s.logger.Errorf("Error retrieving tenant: %v", err)
		return nil, customErrors.NewInternalError("Failed to get workspace")
	}
	if tenant == nil || tenant.OwnerID != userID {
		return nil, customErrors.NewForbiddenError("Workspace not found or unauthorized")
	}

	return workspaceToResponse(workspace), nil
}

func (s *workspaceService) UpdateWorkspace(ctx context.Context, userID, id string, req *dtos.UpdateWorkspaceRequest) (*dtos.WorkspaceResponse, error) {
	workspace, err := s.repo.GetByID(ctx, id)
	if err != nil {
		s.logger.Errorf("Error retrieving workspace: %v", err)
		return nil, customErrors.NewInternalError("Failed to update workspace")
	}
	if workspace == nil {
		return nil, customErrors.NewNotFoundError("Workspace")
	}

	tenant, err := s.tenantRepo.GetByID(ctx, workspace.TenantID)
	if err != nil {
		s.logger.Errorf("Error retrieving tenant: %v", err)
		return nil, customErrors.NewInternalError("Failed to update workspace")
	}
	if tenant == nil || tenant.OwnerID != userID {
		return nil, customErrors.NewForbiddenError("Workspace not found or unauthorized")
	}

	if req.Name != "" {
		workspace.Name = req.Name
	}
	if req.Description != "" {
		workspace.Description = req.Description
	}
	if req.IsActive != nil {
		workspace.IsActive = *req.IsActive
	}

	if err := s.repo.Update(ctx, workspace); err != nil {
		s.logger.Errorf("Error updating workspace: %v", err)
		return nil, customErrors.NewInternalError("Failed to update workspace")
	}

	return workspaceToResponse(workspace), nil
}

func (s *workspaceService) DeleteWorkspace(ctx context.Context, userID, id string) error {
	workspace, err := s.repo.GetByID(ctx, id)
	if err != nil {
		s.logger.Errorf("Error retrieving workspace: %v", err)
		return customErrors.NewInternalError("Failed to delete workspace")
	}
	if workspace == nil {
		return customErrors.NewNotFoundError("Workspace")
	}

	tenant, err := s.tenantRepo.GetByID(ctx, workspace.TenantID)
	if err != nil {
		s.logger.Errorf("Error retrieving tenant: %v", err)
		return customErrors.NewInternalError("Failed to delete workspace")
	}
	if tenant == nil || tenant.OwnerID != userID {
		return customErrors.NewForbiddenError("Workspace not found or unauthorized")
	}

	if err := s.repo.Delete(ctx, id); err != nil {
		s.logger.Errorf("Error deleting workspace: %v", err)
		return customErrors.NewInternalError("Failed to delete workspace")
	}

	return nil
}

func (s *workspaceService) ListWorkspaces(ctx context.Context, userID, tenantID string, skip, limit int) ([]*dtos.WorkspaceResponse, error) {
	tenant, err := s.tenantRepo.GetByID(ctx, tenantID)
	if err != nil {
		s.logger.Errorf("Error retrieving tenant: %v", err)
		return nil, customErrors.NewInternalError("Failed to list workspaces")
	}
	if tenant == nil || tenant.OwnerID != userID {
		return nil, customErrors.NewForbiddenError("Tenant not found or unauthorized")
	}

	workspaces, err := s.repo.List(ctx, tenantID, skip, limit)
	if err != nil {
		s.logger.Errorf("Error listing workspaces: %v", err)
		return nil, customErrors.NewInternalError("Failed to list workspaces")
	}

	var responses []*dtos.WorkspaceResponse
	for _, workspace := range workspaces {
		responses = append(responses, workspaceToResponse(workspace))
	}
	return responses, nil
}

func workspaceToResponse(workspace *models.Workspace) *dtos.WorkspaceResponse {
	return &dtos.WorkspaceResponse{
		ID:          workspace.ID,
		TenantID:    workspace.TenantID,
		Name:        workspace.Name,
		Description: workspace.Description,
		IsActive:    workspace.IsActive,
	}
}
