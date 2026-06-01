package handlers

import (
	"github.com/gofiber/fiber/v2"
	"go.uber.org/zap"

	"github.com/velora-teknologi/velora/internal/application/dtos"
	"github.com/velora-teknologi/velora/internal/application/services"
	customErrors "github.com/velora-teknologi/velora/pkg/errors"
)

type APIKeyHandler struct {
	service services.APIKeyService
	logger  *zap.SugaredLogger
}

func NewAPIKeyHandler(service services.APIKeyService, logger *zap.SugaredLogger) *APIKeyHandler {
	return &APIKeyHandler{service: service, logger: logger}
}

// CreateAPIKey creates a new API key for the authenticated user
// POST /api/v1/api-keys
func (h *APIKeyHandler) CreateAPIKey(c *fiber.Ctx) error {
	userID, err := getUserIDFromContext(c)
	if err != nil {
		return respondCustomError(c, err)
	}

	var req dtos.CreateAPIKeyRequest
	if err := c.BodyParser(&req); err != nil {
		return respondError(c, fiber.StatusBadRequest, "Invalid request body")
	}

	apiKey, svcErr := h.service.CreateAPIKey(c.Context(), userID, &req)
	if svcErr != nil {
		return handleServiceError(c, svcErr)
	}

	return respondSuccess(c, fiber.StatusCreated, apiKey)
}

// GetAPIKey returns an API key by ID
// GET /api/v1/api-keys/:id
func (h *APIKeyHandler) GetAPIKey(c *fiber.Ctx) error {
	userID, err := getUserIDFromContext(c)
	if err != nil {
		return respondCustomError(c, err)
	}

	id := c.Params("id")
	apiKey, svcErr := h.service.GetAPIKey(c.Context(), userID, id)
	if svcErr != nil {
		return handleServiceError(c, svcErr)
	}

	return respondSuccess(c, fiber.StatusOK, apiKey)
}

// ListAPIKeys lists API keys for the authenticated user
// GET /api/v1/api-keys
func (h *APIKeyHandler) ListAPIKeys(c *fiber.Ctx) error {
	userID, err := getUserIDFromContext(c)
	if err != nil {
		return respondCustomError(c, err)
	}

	skip := c.QueryInt("skip", 0)
	limit := c.QueryInt("limit", 10)
	if skip < 0 {
		skip = 0
	}
	if limit < 1 || limit > 100 {
		limit = 10
	}

	apiKeys, svcErr := h.service.ListAPIKeys(c.Context(), userID, skip, limit)
	if svcErr != nil {
		return handleServiceError(c, svcErr)
	}

	return respondSuccess(c, fiber.StatusOK, apiKeys)
}

// RevokeAPIKey deactivates an existing API key
// DELETE /api/v1/api-keys/:id
func (h *APIKeyHandler) RevokeAPIKey(c *fiber.Ctx) error {
	userID, err := getUserIDFromContext(c)
	if err != nil {
		return respondCustomError(c, err)
	}

	id := c.Params("id")
	if svcErr := h.service.RevokeAPIKey(c.Context(), userID, id); svcErr != nil {
		return handleServiceError(c, svcErr)
	}

	return c.SendStatus(fiber.StatusNoContent)
}

func handleServiceError(c *fiber.Ctx, err error) error {
	if customErr, ok := err.(*customErrors.CustomError); ok {
		return respondCustomError(c, customErr)
	}
	return respondError(c, fiber.StatusInternalServerError, "Internal server error")
}
