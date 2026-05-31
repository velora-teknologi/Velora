package persistence

import (
	"context"

	"github.com/velora-teknologi/velora/internal/domain/models"
	"github.com/velora-teknologi/velora/internal/domain/repositories"
	"go.uber.org/zap"
	"gorm.io/gorm"
)

type postgresAgentRepository struct {
	db     *gorm.DB
	logger *zap.SugaredLogger
}

func NewPostgresAgentRepository(db *gorm.DB, logger *zap.SugaredLogger) repositories.AgentRepository {
	return &postgresAgentRepository{db: db, logger: logger}
}

func (r *postgresAgentRepository) Create(ctx context.Context, agent *models.Agent) error {
	return r.db.WithContext(ctx).Create(agent).Error
}

func (r *postgresAgentRepository) GetByID(ctx context.Context, id string) (*models.Agent, error) {
	var agent models.Agent
	if err := r.db.WithContext(ctx).Where("id = ?", id).First(&agent).Error; err != nil {
		if err == gorm.ErrRecordNotFound {
			return nil, nil
		}
		return nil, err
	}
	return &agent, nil
}

func (r *postgresAgentRepository) GetByUserID(ctx context.Context, userID string) ([]*models.Agent, error) {
	var agents []*models.Agent
	if err := r.db.WithContext(ctx).Where("user_id = ?", userID).Find(&agents).Error; err != nil {
		return nil, err
	}
	return agents, nil
}

func (r *postgresAgentRepository) Update(ctx context.Context, agent *models.Agent) error {
	return r.db.WithContext(ctx).Save(agent).Error
}

func (r *postgresAgentRepository) Delete(ctx context.Context, id string) error {
	return r.db.WithContext(ctx).Delete(&models.Agent{}, "id = ?", id).Error
}

func (r *postgresAgentRepository) List(ctx context.Context, userID string, skip, limit int) ([]*models.Agent, error) {
	var agents []*models.Agent
	if err := r.db.WithContext(ctx).
		Where("user_id = ?", userID).
		Offset(skip).
		Limit(limit).
		Find(&agents).Error; err != nil {
		return nil, err
	}
	return agents, nil
}
