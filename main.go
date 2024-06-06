package main

import (
	"os"

	"app2/config"
	"app2/service"

	"github.com/gofiber/fiber/v2"
	"github.com/gofiber/fiber/v2/middleware/cors"
	"github.com/gofiber/fiber/v2/middleware/logger"
	"github.com/gofiber/fiber/v2/middleware/recover"
)

func main() {
	err := run()

	if err != nil {
		panic(err)
	}
}

func run() error {
	err := config.LoadEnv()

	if err != nil {
		return err
	}

	err = config.ConnectDB()

	if err != nil {
		return err
	}

	defer config.CloseDB()

	app := fiber.New()

	app.Use(logger.New())
	app.Use(recover.New())
	app.Use(cors.New())

	v1 := app.Group("/api/v1")

	service.CategoryGroup(v1)
	service.ProductGroup(v1)
	service.UserGroup(v1)

	var port string
	
	if port = os.Getenv("PORT"); port == "" {
		port = "8080"
	}

	app.Listen(":" + port)

	return nil
}