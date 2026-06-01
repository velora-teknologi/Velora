package services

import (
	"context"
	"crypto/rand"
	"encoding/hex"
	"encoding/json"
	"strings"
	"time"

	"github.com/velora-teknologi/velora/internal/application/dtos"
	"github.com/velora-teknologi/velora/internal/domain/models"
	"github.com/velora-teknologi/velora/internal/domain/repositories"
	customErrors "github.com/velora-teknologi/velora/pkg/errors"
	"go.uber.org/zap"
	"golang.org/x/crypto/bcrypt"
)

type InvitationService interface {
	CreateInvitation(ctx context.Context, userID string, req *dtos.CreateInvitationRequest) (*dtos.InvitationResponse, error)
	ListInvitations(ctx context.Context, userID, tenantID string, skip, limit int) ([]*dtos.InvitationResponse, error)
	AcceptInvitation(ctx context.Context, token string, req *dtos.AcceptInvitationRequest) (*dtos.AcceptInvitationResponse, error)
}

type invitationService struct {
	repo       repositories.InvitationRepository
	tenantRepo repositories.TenantRepository
	userRepo   repositories.UserRepository
	publisher  Publisher
	logger     *zap.SugaredLogger
}

func NewInvitationService(repo repositories.InvitationRepository, tenantRepo repositories.TenantRepository, userRepo repositories.UserRepository, publisher Publisher, logger *zap.SugaredLogger) InvitationService {
	return &invitationService{repo: repo, tenantRepo: tenantRepo, userRepo: userRepo, publisher: publisher, logger: logger}
}

func (s *invitationService) CreateInvitation(ctx context.Context, userID string, req *dtos.CreateInvitationRequest) (*dtos.InvitationResponse, error) {
	if req.Email == "" {
		return nil, customErrors.NewValidationError("Invitation email is required")
	}

	tenant, err := s.tenantRepo.GetByID(ctx, req.TenantID)
	if err != nil {
		s.logger.Errorf("Error retrieving tenant: %v", err)
		return nil, customErrors.NewInternalError("Failed to create invitation")
	}
	if tenant == nil || tenant.OwnerID != userID {
		return nil, customErrors.NewForbiddenError("Tenant not found or unauthorized")
	}

	role := req.Role
	if role == "" {
		role = "member"
	}

	token, err := generateToken(32)
	if err != nil {
		s.logger.Errorf("Error generating invitation token: %v", err)
		return nil, customErrors.NewInternalError("Failed to create invitation")
	}

	invitation := &models.Invitation{
		TenantID:  req.TenantID,
		Email:     strings.ToLower(req.Email),
		Token:     token,
		Role:      role,
		Status:    "pending",
		InvitedBy: userID,
	}

	if req.ExpiresIn > 0 {
		expiresAt := time.Now().Add(time.Duration(req.ExpiresIn) * time.Second)
		invitation.ExpiresAt = &expiresAt
	}

	if err := s.repo.Create(ctx, invitation); err != nil {
		s.logger.Errorf("Error creating invitation: %v", err)
		return nil, customErrors.NewInternalError("Failed to create invitation")
	}

	if s.publisher != nil {
		eventData, err := json.Marshal(map[string]interface{}{
			"invitation_id": invitation.ID,
			"tenant_id":     invitation.TenantID,
			"email":         invitation.Email,
			"role":          invitation.Role,
		})
		if err == nil {
			if err := s.publisher.Publish("tenant.invitation.created", eventData); err != nil {
				s.logger.Warnf("Failed to publish tenant.invitation.created event: %v", err)
			}
		} else {
			s.logger.Warnf("Failed to marshal tenant.invitation.created event: %v", err)
		}
	}

	return invitationToResponse(invitation), nil
}

func (s *invitationService) ListInvitations(ctx context.Context, userID, tenantID string, skip, limit int) ([]*dtos.InvitationResponse, error) {
	tenant, err := s.tenantRepo.GetByID(ctx, tenantID)
	if err != nil {
		s.logger.Errorf("Error retrieving tenant: %v", err)
		return nil, customErrors.NewInternalError("Failed to list invitations")
	}
	if tenant == nil || tenant.OwnerID != userID {
		return nil, customErrors.NewForbiddenError("Tenant not found or unauthorized")
	}

	invitations, err := s.repo.ListByTenantID(ctx, tenantID, skip, limit)
	if err != nil {
		s.logger.Errorf("Error listing invitations: %v", err)
		return nil, customErrors.NewInternalError("Failed to list invitations")
	}

	var responses []*dtos.InvitationResponse
	for _, invitation := range invitations {
		responses = append(responses, invitationToResponse(invitation))
	}
	return responses, nil
}

