package services

import (
	"context"

	"github.com/velora-teknologi/velora/internal/application/dtos"
	"github.com/velora-teknologi/velora/internal/domain/repositories"
	customErrors "github.com/velora-teknologi/velora/pkg/errors"
	"go.uber.org/zap"
)

type TenantMemberService interface {
	ListTenantMembers(ctx context.Context, userID, tenantID string, skip, limit int) ([]*dtos.TenantMemberResponse, error)
}

type tenantMemberService struct {
	userRepo   repositories.UserRepository
	tenantRepo repositories.TenantRepository
	logger     *zap.SugaredLogger
}

func NewTenantMemberService(userRepo repositories.UserRepository, tenantRepo repositories.TenantRepository, logger *zap.SugaredLogger) TenantMemberService {
	return &tenantMemberService{userRepo: userRepo, tenantRepo: tenantRepo, logger: logger}
}

func (s *tenantMemberService) ListTenantMembers(ctx context.Context, userID, tenantID string, skip, limit int) ([]*dtos.TenantMemberResponse, error) {
	if tenantID == "" {
		return nil, customErrors.NewValidationError("Tenant ID is required")
	}

	tenant, err := s.tenantRepo.GetByID(ctx, tenantID)
	if err != nil {
		s.logger.Errorf("Error retrieving tenant: %v", err)
		return nil, customErrors.NewInternalError("Failed to list tenant members")
	}
	if tenant == nil || tenant.OwnerID != userID {
		return nil, customErrors.NewForbiddenError("Tenant not found or unauthorized")
	}

	users, err := s.userRepo.ListByTenantID(ctx, tenantID, skip, limit)
	if err != nil {
		s.logger.Errorf("Error listing tenant users: %v", err)
		return nil, customErrors.NewInternalError("Failed to list tenant members")
	}

	var responses []*dtos.TenantMemberResponse
	for _, user := range users {
		responses = append(responses, &dtos.TenantMemberResponse{
			ID:       user.ID,
			Email:    user.Email,
			FullName: user.FullName,
			Role:     user.Role,
			IsActive: user.IsActive,
		})
	}

	return responses, nil
}
