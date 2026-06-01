package services_test

import (
	"context"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
	"go.uber.org/zap/zaptest"

	"github.com/velora-teknologi/velora/internal/application/dtos"
	"github.com/velora-teknologi/velora/internal/application/services"
	"github.com/velora-teknologi/velora/internal/domain/models"
	customErrors "github.com/velora-teknologi/velora/pkg/errors"
)

type mockTenantRepoForTeam struct {
	mock.Mock
}

func (m *mockTenantRepoForTeam) Create(ctx context.Context, tenant *models.Tenant) error { return nil }
func (m *mockTenantRepoForTeam) GetByID(ctx context.Context, id string) (*models.Tenant, error) {
	args := m.Called(ctx, id)
	if tenant, ok := args.Get(0).(*models.Tenant); ok {
		return tenant, args.Error(1)
	}
	return nil, args.Error(1)
}
func (m *mockTenantRepoForTeam) Update(ctx context.Context, tenant *models.Tenant) error { return nil }
func (m *mockTenantRepoForTeam) Delete(ctx context.Context, id string) error             { return nil }
func (m *mockTenantRepoForTeam) List(ctx context.Context, ownerID string, skip, limit int) ([]*models.Tenant, error) {
	return nil, nil
}

type mockTeamRepository struct {
	mock.Mock
}

func (m *mockTeamRepository) Create(ctx context.Context, team *models.Team) error {
	args := m.Called(ctx, team)
	return args.Error(0)
}
func (m *mockTeamRepository) GetByID(ctx context.Context, id string) (*models.Team, error) {
	args := m.Called(ctx, id)
	return args.Get(0).(*models.Team), args.Error(1)
}
func (m *mockTeamRepository) Update(ctx context.Context, team *models.Team) error {
	args := m.Called(ctx, team)
	return args.Error(0)
}
func (m *mockTeamRepository) Delete(ctx context.Context, id string) error {
	args := m.Called(ctx, id)
	return args.Error(0)
}
func (m *mockTeamRepository) List(ctx context.Context, tenantID string, skip, limit int) ([]*models.Team, error) {
	args := m.Called(ctx, tenantID, skip, limit)
	return args.Get(0).([]*models.Team), args.Error(1)
}

type mockPublisherForTeam struct {
	mock.Mock
}

func (m *mockPublisherForTeam) Publish(subject string, data []byte) error {
	args := m.Called(subject, data)
	return args.Error(0)
}

func TestCreateTeamSuccess(t *testing.T) {
	logger := zaptest.NewLogger(t).Sugar()
	tenantRepo := new(mockTenantRepoForTeam)
	teamRepo := new(mockTeamRepository)
	publisher := new(mockPublisherForTeam)

	tenant := &models.Tenant{ID: "tenant-1", OwnerID: "user-1"}
	tenantRepo.On("GetByID", mock.Anything, "tenant-1").Return(tenant, nil)
	teamRepo.On("Create", mock.Anything, mock.AnythingOfType("*models.Team")).Return(nil)
	publisher.On("Publish", "team.created", mock.Anything).Return(nil)

	service := services.NewTeamService(teamRepo, tenantRepo, publisher, logger)
	req := &dtos.CreateTeamRequest{TenantID: "tenant-1", Name: "Platform Team", Description: "Core operations"}

	result, err := service.CreateTeam(context.Background(), "user-1", req)
	assert.NoError(t, err)
	assert.Equal(t, req.Name, result.Name)
	assert.Equal(t, req.TenantID, result.TenantID)
}

func TestCreateTeamUnauthorized(t *testing.T) {
	logger := zaptest.NewLogger(t).Sugar()
	tenantRepo := new(mockTenantRepoForTeam)
	teamRepo := new(mockTeamRepository)
	publisher := new(mockPublisherForTeam)

	tenant := &models.Tenant{ID: "tenant-1", OwnerID: "user-2"}
	tenantRepo.On("GetByID", mock.Anything, "tenant-1").Return(tenant, nil)

	service := services.NewTeamService(teamRepo, tenantRepo, publisher, logger)
	req := &dtos.CreateTeamRequest{TenantID: "tenant-1", Name: "Platform Team"}

	result, err := service.CreateTeam(context.Background(), "user-1", req)
	assert.Nil(t, result)
	assert.Error(t, err)
	assert.IsType(t, &customErrors.CustomError{}, err)
	assert.Equal(t, customErrors.ForbiddenError, err.(*customErrors.CustomError).Type)
}

func TestListTeamsSuccess(t *testing.T) {
	logger := zaptest.NewLogger(t).Sugar()
	tenantRepo := new(mockTenantRepoForTeam)
	teamRepo := new(mockTeamRepository)
	publisher := new(mockPublisherForTeam)

	tenant := &models.Tenant{ID: "tenant-1", OwnerID: "user-1"}
	team := &models.Team{ID: "team-1", TenantID: "tenant-1", Name: "Platform Team", Description: "Core operations", IsActive: true}

	tenantRepo.On("GetByID", mock.Anything, "tenant-1").Return(tenant, nil)
	teamRepo.On("List", mock.Anything, "tenant-1", 0, 20).Return([]*models.Team{team}, nil)

	service := services.NewTeamService(teamRepo, tenantRepo, publisher, logger)

	result, err := service.ListTeams(context.Background(), "user-1", "tenant-1", 0, 20)
	assert.NoError(t, err)
	assert.Len(t, result, 1)
	assert.Equal(t, "team-1", result[0].ID)
}
