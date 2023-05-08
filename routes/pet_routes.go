package routes

import (
	"github.com/gofiber/fiber/v2"
	"pet-appointments-api/controllers"
)

func PetRoutes(app *fiber.App) {
	app.Post("/pet", controllers.CreatePet)
	app.Get("/pet/:petId", controllers.GetPet)
	app.Put("/pet/:petId", controllers.EditPet)
	app.Delete("/pet/:petId", controllers.DeletePet)
	app.Get("/pets", controllers.GetAllPets)
}
