package config

import (
	"app/middleware"
	"app/routes"
	"app/service"

	"github.com/gofiber/fiber/v2"
	"github.com/gofiber/fiber/v2/middleware/cors"
	"github.com/gofiber/fiber/v2/middleware/logger"
	"github.com/gofiber/fiber/v2/middleware/recover"
)

func NewAppConfig() *fiber.App {
	// load dot env
	LoadEnv()

	app := fiber.New()

	app.Use(logger.New())
	app.Use(recover.New())
	app.Use(cors.New())

	v1 := app.Group("/api/v1")

	db := ConnectDB()

	middleware := middleware.NewMiddleware(db)

	categoryService := service.NewCategoryService(db)
	productService := service.NewProductService(db)
	userService := service.NewUserService(db)

	categoryRoutes := routes.NewCategoryRoutes(v1, categoryService, middleware)
	productRoutes := routes.NewProductRoutes(v1, productService, middleware)
	userRoutes := routes.NewUserRoutes(v1, userService, middleware)

	categoryRoutes.CategoryGroup()
	productRoutes.ProductGroup()
	userRoutes.UserGroup()

	return app
}
