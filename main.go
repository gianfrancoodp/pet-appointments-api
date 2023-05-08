package main

import (
	"github.com/gofiber/fiber/v2"
	"pet-appointments-api/configs"
	"pet-appointments-api/routes"
)

func main() {
	app := fiber.New()

	//run database
	configs.ConnectDB()

	//routes
	routes.AppointmentRoutes(app)
	routes.OwnerRoutes(app)

	app.Listen(":6000")
}
