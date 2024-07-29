package models

import (
	"time"

	"go.mongodb.org/mongo-driver/bson/primitive"
)

// Rental representa un inquilino en el sistema
type Rental struct {
	ID            primitive.ObjectID `bson:"_id,omitempty" json:"id,omitempty"`
	Nombres       string             `bson:"nombres" json:"nombres"`
	Apellidos     string             `bson:"apellidos" json:"apellidos"`
	Correo        string             `bson:"correo" json:"correo"`
	NumeroCelular string             `bson:"numeroCelular" json:"numeroCelular"`
	INEString     string             `bson:"INEString" json:"INEString"`
	RoomNumber    string             `bson:"RoomNumber" json:"RoomNumber"`
	ContratoURL   string             `bson:"contratoUrl" json:"contratoUrl"` // URL del contrato
	INEURL        string             `bson:"ineUrl" json:"ineUrl"`           // URL del INE
	History       []HistoryRecord    `bson:"history" json:"history"`
	CreatedAt     time.Time          `bson:"createdAt" json:"createdAt"`
	UpdatedAt     time.Time          `bson:"updatedAt" json:"updatedAt"`
}
