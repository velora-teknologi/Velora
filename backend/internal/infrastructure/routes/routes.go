package routes

import (
	"github.com/gofiber/fiber/v2"
	"github.com/nats-io/nats.go"
	"github.com/redis/go-redis/v9"
	"go.uber.org/zap"
	"gorm.io/gorm"
)

func SetupRoutes(
	app *fiber.App,
	db *gorm.DB,
	rdb *redis.Client,
	nc *nats.Conn,
	logger *zap.SugaredLogger,
) {
	// Health check
	app.Get("/health", func(c *fiber.Ctx) error {
		return c.JSON(fiber.Map{
			"status": "healthy",
		})
	})

	// API v1 routes
	_ = app.Group("/api/v1")

	// Example: Auth routes would go here
	// setupAuthRoutes(v1, db, logger)

	// Example: User routes would go here
	// setupUserRoutes(v1, db, logger)

	// Example: Agent routes would go here
	// setupAgentRoutes(v1, db, rdb, nc, logger)

	logger.Info("Routes setup complete")
}
