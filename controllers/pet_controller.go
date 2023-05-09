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

var petCollection *mongo.Collection = configs.GetCollection(configs.DB, "pets")
var validatePet = validator.New()

// Create a new Pet
func CreatePet(c *fiber.Ctx) error {
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	var pet models.Pet
	defer cancel()

	//validate the request body
	if err := c.BodyParser(&pet); err != nil {
		return c.Status(http.StatusBadRequest).JSON(responses.Response{Status: http.StatusBadRequest, Message: "Error: the request body is invalid, please check it again.", Data: &fiber.Map{"data": err.Error()}})
	}

	//use the validator library to validate required fields
	if validationErr := validatePet.Struct(&pet); validationErr != nil {
		return c.Status(http.StatusBadRequest).JSON(responses.Response{Status: http.StatusBadRequest, Message: "Error: some fields could be invalid.", Data: &fiber.Map{"data": validationErr.Error()}})
	}

	newPet := models.Pet{
		Id:           primitive.NewObjectID(),
		OwnerId:      pet.OwnerId,
		Name:         pet.Name,
		Age:          pet.Age,
		PetType:      pet.PetType,
		Breed:        pet.Breed,
		CreationDate: time.Now(),
	}

	result, err := petCollection.InsertOne(ctx, newPet)
	if err != nil {
		return c.Status(http.StatusInternalServerError).JSON(responses.Response{Status: http.StatusInternalServerError, Message: "Error: The Pet creation process failed.", Data: &fiber.Map{"data": err.Error()}})
	}

	return c.Status(http.StatusCreated).JSON(responses.Response{Status: http.StatusCreated, Message: "A new Pet was created successfully.", Data: &fiber.Map{"data": result}})
}

// Get a Pet
func GetPet(c *fiber.Ctx) error {
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	petId := c.Params("petId")
	var pet models.Pet
	defer cancel()

	objId, _ := primitive.ObjectIDFromHex(petId)

	//validate if the pet ID exists
	err := petCollection.FindOne(ctx, bson.M{"id": objId}).Decode(&pet)
	if err != nil {
		return c.Status(http.StatusInternalServerError).JSON(responses.Response{Status: http.StatusInternalServerError, Message: "Error: invalid pet ID.", Data: &fiber.Map{"data": err.Error()}})
	}

	return c.Status(http.StatusOK).JSON(responses.Response{Status: http.StatusOK, Message: "The operation was successfully.", Data: &fiber.Map{"data": pet}})
}

// Edit a Pet
func EditPet(c *fiber.Ctx) error {
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	petId := c.Params("petId")
	var pet models.Pet
	defer cancel()

	objId, _ := primitive.ObjectIDFromHex(petId)

	//validate the request body
	if err := c.BodyParser(&pet); err != nil {
		return c.Status(http.StatusBadRequest).JSON(responses.Response{Status: http.StatusBadRequest, Message: "Error: the request body is invalid, please check it again.", Data: &fiber.Map{"data": err.Error()}})
	}

	//use the validator library to validate required fields
	if validationErr := validatePet.Struct(&pet); validationErr != nil {
		return c.Status(http.StatusBadRequest).JSON(responses.Response{Status: http.StatusBadRequest, Message: "Error: some fields could be invalid.", Data: &fiber.Map{"data": validationErr.Error()}})
	}

	update := bson.M{"ownerId": pet.OwnerId, "name": pet.Name, "age": pet.Age, "petType": pet.PetType, "breed": pet.Breed}

	result, err := petCollection.UpdateOne(ctx, bson.M{"id": objId}, bson.M{"$set": update})

	if err != nil {
		return c.Status(http.StatusInternalServerError).JSON(responses.Response{Status: http.StatusInternalServerError, Message: "Error: the Pet edit process failed.", Data: &fiber.Map{"data": err.Error()}})
	}

	//get updated owner details
	var updatedPet models.Pet
	if result.MatchedCount == 1 {
		err := petCollection.FindOne(ctx, bson.M{"id": objId}).Decode(&updatedPet)

		if err != nil {
			return c.Status(http.StatusInternalServerError).JSON(responses.Response{Status: http.StatusInternalServerError, Message: "Error: The Pet edit process failed.", Data: &fiber.Map{"data": err.Error()}})
		}
	}

	return c.Status(http.StatusOK).JSON(responses.Response{Status: http.StatusOK, Message: "The Pet with the ID " + petId + " was edited correctly.", Data: &fiber.Map{"data": updatedPet}})
}

// Delete a Pet
func DeletePet(c *fiber.Ctx) error {
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	petId := c.Params("petId")
	defer cancel()

	objId, _ := primitive.ObjectIDFromHex(petId)

	result, err := petCollection.DeleteOne(ctx, bson.M{"id": objId})

	//validate if the DeleteOne functions returns an Error
	if err != nil {
		return c.Status(http.StatusInternalServerError).JSON(responses.Response{Status: http.StatusInternalServerError, Message: "Error: There is no pet with that ID. ", Data: &fiber.Map{"data": err.Error()}})
	}

	//validate the ID number
	if result.DeletedCount < 1 {
		return c.Status(http.StatusNotFound).JSON(
			responses.Response{Status: http.StatusNotFound, Message: "Error", Data: &fiber.Map{"data": "Error: The Pet with the ID " + petId + " does not exists."}},
		)
	}

	return c.Status(http.StatusOK).JSON(
		responses.Response{Status: http.StatusOK, Message: "Success", Data: &fiber.Map{"data": "The Pet was deleted successfully."}},
	)
}

// Get All Pets
func GetAllPets(c *fiber.Ctx) error {
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	var pets []models.Pet
	defer cancel()

	results, err := petCollection.Find(ctx, bson.M{})

	//validate if the context has a collection
	if err != nil {
		return c.Status(http.StatusInternalServerError).JSON(responses.Response{Status: http.StatusInternalServerError, Message: "Error", Data: &fiber.Map{"data": err.Error()}})
	}

	//reading from the db in an optimal way
	defer results.Close(ctx)
	for results.Next(ctx) {
		var singlePet models.Pet
		if err = results.Decode(&singlePet); err != nil {
			return c.Status(http.StatusInternalServerError).JSON(responses.Response{Status: http.StatusInternalServerError, Message: "Error", Data: &fiber.Map{"data": err.Error()}})
		}

		pets = append(pets, singlePet)
	}

	return c.Status(http.StatusOK).JSON(
		responses.Response{Status: http.StatusOK, Message: "Success", Data: &fiber.Map{"data": pets}},
	)
}
