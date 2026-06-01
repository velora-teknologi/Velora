package handlers

import (
	"github.com/gofiber/fiber/v2"
	"go.uber.org/zap"

	"github.com/velora-teknologi/velora/internal/application/dtos"
	"github.com/velora-teknologi/velora/internal/application/services"
	customErrors "github.com/velora-teknologi/velora/pkg/errors"
)

type WorkspaceHandler struct {
	service services.WorkspaceService
	logger  *zap.SugaredLogger
}

func NewWorkspaceHandler(service services.WorkspaceService, logger *zap.SugaredLogger) *WorkspaceHandler {
	return &WorkspaceHandler{service: service, logger: logger}
}

func (h *WorkspaceHandler) CreateWorkspace(c *fiber.Ctx) error {
	userID, authErr := getUserIDFromContext(c)
	if authErr != nil {
		return respondCustomError(c, authErr)
	}

	var req dtos.CreateWorkspaceRequest
	if err := c.BodyParser(&req); err != nil {
		return respondError(c, fiber.StatusBadRequest, "Invalid request body")
	}

	workspace, err := h.service.CreateWorkspace(c.Context(), userID, &req)
	if err != nil {
		if customErr, ok := err.(*customErrors.CustomError); ok {
			return respondCustomError(c, customErr)
		}
		return respondError(c, fiber.StatusInternalServerError, "Internal server error")
	}

	return respondSuccess(c, fiber.StatusCreated, workspace)
}

func (h *WorkspaceHandler) GetWorkspace(c *fiber.Ctx) error {
	userID, authErr := getUserIDFromContext(c)
	if authErr != nil {
		return respondCustomError(c, authErr)
	}

	id := c.Params("id")
	workspace, err := h.service.GetWorkspace(c.Context(), userID, id)
	if err != nil {
		if customErr, ok := err.(*customErrors.CustomError); ok {
			return respondCustomError(c, customErr)
		}
		return respondError(c, fiber.StatusInternalServerError, "Internal server error")
	}

	return respondSuccess(c, fiber.StatusOK, workspace)
}

func (h *WorkspaceHandler) UpdateWorkspace(c *fiber.Ctx) error {
	userID, authErr := getUserIDFromContext(c)
	if authErr != nil {
		return respondCustomError(c, authErr)
	}

	id := c.Params("id")
	var req dtos.UpdateWorkspaceRequest
	if err := c.BodyParser(&req); err != nil {
		return respondError(c, fiber.StatusBadRequest, "Invalid request body")
	}

	workspace, err := h.service.UpdateWorkspace(c.Context(), userID, id, &req)
	if err != nil {
		if customErr, ok := err.(*customErrors.CustomError); ok {
			return respondCustomError(c, customErr)
		}
		return respondError(c, fiber.StatusInternalServerError, "Internal server error")
	}

	return respondSuccess(c, fiber.StatusOK, workspace)
}

func (h *WorkspaceHandler) DeleteWorkspace(c *fiber.Ctx) error {
	userID, authErr := getUserIDFromContext(c)
	if authErr != nil {
		return respondCustomError(c, authErr)
	}

	id := c.Params("id")
	if err := h.service.DeleteWorkspace(c.Context(), userID, id); err != nil {
		if customErr, ok := err.(*customErrors.CustomError); ok {
			return respondCustomError(c, customErr)
		}
		return respondError(c, fiber.StatusInternalServerError, "Internal server error")
	}

	return c.SendStatus(fiber.StatusNoContent)
}

func (h *WorkspaceHandler) ListWorkspaces(c *fiber.Ctx) error {
	userID, authErr := getUserIDFromContext(c)
	if authErr != nil {
		return respondCustomError(c, authErr)
	}

	tenantID := c.Query("tenant_id")
	if tenantID == "" {
		return respondError(c, fiber.StatusBadRequest, "tenant_id is required")
	}

	skip := c.QueryInt("skip", 0)
	limit := c.QueryInt("limit", 20)
	if skip < 0 {
		skip = 0
	}
	if limit < 1 || limit > 100 {
		limit = 20
	}

	workspaces, err := h.service.ListWorkspaces(c.Context(), userID, tenantID, skip, limit)
	if err != nil {
		if customErr, ok := err.(*customErrors.CustomError); ok {
			return respondCustomError(c, customErr)
		}
		return respondError(c, fiber.StatusInternalServerError, "Internal server error")
	}

	return respondSuccess(c, fiber.StatusOK, workspaces)
}
