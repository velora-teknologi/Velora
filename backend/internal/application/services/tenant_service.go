package services

import (
	"context"
	"encoding/json"
	"strings"

	"github.com/velora-teknologi/velora/internal/application/dtos"
	"github.com/velora-teknologi/velora/internal/domain/models"
	"github.com/velora-teknologi/velora/internal/domain/repositories"
	customErrors "github.com/velora-teknologi/velora/pkg/errors"
	"go.uber.org/zap"
)

type TenantService interface {
	CreateTenant(ctx context.Context, userID string, req *dtos.CreateTenantRequest) (*dtos.TenantResponse, error)
	GetTenant(ctx context.Context, userID, id string) (*dtos.TenantResponse, error)
	UpdateTenant(ctx context.Context, userID, id string, req *dtos.UpdateTenantRequest) (*dtos.TenantResponse, error)
	DeleteTenant(ctx context.Context, userID, id string) error
	ListTenants(ctx context.Context, userID string, skip, limit int) ([]*dtos.TenantResponse, error)
}

type tenantService struct {
	repo      repositories.TenantRepository
	publisher Publisher
	logger    *zap.SugaredLogger
}

func NewTenantService(repo repositories.TenantRepository, publisher Publisher, logger *zap.SugaredLogger) TenantService {
	return &tenantService{repo: repo, publisher: publisher, logger: logger}
}

func (s *tenantService) CreateTenant(ctx context.Context, userID string, req *dtos.CreateTenantRequest) (*dtos.TenantResponse, error) {
	if req.Name == "" {
		return nil, customErrors.NewValidationError("Tenant name is required")
	}

	slug := req.Slug
	if slug == "" {
		slug = generateSlug(req.Name)
	}

	plan := req.Plan
	if plan == "" {
		plan = "free"
	}

	tenant := &models.Tenant{
		Name:        req.Name,
		Slug:        slug,
		Plan:        plan,
		Description: req.Description,
		OwnerID:     userID,
		IsActive:    true,
	}

	if err := s.repo.Create(ctx, tenant); err != nil {
		s.logger.Errorf("Error creating tenant: %v", err)
		return nil, customErrors.NewInternalError("Failed to create tenant")
	}

	if s.publisher != nil {
		eventData, err := json.Marshal(map[string]interface{}{
			"tenant_id": tenant.ID,
			"owner_id":  tenant.OwnerID,
			"slug":      tenant.Slug,
			"plan":      tenant.Plan,
		})
		if err == nil {
			if err := s.publisher.Publish("tenant.created", eventData); err != nil {
				s.logger.Warnf("Failed to publish tenant.created event: %v", err)
			}
		} else {
			s.logger.Warnf("Failed to marshal tenant.created event: %v", err)
		}
	}

	return tenantToResponse(tenant), nil
}

func generateSlug(name string) string {
	slug := strings.ToLower(name)
	slug = strings.TrimSpace(slug)
	slug = strings.ReplaceAll(slug, " ", "-")
	return slug
}

func (s *tenantService) GetTenant(ctx context.Context, userID, id string) (*dtos.TenantResponse, error) {
	tenant, err := s.repo.GetByID(ctx, id)
	if err != nil {
		s.logger.Errorf("Error retrieving tenant: %v", err)
		return nil, customErrors.NewInternalError("Failed to get tenant")
	}
	if tenant == nil || tenant.OwnerID != userID {
		return nil, customErrors.NewNotFoundError("Tenant")
	}
	return tenantToResponse(tenant), nil
}

func (s *tenantService) UpdateTenant(ctx context.Context, userID, id string, req *dtos.UpdateTenantRequest) (*dtos.TenantResponse, error) {
	tenant, err := s.repo.GetByID(ctx, id)
	if err != nil {
		s.logger.Errorf("Error retrieving tenant: %v", err)
		return nil, customErrors.NewInternalError("Failed to update tenant")
	}
	if tenant == nil || tenant.OwnerID != userID {
		return nil, customErrors.NewNotFoundError("Tenant")
	}

	if req.Name != "" {
		tenant.Name = req.Name
	}
	if req.Description != "" {
		tenant.Description = req.Description
	}
	if req.IsActive != nil {
		tenant.IsActive = *req.IsActive
	}

	if err := s.repo.Update(ctx, tenant); err != nil {
		s.logger.Errorf("Error updating tenant: %v", err)
		return nil, customErrors.NewInternalError("Failed to update tenant")
	}

	return tenantToResponse(tenant), nil
}

func (s *tenantService) DeleteTenant(ctx context.Context, userID, id string) error {
	tenant, err := s.repo.GetByID(ctx, id)
	if err != nil {
		s.logger.Errorf("Error retrieving tenant: %v", err)
		return customErrors.NewInternalError("Failed to delete tenant")
	}
	if tenant == nil || tenant.OwnerID != userID {
		return customErrors.NewNotFoundError("Tenant")
	}
	if err := s.repo.Delete(ctx, id); err != nil {
		s.logger.Errorf("Error deleting tenant: %v", err)
		return customErrors.NewInternalError("Failed to delete tenant")
	}
	return nil
}

func (s *tenantService) ListTenants(ctx context.Context, userID string, skip, limit int) ([]*dtos.TenantResponse, error) {
	tenants, err := s.repo.List(ctx, userID, skip, limit)
	if err != nil {
		s.logger.Errorf("Error listing tenants: %v", err)
		return nil, customErrors.NewInternalError("Failed to list tenants")
	}

	var responses []*dtos.TenantResponse
	for _, tenant := range tenants {
		responses = append(responses, tenantToResponse(tenant))
	}
	return responses, nil
}

func tenantToResponse(tenant *models.Tenant) *dtos.TenantResponse {
	return &dtos.TenantResponse{
		ID:          tenant.ID,
		Name:        tenant.Name,
		Slug:        tenant.Slug,
		Plan:        tenant.Plan,
		Description: tenant.Description,
		OwnerID:     tenant.OwnerID,
		IsActive:    tenant.IsActive,
	}
}
