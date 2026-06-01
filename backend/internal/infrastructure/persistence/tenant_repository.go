package persistence

import (
	"context"

	"github.com/velora-teknologi/velora/internal/domain/models"
	"github.com/velora-teknologi/velora/internal/domain/repositories"
	"go.uber.org/zap"
	"gorm.io/gorm"
)

type postgresTenantRepository struct {
	db     *gorm.DB
	logger *zap.SugaredLogger
}

func NewPostgresTenantRepository(db *gorm.DB, logger *zap.SugaredLogger) repositories.TenantRepository {
	return &postgresTenantRepository{db: db, logger: logger}
}

func (r *postgresTenantRepository) Create(ctx context.Context, tenant *models.Tenant) error {
	return r.db.WithContext(ctx).Create(tenant).Error
}

func (r *postgresTenantRepository) GetByID(ctx context.Context, id string) (*models.Tenant, error) {
	var tenant models.Tenant
	if err := r.db.WithContext(ctx).Where("id = ?", id).First(&tenant).Error; err != nil {
		if err == gorm.ErrRecordNotFound {
			return nil, nil
		}
		return nil, err
	}
	return &tenant, nil
}

func (r *postgresTenantRepository) Update(ctx context.Context, tenant *models.Tenant) error {
	return r.db.WithContext(ctx).Save(tenant).Error
}

func (r *postgresTenantRepository) Delete(ctx context.Context, id string) error {
	return r.db.WithContext(ctx).Delete(&models.Tenant{}, "id = ?", id).Error
}

func (r *postgresTenantRepository) List(ctx context.Context, ownerID string, skip, limit int) ([]*models.Tenant, error) {
	var tenants []*models.Tenant
	if err := r.db.WithContext(ctx).
		Where("owner_id = ?", ownerID).
		Offset(skip).
		Limit(limit).
		Find(&tenants).Error; err != nil {
		return nil, err
	}
	return tenants, nil
}
