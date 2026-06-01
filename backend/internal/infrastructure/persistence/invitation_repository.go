package persistence

import (
	"context"

	"github.com/velora-teknologi/velora/internal/domain/models"
	"github.com/velora-teknologi/velora/internal/domain/repositories"
	"go.uber.org/zap"
	"gorm.io/gorm"
)

type postgresInvitationRepository struct {
	db     *gorm.DB
	logger *zap.SugaredLogger
}

func NewPostgresInvitationRepository(db *gorm.DB, logger *zap.SugaredLogger) repositories.InvitationRepository {
	return &postgresInvitationRepository{db: db, logger: logger}
}

func (r *postgresInvitationRepository) Create(ctx context.Context, invitation *models.Invitation) error {
	return r.db.WithContext(ctx).Create(invitation).Error
}

func (r *postgresInvitationRepository) GetByToken(ctx context.Context, token string) (*models.Invitation, error) {
	var invitation models.Invitation
	if err := r.db.WithContext(ctx).Where("token = ?", token).First(&invitation).Error; err != nil {
		if err == gorm.ErrRecordNotFound {
			return nil, nil
		}
		return nil, err
	}
	return &invitation, nil
}

func (r *postgresInvitationRepository) ListByTenantID(ctx context.Context, tenantID string, skip, limit int) ([]*models.Invitation, error) {
	var invitations []*models.Invitation
	if err := r.db.WithContext(ctx).
		Where("tenant_id = ?", tenantID).
		Offset(skip).
		Limit(limit).
		Find(&invitations).Error; err != nil {
		return nil, err
	}
	return invitations, nil
}

func (r *postgresInvitationRepository) Update(ctx context.Context, invitation *models.Invitation) error {
	return r.db.WithContext(ctx).Save(invitation).Error
}
