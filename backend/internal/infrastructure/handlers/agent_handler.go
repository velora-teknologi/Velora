package handlers

import (
	"github.com/gofiber/fiber/v2"
	"go.uber.org/zap"

	"github.com/velora-teknologi/velora/internal/application/dtos"
	"github.com/velora-teknologi/velora/internal/application/services"
	customErrors "github.com/velora-teknologi/velora/pkg/errors"
)

type AgentHandler struct {
	service         services.AgentService
	workflowService services.WorkflowService
	logger          *zap.SugaredLogger
}

func NewAgentHandler(service services.AgentService, workflowService services.WorkflowService, logger *zap.SugaredLogger) *AgentHandler {
	return &AgentHandler{service: service, workflowService: workflowService, logger: logger}
}

func (h *AgentHandler) CreateAgent(c *fiber.Ctx) error {
	userID, authErr := getUserIDFromContext(c)
	if authErr != nil {
		return respondCustomError(c, authErr)
	}

	var req dtos.CreateAgentRequest
	if err := c.BodyParser(&req); err != nil {
		return respondError(c, fiber.StatusBadRequest, "Invalid request body")
	}

	agent, err := h.service.CreateAgent(c.Context(), userID, &req)
	if err != nil {
		if customErr, ok := err.(*customErrors.CustomError); ok {
			return respondCustomError(c, customErr)
		}
		return respondError(c, fiber.StatusInternalServerError, "Internal server error")
	}

	return respondSuccess(c, fiber.StatusCreated, agent)
}

func (h *AgentHandler) GetAgent(c *fiber.Ctx) error {
	userID, authErr := getUserIDFromContext(c)
	if authErr != nil {
		return respondCustomError(c, authErr)
	}

	id := c.Params("id")
	agent, err := h.service.GetAgent(c.Context(), userID, id)
	if err != nil {
		if customErr, ok := err.(*customErrors.CustomError); ok {
			return respondCustomError(c, customErr)
		}
		return respondError(c, fiber.StatusInternalServerError, "Internal server error")
	}

	return respondSuccess(c, fiber.StatusOK, agent)
}

func (h *AgentHandler) UpdateAgent(c *fiber.Ctx) error {
	userID, authErr := getUserIDFromContext(c)
	if authErr != nil {
		return respondCustomError(c, authErr)
	}

	id := c.Params("id")
	var req dtos.UpdateAgentRequest
	if err := c.BodyParser(&req); err != nil {
		return respondError(c, fiber.StatusBadRequest, "Invalid request body")
	}

	agent, err := h.service.UpdateAgent(c.Context(), userID, id, &req)
	if err != nil {
		if customErr, ok := err.(*customErrors.CustomError); ok {
			return respondCustomError(c, customErr)
		}
		return respondError(c, fiber.StatusInternalServerError, "Internal server error")
	}

	return respondSuccess(c, fiber.StatusOK, agent)
}

func (h *AgentHandler) DeleteAgent(c *fiber.Ctx) error {
	userID, authErr := getUserIDFromContext(c)
	if authErr != nil {
		return respondCustomError(c, authErr)
	}

	id := c.Params("id")
	if err := h.service.DeleteAgent(c.Context(), userID, id); err != nil {
		if customErr, ok := err.(*customErrors.CustomError); ok {
			return respondCustomError(c, customErr)
		}
		return respondError(c, fiber.StatusInternalServerError, "Internal server error")
	}

	return c.SendStatus(fiber.StatusNoContent)
}

func (h *AgentHandler) ListAgents(c *fiber.Ctx) error {
	userID, authErr := getUserIDFromContext(c)
	if authErr != nil {
		return respondCustomError(c, authErr)
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
			return respondCustomError(c, customErr)
		}
		return respondError(c, fiber.StatusInternalServerError, "Internal server error")
	}

	return respondSuccess(c, fiber.StatusOK, agents)
}

func (h *AgentHandler) ListAgentWorkflows(c *fiber.Ctx) error {
	if h.workflowService == nil {
		return respondError(c, fiber.StatusInternalServerError, "workflow service unavailable")
	}

	userID, authErr := getUserIDFromContext(c)
	if authErr != nil {
		return respondCustomError(c, authErr)
	}

	agentID := c.Params("id")
	skip := c.QueryInt("skip", 0)
	limit := c.QueryInt("limit", 10)
	if skip < 0 {
		skip = 0
	}
	if limit < 1 || limit > 100 {
		limit = 10
	}

	workflows, err := h.workflowService.ListWorkflows(c.Context(), userID, agentID, skip, limit)
	if err != nil {
		if customErr, ok := err.(*customErrors.CustomError); ok {
			return respondCustomError(c, customErr)
		}
		return respondError(c, fiber.StatusInternalServerError, "Internal server error")
	}

	return respondSuccess(c, fiber.StatusOK, workflows)
}

func (h *AgentHandler) ExecuteAgent(c *fiber.Ctx) error {
	userID, authErr := getUserIDFromContext(c)
	if authErr != nil {
		return respondCustomError(c, authErr)
	}

	agentID := c.Params("id")
	var req dtos.AgentExecutionRequest
	if err := c.BodyParser(&req); err != nil {
		return respondError(c, fiber.StatusBadRequest, "Invalid request body")
	}

	resp, err := h.service.ExecuteAgent(c.Context(), userID, agentID, req.Input)
	if err != nil {
		if customErr, ok := err.(*customErrors.CustomError); ok {
			return respondCustomError(c, customErr)
		}
		return respondError(c, fiber.StatusInternalServerError, "Internal server error")
	}

	return respondSuccess(c, fiber.StatusAccepted, resp)
}
