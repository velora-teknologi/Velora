package handlers

import (
	"github.com/gofiber/fiber/v2"
	"go.uber.org/zap"

	"github.com/velora-teknologi/velora/internal/application/dtos"
	"github.com/velora-teknologi/velora/internal/application/services"
	customErrors "github.com/velora-teknologi/velora/pkg/errors"
)

type TeamHandler struct {
	service services.TeamService
	logger  *zap.SugaredLogger
}

func NewTeamHandler(service services.TeamService, logger *zap.SugaredLogger) *TeamHandler {
	return &TeamHandler{service: service, logger: logger}
}

func (h *TeamHandler) CreateTeam(c *fiber.Ctx) error {
	userID, authErr := getUserIDFromContext(c)
	if authErr != nil {
		return respondCustomError(c, authErr)
	}

	var req dtos.CreateTeamRequest
	if err := c.BodyParser(&req); err != nil {
		return respondError(c, fiber.StatusBadRequest, "Invalid request body")
	}

	team, err := h.service.CreateTeam(c.Context(), userID, &req)
	if err != nil {
		if customErr, ok := err.(*customErrors.CustomError); ok {
			return respondCustomError(c, customErr)
		}
		return respondError(c, fiber.StatusInternalServerError, "Internal server error")
	}

	return respondSuccess(c, fiber.StatusCreated, team)
}

func (h *TeamHandler) GetTeam(c *fiber.Ctx) error {
	userID, authErr := getUserIDFromContext(c)
	if authErr != nil {
		return respondCustomError(c, authErr)
	}

	id := c.Params("id")
	team, err := h.service.GetTeam(c.Context(), userID, id)
	if err != nil {
		if customErr, ok := err.(*customErrors.CustomError); ok {
			return respondCustomError(c, customErr)
		}
		return respondError(c, fiber.StatusInternalServerError, "Internal server error")
	}

	return respondSuccess(c, fiber.StatusOK, team)
}

func (h *TeamHandler) UpdateTeam(c *fiber.Ctx) error {
	userID, authErr := getUserIDFromContext(c)
	if authErr != nil {
		return respondCustomError(c, authErr)
	}

	id := c.Params("id")
	var req dtos.UpdateTeamRequest
	if err := c.BodyParser(&req); err != nil {
		return respondError(c, fiber.StatusBadRequest, "Invalid request body")
	}

	team, err := h.service.UpdateTeam(c.Context(), userID, id, &req)
	if err != nil {
		if customErr, ok := err.(*customErrors.CustomError); ok {
			return respondCustomError(c, customErr)
		}
		return respondError(c, fiber.StatusInternalServerError, "Internal server error")
	}

	return respondSuccess(c, fiber.StatusOK, team)
}

func (h *TeamHandler) DeleteTeam(c *fiber.Ctx) error {
	userID, authErr := getUserIDFromContext(c)
	if authErr != nil {
		return respondCustomError(c, authErr)
	}

	id := c.Params("id")
	if err := h.service.DeleteTeam(c.Context(), userID, id); err != nil {
		if customErr, ok := err.(*customErrors.CustomError); ok {
			return respondCustomError(c, customErr)
		}
		return respondError(c, fiber.StatusInternalServerError, "Internal server error")
	}

	return c.SendStatus(fiber.StatusNoContent)
}

func (h *TeamHandler) ListTeams(c *fiber.Ctx) error {
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

	teams, err := h.service.ListTeams(c.Context(), userID, tenantID, skip, limit)
	if err != nil {
		if customErr, ok := err.(*customErrors.CustomError); ok {
			return respondCustomError(c, customErr)
		}
		return respondError(c, fiber.StatusInternalServerError, "Internal server error")
	}

	return respondSuccess(c, fiber.StatusOK, teams)
}
