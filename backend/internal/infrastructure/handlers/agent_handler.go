package handlers

import (
	"github.com/gofiber/fiber/v2"
	"go.uber.org/zap"

	"github.com/velora-teknologi/velora/internal/application/dtos"
	"github.com/velora-teknologi/velora/internal/application/services"
	customErrors "github.com/velora-teknologi/velora/pkg/errors"
)

type AgentHandler struct {
	service services.AgentService
	logger  *zap.SugaredLogger
}

func NewAgentHandler(service services.AgentService, logger *zap.SugaredLogger) *AgentHandler {
	return &AgentHandler{service: service, logger: logger}
}

func (h *AgentHandler) CreateAgent(c *fiber.Ctx) error {
	userID, authErr := getUserIDFromContext(c)
	if authErr != nil {
		return c.Status(authErr.StatusCode).JSON(fiber.Map{"error": authErr.Message})
	}

	var req dtos.CreateAgentRequest
	if err := c.BodyParser(&req); err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"error": "Invalid request body"})
	}

	agent, err := h.service.CreateAgent(c.Context(), userID, &req)
	if err != nil {
		if customErr, ok := err.(*customErrors.CustomError); ok {
			return c.Status(customErr.StatusCode).JSON(fiber.Map{"error": customErr.Message})
		}
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{"error": "Internal server error"})
	}

	return c.Status(fiber.StatusCreated).JSON(agent)
}

func (h *AgentHandler) GetAgent(c *fiber.Ctx) error {
	userID, authErr := getUserIDFromContext(c)
	if authErr != nil {
		return c.Status(authErr.StatusCode).JSON(fiber.Map{"error": authErr.Message})
	}

	id := c.Params("id")
	agent, err := h.service.GetAgent(c.Context(), userID, id)
	if err != nil {
		if customErr, ok := err.(*customErrors.CustomError); ok {
			return c.Status(customErr.StatusCode).JSON(fiber.Map{"error": customErr.Message})
		}
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{"error": "Internal server error"})
	}

	return c.JSON(agent)
}

func (h *AgentHandler) UpdateAgent(c *fiber.Ctx) error {
	userID, authErr := getUserIDFromContext(c)
	if authErr != nil {
		return c.Status(authErr.StatusCode).JSON(fiber.Map{"error": authErr.Message})
	}

	id := c.Params("id")
	var req dtos.UpdateAgentRequest
	if err := c.BodyParser(&req); err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"error": "Invalid request body"})
	}

	agent, err := h.service.UpdateAgent(c.Context(), userID, id, &req)
	if err != nil {
		if customErr, ok := err.(*customErrors.CustomError); ok {
			return c.Status(customErr.StatusCode).JSON(fiber.Map{"error": customErr.Message})
		}
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{"error": "Internal server error"})
	}

	return c.JSON(agent)
}

func (h *AgentHandler) DeleteAgent(c *fiber.Ctx) error {
	userID, authErr := getUserIDFromContext(c)
	if authErr != nil {
		return c.Status(authErr.StatusCode).JSON(fiber.Map{"error": authErr.Message})
	}

	id := c.Params("id")
	if err := h.service.DeleteAgent(c.Context(), userID, id); err != nil {
		if customErr, ok := err.(*customErrors.CustomError); ok {
			return c.Status(customErr.StatusCode).JSON(fiber.Map{"error": customErr.Message})
		}
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{"error": "Internal server error"})
	}

	return c.SendStatus(fiber.StatusNoContent)
}

func (h *AgentHandler) ListAgents(c *fiber.Ctx) error {
	userID, authErr := getUserIDFromContext(c)
	if authErr != nil {
		return c.Status(authErr.StatusCode).JSON(fiber.Map{"error": authErr.Message})
	}

	skip := c.QueryInt("skip", 0)
	limit := c.QueryInt("limit", 10)
	if skip < 0 {
		skip = 0
	}
	if limit < 1 || limit > 100 {
		limit = 10
	}

	agents, err := h.service.ListAgents(c.Context(), userID, skip, limit)
	if err != nil {
		if customErr, ok := err.(*customErrors.CustomError); ok {
			return c.Status(customErr.StatusCode).JSON(fiber.Map{"error": customErr.Message})
		}
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{"error": "Internal server error"})
	}

	return c.JSON(fiber.Map{"data": agents})
}
