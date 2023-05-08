package controllers

import (
	"github.com/go-playground/validator/v10"
	"github.com/gofiber/fiber/v2"
	"go.mongodb.org/mongo-driver/mongo"
	"pet-appointments-api/configs"
)

var petCollection *mongo.Collection = configs.GetCollection(configs.DB, "pets")
var validatePet = validator.New()

// Create a new Appointment
func CreatePet(c *fiber.Ctx) error {
	return nil
}

// Get an Appointment
func GetPet(c *fiber.Ctx) error {
	return nil
}

// Edit an Appointment
func EditPet(c *fiber.Ctx) error {
	return nil
}

// Delete an Appointment
func DeletePet(c *fiber.Ctx) error {
	return nil
}

// Get All Appointments
func GetAllPets(c *fiber.Ctx) error {
	return nil
}
