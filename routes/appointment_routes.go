package routes

import (
	"github.com/gofiber/fiber/v2"
	"pet-appointments-api/controllers"
)

func AppointmentRoutes(app *fiber.App) {
	app.Post("/appointment", controllers.CreateAppointment)
	app.Get("/appointment/:appointmentId", controllers.GetAppointment)
	app.Put("/appointment/:appointmentId", controllers.EditAppointment)
	app.Delete("/appointment/:appointmentId", controllers.DeleteAppointment)
	app.Get("/appointments", controllers.GetAllAppointments)
}
