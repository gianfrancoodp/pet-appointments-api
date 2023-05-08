package controllers

import (
	"context"
	"go.mongodb.org/mongo-driver/bson"
	"net/http"
	"pet-appointments-api/configs"
	"pet-appointments-api/models"
	"pet-appointments-api/responses"
	"time"

	"github.com/go-playground/validator/v10"
	"github.com/gofiber/fiber/v2"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
)

var appointmentCollection *mongo.Collection = configs.GetCollection(configs.DB, "appointments")
var validate = validator.New()

// Create a new Appointment
func CreateAppointment(c *fiber.Ctx) error {
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	var appointment models.Appointment
	defer cancel()

	//validate the request body
	if err := c.BodyParser(&appointment); err != nil {
		return c.Status(http.StatusBadRequest).JSON(responses.AppointmentResponse{Status: http.StatusBadRequest, Message: "Error: the request body is invalid, please check it again.", Data: &fiber.Map{"data": err.Error()}})
	}

	//use the validator library to validate required fields
	if validationErr := validate.Struct(&appointment); validationErr != nil {
		return c.Status(http.StatusBadRequest).JSON(responses.AppointmentResponse{Status: http.StatusBadRequest, Message: "Error: some fields could be invalid.", Data: &fiber.Map{"data": validationErr.Error()}})
	}

	newAppointment := models.Appointment{
		Id:          primitive.NewObjectID(),
		OwnerId:     appointment.OwnerId,
		PetId:       appointment.PetId,
		PartnerId:   appointment.PartnerId,
		Service:     appointment.Service,
		Amount:      appointment.Amount,
		PaymentType: appointment.PaymentType,
		Date:        time.Now(),
	}

	result, err := appointmentCollection.InsertOne(ctx, newAppointment)
	if err != nil {
		return c.Status(http.StatusInternalServerError).JSON(responses.AppointmentResponse{Status: http.StatusInternalServerError, Message: "Error: the Appointment creation process failed.", Data: &fiber.Map{"data": err.Error()}})
	}

	return c.Status(http.StatusCreated).JSON(responses.AppointmentResponse{Status: http.StatusCreated, Message: "A new Appointment was created successfully.", Data: &fiber.Map{"data": result}})
}

// Get an Appointment
func GetAppointment(c *fiber.Ctx) error {
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	appointmentId := c.Params("appointmentId")
	var appointment models.Appointment
	defer cancel()

	objId, _ := primitive.ObjectIDFromHex(appointmentId)

	//validate if the appointmentId ID exists
	err := appointmentCollection.FindOne(ctx, bson.M{"id": objId}).Decode(&appointment)
	if err != nil {
		return c.Status(http.StatusInternalServerError).JSON(responses.AppointmentResponse{Status: http.StatusInternalServerError, Message: "Error: invalid appointment ID.", Data: &fiber.Map{"data": err.Error()}})
	}

	return c.Status(http.StatusOK).JSON(responses.AppointmentResponse{Status: http.StatusOK, Message: "The operation was successfully.", Data: &fiber.Map{"data": appointment}})
}

// Edit an Appointment
func EditAppointment(c *fiber.Ctx) error {
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	appointmentId := c.Params("appointmentId")
	var appointment models.Appointment
	defer cancel()

	objId, _ := primitive.ObjectIDFromHex(appointmentId)

	//validate the request body
	if err := c.BodyParser(&appointment); err != nil {
		return c.Status(http.StatusBadRequest).JSON(responses.AppointmentResponse{Status: http.StatusBadRequest, Message: "Error: the request body is invalid, please check it again.", Data: &fiber.Map{"data": err.Error()}})
	}

	//use the validator library to validate required fields
	if validationErr := validate.Struct(&appointment); validationErr != nil {
		return c.Status(http.StatusBadRequest).JSON(responses.AppointmentResponse{Status: http.StatusBadRequest, Message: "Error: some fields could be invalid.", Data: &fiber.Map{"data": validationErr.Error()}})
	}

	update := bson.M{"ownerId": appointment.OwnerId, "petId": appointment.PetId, "partnerId": appointment.PartnerId, "service": appointment.Service, "amount": appointment.Amount, "paymentType": appointment.PaymentType}

	result, err := appointmentCollection.UpdateOne(ctx, bson.M{"id": objId}, bson.M{"$set": update})

	if err != nil {
		return c.Status(http.StatusInternalServerError).JSON(responses.AppointmentResponse{Status: http.StatusInternalServerError, Message: "Error: the Appointment edit process failed.", Data: &fiber.Map{"data": err.Error()}})
	}

	//get updated appointment details
	var updatedAppointment models.Appointment
	if result.MatchedCount == 1 {
		err := appointmentCollection.FindOne(ctx, bson.M{"id": objId}).Decode(&updatedAppointment)

		if err != nil {
			return c.Status(http.StatusInternalServerError).JSON(responses.AppointmentResponse{Status: http.StatusInternalServerError, Message: "Error: the Appointment edit process failed.", Data: &fiber.Map{"data": err.Error()}})
		}
	}

	return c.Status(http.StatusOK).JSON(responses.AppointmentResponse{Status: http.StatusOK, Message: "The Appointment with the ID " + appointmentId + " was edited correctly.", Data: &fiber.Map{"data": updatedAppointment}})
}

// Delete an Appointment
func DeleteAppointment(c *fiber.Ctx) error {
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	appointmentId := c.Params("appointmentId")
	defer cancel()

	objId, _ := primitive.ObjectIDFromHex(appointmentId)

	result, err := appointmentCollection.DeleteOne(ctx, bson.M{"id": objId})

	//validate if the DeleteOne functions returns an Error
	if err != nil {
		return c.Status(http.StatusInternalServerError).JSON(responses.AppointmentResponse{Status: http.StatusInternalServerError, Message: "Error: There is no appointment with that ID. ", Data: &fiber.Map{"data": err.Error()}})
	}

	//validate the ID number
	if result.DeletedCount < 1 {
		return c.Status(http.StatusNotFound).JSON(
			responses.AppointmentResponse{Status: http.StatusNotFound, Message: "Error", Data: &fiber.Map{"data": "Error: The appointment with the ID " + appointmentId + " does not exists."}},
		)
	}

	return c.Status(http.StatusOK).JSON(
		responses.AppointmentResponse{Status: http.StatusOK, Message: "Success", Data: &fiber.Map{"data": "The appointment was deleted successfully."}},
	)
}

// Get All Appointments
func GetAllAppointments(c *fiber.Ctx) error {
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	var appointments []models.Appointment
	defer cancel()

	results, err := appointmentCollection.Find(ctx, bson.M{})

	//validate if the context has a collection
	if err != nil {
		return c.Status(http.StatusInternalServerError).JSON(responses.AppointmentResponse{Status: http.StatusInternalServerError, Message: "Error", Data: &fiber.Map{"data": err.Error()}})
	}

	//reading from the db in an optimal way
	defer results.Close(ctx)
	for results.Next(ctx) {
		var singleAppointment models.Appointment
		if err = results.Decode(&singleAppointment); err != nil {
			return c.Status(http.StatusInternalServerError).JSON(responses.AppointmentResponse{Status: http.StatusInternalServerError, Message: "Error", Data: &fiber.Map{"data": err.Error()}})
		}

		appointments = append(appointments, singleAppointment)
	}

	return c.Status(http.StatusOK).JSON(
		responses.AppointmentResponse{Status: http.StatusOK, Message: "Success", Data: &fiber.Map{"data": appointments}},
	)
}
