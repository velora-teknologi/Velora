package handlers

import (
	"github.com/gofiber/fiber/v2"
	"go.uber.org/zap"

	"github.com/velora-teknologi/velora/internal/application/dtos"
	"github.com/velora-teknologi/velora/internal/application/services"
	customErrors "github.com/velora-teknologi/velora/pkg/errors"
)

type InvitationHandler struct {
	service services.InvitationService
	logger  *zap.SugaredLogger
}

func NewInvitationHandler(service services.InvitationService, logger *zap.SugaredLogger) *InvitationHandler {
	return &InvitationHandler{service: service, logger: logger}
}

func (h *InvitationHandler) CreateInvitation(c *fiber.Ctx) error {
	userID, authErr := getUserIDFromContext(c)
	if authErr != nil {
		return respondCustomError(c, authErr)
	}

	var req dtos.CreateInvitationRequest
	if err := c.BodyParser(&req); err != nil {
		return respondError(c, fiber.StatusBadRequest, "Invalid request body")
	}

	invitation, err := h.service.CreateInvitation(c.Context(), userID, &req)
	if err != nil {
		if customErr, ok := err.(*customErrors.CustomError); ok {
			return respondCustomError(c, customErr)
		}
		return respondError(c, fiber.StatusInternalServerError, "Internal server error")
	}

	return respondSuccess(c, fiber.StatusCreated, invitation)
}

func (h *InvitationHandler) ListInvitations(c *fiber.Ctx) error {
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

	invitations, err := h.service.ListInvitations(c.Context(), userID, tenantID, skip, limit)
	if err != nil {
		if customErr, ok := err.(*customErrors.CustomError); ok {
			return respondCustomError(c, customErr)
		}
		return respondError(c, fiber.StatusInternalServerError, "Internal server error")
	}

	return respondSuccess(c, fiber.StatusOK, invitations)
}

func (h *InvitationHandler) AcceptInvitation(c *fiber.Ctx) error {
	token := c.Params("token")
	var req dtos.AcceptInvitationRequest
	if err := c.BodyParser(&req); err != nil {
		return respondError(c, fiber.StatusBadRequest, "Invalid request body")
	}

	resp, err := h.service.AcceptInvitation(c.Context(), token, &req)
	if err != nil {
		if customErr, ok := err.(*customErrors.CustomError); ok {
			return respondCustomError(c, customErr)
		}
		return respondError(c, fiber.StatusInternalServerError, "Internal server error")
	}

	return respondSuccess(c, fiber.StatusOK, resp)
}
