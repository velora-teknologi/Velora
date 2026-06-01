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

type mockWorkspaceMemberRepo struct {
	mock.Mock
}

func (m *mockWorkspaceMemberRepo) Create(ctx context.Context, member *models.WorkspaceMember) error {
	return m.Called(ctx, member).Error(0)
}
func (m *mockWorkspaceMemberRepo) GetByID(ctx context.Context, id string) (*models.WorkspaceMember, error) {
	args := m.Called(ctx, id)
	if member, ok := args.Get(0).(*models.WorkspaceMember); ok {
		return member, args.Error(1)
	}
	return nil, args.Error(1)
}
func (m *mockWorkspaceMemberRepo) GetByWorkspaceAndUser(ctx context.Context, workspaceID, userID string) (*models.WorkspaceMember, error) {
	args := m.Called(ctx, workspaceID, userID)
	if member, ok := args.Get(0).(*models.WorkspaceMember); ok {
		return member, args.Error(1)
	}
	return nil, args.Error(1)
}
func (m *mockWorkspaceMemberRepo) Delete(ctx context.Context, id string) error {
	return m.Called(ctx, id).Error(0)
}
func (m *mockWorkspaceMemberRepo) ListByWorkspaceID(ctx context.Context, workspaceID string, skip, limit int) ([]*models.WorkspaceMember, error) {
	args := m.Called(ctx, workspaceID, skip, limit)
	if members, ok := args.Get(0).([]*models.WorkspaceMember); ok {
		return members, args.Error(1)
	}
	return nil, args.Error(1)
}

type mockWorkspaceRepoForMember struct {
	mock.Mock
}

func (m *mockWorkspaceRepoForMember) Create(ctx context.Context, workspace *models.Workspace) error {
	return nil
}
func (m *mockWorkspaceRepoForMember) GetByID(ctx context.Context, id string) (*models.Workspace, error) {
	args := m.Called(ctx, id)
	if workspace, ok := args.Get(0).(*models.Workspace); ok {
		return workspace, args.Error(1)
	}
	return nil, args.Error(1)
}
func (m *mockWorkspaceRepoForMember) Update(ctx context.Context, workspace *models.Workspace) error {
	return nil
}
func (m *mockWorkspaceRepoForMember) Delete(ctx context.Context, id string) error { return nil }
func (m *mockWorkspaceRepoForMember) List(ctx context.Context, tenantID string, skip, limit int) ([]*models.Workspace, error) {
	return nil, nil
}

type mockTenantRepoForWorkspaceMember struct {
	mock.Mock
}

func (m *mockTenantRepoForWorkspaceMember) Create(ctx context.Context, tenant *models.Tenant) error {
	return nil
}
func (m *mockTenantRepoForWorkspaceMember) GetByID(ctx context.Context, id string) (*models.Tenant, error) {
	args := m.Called(ctx, id)
	if tenant, ok := args.Get(0).(*models.Tenant); ok {
		return tenant, args.Error(1)
	}
	return nil, args.Error(1)
}
func (m *mockTenantRepoForWorkspaceMember) Update(ctx context.Context, tenant *models.Tenant) error {
	return nil
}
func (m *mockTenantRepoForWorkspaceMember) Delete(ctx context.Context, id string) error { return nil }
func (m *mockTenantRepoForWorkspaceMember) List(ctx context.Context, ownerID string, skip, limit int) ([]*models.Tenant, error) {
	return nil, nil
}

type mockUserRepoForWorkspaceMember struct {
	mock.Mock
}

