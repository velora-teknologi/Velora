package handlers

import (
	"github.com/gofiber/fiber/v2"
	"go.uber.org/zap"

	"github.com/velora-teknologi/velora/internal/application/dtos"
	"github.com/velora-teknologi/velora/internal/application/services"
	customErrors "github.com/velora-teknologi/velora/pkg/errors"
)

type TenantHandler struct {
	service services.TenantService
	logger  *zap.SugaredLogger
}

func NewTenantHandler(service services.TenantService, logger *zap.SugaredLogger) *TenantHandler {
	return &TenantHandler{service: service, logger: logger}
}

func (h *TenantHandler) CreateTenant(c *fiber.Ctx) error {
	userID, authErr := getUserIDFromContext(c)
	if authErr != nil {
		return respondCustomError(c, authErr)
	}
	var req dtos.CreateTenantRequest
	if err := c.BodyParser(&req); err != nil {
		h.logger.Warnf("Invalid create tenant payload: %v", err)
		return respondError(c, fiber.StatusBadRequest, "Invalid request payload")
	}

	tenant, err := h.service.CreateTenant(c.Context(), userID, &req)
	if err != nil {
		if customErr, ok := err.(*customErrors.CustomError); ok {
			return respondCustomError(c, customErr)
		}
		return respondError(c, fiber.StatusInternalServerError, "Internal server error")
	}

	return respondSuccess(c, fiber.StatusCreated, tenant)
}

func (h *TenantHandler) GetTenant(c *fiber.Ctx) error {
	userID, authErr := getUserIDFromContext(c)
	if authErr != nil {
		return respondCustomError(c, authErr)
	}
	id := c.Params("id")

	tenant, err := h.service.GetTenant(c.Context(), userID, id)
	if err != nil {
		if customErr, ok := err.(*customErrors.CustomError); ok {
			return respondCustomError(c, customErr)
		}
		return respondError(c, fiber.StatusInternalServerError, "Internal server error")
	}

	return respondSuccess(c, fiber.StatusOK, tenant)
}

func (h *TenantHandler) UpdateTenant(c *fiber.Ctx) error {
	userID, authErr := getUserIDFromContext(c)
	if authErr != nil {
		return respondCustomError(c, authErr)
	}
	id := c.Params("id")
	var req dtos.UpdateTenantRequest
	if err := c.BodyParser(&req); err != nil {
		h.logger.Warnf("Invalid update tenant payload: %v", err)
		return respondError(c, fiber.StatusBadRequest, "Invalid request payload")
	}

	tenant, err := h.service.UpdateTenant(c.Context(), userID, id, &req)
	if err != nil {
		if customErr, ok := err.(*customErrors.CustomError); ok {
			return respondCustomError(c, customErr)
		}
		return respondError(c, fiber.StatusInternalServerError, "Internal server error")
	}

	return respondSuccess(c, fiber.StatusOK, tenant)
}

func (h *TenantHandler) DeleteTenant(c *fiber.Ctx) error {
	userID, authErr := getUserIDFromContext(c)
	if authErr != nil {
		return respondCustomError(c, authErr)
	}
	id := c.Params("id")

	if err := h.service.DeleteTenant(c.Context(), userID, id); err != nil {
		if customErr, ok := err.(*customErrors.CustomError); ok {
			return respondCustomError(c, customErr)
		}
		return respondError(c, fiber.StatusInternalServerError, "Internal server error")
	}

	return c.SendStatus(fiber.StatusNoContent)
}

func (h *TenantHandler) ListTenants(c *fiber.Ctx) error {
	userID, authErr := getUserIDFromContext(c)
	if authErr != nil {
		return respondCustomError(c, authErr)
	}
	skip, limit := getPagination(c)

	tenants, err := h.service.ListTenants(c.Context(), userID, skip, limit)
	if err != nil {
		if customErr, ok := err.(*customErrors.CustomError); ok {
			return respondCustomError(c, customErr)
		}
		return respondError(c, fiber.StatusInternalServerError, "Internal server error")
	}

	return respondSuccess(c, fiber.StatusOK, tenants)
}

func getPagination(c *fiber.Ctx) (int, int) {
	page := c.QueryInt("page", 1)
	if page < 1 {
		page = 1
	}
	limit := c.QueryInt("limit", 20)
	if limit < 1 {
		limit = 20
	}
	return (page - 1) * limit, limit
}
