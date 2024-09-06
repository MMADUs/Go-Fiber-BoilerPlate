package routes

import (
	"app/middleware"
	"app/service"

	"github.com/gofiber/fiber/v2"
)

type ProductRoutes interface {
	ProductGroup()
}

type implProductRoutes struct {
	router     fiber.Router
	service    service.ProductService
	middleware middleware.Middleware
}

func NewProductRoutes(router fiber.Router, service service.ProductService, middleware middleware.Middleware) ProductRoutes {
	return &implProductRoutes{
		router:     router,
		service:    service,
		middleware: middleware,
	}
}

func (r *implProductRoutes) ProductGroup() {
	ProductGroup := r.router.Group("/product")

	ProductGroup.Post("/", r.middleware.Authenticate, r.middleware.GetCredential, r.middleware.Authorize(0, 1), r.service.CreateProduct)
	ProductGroup.Get("/", r.service.GetAllProducts)
	ProductGroup.Get("/page", r.service.PaginatedProduct)
	ProductGroup.Get("/:id", r.service.GetProductById)
	ProductGroup.Put("/:id", r.middleware.Authenticate, r.middleware.GetCredential, r.service.UpdateProduct)
	ProductGroup.Delete("/:id", r.middleware.Authenticate, r.middleware.GetCredential, r.service.DeleteProduct)
}
