package models

import (
	"go.mongodb.org/mongo-driver/bson/primitive"
	"time"
)

type Owner struct {
	Id           primitive.ObjectID `json:"id,omitempty"`
	Name         string             `json:"name,omitempty" validate:"required"`
	LastName     string             `json:"lastName,omitempty" validate:"required"`
	IdNumber     int                `json:"idNumber,omitempty" validate:"required"`
	Phone        int                `json:"phone,omitempty" validate:"required"`
	Email        string             `json:"email,omitempty" validate:"required"`
	CreationDate time.Time          `json:"creationDate,omitempty" form:"date"`
	Pets         []string           `json:"pets,omitempty"`
}
