package persistence

import (
	"context"

	"github.com/velora-teknologi/velora/internal/domain/models"
	"github.com/velora-teknologi/velora/internal/domain/repositories"
	"go.uber.org/zap"
	"gorm.io/gorm"
)

type postgresWorkflowRepository struct {
	db     *gorm.DB
	logger *zap.SugaredLogger
}

func NewPostgresWorkflowRepository(db *gorm.DB, logger *zap.SugaredLogger) repositories.WorkflowRepository {
	return &postgresWorkflowRepository{db: db, logger: logger}
}

func (r *postgresWorkflowRepository) Create(ctx context.Context, workflow *models.Workflow) error {
	return r.db.WithContext(ctx).Create(workflow).Error
}

func (r *postgresWorkflowRepository) GetByID(ctx context.Context, id string) (*models.Workflow, error) {
	var workflow models.Workflow
	if err := r.db.WithContext(ctx).Where("id = ?", id).First(&workflow).Error; err != nil {
		if err == gorm.ErrRecordNotFound {
			return nil, nil
		}
		return nil, err
	}
	return &workflow, nil
}

func (r *postgresWorkflowRepository) GetByAgentID(ctx context.Context, agentID string) ([]*models.Workflow, error) {
	var workflows []*models.Workflow
	if err := r.db.WithContext(ctx).Where("agent_id = ?", agentID).Find(&workflows).Error; err != nil {
		return nil, err
	}
	return workflows, nil
}

func (r *postgresWorkflowRepository) Update(ctx context.Context, workflow *models.Workflow) error {
	return r.db.WithContext(ctx).Save(workflow).Error
}

func (r *postgresWorkflowRepository) Delete(ctx context.Context, id string) error {
	return r.db.WithContext(ctx).Delete(&models.Workflow{}, "id = ?", id).Error
}

func (r *postgresWorkflowRepository) List(ctx context.Context, agentID string, skip, limit int) ([]*models.Workflow, error) {
	var workflows []*models.Workflow
	if err := r.db.WithContext(ctx).
		Where("agent_id = ?", agentID).
		Offset(skip).
		Limit(limit).
		Find(&workflows).Error; err != nil {
		return nil, err
	}
	return workflows, nil
}
