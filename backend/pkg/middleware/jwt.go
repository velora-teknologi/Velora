package middleware

import (
	"strings"

	"github.com/gofiber/fiber/v2"
	"github.com/golang-jwt/jwt/v5"
	customErrors "github.com/velora-teknologi/velora/pkg/errors"
)

func JWTMiddleware(secret string) fiber.Handler {
	return func(c *fiber.Ctx) error {
		authHeader := c.Get("Authorization")
		if authHeader == "" {
			return c.Status(fiber.StatusUnauthorized).JSON(fiber.Map{
				"error": "missing authorization header",
			})
		}

		parts := strings.SplitN(authHeader, " ", 2)
		if len(parts) != 2 || parts[0] != "Bearer" {
			return c.Status(fiber.StatusUnauthorized).JSON(fiber.Map{
				"error": "invalid authorization header format",
			})
		}

		tokenString := parts[1]

		token, err := jwt.Parse(tokenString, func(token *jwt.Token) (interface{}, error) {
			return []byte(secret), nil
		})

		if err != nil || !token.Valid {
			return c.Status(fiber.StatusUnauthorized).JSON(fiber.Map{
				"error": "invalid token",
			})
		}

		claims := token.Claims.(jwt.MapClaims)
		c.Locals("user", claims)

		return c.Next()
	}
}

func ErrorHandler(c *fiber.Ctx, err error) error {
	if customErr, ok := err.(*customErrors.CustomError); ok {
		return c.Status(customErr.StatusCode).JSON(fiber.Map{
			"type":    customErr.Type,
			"message": customErr.Message,
			"details": customErr.Details,
		})
	}

	return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
		"type":    "INTERNAL_ERROR",
		"message": "internal server error",
	})
}
