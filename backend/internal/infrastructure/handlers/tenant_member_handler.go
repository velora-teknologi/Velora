package handlers

import (
	"github.com/gofiber/fiber/v2"
	"go.uber.org/zap"

	"github.com/velora-teknologi/velora/internal/application/services"
	customErrors "github.com/velora-teknologi/velora/pkg/errors"
)

type TenantMemberHandler struct {
	service services.TenantMemberService
	logger  *zap.SugaredLogger
}

func NewTenantMemberHandler(service services.TenantMemberService, logger *zap.SugaredLogger) *TenantMemberHandler {
	return &TenantMemberHandler{service: service, logger: logger}
}

func (h *TenantMemberHandler) ListTenantMembers(c *fiber.Ctx) error {
	userID, authErr := getUserIDFromContext(c)
	if authErr != nil {
		return respondCustomError(c, authErr)
	}

	tenantID := c.Params("id")
	skip, limit := getPagination(c)

	members, err := h.service.ListTenantMembers(c.Context(), userID, tenantID, skip, limit)
	if err != nil {
		if customErr, ok := err.(*customErrors.CustomError); ok {
			return respondCustomError(c, customErr)
		}
		return respondError(c, fiber.StatusInternalServerError, "Internal server error")
	}

	return respondSuccess(c, fiber.StatusOK, members)
}