func (m *mockUserRepoForWorkspaceMember) Create(ctx context.Context, user *models.User) error {
	return nil
}
func (m *mockUserRepoForWorkspaceMember) GetByID(ctx context.Context, id string) (*models.User, error) {
	args := m.Called(ctx, id)
	if user, ok := args.Get(0).(*models.User); ok {
		return user, args.Error(1)
	}
	return nil, args.Error(1)
}
func (m *mockUserRepoForWorkspaceMember) GetByEmail(ctx context.Context, email string) (*models.User, error) {
	return nil, nil
}
func (m *mockUserRepoForWorkspaceMember) Update(ctx context.Context, user *models.User) error {
	return nil
}
func (m *mockUserRepoForWorkspaceMember) Delete(ctx context.Context, id string) error { return nil }
func (m *mockUserRepoForWorkspaceMember) List(ctx context.Context, skip, limit int) ([]*models.User, error) {
	return nil, nil
}
func (m *mockUserRepoForWorkspaceMember) ListByTenantID(ctx context.Context, tenantID string, skip, limit int) ([]*models.User, error) {
	return nil, nil
}

type mockPublisherForWorkspaceMember struct {
	mock.Mock
}

func (m *mockPublisherForWorkspaceMember) Publish(subject string, data []byte) error {
	return m.Called(subject, data).Error(0)
}

func TestAddWorkspaceMemberSuccess(t *testing.T) {
	logger := zaptest.NewLogger(t).Sugar()
	workspaceMemberRepo := new(mockWorkspaceMemberRepo)
	workspaceRepo := new(mockWorkspaceRepoForMember)
	tenantRepo := new(mockTenantRepoForWorkspaceMember)
	userRepo := new(mockUserRepoForWorkspaceMember)
	publisher := new(mockPublisherForWorkspaceMember)

	workspace := &models.Workspace{ID: "workspace-1", TenantID: "tenant-1"}
	tenant := &models.Tenant{ID: "tenant-1", OwnerID: "owner-1"}
	user := &models.User{ID: "user-2", TenantID: "tenant-1"}

	workspaceRepo.On("GetByID", mock.Anything, "workspace-1").Return(workspace, nil)
	tenantRepo.On("GetByID", mock.Anything, "tenant-1").Return(tenant, nil)
	userRepo.On("GetByID", mock.Anything, "user-2").Return(user, nil)
	workspaceMemberRepo.On("GetByWorkspaceAndUser", mock.Anything, "workspace-1", "user-2").Return(nil, nil)
	workspaceMemberRepo.On("Create", mock.Anything, mock.AnythingOfType("*models.WorkspaceMember")).Return(nil)
	publisher.On("Publish", "workspace.member.added", mock.Anything).Return(nil)

	service := services.NewWorkspaceMemberService(workspaceMemberRepo, workspaceRepo, tenantRepo, userRepo, publisher, logger)
	req := &dtos.CreateWorkspaceMemberRequest{WorkspaceID: "workspace-1", UserID: "user-2", Role: "member"}

	result, err := service.AddWorkspaceMember(context.Background(), "owner-1", req)
	assert.NoError(t, err)
	assert.Equal(t, "workspace-1", result.WorkspaceID)
	assert.Equal(t, "user-2", result.UserID)
	assert.Equal(t, "member", result.Role)
}

func TestAddWorkspaceMemberDuplicate(t *testing.T) {
	logger := zaptest.NewLogger(t).Sugar()
	workspaceMemberRepo := new(mockWorkspaceMemberRepo)
	workspaceRepo := new(mockWorkspaceRepoForMember)
	tenantRepo := new(mockTenantRepoForWorkspaceMember)
	userRepo := new(mockUserRepoForWorkspaceMember)
	publisher := new(mockPublisherForWorkspaceMember)

	workspace := &models.Workspace{ID: "workspace-1", TenantID: "tenant-1"}
	tenant := &models.Tenant{ID: "tenant-1", OwnerID: "owner-1"}
	user := &models.User{ID: "user-2", TenantID: "tenant-1"}
	existing := &models.WorkspaceMember{ID: "member-1", WorkspaceID: "workspace-1", UserID: "user-2"}

	workspaceRepo.On("GetByID", mock.Anything, "workspace-1").Return(workspace, nil)
	tenantRepo.On("GetByID", mock.Anything, "tenant-1").Return(tenant, nil)
	userRepo.On("GetByID", mock.Anything, "user-2").Return(user, nil)
	workspaceMemberRepo.On("GetByWorkspaceAndUser", mock.Anything, "workspace-1", "user-2").Return(existing, nil)

	service := services.NewWorkspaceMemberService(workspaceMemberRepo, workspaceRepo, tenantRepo, userRepo, publisher, logger)
	req := &dtos.CreateWorkspaceMemberRequest{WorkspaceID: "workspace-1", UserID: "user-2", Role: "member"}

	result, err := service.AddWorkspaceMember(context.Background(), "owner-1", req)
	assert.Nil(t, result)
	assert.Error(t, err)
	assert.Equal(t, customErrors.ConflictError, err.(*customErrors.CustomError).Type)
}

