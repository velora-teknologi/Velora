package handlers

import (
	"github.com/gofiber/fiber/v2"
	"go.uber.org/zap"

	"github.com/velora-teknologi/velora/internal/application/dtos"
	"github.com/velora-teknologi/velora/internal/application/services"
	customErrors "github.com/velora-teknologi/velora/pkg/errors"
)

type UserHandler struct {
	service services.UserService
	logger  *zap.SugaredLogger
}

func NewUserHandler(service services.UserService, logger *zap.SugaredLogger) *UserHandler {
	return &UserHandler{
		service: service,
		logger:  logger,
	}
}

// CreateUser creates a new user
// POST /api/v1/users
func (h *UserHandler) CreateUser(c *fiber.Ctx) error {
	var req dtos.CreateUserRequest

	if err := c.BodyParser(&req); err != nil {
		return respondError(c, fiber.StatusBadRequest, "Invalid request body")
	}

	// Validate request
	if req.Email == "" || req.Password == "" || req.FullName == "" {
		return respondError(c, fiber.StatusBadRequest, "Missing required fields")
	}

	user, err := h.service.CreateUser(c.Context(), &req)
	if err != nil {
		if customErr, ok := err.(*customErrors.CustomError); ok {
			return respondCustomError(c, customErr)
		}
		return respondError(c, fiber.StatusInternalServerError, "Internal server error")
	}

	return respondSuccess(c, fiber.StatusCreated, user)
}

// RegisterUser registers a new user
// POST /api/v1/auth/register
func (h *UserHandler) RegisterUser(c *fiber.Ctx) error {
	return h.CreateUser(c)
}

// GetUser gets a user by ID
// GET /api/v1/users/:id
func (h *UserHandler) GetUser(c *fiber.Ctx) error {
	id := c.Params("id")

	if ok, err := canManageUser(c, id); !ok {
		return respondCustomError(c, err)
	}

	user, err := h.service.GetUser(c.Context(), id)
	if err != nil {
		if customErr, ok := err.(*customErrors.CustomError); ok {
			return respondCustomError(c, customErr)
		}
		return respondError(c, fiber.StatusInternalServerError, "Internal server error")
	}

	return respondSuccess(c, fiber.StatusOK, user)
}

// GetCurrentUser returns the authenticated user's profile
func (h *UserHandler) GetCurrentUser(c *fiber.Ctx) error {
	claims, ok := c.Locals("user").(map[string]interface{})
	if !ok {
		return respondError(c, fiber.StatusUnauthorized, "invalid authentication claims")
	}

	userID, ok := claims["sub"].(string)
	if !ok || userID == "" {
		return respondError(c, fiber.StatusUnauthorized, "missing user id in token")
	}

	user, err := h.service.GetUser(c.Context(), userID)
	if err != nil {
		if customErr, ok := err.(*customErrors.CustomError); ok {
			return respondCustomError(c, customErr)
		}
		return respondError(c, fiber.StatusInternalServerError, "Internal server error")
	}

	return respondSuccess(c, fiber.StatusOK, user)
}

// UpdateUser updates a user
// PUT /api/v1/users/:id
func (h *UserHandler) UpdateUser(c *fiber.Ctx) error {
	id := c.Params("id")

	var req dtos.UpdateUserRequest
	if err := c.BodyParser(&req); err != nil {
		return respondError(c, fiber.StatusBadRequest, "Invalid request body")
	}

	if req.Role != "" && !isAdmin(c) {
		return respondCustomError(c, customErrors.NewForbiddenError("insufficient permissions"))
	}

	if ok, err := canManageUser(c, id); !ok {
		return respondCustomError(c, err)
	}

	user, err := h.service.UpdateUser(c.Context(), id, &req)
	if err != nil {
		if customErr, ok := err.(*customErrors.CustomError); ok {
			return respondCustomError(c, customErr)
		}
		return respondError(c, fiber.StatusInternalServerError, "Internal server error")
	}

	return respondSuccess(c, fiber.StatusOK, user)
}

// DeleteUser deletes a user
// DELETE /api/v1/users/:id
func (h *UserHandler) DeleteUser(c *fiber.Ctx) error {
	id := c.Params("id")

	if ok, err := canManageUser(c, id); !ok {
		return respondCustomError(c, err)
	}

	err := h.service.DeleteUser(c.Context(), id)
	if err != nil {
		if customErr, ok := err.(*customErrors.CustomError); ok {
			return respondCustomError(c, customErr)
		}
		return respondError(c, fiber.StatusInternalServerError, "Internal server error")
	}

	return c.SendStatus(fiber.StatusNoContent)
}

// ListUsers lists all users
// GET /api/v1/users
func (h *UserHandler) ListUsers(c *fiber.Ctx) error {
	if !isAdmin(c) {
		return respondCustomError(c, customErrors.NewForbiddenError("insufficient permissions"))
	}

	skip := c.QueryInt("skip", 0)
	limit := c.QueryInt("limit", 10)

	// Validate pagination
	if skip < 0 {
		skip = 0
	}
	if limit < 1 || limit > 100 {
		limit = 10
	}

	users, err := h.service.ListUsers(c.Context(), skip, limit)
	if err != nil {
		if customErr, ok := err.(*customErrors.CustomError); ok {
			return respondCustomError(c, customErr)
		}
		return respondError(c, fiber.StatusInternalServerError, "Internal server error")
	}

	return respondSuccess(c, fiber.StatusOK, users)
}

// LoginUser authenticates a user
// POST /api/v1/auth/login
func (h *UserHandler) LoginUser(c *fiber.Ctx) error {
	var req dtos.LoginRequest

	if err := c.BodyParser(&req); err != nil {
		return respondError(c, fiber.StatusBadRequest, "Invalid request body")
	}

	// Validate request
	if req.Email == "" || req.Password == "" {
		return respondError(c, fiber.StatusBadRequest, "Missing required fields")
	}

	resp, err := h.service.LoginUser(c.Context(), &req)
	if err != nil {
		if customErr, ok := err.(*customErrors.CustomError); ok {
			return respondCustomError(c, customErr)
		}
		return respondError(c, fiber.StatusInternalServerError, "Internal server error")
	}

	return respondSuccess(c, fiber.StatusOK, resp)
}
