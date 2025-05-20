package routes

import (
	"go-fiber-api/controllers"
	"go-fiber-api/middleware"

	"github.com/gofiber/fiber/v2"
)

func Setup(app *fiber.App) {
	// Auth
	app.Post("/login", controllers.Login)

	// Protected API group
	api := app.Group("/api", middleware.Protected())
	// Upload URL
	api.Put("/presigned_url", controllers.GetUploadUrl)

	// Person routes
	// persons := api.Group("/persons")
	// persons.Post("/", controllers.CreatePerson)
	// persons.Get("/search", controllers.SearchPersons)
	// persons.Get("/family", controllers.GetFamilyInfo)
	// persons.Get("/", controllers.GetPersonByID)
	// persons.Put("/", controllers.UpdatePerson)
	// persons.Delete("/", controllers.DeletePerson)
}
