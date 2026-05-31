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
}

type AgentRepository interface {
	Create(ctx context.Context, agent *models.Agent) error
	GetByID(ctx context.Context, id string) (*models.Agent, error)
	GetByUserID(ctx context.Context, userID string) ([]*models.Agent, error)
	Update(ctx context.Context, agent *models.Agent) error
	Delete(ctx context.Context, id string) error
	List(ctx context.Context, userID string, skip, limit int) ([]*models.Agent, error)
}

type WorkflowRepository interface {
	Create(ctx context.Context, workflow *models.Workflow) error
	GetByID(ctx context.Context, id string) (*models.Workflow, error)
	GetByAgentID(ctx context.Context, agentID string) ([]*models.Workflow, error)
	Update(ctx context.Context, workflow *models.Workflow) error
	Delete(ctx context.Context, id string) error
	List(ctx context.Context, agentID string, skip, limit int) ([]*models.Workflow, error)
}
