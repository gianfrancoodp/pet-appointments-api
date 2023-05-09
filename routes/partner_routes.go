package routes

import (
	"github.com/gofiber/fiber/v2"
	"pet-appointments-api/controllers"
)

func PartnerRoutes(app *fiber.App) {
	app.Post("/partner", controllers.CreatePartner)
	app.Get("/partner/:partnerId", controllers.GetPartner)
	app.Put("/partner/:partnerId", controllers.EditPartner)
	app.Delete("/partner/:partnerId", controllers.DeletePartner)
	app.Get("/partners", controllers.GetAllPartners)
}
