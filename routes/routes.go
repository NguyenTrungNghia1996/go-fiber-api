package routes

import (
	"go-fiber-api/controllers"

	"github.com/gofiber/fiber/v2"
)

func Setup(app *fiber.App) {
	// Auth
	app.Post("/login", controllers.Login)

	// Protected API group
	// api := app.Group("/api", middleware.Protected())
	api := app.Group("/api")
	// Upload URL
	api.Put("/presigned_url", controllers.GetUploadUrl)

	// Person routes
	persons := api.Group("/persons")
	persons.Post("/", controllers.CreatePerson)
	persons.Get("/search", controllers.SearchPersons)
	persons.Get("/:id/family", controllers.GetFamilyInfo)
	persons.Get("/:id", controllers.GetPersonByID)
	persons.Put("/:id", controllers.UpdatePerson)
	persons.Delete("/:id", controllers.DeletePerson)
}
