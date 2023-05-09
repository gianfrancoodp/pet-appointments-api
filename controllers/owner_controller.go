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

var ownerCollection *mongo.Collection = configs.GetCollection(configs.DB, "owners")
var validateOwner = validator.New()

// Create a new Appointment
func CreateOwner(c *fiber.Ctx) error {
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	var owner models.Owner
	defer cancel()

	//validate the request body
	if err := c.BodyParser(&owner); err != nil {
		return c.Status(http.StatusBadRequest).JSON(responses.OwnerResponse{Status: http.StatusBadRequest, Message: "Error: the request body is invalid, please check it again.", Data: &fiber.Map{"data": err.Error()}})
	}

	//use the validator library to validate required fields
	if validationErr := validateOwner.Struct(&owner); validationErr != nil {
		return c.Status(http.StatusBadRequest).JSON(responses.OwnerResponse{Status: http.StatusBadRequest, Message: "Error: some fields could be invalid.", Data: &fiber.Map{"data": validationErr.Error()}})
	}

	newOwner := models.Owner{
		Id:           primitive.NewObjectID(),
		Name:         owner.Name,
		LastName:     owner.LastName,
		IdNumber:     owner.IdNumber,
		Phone:        owner.Phone,
		Email:        owner.Email,
		CreationDate: time.Now(),
		Pets:         owner.Pets,
	}

	result, err := ownerCollection.InsertOne(ctx, newOwner)
	if err != nil {
		return c.Status(http.StatusInternalServerError).JSON(responses.OwnerResponse{Status: http.StatusInternalServerError, Message: "Error: The Owner creation process failed.", Data: &fiber.Map{"data": err.Error()}})
	}

	return c.Status(http.StatusCreated).JSON(responses.OwnerResponse{Status: http.StatusCreated, Message: "A new Owner was created successfully.", Data: &fiber.Map{"data": result}})
}

// Get an Appointment
func GetOwner(c *fiber.Ctx) error {
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	ownerId := c.Params("ownerId")
	var owner models.Owner
	defer cancel()

	objId, _ := primitive.ObjectIDFromHex(ownerId)

	//validate if the owner ID exists
	err := ownerCollection.FindOne(ctx, bson.M{"id": objId}).Decode(&owner)
	if err != nil {
		return c.Status(http.StatusInternalServerError).JSON(responses.OwnerResponse{Status: http.StatusInternalServerError, Message: "Error: invalid owner ID.", Data: &fiber.Map{"data": err.Error()}})
	}

	return c.Status(http.StatusOK).JSON(responses.OwnerResponse{Status: http.StatusOK, Message: "The operation was successfully.", Data: &fiber.Map{"data": owner}})
}

// Edit an Appointment
func EditOwner(c *fiber.Ctx) error {
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	ownerId := c.Params("ownerId")
	var owner models.Owner
	defer cancel()

	objId, _ := primitive.ObjectIDFromHex(ownerId)

	//validate the request body
	if err := c.BodyParser(&owner); err != nil {
		return c.Status(http.StatusBadRequest).JSON(responses.OwnerResponse{Status: http.StatusBadRequest, Message: "Error: the request body is invalid, please check it again.", Data: &fiber.Map{"data": err.Error()}})
	}

	//use the validator library to validate required fields
	if validationErr := validateOwner.Struct(&owner); validationErr != nil {
		return c.Status(http.StatusBadRequest).JSON(responses.OwnerResponse{Status: http.StatusBadRequest, Message: "Error: some fields could be invalid.", Data: &fiber.Map{"data": validationErr.Error()}})
	}

	update := bson.M{"name": owner.Name, "lastName": owner.LastName, "idNumber": owner.IdNumber, "phone": owner.Phone, "email": owner.Email, "pets": owner.Pets}

	result, err := ownerCollection.UpdateOne(ctx, bson.M{"id": objId}, bson.M{"$set": update})

	if err != nil {
		return c.Status(http.StatusInternalServerError).JSON(responses.OwnerResponse{Status: http.StatusInternalServerError, Message: "Error: The Owner edit process failed.", Data: &fiber.Map{"data": err.Error()}})
	}

	//get updated owner details
	var updatedOwner models.Owner
	if result.MatchedCount == 1 {
		err := ownerCollection.FindOne(ctx, bson.M{"id": objId}).Decode(&updatedOwner)

		if err != nil {
			return c.Status(http.StatusInternalServerError).JSON(responses.OwnerResponse{Status: http.StatusInternalServerError, Message: "Error: The Owner edit process failed.", Data: &fiber.Map{"data": err.Error()}})
		}
	}

	return c.Status(http.StatusOK).JSON(responses.OwnerResponse{Status: http.StatusOK, Message: "The Owner with the ID " + ownerId + " was edited correctly.", Data: &fiber.Map{"data": updatedOwner}})
}

// Delete an Appointment
func DeleteOwner(c *fiber.Ctx) error {
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	ownerId := c.Params("ownerId")
	defer cancel()

	objId, _ := primitive.ObjectIDFromHex(ownerId)

	result, err := ownerCollection.DeleteOne(ctx, bson.M{"id": objId})

	//validate if the DeleteOne functions returns an Error
	if err != nil {
		return c.Status(http.StatusInternalServerError).JSON(responses.OwnerResponse{Status: http.StatusInternalServerError, Message: "Error: There is no owner with that ID. ", Data: &fiber.Map{"data": err.Error()}})
	}

	//validate the ID number
	if result.DeletedCount < 1 {
		return c.Status(http.StatusNotFound).JSON(
			responses.OwnerResponse{Status: http.StatusNotFound, Message: "Error", Data: &fiber.Map{"data": "Error: The Owner with the ID " + ownerId + " does not exists."}},
		)
	}

	return c.Status(http.StatusOK).JSON(
		responses.OwnerResponse{Status: http.StatusOK, Message: "Success", Data: &fiber.Map{"data": "The Owner was deleted successfully."}},
	)
}

// Get All Appointments
func GetAllOwners(c *fiber.Ctx) error {
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	var owners []models.Owner
	defer cancel()

	results, err := ownerCollection.Find(ctx, bson.M{})

	//validate if the context has a collection
	if err != nil {
		return c.Status(http.StatusInternalServerError).JSON(responses.OwnerResponse{Status: http.StatusInternalServerError, Message: "Error", Data: &fiber.Map{"data": err.Error()}})
	}

	//reading from the db in an optimal way
	defer results.Close(ctx)
	for results.Next(ctx) {
		var singleOwner models.Owner
		if err = results.Decode(&singleOwner); err != nil {
			return c.Status(http.StatusInternalServerError).JSON(responses.OwnerResponse{Status: http.StatusInternalServerError, Message: "Error", Data: &fiber.Map{"data": err.Error()}})
		}

		owners = append(owners, singleOwner)
	}

	return c.Status(http.StatusOK).JSON(
		responses.OwnerResponse{Status: http.StatusOK, Message: "Success", Data: &fiber.Map{"data": owners}},
	)
}
