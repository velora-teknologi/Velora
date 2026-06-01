package database

import (
	"fmt"

	"go.uber.org/zap"
	"gorm.io/gorm"

	"github.com/velora-teknologi/velora/internal/domain/models"
)

// Migrate runs all database migrations
func Migrate(db *gorm.DB, logger *zap.SugaredLogger) error {
	logger.Info("Running database migrations...")

	// Create extensions
	if err := db.Exec("CREATE EXTENSION IF NOT EXISTS \"uuid-ossp\"").Error; err != nil {
		logger.Warnf("Failed to create uuid-ossp extension: %v", err)
	}

	// Migrate all models
	if err := db.AutoMigrate(
		&models.User{},
		&models.Agent{},
		&models.Workflow{},
		&models.Tenant{},
		&models.Workspace{},
		&models.Team{},
		&models.Invitation{},
		&models.APIKey{},
	); err != nil {
		logger.Errorf("Migration failed: %v", err)
		return fmt.Errorf("migration error: %w", err)
	}

	logger.Info("Database migrations completed successfully")
	return nil
}

// DropAllTables drops all tables (use with caution in development only)
func DropAllTables(db *gorm.DB, logger *zap.SugaredLogger) error {
	logger.Warn("Dropping all tables...")

	return db.Migrator().DropTable(
		&models.User{},
		&models.Agent{},
		&models.Workflow{},
		&models.Tenant{},
		&models.Workspace{},
		&models.Team{},
		&models.Invitation{},
		&models.APIKey{},
	)
}

// CreateIndexes creates custom indexes for better query performance
func CreateIndexes(db *gorm.DB, logger *zap.SugaredLogger) error {
	logger.Info("Creating database indexes...")

	// User indexes
	if err := db.Migrator().CreateIndex(&models.User{}, "email"); err != nil {
		logger.Warnf("Failed to create email index: %v", err)
	}

	if err := db.Migrator().CreateIndex(&models.User{}, "is_active"); err != nil {
		logger.Warnf("Failed to create is_active index: %v", err)
	}

	// Agent indexes
	if err := db.Migrator().CreateIndex(&models.Agent{}, "user_id"); err != nil {
		logger.Warnf("Failed to create agent user_id index: %v", err)
	}

	if err := db.Migrator().CreateIndex(&models.Agent{}, "status"); err != nil {
		logger.Warnf("Failed to create agent status index: %v", err)
	}

	// Workflow indexes
	if err := db.Migrator().CreateIndex(&models.Workflow{}, "agent_id"); err != nil {
		logger.Warnf("Failed to create workflow agent_id index: %v", err)
	}

	if err := db.Migrator().CreateIndex(&models.Workflow{}, "is_active"); err != nil {
		logger.Warnf("Failed to create workflow is_active index: %v", err)
	}

	// Tenant indexes
	if err := db.Migrator().CreateIndex(&models.Tenant{}, "owner_id"); err != nil {
		logger.Warnf("Failed to create tenant owner_id index: %v", err)
	}

	if err := db.Migrator().CreateIndex(&models.Tenant{}, "is_active"); err != nil {
		logger.Warnf("Failed to create tenant is_active index: %v", err)
	}

	if err := db.Migrator().CreateIndex(&models.Team{}, "tenant_id"); err != nil {
		logger.Warnf("Failed to create team tenant_id index: %v", err)
	}

	if err := db.Migrator().CreateIndex(&models.Team{}, "is_active"); err != nil {
		logger.Warnf("Failed to create team is_active index: %v", err)
	}

	logger.Info("Database indexes created successfully")
	return nil
}
