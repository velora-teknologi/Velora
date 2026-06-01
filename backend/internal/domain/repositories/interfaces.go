package repositories

import (
	"context"

	"github.com/velora-teknologi/velora/internal/domain/models"
)

type UserRepository interface {
	Create(ctx context.Context, user *models.User) error
	GetByID(ctx context.Context, id string) (*models.User, error)
	GetByEmail(ctx context.Context, email string) (*models.User, error)
	Update(ctx context.Context, user *models.User) error
	Delete(ctx context.Context, id string) error
	List(ctx context.Context, skip, limit int) ([]*models.User, error)
	ListByTenantID(ctx context.Context, tenantID string, skip, limit int) ([]*models.User, error)
}

type AgentRepository interface {
	Create(ctx context.Context, agent *models.Agent) error
	GetByID(ctx context.Context, id string) (*models.Agent, error)
	GetByUserID(ctx context.Context, userID string) ([]*models.Agent, error)
	Update(ctx context.Context, agent *models.Agent) error
	Delete(ctx context.Context, id string) error
	List(ctx context.Context, userID string, skip, limit int) ([]*models.Agent, error)
}

type TenantRepository interface {
	Create(ctx context.Context, tenant *models.Tenant) error
	GetByID(ctx context.Context, id string) (*models.Tenant, error)
	Update(ctx context.Context, tenant *models.Tenant) error
	Delete(ctx context.Context, id string) error
	List(ctx context.Context, ownerID string, skip, limit int) ([]*models.Tenant, error)
}

type WorkspaceRepository interface {
	Create(ctx context.Context, workspace *models.Workspace) error
	GetByID(ctx context.Context, id string) (*models.Workspace, error)
	Update(ctx context.Context, workspace *models.Workspace) error
	Delete(ctx context.Context, id string) error
	List(ctx context.Context, tenantID string, skip, limit int) ([]*models.Workspace, error)
}

type TeamRepository interface {
	Create(ctx context.Context, team *models.Team) error
	GetByID(ctx context.Context, id string) (*models.Team, error)
	Update(ctx context.Context, team *models.Team) error
	Delete(ctx context.Context, id string) error
	List(ctx context.Context, tenantID string, skip, limit int) ([]*models.Team, error)
}

type TeamMemberRepository interface {
	Create(ctx context.Context, member *models.TeamMember) error
	GetByID(ctx context.Context, id string) (*models.TeamMember, error)
	GetByTeamAndUser(ctx context.Context, teamID, userID string) (*models.TeamMember, error)
	Delete(ctx context.Context, id string) error
	ListByTeamID(ctx context.Context, teamID string, skip, limit int) ([]*models.TeamMember, error)
}

type WorkspaceMemberRepository interface {
	Create(ctx context.Context, member *models.WorkspaceMember) error
	GetByID(ctx context.Context, id string) (*models.WorkspaceMember, error)
	GetByWorkspaceAndUser(ctx context.Context, workspaceID, userID string) (*models.WorkspaceMember, error)
	Delete(ctx context.Context, id string) error
	ListByWorkspaceID(ctx context.Context, workspaceID string, skip, limit int) ([]*models.WorkspaceMember, error)
}

type InvitationRepository interface {
	Create(ctx context.Context, invitation *models.Invitation) error
	GetByToken(ctx context.Context, token string) (*models.Invitation, error)
	ListByTenantID(ctx context.Context, tenantID string, skip, limit int) ([]*models.Invitation, error)
	Update(ctx context.Context, invitation *models.Invitation) error
}

type APIKeyRepository interface {
	Create(ctx context.Context, apiKey *models.APIKey) error
	GetByID(ctx context.Context, id string) (*models.APIKey, error)
	GetByKeyID(ctx context.Context, keyID string) (*models.APIKey, error)
	Update(ctx context.Context, apiKey *models.APIKey) error
	Delete(ctx context.Context, id string) error
	List(ctx context.Context, userID string, skip, limit int) ([]*models.APIKey, error)
}

type WorkflowRepository interface {
	Create(ctx context.Context, workflow *models.Workflow) error
	GetByID(ctx context.Context, id string) (*models.Workflow, error)
	GetByAgentID(ctx context.Context, agentID string) ([]*models.Workflow, error)
	Update(ctx context.Context, workflow *models.Workflow) error
	Delete(ctx context.Context, id string) error
	List(ctx context.Context, agentID string, skip, limit int) ([]*models.Workflow, error)
}
