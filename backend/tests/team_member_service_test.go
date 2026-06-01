package services_test

import (
	"context"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
	"github.com/velora-teknologi/velora/internal/application/dtos"
	"github.com/velora-teknologi/velora/internal/application/services"
	"github.com/velora-teknologi/velora/internal/domain/models"
	customErrors "github.com/velora-teknologi/velora/pkg/errors"
	"go.uber.org/zap/zaptest"
)

type mockTeamMemberRepo struct {
	mock.Mock
}

func (m *mockTeamMemberRepo) Create(ctx context.Context, member *models.TeamMember) error {
	return m.Called(ctx, member).Error(0)
}
func (m *mockTeamMemberRepo) GetByID(ctx context.Context, id string) (*models.TeamMember, error) {
	args := m.Called(ctx, id)
	if member, ok := args.Get(0).(*models.TeamMember); ok {
		return member, args.Error(1)
	}
	return nil, args.Error(1)
}
func (m *mockTeamMemberRepo) GetByTeamAndUser(ctx context.Context, teamID, userID string) (*models.TeamMember, error) {
	args := m.Called(ctx, teamID, userID)
	if member, ok := args.Get(0).(*models.TeamMember); ok {
		return member, args.Error(1)
	}
	return nil, args.Error(1)
}
func (m *mockTeamMemberRepo) Delete(ctx context.Context, id string) error {
	return m.Called(ctx, id).Error(0)
}
func (m *mockTeamMemberRepo) ListByTeamID(ctx context.Context, teamID string, skip, limit int) ([]*models.TeamMember, error) {
	args := m.Called(ctx, teamID, skip, limit)
	if members, ok := args.Get(0).([]*models.TeamMember); ok {
		return members, args.Error(1)
	}
	return nil, args.Error(1)
}

type mockTeamRepoForMember struct {
	mock.Mock
}

func (m *mockTeamRepoForMember) Create(ctx context.Context, team *models.Team) error { return nil }
func (m *mockTeamRepoForMember) GetByID(ctx context.Context, id string) (*models.Team, error) {
	args := m.Called(ctx, id)
	if team, ok := args.Get(0).(*models.Team); ok {
		return team, args.Error(1)
	}
	return nil, args.Error(1)
}
func (m *mockTeamRepoForMember) Update(ctx context.Context, team *models.Team) error { return nil }
func (m *mockTeamRepoForMember) Delete(ctx context.Context, id string) error         { return nil }
func (m *mockTeamRepoForMember) List(ctx context.Context, tenantID string, skip, limit int) ([]*models.Team, error) {
	return nil, nil
}

type mockTenantRepoForMember struct {
	mock.Mock
}

func (m *mockTenantRepoForMember) Create(ctx context.Context, tenant *models.Tenant) error {
	return nil
}
func (m *mockTenantRepoForMember) GetByID(ctx context.Context, id string) (*models.Tenant, error) {
	args := m.Called(ctx, id)
	if tenant, ok := args.Get(0).(*models.Tenant); ok {
		return tenant, args.Error(1)
	}
	return nil, args.Error(1)
}
func (m *mockTenantRepoForMember) Update(ctx context.Context, tenant *models.Tenant) error {
	return nil
}
func (m *mockTenantRepoForMember) Delete(ctx context.Context, id string) error { return nil }
func (m *mockTenantRepoForMember) List(ctx context.Context, ownerID string, skip, limit int) ([]*models.Tenant, error) {
	return nil, nil
}

type mockUserRepoForMember struct {
	mock.Mock
}

func (m *mockUserRepoForMember) Create(ctx context.Context, user *models.User) error { return nil }
func (m *mockUserRepoForMember) GetByID(ctx context.Context, id string) (*models.User, error) {
	args := m.Called(ctx, id)
	if user, ok := args.Get(0).(*models.User); ok {
		return user, args.Error(1)
	}
	return nil, args.Error(1)
}
func (m *mockUserRepoForMember) GetByEmail(ctx context.Context, email string) (*models.User, error) {
	return nil, nil
}
func (m *mockUserRepoForMember) Update(ctx context.Context, user *models.User) error { return nil }
func (m *mockUserRepoForMember) Delete(ctx context.Context, id string) error         { return nil }
func (m *mockUserRepoForMember) List(ctx context.Context, skip, limit int) ([]*models.User, error) {
	return nil, nil
}
func (m *mockUserRepoForMember) ListByTenantID(ctx context.Context, tenantID string, skip, limit int) ([]*models.User, error) {
	return nil, nil
}

type mockPublisherForMember struct {
	mock.Mock
}

func (m *mockPublisherForMember) Publish(subject string, data []byte) error {
	return m.Called(subject, data).Error(0)
}

