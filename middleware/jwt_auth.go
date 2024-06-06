package middleware

import (
	"fmt"
	"os"

	"github.com/gofiber/fiber/v2"
	"github.com/golang-jwt/jwt/v4"
)

func Authenticate(c *fiber.Ctx) error {
	userToken := c.Cookies("token")

	if userToken == "" {
		return c.Status(fiber.StatusUnauthorized).JSON(fiber.Map{
			"message": "Cannot access resources",
		})
	}

	token, err := jwt.Parse(userToken, func(token *jwt.Token) (interface{}, error) {
		if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
			return nil, fmt.Errorf("unexpected signing method: %v", token.Header["alg"])
		}
		return []byte(os.Getenv("SECRET_KEY")), nil
	})

	if err != nil {
		return c.Status(fiber.StatusUnauthorized).JSON(fiber.Map{
			"message": "your session has expired",
		})
	}

	if tokenValue, ok := token.Claims.(jwt.MapClaims); ok && token.Valid {
		userID, ok := tokenValue["user_id"].(string)
		if !ok {
			return c.Status(fiber.StatusUnauthorized).JSON(fiber.Map{
				"message": "unexpected value",
			})
		}

		c.Locals("user_id", userID)
		return c.Next()
	}

	return c.Status(fiber.StatusUnauthorized).JSON(fiber.Map{
		"message": "invalid session",
	})
}
