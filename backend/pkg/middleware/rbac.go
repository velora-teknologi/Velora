package middleware

import (
	"strings"

	"github.com/gofiber/fiber/v2"
)

func RBACMiddleware(allowedRoles ...string) fiber.Handler {
	allowed := make(map[string]struct{}, len(allowedRoles))
	for _, role := range allowedRoles {
		allowed[strings.ToLower(role)] = struct{}{}
	}

	return func(c *fiber.Ctx) error {
		user, ok := c.Locals("user").(map[string]interface{})
		if !ok {
			return c.Status(fiber.StatusUnauthorized).JSON(fiber.Map{
				"success": false,
				"message": "invalid authentication claims",
				"errors":  []string{"invalid authentication claims"},
			})
		}

		roleValue, ok := user["role"].(string)
		if !ok || roleValue == "" {
			return c.Status(fiber.StatusForbidden).JSON(fiber.Map{
				"success": false,
				"message": "insufficient permissions",
				"errors":  []string{"insufficient permissions"},
			})
		}

		if _, ok := allowed[strings.ToLower(roleValue)]; !ok {
			return c.Status(fiber.StatusForbidden).JSON(fiber.Map{
				"success": false,
				"message": "insufficient permissions",
				"errors":  []string{"insufficient permissions"},
			})
		}

		return c.Next()
	}
}
