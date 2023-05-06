package main

import (
	"github.com/gofiber/fiber/v2"
	"pet-appointments-api/configs"
	"pet-appointments-api/routes"
)

func main() {
	app := fiber.New()

	//run database connection
	configs.ConnectDB()

	//REST API routes
	routes.AppointmentRoute(app)

	//Port defined
	app.Listen(":6000")
}
