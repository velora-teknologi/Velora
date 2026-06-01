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

type TeamService interface {
	CreateTeam(ctx context.Context, userID string, req *dtos.CreateTeamRequest) (*dtos.TeamResponse, error)
	GetTeam(ctx context.Context, userID, id string) (*dtos.TeamResponse, error)
	UpdateTeam(ctx context.Context, userID, id string, req *dtos.UpdateTeamRequest) (*dtos.TeamResponse, error)
	DeleteTeam(ctx context.Context, userID, id string) error
	ListTeams(ctx context.Context, userID, tenantID string, skip, limit int) ([]*dtos.TeamResponse, error)
}

type teamService struct {
	repo       repositories.TeamRepository
	tenantRepo repositories.TenantRepository
	publisher  Publisher
	logger     *zap.SugaredLogger
}

func NewTeamService(repo repositories.TeamRepository, tenantRepo repositories.TenantRepository, publisher Publisher, logger *zap.SugaredLogger) TeamService {
	return &teamService{repo: repo, tenantRepo: tenantRepo, publisher: publisher, logger: logger}
}

func (s *teamService) CreateTeam(ctx context.Context, userID string, req *dtos.CreateTeamRequest) (*dtos.TeamResponse, error) {
	if req.Name == "" {
		return nil, customErrors.NewValidationError("Team name is required")
	}
	if req.TenantID == "" {
		return nil, customErrors.NewValidationError("Tenant ID is required")
	}

	tenant, err := s.tenantRepo.GetByID(ctx, req.TenantID)
	if err != nil {
		s.logger.Errorf("Error retrieving tenant: %v", err)
		return nil, customErrors.NewInternalError("Failed to create team")
	}
	if tenant == nil || tenant.OwnerID != userID {
		return nil, customErrors.NewForbiddenError("Tenant not found or unauthorized")
	}

	team := &models.Team{
		TenantID:    req.TenantID,
		Name:        req.Name,
		Description: req.Description,
		IsActive:    true,
	}

	if err := s.repo.Create(ctx, team); err != nil {
		s.logger.Errorf("Error creating team: %v", err)
		return nil, customErrors.NewInternalError("Failed to create team")
	}

	if s.publisher != nil {
		eventData, err := json.Marshal(map[string]interface{}{
			"team_id":   team.ID,
			"tenant_id": team.TenantID,
			"name":      team.Name,
		})
		if err == nil {
			if err := s.publisher.Publish("team.created", eventData); err != nil {
				s.logger.Warnf("Failed to publish team.created event: %v", err)
			}
		} else {
			s.logger.Warnf("Failed to marshal team.created event: %v", err)
		}
	}

	return teamToResponse(team), nil
}

func (s *teamService) GetTeam(ctx context.Context, userID, id string) (*dtos.TeamResponse, error) {
	team, err := s.repo.GetByID(ctx, id)
	if err != nil {
		s.logger.Errorf("Error retrieving team: %v", err)
		return nil, customErrors.NewInternalError("Failed to get team")
	}
	if team == nil {
		return nil, customErrors.NewNotFoundError("Team")
	}

	tenant, err := s.tenantRepo.GetByID(ctx, team.TenantID)
	if err != nil {
		s.logger.Errorf("Error retrieving tenant: %v", err)
		return nil, customErrors.NewInternalError("Failed to get team")
	}
	if tenant == nil || tenant.OwnerID != userID {
		return nil, customErrors.NewForbiddenError("Team not found or unauthorized")
	}

	return teamToResponse(team), nil
}

func (s *teamService) UpdateTeam(ctx context.Context, userID, id string, req *dtos.UpdateTeamRequest) (*dtos.TeamResponse, error) {
	team, err := s.repo.GetByID(ctx, id)
	if err != nil {
		s.logger.Errorf("Error retrieving team: %v", err)
		return nil, customErrors.NewInternalError("Failed to update team")
	}
	if team == nil {
		return nil, customErrors.NewNotFoundError("Team")
	}

	tenant, err := s.tenantRepo.GetByID(ctx, team.TenantID)
	if err != nil {
		s.logger.Errorf("Error retrieving tenant: %v", err)
		return nil, customErrors.NewInternalError("Failed to update team")
	}
	if tenant == nil || tenant.OwnerID != userID {
		return nil, customErrors.NewForbiddenError("Team not found or unauthorized")
	}

	if req.Name != "" {
		team.Name = req.Name
	}
	if req.Description != "" {
		team.Description = req.Description
	}
	if req.IsActive != nil {
		team.IsActive = *req.IsActive
	}

	if err := s.repo.Update(ctx, team); err != nil {
		s.logger.Errorf("Error updating team: %v", err)
		return nil, customErrors.NewInternalError("Failed to update team")
	}

	return teamToResponse(team), nil
}

func (s *teamService) DeleteTeam(ctx context.Context, userID, id string) error {
	team, err := s.repo.GetByID(ctx, id)
	if err != nil {
		s.logger.Errorf("Error retrieving team: %v", err)
		return customErrors.NewInternalError("Failed to delete team")
	}
	if team == nil {
		return customErrors.NewNotFoundError("Team")
	}

	tenant, err := s.tenantRepo.GetByID(ctx, team.TenantID)
	if err != nil {
		s.logger.Errorf("Error retrieving tenant: %v", err)
		return customErrors.NewInternalError("Failed to delete team")
	}
	if tenant == nil || tenant.OwnerID != userID {
		return customErrors.NewForbiddenError("Team not found or unauthorized")
	}

	if err := s.repo.Delete(ctx, id); err != nil {
		s.logger.Errorf("Error deleting team: %v", err)
		return customErrors.NewInternalError("Failed to delete team")
	}

	return nil
}

func (s *teamService) ListTeams(ctx context.Context, userID, tenantID string, skip, limit int) ([]*dtos.TeamResponse, error) {
	if tenantID == "" {
		return nil, customErrors.NewValidationError("Tenant ID is required")
	}

	tenant, err := s.tenantRepo.GetByID(ctx, tenantID)
	if err != nil {
		s.logger.Errorf("Error retrieving tenant: %v", err)
		return nil, customErrors.NewInternalError("Failed to list teams")
	}
	if tenant == nil || tenant.OwnerID != userID {
		return nil, customErrors.NewForbiddenError("Tenant not found or unauthorized")
	}

	teams, err := s.repo.List(ctx, tenantID, skip, limit)
	if err != nil {
		s.logger.Errorf("Error listing teams: %v", err)
		return nil, customErrors.NewInternalError("Failed to list teams")
	}

	var responses []*dtos.TeamResponse
	for _, team := range teams {
		responses = append(responses, teamToResponse(team))
	}

	return responses, nil
}

func teamToResponse(team *models.Team) *dtos.TeamResponse {
	return &dtos.TeamResponse{
		ID:          team.ID,
		TenantID:    team.TenantID,
		Name:        team.Name,
		Description: team.Description,
		IsActive:    team.IsActive,
	}
}
