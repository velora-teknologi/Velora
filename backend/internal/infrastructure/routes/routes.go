package routes

import (
	"github.com/gofiber/fiber/v2"
	"github.com/nats-io/nats.go"
	"github.com/redis/go-redis/v9"
	"go.uber.org/zap"
	"gorm.io/gorm"

	"github.com/velora-teknologi/velora/internal/application/services"
	"github.com/velora-teknologi/velora/internal/infrastructure/config"
	"github.com/velora-teknologi/velora/internal/infrastructure/handlers"
	"github.com/velora-teknologi/velora/internal/infrastructure/persistence"
	"github.com/velora-teknologi/velora/pkg/middleware"
)

func SetupRoutes(
	app *fiber.App,
	db *gorm.DB,
	rdb *redis.Client,
	nc *nats.Conn,
	cfg *config.Config,
	logger *zap.SugaredLogger,
) {
	// Health check
	app.Get("/health", func(c *fiber.Ctx) error {
		return c.JSON(fiber.Map{
			"status": "healthy",
		})
	})

	v1 := app.Group("/api/v1")

	userRepo := persistence.NewPostgresUserRepository(db, logger)
	userService := services.NewUserService(userRepo, cfg.JWTSecret, logger)
	userHandler := handlers.NewUserHandler(userService, logger)

	agentRepo := persistence.NewPostgresAgentRepository(db, logger)
	agentService := services.NewAgentService(agentRepo, logger)
	agentHandler := handlers.NewAgentHandler(agentService, logger)

	workflowRepo := persistence.NewPostgresWorkflowRepository(db, logger)
	workflowService := services.NewWorkflowService(workflowRepo, agentRepo, logger)
	workflowHandler := handlers.NewWorkflowHandler(workflowService, logger)

	authGroup := v1.Group("/auth")
	authGroup.Post("/login", userHandler.LoginUser)

	userPublic := v1.Group("/users")
	userPublic.Post("", userHandler.CreateUser)

	userGroup := v1.Group("/users", middleware.JWTMiddleware(cfg.JWTSecret))
	userGroup.Get("/me", userHandler.GetCurrentUser)
	userGroup.Get(":id", userHandler.GetUser)
	userGroup.Put(":id", userHandler.UpdateUser)
	userGroup.Delete(":id", userHandler.DeleteUser)
	userGroup.Get("", userHandler.ListUsers)

	agentGroup := v1.Group("/agents", middleware.JWTMiddleware(cfg.JWTSecret))
	agentGroup.Post("", agentHandler.CreateAgent)
	agentGroup.Get("", agentHandler.ListAgents)
	agentGroup.Get(":id", agentHandler.GetAgent)
	agentGroup.Put(":id", agentHandler.UpdateAgent)
	agentGroup.Delete(":id", agentHandler.DeleteAgent)

	workflowGroup := v1.Group("/workflows", middleware.JWTMiddleware(cfg.JWTSecret))
	workflowGroup.Post("", workflowHandler.CreateWorkflow)
	workflowGroup.Get("", workflowHandler.ListWorkflows)
	workflowGroup.Get(":id", workflowHandler.GetWorkflow)
	workflowGroup.Put(":id", workflowHandler.UpdateWorkflow)
	workflowGroup.Delete(":id", workflowHandler.DeleteWorkflow)

	logger.Info("Routes setup complete")
}
