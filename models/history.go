package models

import (
	"time"
)

type HistoryRecord struct {
	Action   string    `bson:"action" json:"action"`     // "checkIn" o "checkOut"
	DateTime time.Time `bson:"dateTime" json:"dateTime"` // Fecha y hora del evento
}
