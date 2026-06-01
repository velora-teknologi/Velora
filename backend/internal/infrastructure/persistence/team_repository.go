package persistence

import (
	"context"

	"github.com/velora-teknologi/velora/internal/domain/models"
	"github.com/velora-teknologi/velora/internal/domain/repositories"
	"go.uber.org/zap"
	"gorm.io/gorm"
)

type postgresTeamRepository struct {
	db     *gorm.DB
	logger *zap.SugaredLogger
}

func NewPostgresTeamRepository(db *gorm.DB, logger *zap.SugaredLogger) repositories.TeamRepository {
	return &postgresTeamRepository{db: db, logger: logger}
}

func (r *postgresTeamRepository) Create(ctx context.Context, team *models.Team) error {
	return r.db.WithContext(ctx).Create(team).Error
}

func (r *postgresTeamRepository) GetByID(ctx context.Context, id string) (*models.Team, error) {
	var team models.Team
	if err := r.db.WithContext(ctx).Where("id = ?", id).First(&team).Error; err != nil {
		if err == gorm.ErrRecordNotFound {
			return nil, nil
		}
		return nil, err
	}
	return &team, nil
}

func (r *postgresTeamRepository) Update(ctx context.Context, team *models.Team) error {
	return r.db.WithContext(ctx).Save(team).Error
}

func (r *postgresTeamRepository) Delete(ctx context.Context, id string) error {
	return r.db.WithContext(ctx).Delete(&models.Team{}, "id = ?", id).Error
}

func (r *postgresTeamRepository) List(ctx context.Context, tenantID string, skip, limit int) ([]*models.Team, error) {
	var teams []*models.Team
	if err := r.db.WithContext(ctx).
		Where("tenant_id = ?", tenantID).
		Offset(skip).
		Limit(limit).
		Find(&teams).Error; err != nil {
		return nil, err
	}
	return teams, nil
}
