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

type TeamMemberService interface {
	AddTeamMember(ctx context.Context, userID string, req *dtos.CreateTeamMemberRequest) (*dtos.TeamMemberResponse, error)
	ListTeamMembers(ctx context.Context, userID, teamID string, skip, limit int) ([]*dtos.TeamMemberResponse, error)
	RemoveTeamMember(ctx context.Context, userID, teamID, memberID string) error
}

type teamMemberService struct {
	repo       repositories.TeamMemberRepository
	teamRepo   repositories.TeamRepository
	tenantRepo repositories.TenantRepository
	userRepo   repositories.UserRepository
	publisher  Publisher
	logger     *zap.SugaredLogger
}

func NewTeamMemberService(
	repo repositories.TeamMemberRepository,
	teamRepo repositories.TeamRepository,
	tenantRepo repositories.TenantRepository,
	userRepo repositories.UserRepository,
	publisher Publisher,
	logger *zap.SugaredLogger,
) TeamMemberService {
	return &teamMemberService{
		repo:       repo,
		teamRepo:   teamRepo,
		tenantRepo: tenantRepo,
		userRepo:   userRepo,
		publisher:  publisher,
		logger:     logger,
	}
}

func (s *teamMemberService) AddTeamMember(ctx context.Context, userID string, req *dtos.CreateTeamMemberRequest) (*dtos.TeamMemberResponse, error) {
	if req.TeamID == "" {
		return nil, customErrors.NewValidationError("Team ID is required")
	}
	if req.UserID == "" {
		return nil, customErrors.NewValidationError("User ID is required")
	}

	role := req.Role
	if role == "" {
		role = "member"
	}
	if role != "member" && role != "lead" {
		return nil, customErrors.NewValidationError("Invalid role")
	}

	team, err := s.teamRepo.GetByID(ctx, req.TeamID)
	if err != nil {
		s.logger.Errorf("Error retrieving team: %v", err)
		return nil, customErrors.NewInternalError("Failed to add team member")
	}
	if team == nil {
		return nil, customErrors.NewNotFoundError("Team")
	}

	tenant, err := s.tenantRepo.GetByID(ctx, team.TenantID)
	if err != nil {
		s.logger.Errorf("Error retrieving tenant: %v", err)
		return nil, customErrors.NewInternalError("Failed to add team member")
	}
	if tenant == nil || tenant.OwnerID != userID {
		return nil, customErrors.NewForbiddenError("Team not found or unauthorized")
	}

	user, err := s.userRepo.GetByID(ctx, req.UserID)
	if err != nil {
		s.logger.Errorf("Error retrieving user: %v", err)
		return nil, customErrors.NewInternalError("Failed to add team member")
	}
	if user == nil {
		return nil, customErrors.NewNotFoundError("User")
	}
	if user.TenantID != "" && user.TenantID != team.TenantID {
		return nil, customErrors.NewConflictError("User belongs to a different tenant")
	}

	existing, err := s.repo.GetByTeamAndUser(ctx, req.TeamID, req.UserID)
	if err != nil {
		s.logger.Errorf("Error checking existing team membership: %v", err)
		return nil, customErrors.NewInternalError("Failed to add team member")
	}
	if existing != nil {
		return nil, customErrors.NewConflictError("User already belongs to the team")
	}

	member := &models.TeamMember{
		TeamID: req.TeamID,
		UserID: req.UserID,
		Role:   role,
	}

	if err := s.repo.Create(ctx, member); err != nil {
		s.logger.Errorf("Error creating team member: %v", err)
		return nil, customErrors.NewInternalError("Failed to add team member")
	}

	if s.publisher != nil {
		eventData, err := json.Marshal(map[string]interface{}{
			"team_member_id": member.ID,
			"team_id":        member.TeamID,
			"user_id":        member.UserID,
			"role":           member.Role,
		})
		if err == nil {
			if err := s.publisher.Publish("team.member.added", eventData); err != nil {
				s.logger.Warnf("Failed to publish team.member.added event: %v", err)
			}
		} else {
			s.logger.Warnf("Failed to marshal team.member.added event: %v", err)
		}
	}

	return teamMemberToResponse(member), nil
}

func (s *teamMemberService) ListTeamMembers(ctx context.Context, userID, teamID string, skip, limit int) ([]*dtos.TeamMemberResponse, error) {
	if teamID == "" {
		return nil, customErrors.NewValidationError("Team ID is required")
	}

	team, err := s.teamRepo.GetByID(ctx, teamID)
	if err != nil {
		s.logger.Errorf("Error retrieving team: %v", err)
		return nil, customErrors.NewInternalError("Failed to list team members")
	}
	if team == nil {
		return nil, customErrors.NewNotFoundError("Team")
	}

	tenant, err := s.tenantRepo.GetByID(ctx, team.TenantID)
	if err != nil {
		s.logger.Errorf("Error retrieving tenant: %v", err)
		return nil, customErrors.NewInternalError("Failed to list team members")
	}
	if tenant == nil || tenant.OwnerID != userID {
		return nil, customErrors.NewForbiddenError("Team not found or unauthorized")
	}

	members, err := s.repo.ListByTeamID(ctx, teamID, skip, limit)
	if err != nil {
		s.logger.Errorf("Error listing team members: %v", err)
		return nil, customErrors.NewInternalError("Failed to list team members")
	}

	var responses []*dtos.TeamMemberResponse
	for _, member := range members {
		responses = append(responses, teamMemberToResponse(member))
	}

	return responses, nil
}

func (s *teamMemberService) RemoveTeamMember(ctx context.Context, userID, teamID, memberID string) error {
	team, err := s.teamRepo.GetByID(ctx, teamID)
	if err != nil {
		s.logger.Errorf("Error retrieving team: %v", err)
		return customErrors.NewInternalError("Failed to remove team member")
	}
	if team == nil {
		return customErrors.NewNotFoundError("Team")
	}

	tenant, err := s.tenantRepo.GetByID(ctx, team.TenantID)
	if err != nil {
		s.logger.Errorf("Error retrieving tenant: %v", err)
		return customErrors.NewInternalError("Failed to remove team member")
	}
	if tenant == nil || tenant.OwnerID != userID {
		return customErrors.NewForbiddenError("Team not found or unauthorized")
	}

	member, err := s.repo.GetByID(ctx, memberID)
	if err != nil {
		s.logger.Errorf("Error retrieving team member: %v", err)
		return customErrors.NewInternalError("Failed to remove team member")
	}
	if member == nil || member.TeamID != teamID {
		return customErrors.NewNotFoundError("Team member")
	}

	if err := s.repo.Delete(ctx, memberID); err != nil {
		s.logger.Errorf("Error deleting team member: %v", err)
		return customErrors.NewInternalError("Failed to remove team member")
	}

	return nil
}

func teamMemberToResponse(member *models.TeamMember) *dtos.TeamMemberResponse {
	return &dtos.TeamMemberResponse{
		ID:     member.ID,
		TeamID: member.TeamID,
		UserID: member.UserID,
		Role:   member.Role,
	}
}
