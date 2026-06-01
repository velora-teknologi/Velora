package persistence

import (
	"context"

	"go.uber.org/zap"
	"gorm.io/gorm"

	"github.com/velora-teknologi/velora/internal/domain/models"
	"github.com/velora-teknologi/velora/internal/domain/repositories"
)

type PostgresAPIKeyRepository struct {
	db     *gorm.DB
	logger *zap.SugaredLogger
}

func NewPostgresAPIKeyRepository(db *gorm.DB, logger *zap.SugaredLogger) repositories.APIKeyRepository {
	return &PostgresAPIKeyRepository{db: db, logger: logger}
}

func (r *PostgresAPIKeyRepository) Create(ctx context.Context, apiKey *models.APIKey) error {
	if err := r.db.WithContext(ctx).Create(apiKey).Error; err != nil {
		r.logger.Errorf("Failed to create api key: %v", err)
		return err
	}
	return nil
}

func (r *PostgresAPIKeyRepository) GetByID(ctx context.Context, id string) (*models.APIKey, error) {
	var apiKey models.APIKey
	if err := r.db.WithContext(ctx).First(&apiKey, "id = ?", id).Error; err != nil {
		if err == gorm.ErrRecordNotFound {
			return nil, nil
		}
		r.logger.Errorf("Failed to get api key by ID: %v", err)
		return nil, err
	}
	return &apiKey, nil
}

func (r *PostgresAPIKeyRepository) GetByKeyID(ctx context.Context, keyID string) (*models.APIKey, error) {
	var apiKey models.APIKey
	if err := r.db.WithContext(ctx).First(&apiKey, "key_id = ?", keyID).Error; err != nil {
		if err == gorm.ErrRecordNotFound {
			return nil, nil
		}
		r.logger.Errorf("Failed to get api key by key id: %v", err)
		return nil, err
	}
	return &apiKey, nil
}

func (r *PostgresAPIKeyRepository) Update(ctx context.Context, apiKey *models.APIKey) error {
	if err := r.db.WithContext(ctx).Save(apiKey).Error; err != nil {
		r.logger.Errorf("Failed to update api key: %v", err)
		return err
	}
	return nil
}

func (r *PostgresAPIKeyRepository) Delete(ctx context.Context, id string) error {
	if err := r.db.WithContext(ctx).Delete(&models.APIKey{}, "id = ?", id).Error; err != nil {
		r.logger.Errorf("Failed to delete api key: %v", err)
		return err
	}
	return nil
}

func (r *PostgresAPIKeyRepository) List(ctx context.Context, userID string, skip, limit int) ([]*models.APIKey, error) {
	var apiKeys []*models.APIKey
	if err := r.db.WithContext(ctx).Where("user_id = ?", userID).Offset(skip).Limit(limit).Find(&apiKeys).Error; err != nil {
		r.logger.Errorf("Failed to list api keys: %v", err)
		return nil, err
	}
	return apiKeys, nil
}
