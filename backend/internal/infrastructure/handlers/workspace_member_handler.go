package handlers

import (
	"github.com/gofiber/fiber/v2"
	"go.uber.org/zap"

	"github.com/velora-teknologi/velora/internal/application/dtos"
	"github.com/velora-teknologi/velora/internal/application/services"
	customErrors "github.com/velora-teknologi/velora/pkg/errors"
)

type WorkspaceMemberHandler struct {
	service services.WorkspaceMemberService
	logger  *zap.SugaredLogger
}

func NewWorkspaceMemberHandler(service services.WorkspaceMemberService, logger *zap.SugaredLogger) *WorkspaceMemberHandler {
	return &WorkspaceMemberHandler{service: service, logger: logger}
}

func (h *WorkspaceMemberHandler) AddWorkspaceMember(c *fiber.Ctx) error {
	userID, authErr := getUserIDFromContext(c)
	if authErr != nil {
		return respondCustomError(c, authErr)
	}

	workspaceID := c.Params("id")
	var req dtos.CreateWorkspaceMemberRequest
	if err := c.BodyParser(&req); err != nil {
		return respondError(c, fiber.StatusBadRequest, "Invalid request body")
	}
	if req.WorkspaceID != "" && req.WorkspaceID != workspaceID {
		return respondError(c, fiber.StatusBadRequest, "Workspace ID in path and body must match")
	}
	req.WorkspaceID = workspaceID

	member, err := h.service.AddWorkspaceMember(c.Context(), userID, &req)
	if err != nil {
		if customErr, ok := err.(*customErrors.CustomError); ok {
			return respondCustomError(c, customErr)
		}
		return respondError(c, fiber.StatusInternalServerError, "Internal server error")
	}

	return respondSuccess(c, fiber.StatusCreated, member)
}

func (h *WorkspaceMemberHandler) ListWorkspaceMembers(c *fiber.Ctx) error {
	userID, authErr := getUserIDFromContext(c)
	if authErr != nil {
		return respondCustomError(c, authErr)
	}

	workspaceID := c.Params("id")
	skip := c.QueryInt("skip", 0)
	limit := c.QueryInt("limit", 20)
	if skip < 0 {
		skip = 0
	}
	if limit < 1 || limit > 100 {
		limit = 20
	}

	members, err := h.service.ListWorkspaceMembers(c.Context(), userID, workspaceID, skip, limit)
	if err != nil {
		if customErr, ok := err.(*customErrors.CustomError); ok {
			return respondCustomError(c, customErr)
		}
		return respondError(c, fiber.StatusInternalServerError, "Internal server error")
	}

	return respondSuccess(c, fiber.StatusOK, members)
}

func (h *WorkspaceMemberHandler) RemoveWorkspaceMember(c *fiber.Ctx) error {
	userID, authErr := getUserIDFromContext(c)
	if authErr != nil {
		return respondCustomError(c, authErr)
	}

	workspaceID := c.Params("id")
	memberID := c.Params("member_id")

	if err := h.service.RemoveWorkspaceMember(c.Context(), userID, workspaceID, memberID); err != nil {
		if customErr, ok := err.(*customErrors.CustomError); ok {
			return respondCustomError(c, customErr)
		}
		return respondError(c, fiber.StatusInternalServerError, "Internal server error")
	}

	return c.SendStatus(fiber.StatusNoContent)
}
