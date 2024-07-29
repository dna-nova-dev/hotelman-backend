package models

import (
	"time"

	"go.mongodb.org/mongo-driver/bson/primitive"
)

// Room representa una habitación en el sistema
type Room struct {
	ID          primitive.ObjectID `bson:"_id,omitempty" json:"id,omitempty"`
	RoomType    string             `bson:"roomType" json:"roomType"` // Nuevo campo para el tipo de habitación
	RoomNumber  string             `bson:"roomNumber" json:"roomNumber"`
	Name        string             `bson:"name" json:"name"`
	Description string             `bson:"description" json:"description"`
	Status      string             `bson:"status" json:"status"`

	OccupantID *primitive.ObjectID `bson:"occupantId,omitempty" json:"occupantId,omitempty"`
	CreatedAt  time.Time           `bson:"createdAt" json:"createdAt"`
	UpdatedAt  time.Time           `bson:"updatedAt" json:"updatedAt"`
}
