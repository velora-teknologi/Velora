package handlers

import (
	"github.com/gofiber/fiber/v2"
	"go.uber.org/zap"

	"github.com/velora-teknologi/velora/internal/application/dtos"
	"github.com/velora-teknologi/velora/internal/application/services"
	customErrors "github.com/velora-teknologi/velora/pkg/errors"
)

type WorkflowHandler struct {
	service services.WorkflowService
	logger  *zap.SugaredLogger
}

func NewWorkflowHandler(service services.WorkflowService, logger *zap.SugaredLogger) *WorkflowHandler {
	return &WorkflowHandler{service: service, logger: logger}
}

func (h *WorkflowHandler) CreateWorkflow(c *fiber.Ctx) error {
	userID, authErr := getUserIDFromContext(c)
	if authErr != nil {
		return respondCustomError(c, authErr)
	}

	var req dtos.CreateWorkflowRequest
	if err := c.BodyParser(&req); err != nil {
		return respondError(c, fiber.StatusBadRequest, "Invalid request body")
	}

	workflow, err := h.service.CreateWorkflow(c.Context(), userID, &req)
	if err != nil {
		if customErr, ok := err.(*customErrors.CustomError); ok {
			return respondCustomError(c, customErr)
		}
		return respondError(c, fiber.StatusInternalServerError, "Internal server error")
	}

	return respondSuccess(c, fiber.StatusCreated, workflow)
}

func (h *WorkflowHandler) GetWorkflow(c *fiber.Ctx) error {
	userID, authErr := getUserIDFromContext(c)
	if authErr != nil {
		return respondCustomError(c, authErr)
	}

	id := c.Params("id")
	workflow, err := h.service.GetWorkflow(c.Context(), userID, id)
	if err != nil {
		if customErr, ok := err.(*customErrors.CustomError); ok {
			return respondCustomError(c, customErr)
		}
		return respondError(c, fiber.StatusInternalServerError, "Internal server error")
	}

	return respondSuccess(c, fiber.StatusOK, workflow)
}

func (h *WorkflowHandler) UpdateWorkflow(c *fiber.Ctx) error {
	userID, authErr := getUserIDFromContext(c)
	if authErr != nil {
		return respondCustomError(c, authErr)
	}

	id := c.Params("id")
	var req dtos.UpdateWorkflowRequest
	if err := c.BodyParser(&req); err != nil {
		return respondError(c, fiber.StatusBadRequest, "Invalid request body")
	}

	workflow, err := h.service.UpdateWorkflow(c.Context(), userID, id, &req)
	if err != nil {
		if customErr, ok := err.(*customErrors.CustomError); ok {
			return respondCustomError(c, customErr)
		}
		return respondError(c, fiber.StatusInternalServerError, "Internal server error")
	}

	return respondSuccess(c, fiber.StatusOK, workflow)
}

func (h *WorkflowHandler) DeleteWorkflow(c *fiber.Ctx) error {
	userID, authErr := getUserIDFromContext(c)
	if authErr != nil {
		return respondCustomError(c, authErr)
	}

	id := c.Params("id")
	if err := h.service.DeleteWorkflow(c.Context(), userID, id); err != nil {
		if customErr, ok := err.(*customErrors.CustomError); ok {
			return respondCustomError(c, customErr)
		}
		return respondError(c, fiber.StatusInternalServerError, "Internal server error")
	}

	return c.SendStatus(fiber.StatusNoContent)
}

func (h *WorkflowHandler) ListWorkflows(c *fiber.Ctx) error {
	userID, authErr := getUserIDFromContext(c)
	if authErr != nil {
		return respondCustomError(c, authErr)
	}

	agentID := c.Query("agent_id")
	skip := c.QueryInt("skip", 0)
	limit := c.QueryInt("limit", 10)
	if skip < 0 {
		skip = 0
	}
	if limit < 1 || limit > 100 {
		limit = 10
	}

	workflows, err := h.service.ListWorkflows(c.Context(), userID, agentID, skip, limit)
	if err != nil {
		if customErr, ok := err.(*customErrors.CustomError); ok {
			return respondCustomError(c, customErr)
		}
		return respondError(c, fiber.StatusInternalServerError, "Internal server error")
	}

	return respondSuccess(c, fiber.StatusOK, workflows)
}

func (h *WorkflowHandler) ExecuteWorkflow(c *fiber.Ctx) error {
	userID, authErr := getUserIDFromContext(c)
	if authErr != nil {
		return respondCustomError(c, authErr)
	}

	id := c.Params("id")
	execution, err := h.service.ExecuteWorkflow(c.Context(), userID, id)
	if err != nil {
		if customErr, ok := err.(*customErrors.CustomError); ok {
			return respondCustomError(c, customErr)
		}
		return respondError(c, fiber.StatusInternalServerError, "Internal server error")
	}

	return respondSuccess(c, fiber.StatusAccepted, execution)
}
