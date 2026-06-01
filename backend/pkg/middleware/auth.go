package middleware

import (
	"strings"
	"time"

	"github.com/gofiber/fiber/v2"
	"github.com/golang-jwt/jwt/v5"
	"golang.org/x/crypto/bcrypt"

	"github.com/velora-teknologi/velora/internal/domain/repositories"
)

func AuthMiddleware(secret string, apiKeyRepo repositories.APIKeyRepository, userRepo repositories.UserRepository) fiber.Handler {
	return func(c *fiber.Ctx) error {
		authHeader := c.Get("Authorization")
		apiKeyHeader := c.Get("X-API-Key")

		if authHeader == "" && apiKeyHeader == "" {
			return c.Status(fiber.StatusUnauthorized).JSON(fiber.Map{
				"error": "missing authorization header",
			})
		}

		if authHeader != "" {
			parts := strings.SplitN(authHeader, " ", 2)
			if len(parts) == 2 && strings.EqualFold(parts[0], "Bearer") {
				return jwtAuth(c, secret)
			}
			if len(parts) == 2 && strings.EqualFold(parts[0], "ApiKey") {
				return apiKeyAuth(c, apiKeyRepo, userRepo, parts[1])
			}
		}

		if apiKeyHeader != "" {
			return apiKeyAuth(c, apiKeyRepo, userRepo, apiKeyHeader)
		}

		return c.Status(fiber.StatusUnauthorized).JSON(fiber.Map{
			"error": "invalid authorization header format",
		})
	}
}

func jwtAuth(c *fiber.Ctx, secret string) error {
	tokenString := strings.TrimSpace(strings.TrimPrefix(c.Get("Authorization"), "Bearer"))
	if tokenString == "" {
		return c.Status(fiber.StatusUnauthorized).JSON(fiber.Map{
			"error": "invalid authorization header format",
		})
	}

	token, err := jwt.Parse(tokenString, func(token *jwt.Token) (interface{}, error) {
		return []byte(secret), nil
	})

	if err != nil || !token.Valid {
		return c.Status(fiber.StatusUnauthorized).JSON(fiber.Map{
			"error": "invalid token",
		})
	}

	claims, ok := token.Claims.(jwt.MapClaims)
	if !ok {
		return c.Status(fiber.StatusUnauthorized).JSON(fiber.Map{
			"error": "invalid token claims",
		})
	}

	c.Locals("user", claims)
	return c.Next()
}

func apiKeyAuth(c *fiber.Ctx, apiKeyRepo repositories.APIKeyRepository, userRepo repositories.UserRepository, rawKey string) error {
	parts := strings.SplitN(rawKey, ".", 2)
	if len(parts) != 2 || parts[0] == "" || parts[1] == "" {
		return c.Status(fiber.StatusUnauthorized).JSON(fiber.Map{
			"error": "invalid api key format",
		})
	}

	keyID := parts[0]
	secret := parts[1]

	apiKey, err := apiKeyRepo.GetByKeyID(c.Context(), keyID)
	if err != nil {
		return c.Status(fiber.StatusUnauthorized).JSON(fiber.Map{
			"error": "invalid api key",
		})
	}
	if apiKey == nil || !apiKey.IsActive {
		return c.Status(fiber.StatusUnauthorized).JSON(fiber.Map{
			"error": "invalid api key",
		})
	}
	if apiKey.ExpiresAt != nil && time.Now().After(*apiKey.ExpiresAt) {
		return c.Status(fiber.StatusUnauthorized).JSON(fiber.Map{
			"error": "expired api key",
		})
	}
	if err := bcrypt.CompareHashAndPassword([]byte(apiKey.KeyHash), []byte(secret)); err != nil {
		return c.Status(fiber.StatusUnauthorized).JSON(fiber.Map{
			"error": "invalid api key",
		})
	}

	user, err := userRepo.GetByID(c.Context(), apiKey.UserID)
	if err != nil || user == nil {
		return c.Status(fiber.StatusUnauthorized).JSON(fiber.Map{
			"error": "invalid api key",
		})
	}

	claims := map[string]interface{}{
		"sub":  user.ID,
		"role": user.Role,
	}
	c.Locals("user", claims)
	c.Locals("api_key_id", apiKey.ID)

	return c.Next()
}
