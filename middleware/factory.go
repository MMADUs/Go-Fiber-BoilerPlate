package middleware

import (
	"github.com/gofiber/fiber/v2"
	"gorm.io/gorm"
)

type Middleware interface {
	Authenticate(c *fiber.Ctx) error
	Authorize(allowedRoles ...int) func(*fiber.Ctx) error
	GetCredential(c *fiber.Ctx) error
}

type implMiddleware struct {
	db *gorm.DB
}

func NewMiddleware(db *gorm.DB) Middleware {
	return &implMiddleware{
		db: db,
	}
}
