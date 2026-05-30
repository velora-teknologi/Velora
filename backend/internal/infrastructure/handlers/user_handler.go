package handlers
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
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"error": "Invalid request body",
		})
	}

	// Validate request
	if req.Email == "" || req.Password == "" || req.FullName == "" {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"error": "Missing required fields",
		})
	}

	user, err := h.service.CreateUser(c.Context(), &req)
	if err != nil {
		if customErr, ok := err.(*customErrors.CustomError); ok {
			return c.Status(customErr.StatusCode).JSON(fiber.Map{
				"error": customErr.Message,
			})
		}
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"error": "Internal server error",
		})
	}

	return c.Status(fiber.StatusCreated).JSON(user)
}

// GetUser gets a user by ID
// GET /api/v1/users/:id
func (h *UserHandler) GetUser(c *fiber.Ctx) error {
	id := c.Params("id")

	user, err := h.service.GetUser(c.Context(), id)
	if err != nil {
		if customErr, ok := err.(*customErrors.CustomError); ok {
			return c.Status(customErr.StatusCode).JSON(fiber.Map{
				"error": customErr.Message,
			})
		}
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"error": "Internal server error",
		})
	}

	return c.JSON(user)
}

// UpdateUser updates a user
// PUT /api/v1/users/:id
func (h *UserHandler) UpdateUser(c *fiber.Ctx) error {
	id := c.Params("id")

	var req dtos.UpdateUserRequest
	if err := c.BodyParser(&req); err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"error": "Invalid request body",
		})
	}

	user, err := h.service.UpdateUser(c.Context(), id, &req)
	if err != nil {
		if customErr, ok := err.(*customErrors.CustomError); ok {
			return c.Status(customErr.StatusCode).JSON(fiber.Map{
				"error": customErr.Message,
			})
		}
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"error": "Internal server error",
		})
	}

	return c.JSON(user)
}

// DeleteUser deletes a user
// DELETE /api/v1/users/:id
func (h *UserHandler) DeleteUser(c *fiber.Ctx) error {
	id := c.Params("id")

	err := h.service.DeleteUser(c.Context(), id)
	if err != nil {
		if customErr, ok := err.(*customErrors.CustomError); ok {
			return c.Status(customErr.StatusCode).JSON(fiber.Map{
				"error": customErr.Message,
			})
		}
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"error": "Internal server error",
		})
	}

	return c.SendStatus(fiber.StatusNoContent)
}

// ListUsers lists all users
// GET /api/v1/users
func (h *UserHandler) ListUsers(c *fiber.Ctx) error {
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
			return c.Status(customErr.StatusCode).JSON(fiber.Map{
				"error": customErr.Message,
			})
		}
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"error": "Internal server error",
		})
	}

	return c.JSON(fiber.Map{
		"data": users,
	})
}

// LoginUser authenticates a user
// POST /api/v1/auth/login
func (h *UserHandler) LoginUser(c *fiber.Ctx) error {
	var req dtos.LoginRequest

	if err := c.BodyParser(&req); err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"error": "Invalid request body",
		})
	}

	// Validate request
	if req.Email == "" || req.Password == "" {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"error": "Missing required fields",
		})
	}

	resp, err := h.service.LoginUser(c.Context(), &req)
	if err != nil {
		if customErr, ok := err.(*customErrors.CustomError); ok {
			return c.Status(customErr.StatusCode).JSON(fiber.Map{
				"error": customErr.Message,
			})
		}
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"error": "Internal server error",
		})
	}

	return c.JSON(resp)
}
