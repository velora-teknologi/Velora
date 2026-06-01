package persistence

import (
	"context"

	"github.com/velora-teknologi/velora/internal/domain/models"
	"github.com/velora-teknologi/velora/internal/domain/repositories"
	"go.uber.org/zap"
	"gorm.io/gorm"
)

type postgresTeamMemberRepository struct {
	db     *gorm.DB
	logger *zap.SugaredLogger
}

func NewPostgresTeamMemberRepository(db *gorm.DB, logger *zap.SugaredLogger) repositories.TeamMemberRepository {
	return &postgresTeamMemberRepository{db: db, logger: logger}
}

func (r *postgresTeamMemberRepository) Create(ctx context.Context, member *models.TeamMember) error {
	return r.db.WithContext(ctx).Create(member).Error
}

func (r *postgresTeamMemberRepository) GetByID(ctx context.Context, id string) (*models.TeamMember, error) {
	var member models.TeamMember
	if err := r.db.WithContext(ctx).Where("id = ?", id).First(&member).Error; err != nil {
		if err == gorm.ErrRecordNotFound {
			return nil, nil
		}
		return nil, err
	}
	return &member, nil
}

func (r *postgresTeamMemberRepository) GetByTeamAndUser(ctx context.Context, teamID, userID string) (*models.TeamMember, error) {
	var member models.TeamMember
	if err := r.db.WithContext(ctx).
		Where("team_id = ? AND user_id = ?", teamID, userID).
		First(&member).Error; err != nil {
		if err == gorm.ErrRecordNotFound {
			return nil, nil
		}
		return nil, err
	}
	return &member, nil
}

func (r *postgresTeamMemberRepository) Delete(ctx context.Context, id string) error {
	return r.db.WithContext(ctx).Delete(&models.TeamMember{}, "id = ?", id).Error
}

func (r *postgresTeamMemberRepository) ListByTeamID(ctx context.Context, teamID string, skip, limit int) ([]*models.TeamMember, error) {
	var members []*models.TeamMember
	if err := r.db.WithContext(ctx).
		Where("team_id = ?", teamID).
		Offset(skip).
		Limit(limit).
		Find(&members).Error; err != nil {
		return nil, err
	}
	return members, nil
}
