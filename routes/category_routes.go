package routes

import (
	"app/middleware"
	"app/service"

	"github.com/gofiber/fiber/v2"
)

type CategoryRoutes interface {
	CategoryGroup()
}

type implCategoryRoutes struct {
	router     fiber.Router
	service    service.CategoryService
	middleware middleware.Middleware
}

func NewCategoryRoutes(router fiber.Router, service service.CategoryService, middleware middleware.Middleware) CategoryRoutes {
	return &implCategoryRoutes{
		router:     router,
		service:    service,
		middleware: middleware,
	}
}

func (r *implCategoryRoutes) CategoryGroup() {
	categoryRoutes := r.router.Group("/category")

	categoryRoutes.Post("/", r.middleware.Authenticate, r.middleware.GetCredential, r.service.CreateCategory)
	categoryRoutes.Get("/", r.service.GetAllCategory)
	categoryRoutes.Get("/:id", r.service.GetCategoryById)
	categoryRoutes.Put("/:id", r.middleware.Authenticate, r.middleware.GetCredential, r.service.UpdateCategory)
	categoryRoutes.Delete("/:id", r.middleware.Authenticate, r.middleware.GetCredential, r.service.DeleteCategory)
}
