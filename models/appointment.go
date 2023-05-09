package models

import (
	"go.mongodb.org/mongo-driver/bson/primitive"
	"time"
)

type Appointment struct {
	Id          primitive.ObjectID `json:"id,omitempty"`
	OwnerId     string             `json:"ownerId,omitempty" validate:"required"`
	PetId       string             `json:"petId,omitempty" validate:"required"`
	PartnerId   string             `json:"partnerId,omitempty" validate:"required"`
	Service     string             `json:"service,omitempty" validate:"required"`
	Amount      float64            `json:"amount,omitempty" validate:"required"`
	PaymentType string             `json:"paymentType,omitempty" validate:"required"`
	Date        time.Time          `json:"date,omitempty" form:"date"`
}