func TestAddTeamMemberSuccess(t *testing.T) {
	logger := zaptest.NewLogger(t).Sugar()
	teamMemberRepo := new(mockTeamMemberRepo)
	teamRepo := new(mockTeamRepoForMember)
	tenantRepo := new(mockTenantRepoForMember)
	userRepo := new(mockUserRepoForMember)
	publisher := new(mockPublisherForMember)

	team := &models.Team{ID: "team-1", TenantID: "tenant-1"}
	tenant := &models.Tenant{ID: "tenant-1", OwnerID: "owner-1"}
	user := &models.User{ID: "user-2", TenantID: "tenant-1"}

	teamRepo.On("GetByID", mock.Anything, "team-1").Return(team, nil)
	tenantRepo.On("GetByID", mock.Anything, "tenant-1").Return(tenant, nil)
	userRepo.On("GetByID", mock.Anything, "user-2").Return(user, nil)
	teamMemberRepo.On("GetByTeamAndUser", mock.Anything, "team-1", "user-2").Return(nil, nil)
	teamMemberRepo.On("Create", mock.Anything, mock.AnythingOfType("*models.TeamMember")).Return(nil)
	publisher.On("Publish", "team.member.added", mock.Anything).Return(nil)

	service := services.NewTeamMemberService(teamMemberRepo, teamRepo, tenantRepo, userRepo, publisher, logger)
	req := &dtos.CreateTeamMemberRequest{TeamID: "team-1", UserID: "user-2", Role: "member"}

	result, err := service.AddTeamMember(context.Background(), "owner-1", req)
	assert.NoError(t, err)
	assert.Equal(t, "team-1", result.TeamID)
	assert.Equal(t, "user-2", result.UserID)
	assert.Equal(t, "member", result.Role)
}

func TestAddTeamMemberDuplicate(t *testing.T) {
	logger := zaptest.NewLogger(t).Sugar()
	teamMemberRepo := new(mockTeamMemberRepo)
	teamRepo := new(mockTeamRepoForMember)
	tenantRepo := new(mockTenantRepoForMember)
	userRepo := new(mockUserRepoForMember)
	publisher := new(mockPublisherForMember)

	team := &models.Team{ID: "team-1", TenantID: "tenant-1"}
	tenant := &models.Tenant{ID: "tenant-1", OwnerID: "owner-1"}
	user := &models.User{ID: "user-2", TenantID: "tenant-1"}
	existing := &models.TeamMember{ID: "member-1", TeamID: "team-1", UserID: "user-2"}

	teamRepo.On("GetByID", mock.Anything, "team-1").Return(team, nil)
	tenantRepo.On("GetByID", mock.Anything, "tenant-1").Return(tenant, nil)
	userRepo.On("GetByID", mock.Anything, "user-2").Return(user, nil)
	teamMemberRepo.On("GetByTeamAndUser", mock.Anything, "team-1", "user-2").Return(existing, nil)

	service := services.NewTeamMemberService(teamMemberRepo, teamRepo, tenantRepo, userRepo, publisher, logger)
	req := &dtos.CreateTeamMemberRequest{TeamID: "team-1", UserID: "user-2", Role: "member"}

	result, err := service.AddTeamMember(context.Background(), "owner-1", req)
	assert.Nil(t, result)
	assert.Error(t, err)
	assert.Equal(t, customErrors.ConflictError, err.(*customErrors.CustomError).Type)
}

func TestListTeamMembersSuccess(t *testing.T) {
	logger := zaptest.NewLogger(t).Sugar()
	teamMemberRepo := new(mockTeamMemberRepo)
	teamRepo := new(mockTeamRepoForMember)
	tenantRepo := new(mockTenantRepoForMember)
	userRepo := new(mockUserRepoForMember)
	publisher := new(mockPublisherForMember)

	team := &models.Team{ID: "team-1", TenantID: "tenant-1"}
	tenant := &models.Tenant{ID: "tenant-1", OwnerID: "owner-1"}
	members := []*models.TeamMember{{ID: "member-1", TeamID: "team-1", UserID: "user-2", Role: "member"}}

	teamRepo.On("GetByID", mock.Anything, "team-1").Return(team, nil)
	tenantRepo.On("GetByID", mock.Anything, "tenant-1").Return(tenant, nil)
	teamMemberRepo.On("ListByTeamID", mock.Anything, "team-1", 0, 20).Return(members, nil)

	service := services.NewTeamMemberService(teamMemberRepo, teamRepo, tenantRepo, userRepo, publisher, logger)

	result, err := service.ListTeamMembers(context.Background(), "owner-1", "team-1", 0, 20)
	assert.NoError(t, err)
	assert.Len(t, result, 1)
	assert.Equal(t, "member-1", result[0].ID)
}

func TestRemoveTeamMemberSuccess(t *testing.T) {
	logger := zaptest.NewLogger(t).Sugar()
	teamMemberRepo := new(mockTeamMemberRepo)
	teamRepo := new(mockTeamRepoForMember)
	tenantRepo := new(mockTenantRepoForMember)
	userRepo := new(mockUserRepoForMember)
	publisher := new(mockPublisherForMember)

	team := &models.Team{ID: "team-1", TenantID: "tenant-1"}
	tenant := &models.Tenant{ID: "tenant-1", OwnerID: "owner-1"}
	member := &models.TeamMember{ID: "member-1", TeamID: "team-1", UserID: "user-2"}

	teamRepo.On("GetByID", mock.Anything, "team-1").Return(team, nil)
	tenantRepo.On("GetByID", mock.Anything, "tenant-1").Return(tenant, nil)
	teamMemberRepo.On("GetByID", mock.Anything, "member-1").Return(member, nil)
	teamMemberRepo.On("Delete", mock.Anything, "member-1").Return(nil)

	service := services.NewTeamMemberService(teamMemberRepo, teamRepo, tenantRepo, userRepo, publisher, logger)
	err := service.RemoveTeamMember(context.Background(), "owner-1", "team-1", "member-1")
	assert.NoError(t, err)
}
