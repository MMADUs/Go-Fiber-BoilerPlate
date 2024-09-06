package routes

import (
	"app/middleware"
	"app/service"

	"github.com/gofiber/fiber/v2"
)

type UserRoutes interface {
	UserGroup()
}

type implUserRoutes struct {
	router     fiber.Router
	service    service.UserService
	middleware middleware.Middleware
}

func NewUserRoutes(router fiber.Router, service service.UserService, middleware middleware.Middleware) UserRoutes {
	return &implUserRoutes{
		router:     router,
		service:    service,
		middleware: middleware,
	}
}

func (r *implUserRoutes) UserGroup() {
	UserGroup := r.router.Group("/user")

	UserGroup.Post("/register", r.service.Register)
	UserGroup.Post("/Login", r.service.Login)
	UserGroup.Put("/update-password", r.middleware.Authenticate, r.service.UpdatePassword)
}
