package middleware

import (
	"app/model"

	"github.com/gofiber/fiber/v2"
	"gorm.io/gorm"
)

func (m *implMiddleware) GetCredential(c *fiber.Ctx) error {
	user := &model.User{}

	id := c.Locals("user_id")

	if id == nil {
		return c.Status(fiber.StatusUnauthorized).JSON(fiber.Map{
			"message": "Invalid session",
		})
	}

	if err := m.db.Select("id", "username", "email", "phone_number").Where("id = ?", id).First(user).Error; err != nil {
		if err == gorm.ErrRecordNotFound {
			return c.Status(fiber.StatusNotFound).JSON(fiber.Map{
				"message": "Invalid session",
			})
		}
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"message": "Database error",
		})
	}

	c.Locals("user", user)
	return c.Next()
}
