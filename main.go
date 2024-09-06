package main

import (
	"os"

	"app/config"
)

func main() {
	app := config.NewAppConfig()

	var port string
	
	if port = os.Getenv("PORT"); port == "" {
		port = "8080"
	}

	app.Listen(":" + port)
}