func TestListWorkspaceMembersSuccess(t *testing.T) {
	logger := zaptest.NewLogger(t).Sugar()
	workspaceMemberRepo := new(mockWorkspaceMemberRepo)
	workspaceRepo := new(mockWorkspaceRepoForMember)
	tenantRepo := new(mockTenantRepoForWorkspaceMember)
	userRepo := new(mockUserRepoForWorkspaceMember)
	publisher := new(mockPublisherForWorkspaceMember)

	workspace := &models.Workspace{ID: "workspace-1", TenantID: "tenant-1"}
	tenant := &models.Tenant{ID: "tenant-1", OwnerID: "owner-1"}
	members := []*models.WorkspaceMember{{ID: "member-1", WorkspaceID: "workspace-1", UserID: "user-2", Role: "member"}}

	workspaceRepo.On("GetByID", mock.Anything, "workspace-1").Return(workspace, nil)
	tenantRepo.On("GetByID", mock.Anything, "tenant-1").Return(tenant, nil)
	workspaceMemberRepo.On("ListByWorkspaceID", mock.Anything, "workspace-1", 0, 20).Return(members, nil)

	service := services.NewWorkspaceMemberService(workspaceMemberRepo, workspaceRepo, tenantRepo, userRepo, publisher, logger)

	result, err := service.ListWorkspaceMembers(context.Background(), "owner-1", "workspace-1", 0, 20)
	assert.NoError(t, err)
	assert.Len(t, result, 1)
	assert.Equal(t, "member-1", result[0].ID)
}

func TestRemoveWorkspaceMemberSuccess(t *testing.T) {
	logger := zaptest.NewLogger(t).Sugar()
	workspaceMemberRepo := new(mockWorkspaceMemberRepo)
	workspaceRepo := new(mockWorkspaceRepoForMember)
	tenantRepo := new(mockTenantRepoForWorkspaceMember)
	userRepo := new(mockUserRepoForWorkspaceMember)
	publisher := new(mockPublisherForWorkspaceMember)

	workspace := &models.Workspace{ID: "workspace-1", TenantID: "tenant-1"}
	tenant := &models.Tenant{ID: "tenant-1", OwnerID: "owner-1"}
	member := &models.WorkspaceMember{ID: "member-1", WorkspaceID: "workspace-1", UserID: "user-2"}

	workspaceRepo.On("GetByID", mock.Anything, "workspace-1").Return(workspace, nil)
	tenantRepo.On("GetByID", mock.Anything, "tenant-1").Return(tenant, nil)
	workspaceMemberRepo.On("GetByID", mock.Anything, "member-1").Return(member, nil)
	workspaceMemberRepo.On("Delete", mock.Anything, "member-1").Return(nil)

	service := services.NewWorkspaceMemberService(workspaceMemberRepo, workspaceRepo, tenantRepo, userRepo, publisher, logger)
	err := service.RemoveWorkspaceMember(context.Background(), "owner-1", "workspace-1", "member-1")
	assert.NoError(t, err)
}
