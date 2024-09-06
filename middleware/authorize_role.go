package middleware

import (
	"app/model"

	"github.com/gofiber/fiber/v2"
)

func (m *implMiddleware) Authorize(allowedRoles ...int) fiber.Handler {
	return func(c *fiber.Ctx) error {
		user, ok := c.Locals("user").(*model.User)
		if !ok {
			return c.Status(fiber.StatusForbidden).JSON(fiber.Map{
				"message": "Access denied",
			})
		}

		for _, role := range allowedRoles {
			if user.Role == role {
				return c.Next()
			}
		}

		return c.Status(fiber.StatusForbidden).JSON(fiber.Map{
			"message": "Access denied",
		})
	}
}
