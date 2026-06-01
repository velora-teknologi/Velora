package services_test

import (
	"context"
	"testing"

	"github.com/google/uuid"
	"github.com/velora-teknologi/velora/internal/application/dtos"
	"github.com/velora-teknologi/velora/internal/application/services"
	"github.com/velora-teknologi/velora/internal/domain/models"
	"go.uber.org/zap"
)

type mockWorkspaceRepository struct {
	workspaces map[string]*models.Workspace
}

func newMockWorkspaceRepository() *mockWorkspaceRepository {
	return &mockWorkspaceRepository{workspaces: make(map[string]*models.Workspace)}
}

func (m *mockWorkspaceRepository) Create(ctx context.Context, workspace *models.Workspace) error {
	if workspace.ID == "" {
		workspace.ID = uuid.NewString()
	}
	m.workspaces[workspace.ID] = workspace
	return nil
}

func (m *mockWorkspaceRepository) GetByID(ctx context.Context, id string) (*models.Workspace, error) {
	return m.workspaces[id], nil
}

func (m *mockWorkspaceRepository) Update(ctx context.Context, workspace *models.Workspace) error {
	m.workspaces[workspace.ID] = workspace
	return nil
}

func (m *mockWorkspaceRepository) Delete(ctx context.Context, id string) error {
	delete(m.workspaces, id)
	return nil
}

func (m *mockWorkspaceRepository) List(ctx context.Context, tenantID string, skip, limit int) ([]*models.Workspace, error) {
	var workspaces []*models.Workspace
	for _, workspace := range m.workspaces {
		if workspace.TenantID == tenantID {
			workspaces = append(workspaces, workspace)
		}
	}
	return workspaces, nil
}

type mockTenantRepositoryWithOwner struct {
	tenants map[string]*models.Tenant
}

func newMockTenantRepositoryWithOwner() *mockTenantRepositoryWithOwner {
	return &mockTenantRepositoryWithOwner{tenants: make(map[string]*models.Tenant)}
}

func (m *mockTenantRepositoryWithOwner) Create(ctx context.Context, tenant *models.Tenant) error {
	if tenant.ID == "" {
		tenant.ID = uuid.NewString()
	}
	m.tenants[tenant.ID] = tenant
	return nil
}

func (m *mockTenantRepositoryWithOwner) GetByID(ctx context.Context, id string) (*models.Tenant, error) {
	return m.tenants[id], nil
}

func (m *mockTenantRepositoryWithOwner) Update(ctx context.Context, tenant *models.Tenant) error {
	m.tenants[tenant.ID] = tenant
	return nil
}

func (m *mockTenantRepositoryWithOwner) Delete(ctx context.Context, id string) error {
	delete(m.tenants, id)
	return nil
}

func (m *mockTenantRepositoryWithOwner) List(ctx context.Context, ownerID string, skip, limit int) ([]*models.Tenant, error) {
	var tenants []*models.Tenant
	for _, tenant := range m.tenants {
		if tenant.OwnerID == ownerID {
			tenants = append(tenants, tenant)
		}
	}
	return tenants, nil
}

type mockInvitationRepository struct {
	invitations map[string]*models.Invitation
}

func newMockInvitationRepository() *mockInvitationRepository {
	return &mockInvitationRepository{invitations: make(map[string]*models.Invitation)}
}

func (m *mockInvitationRepository) Create(ctx context.Context, invitation *models.Invitation) error {
	if invitation.ID == "" {
		invitation.ID = uuid.NewString()
	}
	m.invitations[invitation.Token] = invitation
	return nil
}

func (m *mockInvitationRepository) GetByToken(ctx context.Context, token string) (*models.Invitation, error) {
	return m.invitations[token], nil
}

func (m *mockInvitationRepository) ListByTenantID(ctx context.Context, tenantID string, skip, limit int) ([]*models.Invitation, error) {
	var invites []*models.Invitation
	for _, invitation := range m.invitations {
		if invitation.TenantID == tenantID {
			invites = append(invites, invitation)
		}
	}
	return invites, nil
}

func (m *mockInvitationRepository) Update(ctx context.Context, invitation *models.Invitation) error {
	m.invitations[invitation.Token] = invitation
	return nil
}

type mockUserRepository struct {
	users map[string]*models.User
}

func newMockUserRepository() *mockUserRepository {
	return &mockUserRepository{users: make(map[string]*models.User)}
}

func (m *mockUserRepository) Create(ctx context.Context, user *models.User) error {
	if user.ID == "" {
		user.ID = uuid.NewString()
	}
	m.users[user.Email] = user
	return nil
}

func (m *mockUserRepository) GetByID(ctx context.Context, id string) (*models.User, error) {
	for _, user := range m.users {
		if user.ID == id {
			return user, nil
		}
	}
	return nil, nil
}

func (m *mockUserRepository) GetByEmail(ctx context.Context, email string) (*models.User, error) {
	return m.users[email], nil
}

