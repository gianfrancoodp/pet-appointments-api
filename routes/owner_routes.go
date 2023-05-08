package routes

import (
	"github.com/gofiber/fiber/v2"
	"pet-appointments-api/controllers"
)

func OwnerRoutes(app *fiber.App) {
	app.Post("/owner", controllers.CreateOwner)
	app.Get("/owner/:ownerId", controllers.GetOwner)
	app.Put("/owner/:ownerId", controllers.EditOwner)
	app.Delete("/owner/:ownerId", controllers.DeleteOwner)
	app.Get("/owners", controllers.GetAllOwners)
}
