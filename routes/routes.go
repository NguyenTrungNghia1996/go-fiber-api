package routes

import (
	"go-fiber-api/controllers"
	"go-fiber-api/middleware"

	"github.com/gofiber/fiber/v2"
)

// Setup thiết lập tất cả các route của ứng dụng
func Setup(app *fiber.App) {
	// Auth
	app.Post("/login", controllers.Login)
	// app.Get("/test", controllers.Hello)

	// Group các route có bảo vệ bằng JWT
	api := app.Group("/api", middleware.Protected())

	api.Get("/protect", controllers.Hello)
	// Person (gia phả)
}
