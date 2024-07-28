package models

import (
	"time"

	"go.mongodb.org/mongo-driver/bson/primitive"
)

// Guest representa un huésped en el sistema
type Guest struct {
	ID               primitive.ObjectID `bson:"_id" json:"id"` // El ID es obligatorio
	Email            string             `bson:"email" json:"email"`
	Phone            string             `bson:"phone" json:"phone"`
	ExtraDescription string             `bson:"extraDescription" json:"extraDescription"`
	Name             string             `bson:"name" json:"name"`
	Height           string             `bson:"height" json:"height"`
	RoomNumber       string             `bson:"roomNumber" json:"roomNumber"`
	Price            float64            `bson:"price" json:"price"`
	Duration         int                `bson:"duration" json:"duration"`
	CreatedAt        time.Time          `bson:"createdAt" json:"createdAt"`
	UpdatedAt        time.Time          `bson:"updatedAt" json:"updatedAt"`
}
