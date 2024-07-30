package models

import (
	"time"

	"go.mongodb.org/mongo-driver/bson/primitive"
)

type Guest struct {
	ID               primitive.ObjectID `bson:"_id" json:"id"`
	CustomID         string             `bson:"customID" json:"customID"` // ID personalizado
	ExtraDescription string             `bson:"extraDescription" json:"extraDescription"`
	Hair             string             `bson:"hair" json:"hair"`
	Height           string             `bson:"height" json:"height"`
	RoomNumber       string             `bson:"roomNumber" json:"roomNumber"`
	Price            float64            `bson:"price" json:"price"`
	Duration         int                `bson:"duration" json:"duration"`
	History          []HistoryRecord    `bson:"history" json:"history"`
	CreatedAt        time.Time          `bson:"createdAt" json:"createdAt"`
	UpdatedAt        time.Time          `bson:"updatedAt" json:"updatedAt"`
}
