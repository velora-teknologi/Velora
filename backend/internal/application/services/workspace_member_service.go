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

type WorkspaceMemberService interface {
	AddWorkspaceMember(ctx context.Context, userID string, req *dtos.CreateWorkspaceMemberRequest) (*dtos.WorkspaceMemberResponse, error)
	ListWorkspaceMembers(ctx context.Context, userID, workspaceID string, skip, limit int) ([]*dtos.WorkspaceMemberResponse, error)
	RemoveWorkspaceMember(ctx context.Context, userID, workspaceID, memberID string) error
}

type workspaceMemberService struct {
	repo          repositories.WorkspaceMemberRepository
	workspaceRepo repositories.WorkspaceRepository
	tenantRepo    repositories.TenantRepository
	userRepo      repositories.UserRepository
	publisher     Publisher
	logger        *zap.SugaredLogger
}

func NewWorkspaceMemberService(
	repo repositories.WorkspaceMemberRepository,
	workspaceRepo repositories.WorkspaceRepository,
	tenantRepo repositories.TenantRepository,
	userRepo repositories.UserRepository,
	publisher Publisher,
	logger *zap.SugaredLogger,
) WorkspaceMemberService {
	return &workspaceMemberService{
		repo:          repo,
		workspaceRepo: workspaceRepo,
		tenantRepo:    tenantRepo,
		userRepo:      userRepo,
		publisher:     publisher,
		logger:        logger,
	}
}

func (s *workspaceMemberService) AddWorkspaceMember(ctx context.Context, userID string, req *dtos.CreateWorkspaceMemberRequest) (*dtos.WorkspaceMemberResponse, error) {
	if req.WorkspaceID == "" {
		return nil, customErrors.NewValidationError("Workspace ID is required")
	}
	if req.UserID == "" {
		return nil, customErrors.NewValidationError("User ID is required")
	}

	role := req.Role
	if role == "" {
		role = "member"
	}
	if role != "member" && role != "manager" {
		return nil, customErrors.NewValidationError("Invalid role")
	}

	workspace, err := s.workspaceRepo.GetByID(ctx, req.WorkspaceID)
	if err != nil {
		s.logger.Errorf("Error retrieving workspace: %v", err)
		return nil, customErrors.NewInternalError("Failed to add workspace member")
	}
	if workspace == nil {
		return nil, customErrors.NewNotFoundError("Workspace")
	}

	tenant, err := s.tenantRepo.GetByID(ctx, workspace.TenantID)
	if err != nil {
		s.logger.Errorf("Error retrieving tenant: %v", err)
		return nil, customErrors.NewInternalError("Failed to add workspace member")
	}
	if tenant == nil || tenant.OwnerID != userID {
		return nil, customErrors.NewForbiddenError("Workspace not found or unauthorized")
	}

	user, err := s.userRepo.GetByID(ctx, req.UserID)
	if err != nil {
		s.logger.Errorf("Error retrieving user: %v", err)
		return nil, customErrors.NewInternalError("Failed to add workspace member")
	}
	if user == nil {
		return nil, customErrors.NewNotFoundError("User")
	}
	if user.TenantID != "" && user.TenantID != workspace.TenantID {
		return nil, customErrors.NewConflictError("User belongs to a different tenant")
	}

	existing, err := s.repo.GetByWorkspaceAndUser(ctx, req.WorkspaceID, req.UserID)
	if err != nil {
		s.logger.Errorf("Error checking existing workspace membership: %v", err)
		return nil, customErrors.NewInternalError("Failed to add workspace member")
	}
	if existing != nil {
		return nil, customErrors.NewConflictError("User already belongs to the workspace")
	}

	member := &models.WorkspaceMember{
		WorkspaceID: req.WorkspaceID,
		UserID:      req.UserID,
		Role:        role,
	}

	if err := s.repo.Create(ctx, member); err != nil {
		s.logger.Errorf("Error creating workspace member: %v", err)
		return nil, customErrors.NewInternalError("Failed to add workspace member")
	}

	if s.publisher != nil {
		eventData, err := json.Marshal(map[string]interface{}{
			"workspace_member_id": member.ID,
			"workspace_id":        member.WorkspaceID,
			"user_id":             member.UserID,
			"role":                member.Role,
		})
		if err == nil {
			if err := s.publisher.Publish("workspace.member.added", eventData); err != nil {
				s.logger.Warnf("Failed to publish workspace.member.added event: %v", err)
			}
		} else {
			s.logger.Warnf("Failed to marshal workspace.member.added event: %v", err)
		}
	}

	return workspaceMemberToResponse(member), nil
}

