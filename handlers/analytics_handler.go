package handlers

import (
	"context"
	"encoding/json"
	"net/http"
	"time"

	"hotelman-backend/constants"

	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
)

// AnalyticsHandler maneja las solicitudes de análisis
type AnalyticsHandler struct {
	Client *mongo.Client
}

// AnalyticsRequest define los parámetros de la solicitud de análisis
type AnalyticsRequest struct {
	StartDate string `json:"startDate"`
	EndDate   string `json:"endDate"`
}

// AnalyticsResponse define la estructura de la respuesta de análisis
type AnalyticsResponse struct {
	TotalPriceGuest float64 `json:"totalPriceGuest"`
	Guest           struct {
		Total    int `json:"total"`
		Occupied int `json:"occupied"`
	} `json:"guest"`
	Rental struct {
		Total    int `json:"total"`
		Occupied int `json:"occupied"`
	} `json:"rental"`
	TotalClients int `json:"totalClients"`
}

// GetAnalyticsHandler maneja la obtención de datos de análisis
func (h *AnalyticsHandler) GetAnalyticsHandler(w http.ResponseWriter, r *http.Request) {
	var req AnalyticsRequest
	err := json.NewDecoder(r.Body).Decode(&req)
	if err != nil {
		http.Error(w, "Error parsing request body", http.StatusBadRequest)
		return
	}

	// Configuración de fechas por defecto (mensuales)
	now := time.Now().UTC()
	var startDate, endDate time.Time

	if req.StartDate != "" && req.EndDate != "" {
		// Convertir las fechas proporcionadas
		startDate, err = time.Parse(time.RFC3339, req.StartDate)
		if err != nil {
			http.Error(w, "Invalid startDate format", http.StatusBadRequest)
			return
		}
		endDate, err = time.Parse(time.RFC3339, req.EndDate)
		if err != nil {
			http.Error(w, "Invalid endDate format", http.StatusBadRequest)
			return
		}
	} else {
		// Establecer el rango mensual por defecto
		startDate = time.Date(now.Year(), now.Month(), 1, 0, 0, 0, 0, time.UTC)
		endDate = time.Date(now.Year(), now.Month()+1, 1, 0, 0, 0, 0, time.UTC).Add(-time.Nanosecond)
	}

	// Conectar a la colección de clientes para calcular el total de precios
	clientCollection := h.Client.Database(constants.MongoDBDatabase).Collection("clients")

	// Obtener el total de precios de los clientes según el período solicitado
	priceSumPipeline := mongo.Pipeline{
		{{"$match", bson.D{
			{"createdAt", bson.D{{"$gte", startDate}, {"$lte", endDate}}},
			{"price", bson.D{{"$exists", true}, {"$type", "double"}}}, // Filtra solo documentos que contienen el campo "price"
		}}},
		{{"$group", bson.D{{"_id", nil}, {"totalPrice", bson.D{{"$sum", "$price"}}}}}},
	}

	cursor, err := clientCollection.Aggregate(context.Background(), priceSumPipeline)
	if err != nil {
		http.Error(w, "Failed to calculate total guest price", http.StatusInternalServerError)
		return
	}

	var priceSumResult struct {
		TotalPrice float64 `bson:"totalPrice"`
	}

	if cursor.Next(context.Background()) {
		if err := cursor.Decode(&priceSumResult); err != nil {
			http.Error(w, "Failed to decode total guest price", http.StatusInternalServerError)
			return
		}
	}

	// Obtener el total de habitaciones de huéspedes
	guestTotalPipeline := mongo.Pipeline{
		{{"$match", bson.D{{"roomType", "guest"}}}},
		{{"$count", "total"}},
	}

	cursor, err = h.Client.Database(constants.MongoDBDatabase).Collection("rooms").Aggregate(context.Background(), guestTotalPipeline)
	if err != nil {
		http.Error(w, "Failed to calculate total guest rooms", http.StatusInternalServerError)
		return
	}

	var guestTotalCountResult struct {
		Total int `bson:"total"`
	}

	if cursor.Next(context.Background()) {
		if err := cursor.Decode(&guestTotalCountResult); err != nil {
			http.Error(w, "Failed to decode total guest rooms", http.StatusInternalServerError)
			return
		}
	}

	// Obtener el total de habitaciones de huéspedes ocupadas
	guestOccupiedPipeline := mongo.Pipeline{
		{{"$match", bson.D{{"roomType", "guest"}, {"status", "occupied"}}}},
		{{"$count", "occupied"}},
	}

	cursor, err = h.Client.Database(constants.MongoDBDatabase).Collection("rooms").Aggregate(context.Background(), guestOccupiedPipeline)
	if err != nil {
		http.Error(w, "Failed to calculate occupied guest rooms", http.StatusInternalServerError)
		return
	}

	var guestOccupiedCountResult struct {
		Occupied int `bson:"occupied"`
	}

	if cursor.Next(context.Background()) {
		if err := cursor.Decode(&guestOccupiedCountResult); err != nil {
			http.Error(w, "Failed to decode occupied guest rooms", http.StatusInternalServerError)
			return
		}
	}

	// Obtener el total de habitaciones de renta
	rentalTotalPipeline := mongo.Pipeline{
		{{"$match", bson.D{{"roomType", "rental"}}}},
		{{"$count", "total"}},
	}

	cursor, err = h.Client.Database(constants.MongoDBDatabase).Collection("rooms").Aggregate(context.Background(), rentalTotalPipeline)
	if err != nil {
		http.Error(w, "Failed to calculate total rental rooms", http.StatusInternalServerError)
		return
	}

	var rentalTotalCountResult struct {
		Total int `bson:"total"`
	}

	if cursor.Next(context.Background()) {
		if err := cursor.Decode(&rentalTotalCountResult); err != nil {
			http.Error(w, "Failed to decode total rental rooms", http.StatusInternalServerError)
			return
		}
	}

	// Obtener el total de habitaciones de renta ocupadas
	rentalOccupiedPipeline := mongo.Pipeline{
		{{"$match", bson.D{{"roomType", "rental"}, {"status", "occupied"}}}},
		{{"$count", "occupied"}},
	}

	cursor, err = h.Client.Database(constants.MongoDBDatabase).Collection("rooms").Aggregate(context.Background(), rentalOccupiedPipeline)
	if err != nil {
		http.Error(w, "Failed to calculate occupied rental rooms", http.StatusInternalServerError)
		return
	}

	var rentalOccupiedCountResult struct {
		Occupied int `bson:"occupied"`
	}

	if cursor.Next(context.Background()) {
		if err := cursor.Decode(&rentalOccupiedCountResult); err != nil {
			http.Error(w, "Failed to decode occupied rental rooms", http.StatusInternalServerError)
			return
		}
	}

	// Obtener el total de clientes
	clientTotalCount, err := clientCollection.CountDocuments(context.Background(), bson.M{})
	if err != nil {
		http.Error(w, "Failed to count total clients", http.StatusInternalServerError)
		return
	}

	// Construir la respuesta
	response := AnalyticsResponse{
		TotalPriceGuest: priceSumResult.TotalPrice,
		Guest: struct {
			Total    int `json:"total"`
			Occupied int `json:"occupied"`
		}{
			Total:    guestTotalCountResult.Total,
			Occupied: guestOccupiedCountResult.Occupied,
		},
		Rental: struct {
			Total    int `json:"total"`
			Occupied int `json:"occupied"`
		}{
			Total:    rentalTotalCountResult.Total,
			Occupied: rentalOccupiedCountResult.Occupied,
		},
		TotalClients: int(clientTotalCount),
	}

	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(response)
}
