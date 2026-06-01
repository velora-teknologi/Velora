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
	"github.com/velora-teknologi/velora/internal/infrastructure/messaging"
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
	// create a NATS publisher adapter and pass it to services that need event publishing
	var publisher services.Publisher
	if nc != nil {
		publisher = messaging.NewNatsPublisher(nc, logger)
	}
	agentService := services.NewAgentService(agentRepo, publisher, logger)
	workflowRepo := persistence.NewPostgresWorkflowRepository(db, logger)
	workflowService := services.NewWorkflowService(workflowRepo, agentRepo, publisher, logger)
	workflowHandler := handlers.NewWorkflowHandler(workflowService, logger)

	tenantRepo := persistence.NewPostgresTenantRepository(db, logger)
	tenantService := services.NewTenantService(tenantRepo, publisher, logger)
	tenantHandler := handlers.NewTenantHandler(tenantService, logger)

	workspaceRepo := persistence.NewPostgresWorkspaceRepository(db, logger)
	workspaceService := services.NewWorkspaceService(workspaceRepo, tenantRepo, publisher, logger)
	workspaceHandler := handlers.NewWorkspaceHandler(workspaceService, logger)

	workspaceMemberRepo := persistence.NewPostgresWorkspaceMemberRepository(db, logger)
	workspaceMemberService := services.NewWorkspaceMemberService(workspaceMemberRepo, workspaceRepo, tenantRepo, userRepo, publisher, logger)
	workspaceMemberHandler := handlers.NewWorkspaceMemberHandler(workspaceMemberService, logger)

	teamRepo := persistence.NewPostgresTeamRepository(db, logger)
	teamService := services.NewTeamService(teamRepo, tenantRepo, publisher, logger)
	teamHandler := handlers.NewTeamHandler(teamService, logger)

	tenantMemberService := services.NewTenantMemberService(userRepo, tenantRepo, logger)
	tenantMemberHandler := handlers.NewTenantMemberHandler(tenantMemberService, logger)

	teamMemberRepo := persistence.NewPostgresTeamMemberRepository(db, logger)
	teamMemberService := services.NewTeamMemberService(teamMemberRepo, teamRepo, tenantRepo, userRepo, publisher, logger)
	teamMemberHandler := handlers.NewTeamMemberHandler(teamMemberService, logger)

	invitationRepo := persistence.NewPostgresInvitationRepository(db, logger)
	invitationService := services.NewInvitationService(invitationRepo, tenantRepo, userRepo, publisher, logger)
	invitationHandler := handlers.NewInvitationHandler(invitationService, logger)

	apiKeyRepo := persistence.NewPostgresAPIKeyRepository(db, logger)
	apiKeyService := services.NewAPIKeyService(apiKeyRepo, userRepo, logger)
	apiKeyHandler := handlers.NewAPIKeyHandler(apiKeyService, logger)

	agentHandler := handlers.NewAgentHandler(agentService, workflowService, logger)

	authGroup := v1.Group("/auth")
	authGroup.Post("/login", userHandler.LoginUser)
	authGroup.Post("/register", userHandler.RegisterUser)

	userPublic := v1.Group("/users")
	userPublic.Post("", userHandler.CreateUser)

	authMiddleware := middleware.AuthMiddleware(cfg.JWTSecret, apiKeyRepo, userRepo)

	userGroup := v1.Group("/users", authMiddleware)
	userGroup.Get("/me", userHandler.GetCurrentUser)
	userGroup.Get(":id", userHandler.GetUser)
	userGroup.Put(":id", userHandler.UpdateUser)
	userGroup.Delete(":id", userHandler.DeleteUser)
	userGroup.Get("", middleware.RBACMiddleware("admin"), userHandler.ListUsers)

	agentGroup := v1.Group("/agents", authMiddleware)
	agentGroup.Post("", agentHandler.CreateAgent)
	agentGroup.Get("", agentHandler.ListAgents)
	agentGroup.Get(":id", agentHandler.GetAgent)
	agentGroup.Post(":id/execute", agentHandler.ExecuteAgent)
	agentGroup.Get(":id/workflows", agentHandler.ListAgentWorkflows)
	agentGroup.Put(":id", agentHandler.UpdateAgent)
	agentGroup.Delete(":id", agentHandler.DeleteAgent)

	tenantGroup := v1.Group("/tenants", authMiddleware)
	tenantGroup.Post("", tenantHandler.CreateTenant)
	tenantGroup.Get("", tenantHandler.ListTenants)
	tenantGroup.Get(":id", tenantHandler.GetTenant)
	tenantGroup.Put(":id", tenantHandler.UpdateTenant)
	tenantGroup.Delete(":id", tenantHandler.DeleteTenant)
	tenantGroup.Get(":id/members", tenantMemberHandler.ListTenantMembers)

	workspaceGroup := v1.Group("/workspaces", authMiddleware)
	workspaceGroup.Post("", workspaceHandler.CreateWorkspace)
	workspaceGroup.Get("", workspaceHandler.ListWorkspaces)
	workspaceGroup.Get(":id", workspaceHandler.GetWorkspace)
	workspaceGroup.Put(":id", workspaceHandler.UpdateWorkspace)
	workspaceGroup.Delete(":id", workspaceHandler.DeleteWorkspace)

	workspaceMembersGroup := workspaceGroup.Group(":id/members")
	workspaceMembersGroup.Post("", workspaceMemberHandler.AddWorkspaceMember)
	workspaceMembersGroup.Get("", workspaceMemberHandler.ListWorkspaceMembers)
	workspaceMembersGroup.Delete(":member_id", workspaceMemberHandler.RemoveWorkspaceMember)

	teamGroup := v1.Group("/teams", authMiddleware)
	teamGroup.Post("", teamHandler.CreateTeam)
	teamGroup.Get("", teamHandler.ListTeams)
	teamGroup.Get(":id", teamHandler.GetTeam)
	teamGroup.Put(":id", teamHandler.UpdateTeam)
	teamGroup.Delete(":id", teamHandler.DeleteTeam)

	teamMembersGroup := teamGroup.Group(":id/members")
	teamMembersGroup.Post("", teamMemberHandler.AddTeamMember)
	teamMembersGroup.Get("", teamMemberHandler.ListTeamMembers)
	teamMembersGroup.Delete(":member_id", teamMemberHandler.RemoveTeamMember)

	apiKeyGroup := v1.Group("/api-keys", authMiddleware)
	apiKeyGroup.Post("", apiKeyHandler.CreateAPIKey)
	apiKeyGroup.Get("", apiKeyHandler.ListAPIKeys)
	apiKeyGroup.Get(":id", apiKeyHandler.GetAPIKey)
	apiKeyGroup.Delete(":id", apiKeyHandler.RevokeAPIKey)

	invitationGroup := v1.Group("/invitations")
	invitationGroup.Post(":token/accept", invitationHandler.AcceptInvitation)

	tenantInvitationGroup := v1.Group("/invitations", authMiddleware)
	tenantInvitationGroup.Post("", invitationHandler.CreateInvitation)
	tenantInvitationGroup.Get("", invitationHandler.ListInvitations)

	workflowGroup := v1.Group("/workflows", authMiddleware)
	workflowGroup.Post("", workflowHandler.CreateWorkflow)
	workflowGroup.Get("", workflowHandler.ListWorkflows)
	workflowGroup.Post(":id/execute", workflowHandler.ExecuteWorkflow)
	workflowGroup.Get(":id", workflowHandler.GetWorkflow)
	workflowGroup.Put(":id", workflowHandler.UpdateWorkflow)
	workflowGroup.Delete(":id", workflowHandler.DeleteWorkflow)

	logger.Info("Routes setup complete")
}
