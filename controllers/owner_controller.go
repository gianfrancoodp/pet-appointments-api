package controllers

import (
	"github.com/go-playground/validator/v10"
	"github.com/gofiber/fiber/v2"
	"go.mongodb.org/mongo-driver/mongo"
	"pet-appointments-api/configs"
)

var ownerCollection *mongo.Collection = configs.GetCollection(configs.DB, "owners")
var validateOwner = validator.New()

// Create a new Appointment
func CreateOwner(c *fiber.Ctx) error {
	return nil
}

// Get an Appointment
func GetOwner(c *fiber.Ctx) error {
	return nil
}

// Edit an Appointment
func EditOwner(c *fiber.Ctx) error {
	return nil
}

// Delete an Appointment
func DeleteOwner(c *fiber.Ctx) error {
	return nil
}

// Get All Appointments
func GetAllOwners(c *fiber.Ctx) error {
	return nil
}
