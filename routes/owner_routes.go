package routes

import (
	"github.com/gofiber/fiber/v2"
	"pet-appointments-api/controllers"
)

func OwnerRoutes(app *fiber.App) {
	app.Post("/appointment", controllers.CreateOwner)
	app.Get("/appointment/:appointmentId", controllers.GetOwner)
	app.Put("/appointment/:appointmentId", controllers.EditOwner)
	app.Delete("/appointment/:appointmentId", controllers.DeleteOwner)
	app.Get("/appointments", controllers.GetAllOwners)
}
