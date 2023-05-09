package controllers

import (
	"context"
	"github.com/go-playground/validator/v10"
	"github.com/gofiber/fiber/v2"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
	"net/http"
	"pet-appointments-api/configs"
	"pet-appointments-api/models"
	"pet-appointments-api/responses"
	"time"
)

var partnerCollection *mongo.Collection = configs.GetCollection(configs.DB, "partners")
var validatePartner = validator.New()

// Create a new Partner
func CreatePartner(c *fiber.Ctx) error {
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	var partner models.Partner
	defer cancel()

	//validate the request body
	if err := c.BodyParser(&partner); err != nil {
		return c.Status(http.StatusBadRequest).JSON(responses.Response{Status: http.StatusBadRequest, Message: "Error: the request body is invalid, please check it again.", Data: &fiber.Map{"data": err.Error()}})
	}

	//use the validator library to validate required fields
	if validationErr := validatePartner.Struct(&partner); validationErr != nil {
		return c.Status(http.StatusBadRequest).JSON(responses.Response{Status: http.StatusBadRequest, Message: "Error: some fields could be invalid.", Data: &fiber.Map{"data": validationErr.Error()}})
	}

	newPartner := models.Partner{
		Id:           primitive.NewObjectID(),
		Name:         partner.Name,
		LastName:     partner.LastName,
		IdNumber:     partner.IdNumber,
		Phone:        partner.Phone,
		Email:        partner.Email,
		CreationDate: time.Now(),
		Services:     partner.Services,
	}

	result, err := partnerCollection.InsertOne(ctx, newPartner)
	if err != nil {
		return c.Status(http.StatusInternalServerError).JSON(responses.Response{Status: http.StatusInternalServerError, Message: "Error: The Partner creation process failed.", Data: &fiber.Map{"data": err.Error()}})
	}

	return c.Status(http.StatusCreated).JSON(responses.Response{Status: http.StatusCreated, Message: "A new Partner was created successfully.", Data: &fiber.Map{"data": result}})
}

// Get a Partner
func GetPartner(c *fiber.Ctx) error {
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	partnerId := c.Params("partnerId")
	var partner models.Partner
	defer cancel()

	objId, _ := primitive.ObjectIDFromHex(partnerId)

	//validate if the partner ID exists
	err := partnerCollection.FindOne(ctx, bson.M{"id": objId}).Decode(&partner)
	if err != nil {
		return c.Status(http.StatusInternalServerError).JSON(responses.Response{Status: http.StatusInternalServerError, Message: "Error: invalid partner ID.", Data: &fiber.Map{"data": err.Error()}})
	}

	return c.Status(http.StatusOK).JSON(responses.Response{Status: http.StatusOK, Message: "The operation was successfully.", Data: &fiber.Map{"data": partner}})
}

// Edit a Partner
func EditPartner(c *fiber.Ctx) error {
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	partnerId := c.Params("partnerId")
	var partner models.Partner
	defer cancel()

	objId, _ := primitive.ObjectIDFromHex(partnerId)

	//validate the request body
	if err := c.BodyParser(&partner); err != nil {
		return c.Status(http.StatusBadRequest).JSON(responses.Response{Status: http.StatusBadRequest, Message: "Error: the request body is invalid, please check it again.", Data: &fiber.Map{"data": err.Error()}})
	}

	//use the validator library to validate required fields
	if validationErr := validatePartner.Struct(&partner); validationErr != nil {
		return c.Status(http.StatusBadRequest).JSON(responses.Response{Status: http.StatusBadRequest, Message: "Error: some fields could be invalid.", Data: &fiber.Map{"data": validationErr.Error()}})
	}

	update := bson.M{"name": partner.Name, "lastName": partner.LastName, "idNumber": partner.IdNumber, "phone": partner.Phone, "email": partner.Email, "services": partner.Services}

	result, err := partnerCollection.UpdateOne(ctx, bson.M{"id": objId}, bson.M{"$set": update})

	if err != nil {
		return c.Status(http.StatusInternalServerError).JSON(responses.Response{Status: http.StatusInternalServerError, Message: "Error: the Partner edit process failed.", Data: &fiber.Map{"data": err.Error()}})
	}

	//get updated partner details
	var updatedPartner models.Partner
	if result.MatchedCount == 1 {
		err := partnerCollection.FindOne(ctx, bson.M{"id": objId}).Decode(&updatedPartner)

		if err != nil {
			return c.Status(http.StatusInternalServerError).JSON(responses.Response{Status: http.StatusInternalServerError, Message: "Error: The Partner edit process failed.", Data: &fiber.Map{"data": err.Error()}})
		}
	}

	return c.Status(http.StatusOK).JSON(responses.Response{Status: http.StatusOK, Message: "The Partner with the ID " + partnerId + " was edited correctly.", Data: &fiber.Map{"data": updatedPartner}})
}

// Delete a Partner
func DeletePartner(c *fiber.Ctx) error {
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	partnerId := c.Params("partnerId")
	defer cancel()

	objId, _ := primitive.ObjectIDFromHex(partnerId)

	result, err := partnerCollection.DeleteOne(ctx, bson.M{"id": objId})

	//validate if the DeleteOne functions returns an Error
	if err != nil {
		return c.Status(http.StatusInternalServerError).JSON(responses.Response{Status: http.StatusInternalServerError, Message: "Error: There is no partner with that ID. ", Data: &fiber.Map{"data": err.Error()}})
	}

	//validate the ID number
	if result.DeletedCount < 1 {
		return c.Status(http.StatusNotFound).JSON(
			responses.Response{Status: http.StatusNotFound, Message: "Error", Data: &fiber.Map{"data": "Error: The Partner with the ID " + partnerId + " does not exists."}},
		)
	}

	return c.Status(http.StatusOK).JSON(
		responses.Response{Status: http.StatusOK, Message: "Success", Data: &fiber.Map{"data": "The Partner was deleted successfully."}},
	)
}

// Get All Partners
func GetAllPartners(c *fiber.Ctx) error {
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	var partners []models.Partner
	defer cancel()

	results, err := partnerCollection.Find(ctx, bson.M{})

	//validate if the context has a collection
	if err != nil {
		return c.Status(http.StatusInternalServerError).JSON(responses.Response{Status: http.StatusInternalServerError, Message: "Error", Data: &fiber.Map{"data": err.Error()}})
	}

	//reading from the db in an optimal way
	defer results.Close(ctx)
	for results.Next(ctx) {
		var singlePartner models.Partner
		if err = results.Decode(&singlePartner); err != nil {
			return c.Status(http.StatusInternalServerError).JSON(responses.Response{Status: http.StatusInternalServerError, Message: "Error", Data: &fiber.Map{"data": err.Error()}})
		}

		partners = append(partners, singlePartner)
	}

	return c.Status(http.StatusOK).JSON(
		responses.Response{Status: http.StatusOK, Message: "Success", Data: &fiber.Map{"data": partners}},
	)
}