func (m *mockUserRepository) Update(ctx context.Context, user *models.User) error {
	m.users[user.Email] = user
	return nil
}

func (m *mockUserRepository) Delete(ctx context.Context, id string) error {
	for email, user := range m.users {
		if user.ID == id {
			delete(m.users, email)
			return nil
		}
	}
	return nil
}

func (m *mockUserRepository) List(ctx context.Context, skip, limit int) ([]*models.User, error) {
	var users []*models.User
	for _, user := range m.users {
		users = append(users, user)
	}
	return users, nil
}
func (m *mockUserRepository) ListByTenantID(ctx context.Context, tenantID string, skip, limit int) ([]*models.User, error) {
	var users []*models.User
	for _, user := range m.users {
		if user.TenantID == tenantID {
			users = append(users, user)
		}
	}
	return users, nil
}

type mockEventPublisher struct {
	subject string
	payload []byte
}

func (m *mockEventPublisher) Publish(subject string, data []byte) error {
	m.subject = subject
	m.payload = append([]byte(nil), data...)
	return nil
}

func TestCreateWorkspace(t *testing.T) {
	tenantRepo := newMockTenantRepositoryWithOwner()
	workspaceRepo := newMockWorkspaceRepository()
	publisher := &mockEventPublisher{}
	logger := zap.NewNop().Sugar()
	service := services.NewWorkspaceService(workspaceRepo, tenantRepo, publisher, logger)

	tenantID := uuid.NewString()
	tenantRepo.Create(context.Background(), &models.Tenant{ID: tenantID, Name: "Acme", Slug: "acme", OwnerID: "owner-1", IsActive: true})

	req := &dtos.CreateWorkspaceRequest{TenantID: tenantID, Name: "East Coast", Description: "Primary workspace"}
	workspace, err := service.CreateWorkspace(context.Background(), "owner-1", req)
	if err != nil {
		t.Fatalf("expected no error, got %v", err)
	}
	if workspace.Name != req.Name || workspace.TenantID != tenantID {
		t.Fatalf("unexpected workspace response: %+v", workspace)
	}
	if publisher.subject != "workspace.created" {
		t.Fatalf("expected workspace.created event, got %s", publisher.subject)
	}
}

func TestCreateInvitationAndAccept(t *testing.T) {
	tenantRepo := newMockTenantRepositoryWithOwner()
	invitationRepo := newMockInvitationRepository()
	userRepo := newMockUserRepository()
	publisher := &mockEventPublisher{}
	logger := zap.NewNop().Sugar()
	service := services.NewInvitationService(invitationRepo, tenantRepo, userRepo, publisher, logger)

	tenantID := uuid.NewString()
	tenantRepo.Create(context.Background(), &models.Tenant{ID: tenantID, Name: "Acme", Slug: "acme", OwnerID: "owner-1", IsActive: true})

	req := &dtos.CreateInvitationRequest{TenantID: tenantID, Email: "invitee@example.com", Role: "member", ExpiresIn: 3600}
	invitation, err := service.CreateInvitation(context.Background(), "owner-1", req)
	if err != nil {
		t.Fatalf("expected no error, got %v", err)
	}
	if invitation.Email != "invitee@example.com" || invitation.Status != "pending" {
		t.Fatalf("unexpected invitation response: %+v", invitation)
	}

	acceptReq := &dtos.AcceptInvitationRequest{Email: "invitee@example.com", Password: "securepass123", FullName: "Invitee User"}
	resp, err := service.AcceptInvitation(context.Background(), invitation.Token, acceptReq)
	if err != nil {
		t.Fatalf("expected no error accepting invitation, got %v", err)
	}
	if resp.Email != "invitee@example.com" || resp.Role != "member" {
		t.Fatalf("unexpected accepted invitation response: %+v", resp)
	}
	if publisher.subject != "tenant.invitation.accepted" {
		t.Fatalf("expected tenant.invitation.accepted event, got %s", publisher.subject)
	}
}

func TestListInvitations(t *testing.T) {
	tenantRepo := newMockTenantRepositoryWithOwner()
	invitationRepo := newMockInvitationRepository()
	userRepo := newMockUserRepository()
	publisher := &mockEventPublisher{}
	logger := zap.NewNop().Sugar()
	service := services.NewInvitationService(invitationRepo, tenantRepo, userRepo, publisher, logger)

	tenantID := uuid.NewString()
	tenantRepo.Create(context.Background(), &models.Tenant{ID: tenantID, Name: "Acme", Slug: "acme", OwnerID: "owner-1", IsActive: true})
	invitationRepo.Create(context.Background(), &models.Invitation{Token: "test-token", TenantID: tenantID, Email: "guest@example.com", Role: "member", Status: "pending", InvitedBy: "owner-1"})

	list, err := service.ListInvitations(context.Background(), "owner-1", tenantID, 0, 10)
	if err != nil {
		t.Fatalf("expected no error, got %v", err)
	}
	if len(list) != 1 {
		t.Fatalf("expected 1 invitation, got %d", len(list))
	}
}
