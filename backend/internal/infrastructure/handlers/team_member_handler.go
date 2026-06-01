package handlers

import (
	"github.com/gofiber/fiber/v2"
	"go.uber.org/zap"

	"github.com/velora-teknologi/velora/internal/application/dtos"
	"github.com/velora-teknologi/velora/internal/application/services"
	customErrors "github.com/velora-teknologi/velora/pkg/errors"
)

type TeamMemberHandler struct {
	service services.TeamMemberService
	logger  *zap.SugaredLogger
}

func NewTeamMemberHandler(service services.TeamMemberService, logger *zap.SugaredLogger) *TeamMemberHandler {
	return &TeamMemberHandler{service: service, logger: logger}
}

func (h *TeamMemberHandler) AddTeamMember(c *fiber.Ctx) error {
	userID, authErr := getUserIDFromContext(c)
	if authErr != nil {
		return respondCustomError(c, authErr)
	}

	teamID := c.Params("id")
	var req dtos.CreateTeamMemberRequest
	if err := c.BodyParser(&req); err != nil {
		return respondError(c, fiber.StatusBadRequest, "Invalid request body")
	}
	if req.TeamID != "" && req.TeamID != teamID {
		return respondError(c, fiber.StatusBadRequest, "Team ID in path and body must match")
	}
	req.TeamID = teamID

	member, err := h.service.AddTeamMember(c.Context(), userID, &req)
	if err != nil {
		if customErr, ok := err.(*customErrors.CustomError); ok {
			return respondCustomError(c, customErr)
		}
		return respondError(c, fiber.StatusInternalServerError, "Internal server error")
	}

	return respondSuccess(c, fiber.StatusCreated, member)
}

func (h *TeamMemberHandler) ListTeamMembers(c *fiber.Ctx) error {
	userID, authErr := getUserIDFromContext(c)
	if authErr != nil {
		return respondCustomError(c, authErr)
	}

	teamID := c.Params("id")
	skip := c.QueryInt("skip", 0)
	limit := c.QueryInt("limit", 20)
	if skip < 0 {
		skip = 0
	}
	if limit < 1 || limit > 100 {
		limit = 20
	}

	members, err := h.service.ListTeamMembers(c.Context(), userID, teamID, skip, limit)
	if err != nil {
		if customErr, ok := err.(*customErrors.CustomError); ok {
			return respondCustomError(c, customErr)
		}
		return respondError(c, fiber.StatusInternalServerError, "Internal server error")
	}

	return respondSuccess(c, fiber.StatusOK, members)
}

func (h *TeamMemberHandler) RemoveTeamMember(c *fiber.Ctx) error {
	userID, authErr := getUserIDFromContext(c)
	if authErr != nil {
		return respondCustomError(c, authErr)
	}

	teamID := c.Params("id")
	memberID := c.Params("member_id")

	if err := h.service.RemoveTeamMember(c.Context(), userID, teamID, memberID); err != nil {
		if customErr, ok := err.(*customErrors.CustomError); ok {
			return respondCustomError(c, customErr)
		}
		return respondError(c, fiber.StatusInternalServerError, "Internal server error")
	}

	return c.SendStatus(fiber.StatusNoContent)
}
