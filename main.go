package main

import (
	"go-fiber-api/config"
	"go-fiber-api/routes"
	"go-fiber-api/seed"
	"log"
	"os"

	"github.com/gofiber/fiber/v2"
	"github.com/gofiber/fiber/v2/middleware/cors"

	// "github.com/gofiber/fiber/v2/middleware/logger"
	"github.com/joho/godotenv"
)

func main() {
	err := godotenv.Load()
	if err != nil {
		log.Fatal("Error loading .env file")
	}

	config.ConnectDB()
	seed.SeedAdminUser()

	app := fiber.New()
	app.Use(cors.New())
	routes.Setup(app)

	// app.Use(logger.New())

	// routes.AuthRoutes(app)

	port := os.Getenv("PORT")
	log.Fatal(app.Listen(":" + port))
}
