package persistence
package persistence

import (
	"context"

	"go.uber.org/zap"
	"gorm.io/gorm"

	"github.com/velora-teknologi/velora/internal/domain/models"
	"github.com/velora-teknologi/velora/internal/domain/repositories"
)

type PostgresUserRepository struct {
	db     *gorm.DB
	logger *zap.SugaredLogger
}

func NewPostgresUserRepository(db *gorm.DB, logger *zap.SugaredLogger) repositories.UserRepository {
	return &PostgresUserRepository{
		db:     db,
		logger: logger,
	}
}

func (r *PostgresUserRepository) Create(ctx context.Context, user *models.User) error {
	if err := r.db.WithContext(ctx).Create(user).Error; err != nil {
		r.logger.Errorf("Failed to create user: %v", err)
		return err
	}
	return nil
}

func (r *PostgresUserRepository) GetByID(ctx context.Context, id string) (*models.User, error) {
	var user models.User
	if err := r.db.WithContext(ctx).First(&user, "id = ?", id).Error; err != nil {
		if err == gorm.ErrRecordNotFound {
			return nil, nil
		}
		r.logger.Errorf("Failed to get user by ID: %v", err)
		return nil, err
	}
	return &user, nil
}

func (r *PostgresUserRepository) GetByEmail(ctx context.Context, email string) (*models.User, error) {
	var user models.User
	if err := r.db.WithContext(ctx).First(&user, "email = ?", email).Error; err != nil {
		if err == gorm.ErrRecordNotFound {
			return nil, nil
		}
		r.logger.Errorf("Failed to get user by email: %v", err)
		return nil, err
	}
	return &user, nil
}

func (r *PostgresUserRepository) Update(ctx context.Context, user *models.User) error {
	if err := r.db.WithContext(ctx).Save(user).Error; err != nil {
		r.logger.Errorf("Failed to update user: %v", err)
		return err
	}
	return nil
}

func (r *PostgresUserRepository) Delete(ctx context.Context, id string) error {
	if err := r.db.WithContext(ctx).Delete(&models.User{}, "id = ?", id).Error; err != nil {
		r.logger.Errorf("Failed to delete user: %v", err)
		return err
	}
	return nil
}

func (r *PostgresUserRepository) List(ctx context.Context, skip, limit int) ([]*models.User, error) {
	var users []*models.User
	if err := r.db.WithContext(ctx).Offset(skip).Limit(limit).Find(&users).Error; err != nil {
		r.logger.Errorf("Failed to list users: %v", err)
		return nil, err
	}
	return users, nil
}
