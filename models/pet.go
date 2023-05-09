package models

import (
	"go.mongodb.org/mongo-driver/bson/primitive"
	"time"
)

type Pet struct {
	Id           primitive.ObjectID `json:"id,omitempty"`
	OwnerId      int                `json:"ownerId,omitempty" validate:"required"`
	Name         string             `json:"name,omitempty" validate:"required"`
	Age          int                `json:"age,omitempty" validate:"required"`
	PetType      string             `json:"petType,omitempty" validate:"required"`
	Breed        string             `json:"breed,omitempty" validate:"required"`
	CreationDate time.Time          `json:"creationDate,omitempty" form:"date"`
}
