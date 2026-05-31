package handlers

import (
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
