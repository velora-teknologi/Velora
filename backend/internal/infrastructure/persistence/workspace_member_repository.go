package persistence

import (
	"context"

	"github.com/velora-teknologi/velora/internal/domain/models"
	"github.com/velora-teknologi/velora/internal/domain/repositories"
	"go.uber.org/zap"
	"gorm.io/gorm"
)

type postgresWorkspaceMemberRepository struct {
	db     *gorm.DB
	logger *zap.SugaredLogger
}

func NewPostgresWorkspaceMemberRepository(db *gorm.DB, logger *zap.SugaredLogger) repositories.WorkspaceMemberRepository {
	return &postgresWorkspaceMemberRepository{db: db, logger: logger}
}

func (r *postgresWorkspaceMemberRepository) Create(ctx context.Context, member *models.WorkspaceMember) error {
	return r.db.WithContext(ctx).Create(member).Error
}

func (r *postgresWorkspaceMemberRepository) GetByID(ctx context.Context, id string) (*models.WorkspaceMember, error) {
	var member models.WorkspaceMember
	if err := r.db.WithContext(ctx).Where("id = ?", id).First(&member).Error; err != nil {
		if err == gorm.ErrRecordNotFound {
			return nil, nil
		}
		return nil, err
	}
	return &member, nil
}

func (r *postgresWorkspaceMemberRepository) GetByWorkspaceAndUser(ctx context.Context, workspaceID, userID string) (*models.WorkspaceMember, error) {
	var member models.WorkspaceMember
	if err := r.db.WithContext(ctx).
		Where("workspace_id = ? AND user_id = ?", workspaceID, userID).
		First(&member).Error; err != nil {
		if err == gorm.ErrRecordNotFound {
			return nil, nil
		}
		return nil, err
	}
	return &member, nil
}

func (r *postgresWorkspaceMemberRepository) Delete(ctx context.Context, id string) error {
	return r.db.WithContext(ctx).Delete(&models.WorkspaceMember{}, "id = ?", id).Error
}

func (r *postgresWorkspaceMemberRepository) ListByWorkspaceID(ctx context.Context, workspaceID string, skip, limit int) ([]*models.WorkspaceMember, error) {
	var members []*models.WorkspaceMember
	if err := r.db.WithContext(ctx).
		Where("workspace_id = ?", workspaceID).
		Offset(skip).
		Limit(limit).
		Find(&members).Error; err != nil {
		return nil, err
	}
	return members, nil
}
