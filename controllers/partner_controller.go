package controllers

import (
	"github.com/go-playground/validator/v10"
	"github.com/gofiber/fiber/v2"
	"go.mongodb.org/mongo-driver/mongo"
	"pet-appointments-api/configs"
)

var partnerCollection *mongo.Collection = configs.GetCollection(configs.DB, "partners")
var validatePartner = validator.New()

// Create a new Partner
func CreatePartner(c *fiber.Ctx) error {
	return nil
}

// Get a Partner
func GetPartner(c *fiber.Ctx) error {
	return nil
}

// Edit a Partner
func EditPartner(c *fiber.Ctx) error {
	return nil
}

// Delete a Partner
func DeletePartner(c *fiber.Ctx) error {
	return nil
}

// Get All Partners
func GetAllPartners(c *fiber.Ctx) error {
	return nil
}
