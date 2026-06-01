package persistence

import (
	"context"

	"github.com/velora-teknologi/velora/internal/domain/models"
	"github.com/velora-teknologi/velora/internal/domain/repositories"
	"go.uber.org/zap"
	"gorm.io/gorm"
)

type postgresWorkspaceRepository struct {
	db     *gorm.DB
	logger *zap.SugaredLogger
}

func NewPostgresWorkspaceRepository(db *gorm.DB, logger *zap.SugaredLogger) repositories.WorkspaceRepository {
	return &postgresWorkspaceRepository{db: db, logger: logger}
}

func (r *postgresWorkspaceRepository) Create(ctx context.Context, workspace *models.Workspace) error {
	return r.db.WithContext(ctx).Create(workspace).Error
}

func (r *postgresWorkspaceRepository) GetByID(ctx context.Context, id string) (*models.Workspace, error) {
	var workspace models.Workspace
	if err := r.db.WithContext(ctx).Where("id = ?", id).First(&workspace).Error; err != nil {
		if err == gorm.ErrRecordNotFound {
			return nil, nil
		}
		return nil, err
	}
	return &workspace, nil
}

func (r *postgresWorkspaceRepository) Update(ctx context.Context, workspace *models.Workspace) error {
	return r.db.WithContext(ctx).Save(workspace).Error
}

func (r *postgresWorkspaceRepository) Delete(ctx context.Context, id string) error {
	return r.db.WithContext(ctx).Delete(&models.Workspace{}, "id = ?", id).Error
}

func (r *postgresWorkspaceRepository) List(ctx context.Context, tenantID string, skip, limit int) ([]*models.Workspace, error) {
	var workspaces []*models.Workspace
	if err := r.db.WithContext(ctx).
		Where("tenant_id = ?", tenantID).
		Offset(skip).
		Limit(limit).
		Find(&workspaces).Error; err != nil {
		return nil, err
	}
	return workspaces, nil
}
