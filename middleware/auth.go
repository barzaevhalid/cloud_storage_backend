package middleware

import (
	"os"
	"strings"

	"github.com/gofiber/fiber/v2"
	"github.com/golang-jwt/jwt/v5"
)

func JWTProtected() fiber.Handler {
	return func(c *fiber.Ctx) error {
		tokenHeader := c.Get("Authorization")

		if tokenHeader == "" {
			return fiber.NewError(fiber.StatusUnauthorized, "missting token")
		}
		parts := strings.Split(tokenHeader, " ")

		if len(parts) != 2 || parts[0] != "Bearer" {
			return fiber.NewError(fiber.StatusUnauthorized, "inalid token format")
		}
		tokenStr := parts[1]
		secret := os.Getenv("JWT_SECRET")

		token, err := jwt.Parse(tokenStr, func(t *jwt.Token) (interface{}, error) {
			return []byte(secret), nil
		})

		if err != nil {
			return fiber.NewError(fiber.StatusUnauthorized, "invalid or expired token---")
		}
		claims := token.Claims.(jwt.MapClaims)

		userId := int64(claims["user_id"].(float64))

		c.Locals("user_id", userId)

		return c.Next()
	}
}
