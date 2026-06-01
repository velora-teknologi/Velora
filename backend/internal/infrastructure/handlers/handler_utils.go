package handlers

import (
	"strings"

	"github.com/gofiber/fiber/v2"
	customErrors "github.com/velora-teknologi/velora/pkg/errors"
)

func getUserIDFromContext(c *fiber.Ctx) (string, *customErrors.CustomError) {
	claims, ok := c.Locals("user").(map[string]interface{})
	if !ok {
		return "", customErrors.NewUnauthorizedError("invalid authentication claims")
	}

	userID, ok := claims["sub"].(string)
	if !ok || userID == "" {
		return "", customErrors.NewUnauthorizedError("missing user id in token")
	}

	return userID, nil
}

func getUserRoleFromContext(c *fiber.Ctx) (string, *customErrors.CustomError) {
	claims, ok := c.Locals("user").(map[string]interface{})
	if !ok {
		return "", customErrors.NewUnauthorizedError("invalid authentication claims")
	}

	role, ok := claims["role"].(string)
	if !ok || role == "" {
		return "", customErrors.NewForbiddenError("insufficient permissions")
	}

	return role, nil
}

func isAdmin(c *fiber.Ctx) bool {
	role, err := getUserRoleFromContext(c)
	return err == nil && strings.EqualFold(role, "admin")
}

func canManageUser(c *fiber.Ctx, targetUserID string) (bool, *customErrors.CustomError) {
	currentUserID, err := getUserIDFromContext(c)
	if err != nil {
		return false, err
	}

	if currentUserID == targetUserID {
		return true, nil
	}

	if isAdmin(c) {
		return true, nil
	}

	return false, customErrors.NewForbiddenError("insufficient permissions")
}

func respondSuccess(c *fiber.Ctx, status int, data any) error {
	return c.Status(status).JSON(fiber.Map{
		"success": true,
		"message": "success",
		"data":    data,
	})
}

func respondError(c *fiber.Ctx, status int, message string) error {
	return c.Status(status).JSON(fiber.Map{
		"success": false,
		"message": message,
		"errors":  []string{message},
	})
}

func respondCustomError(c *fiber.Ctx, err *customErrors.CustomError) error {
	return respondError(c, err.StatusCode, err.Message)
}