func (s *workspaceMemberService) ListWorkspaceMembers(ctx context.Context, userID, workspaceID string, skip, limit int) ([]*dtos.WorkspaceMemberResponse, error) {
	if workspaceID == "" {
		return nil, customErrors.NewValidationError("Workspace ID is required")
	}

	workspace, err := s.workspaceRepo.GetByID(ctx, workspaceID)
	if err != nil {
		s.logger.Errorf("Error retrieving workspace: %v", err)
		return nil, customErrors.NewInternalError("Failed to list workspace members")
	}
	if workspace == nil {
		return nil, customErrors.NewNotFoundError("Workspace")
	}

	tenant, err := s.tenantRepo.GetByID(ctx, workspace.TenantID)
	if err != nil {
		s.logger.Errorf("Error retrieving tenant: %v", err)
		return nil, customErrors.NewInternalError("Failed to list workspace members")
	}
	if tenant == nil || tenant.OwnerID != userID {
		return nil, customErrors.NewForbiddenError("Workspace not found or unauthorized")
	}

	members, err := s.repo.ListByWorkspaceID(ctx, workspaceID, skip, limit)
	if err != nil {
		s.logger.Errorf("Error listing workspace members: %v", err)
		return nil, customErrors.NewInternalError("Failed to list workspace members")
	}

	var responses []*dtos.WorkspaceMemberResponse
	for _, member := range members {
		responses = append(responses, workspaceMemberToResponse(member))
	}

	return responses, nil
}

func (s *workspaceMemberService) RemoveWorkspaceMember(ctx context.Context, userID, workspaceID, memberID string) error {
	workspace, err := s.workspaceRepo.GetByID(ctx, workspaceID)
	if err != nil {
		s.logger.Errorf("Error retrieving workspace: %v", err)
		return customErrors.NewInternalError("Failed to remove workspace member")
	}
	if workspace == nil {
		return customErrors.NewNotFoundError("Workspace")
	}

	tenant, err := s.tenantRepo.GetByID(ctx, workspace.TenantID)
	if err != nil {
		s.logger.Errorf("Error retrieving tenant: %v", err)
		return customErrors.NewInternalError("Failed to remove workspace member")
	}
	if tenant == nil || tenant.OwnerID != userID {
		return customErrors.NewForbiddenError("Workspace not found or unauthorized")
	}

	member, err := s.repo.GetByID(ctx, memberID)
	if err != nil {
		s.logger.Errorf("Error retrieving workspace member: %v", err)
		return customErrors.NewInternalError("Failed to remove workspace member")
	}
	if member == nil || member.WorkspaceID != workspaceID {
		return customErrors.NewNotFoundError("Workspace member")
	}

	if err := s.repo.Delete(ctx, memberID); err != nil {
		s.logger.Errorf("Error deleting workspace member: %v", err)
		return customErrors.NewInternalError("Failed to remove workspace member")
	}

	return nil
}

func workspaceMemberToResponse(member *models.WorkspaceMember) *dtos.WorkspaceMemberResponse {
	return &dtos.WorkspaceMemberResponse{
		ID:          member.ID,
		WorkspaceID: member.WorkspaceID,
		UserID:      member.UserID,
		Role:        member.Role,
	}
}