func (s *invitationService) AcceptInvitation(ctx context.Context, token string, req *dtos.AcceptInvitationRequest) (*dtos.AcceptInvitationResponse, error) {
	invitation, err := s.repo.GetByToken(ctx, token)
	if err != nil {
		s.logger.Errorf("Error retrieving invitation: %v", err)
		return nil, customErrors.NewInternalError("Failed to accept invitation")
	}
	if invitation == nil || invitation.Status != "pending" {
		return nil, customErrors.NewNotFoundError("Invitation")
	}
	if invitation.ExpiresAt != nil && time.Now().After(*invitation.ExpiresAt) {
		return nil, customErrors.NewValidationError("Invitation has expired")
	}
	if strings.ToLower(req.Email) != strings.ToLower(invitation.Email) {
		return nil, customErrors.NewForbiddenError("Invitation email does not match")
	}

	user, err := s.userRepo.GetByEmail(ctx, req.Email)
	if err != nil {
		s.logger.Errorf("Error retrieving user: %v", err)
		return nil, customErrors.NewInternalError("Failed to accept invitation")
	}

	if user == nil {
		hashedPassword, err := bcrypt.GenerateFromPassword([]byte(req.Password), bcrypt.DefaultCost)
		if err != nil {
			s.logger.Errorf("Error hashing password: %v", err)
			return nil, customErrors.NewInternalError("Failed to accept invitation")
		}

		user = &models.User{
			Email:    strings.ToLower(req.Email),
			Password: string(hashedPassword),
			FullName: req.FullName,
			TenantID: invitation.TenantID,
			Role:     invitation.Role,
			IsActive: true,
		}
		if err := s.userRepo.Create(ctx, user); err != nil {
			s.logger.Errorf("Error creating invited user: %v", err)
			return nil, customErrors.NewInternalError("Failed to accept invitation")
		}
	} else {
		if user.TenantID != "" && user.TenantID != invitation.TenantID {
			return nil, customErrors.NewConflictError("User already belongs to another tenant")
		}
		user.TenantID = invitation.TenantID
		user.Role = invitation.Role
		if err := s.userRepo.Update(ctx, user); err != nil {
			s.logger.Errorf("Error updating invited user: %v", err)
			return nil, customErrors.NewInternalError("Failed to accept invitation")
		}
	}

	invitation.Status = "accepted"
	acceptedAt := time.Now()
	invitation.AcceptedAt = &acceptedAt
	if err := s.repo.Update(ctx, invitation); err != nil {
		s.logger.Errorf("Error updating invitation status: %v", err)
		return nil, customErrors.NewInternalError("Failed to accept invitation")
	}

	if s.publisher != nil {
		eventData, err := json.Marshal(map[string]interface{}{
			"invitation_id": invitation.ID,
			"tenant_id":     invitation.TenantID,
			"user_email":    user.Email,
		})
		if err == nil {
			if err := s.publisher.Publish("tenant.invitation.accepted", eventData); err != nil {
				s.logger.Warnf("Failed to publish tenant.invitation.accepted event: %v", err)
			}
		} else {
			s.logger.Warnf("Failed to marshal tenant.invitation.accepted event: %v", err)
		}
	}

	return &dtos.AcceptInvitationResponse{
		UserID: user.ID,
		Email:  user.Email,
		Role:   user.Role,
	}, nil
}

func generateToken(length int) (string, error) {
	buffer := make([]byte, length)
	if _, err := rand.Read(buffer); err != nil {
		return "", err
	}
	return hex.EncodeToString(buffer), nil
}

func invitationToResponse(invitation *models.Invitation) *dtos.InvitationResponse {
	var expiresAt *string
	var acceptedAt *string
	if invitation.ExpiresAt != nil {
		t := invitation.ExpiresAt.UTC().Format(time.RFC3339)
		expiresAt = &t
	}
	if invitation.AcceptedAt != nil {
		t := invitation.AcceptedAt.UTC().Format(time.RFC3339)
		acceptedAt = &t
	}
	return &dtos.InvitationResponse{
		ID:         invitation.ID,
		TenantID:   invitation.TenantID,
		Email:      invitation.Email,
		Role:       invitation.Role,
		Status:     invitation.Status,
		Token:      invitation.Token,
		InvitedBy:  invitation.InvitedBy,
		ExpiresAt:  expiresAt,
		AcceptedAt: acceptedAt,
	}
}
