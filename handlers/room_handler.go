package handlers

import (
	"context"
	"encoding/json"
	"net/http"
	"time"

	"hotelman-backend/constants"
	"hotelman-backend/models"

	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
)

// RoomHandler maneja las solicitudes relacionadas con las habitaciones
type RoomHandler struct {
	Client *mongo.Client
}

// CreateRoomHandler maneja la creación de nuevas habitaciones
func (h *RoomHandler) CreateRoomHandler(w http.ResponseWriter, r *http.Request) {
	var room models.Room
	err := json.NewDecoder(r.Body).Decode(&room)
	if err != nil {
		http.Error(w, "Error parsing request body", http.StatusBadRequest)
		return
	}

	// Validar el tipo de habitación
	if room.RoomType != "rental" && room.RoomType != "guest" {
		http.Error(w, "Invalid room type. Must be either 'rental' or 'guest'", http.StatusBadRequest)
		return
	}

	room.ID = primitive.NewObjectID()
	room.CreatedAt = time.Now()
	room.UpdatedAt = time.Now()

	collection := h.Client.Database(constants.MongoDBDatabase).Collection(constants.CollectionRooms)
	_, err = collection.InsertOne(context.Background(), room)
	if err != nil {
		http.Error(w, "Failed to create room", http.StatusInternalServerError)
		return
	}

	w.WriteHeader(http.StatusCreated)
	json.NewEncoder(w).Encode(room)
}

// UpdateRoomStatusHandler maneja la actualización del estado de una habitación
func (h *RoomHandler) UpdateRoomStatusHandler(w http.ResponseWriter, r *http.Request) {
	var payload struct {
		OccupantID string `json:"occupantId"`
		Status     string `json:"status"`
	}

	err := json.NewDecoder(r.Body).Decode(&payload)
	if err != nil {
		http.Error(w, "Error parsing request body", http.StatusBadRequest)
		return
	}

	occupantID, err := primitive.ObjectIDFromHex(payload.OccupantID)
	if err != nil {
		http.Error(w, "Invalid occupant ID", http.StatusBadRequest)
		return
	}

	collection := h.Client.Database(constants.MongoDBDatabase).Collection(constants.CollectionRooms)
	filter := bson.M{"occupantId": occupantID}
	update := bson.M{"$set": bson.M{"status": payload.Status, "updatedAt": time.Now()}}

	result, err := collection.UpdateOne(context.Background(), filter, update)
	if err != nil {
		http.Error(w, "Failed to update room status", http.StatusInternalServerError)
		return
	}

	if result.MatchedCount == 0 {
		http.Error(w, "No room found with the given occupant ID", http.StatusNotFound)
		return
	}

	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(bson.M{"message": "Room status updated successfully"})
}

// GetRoomOccupantHandler maneja la obtención del inquilino de una habitación
func (h *RoomHandler) GetRoomOccupantHandler(w http.ResponseWriter, r *http.Request) {
	roomNumber := r.URL.Query().Get("roomNumber")
	if roomNumber == "" {
		http.Error(w, "Room number is required", http.StatusBadRequest)
		return
	}

	collection := h.Client.Database(constants.MongoDBDatabase).Collection(constants.CollectionRooms)
	filter := bson.M{"roomNumber": roomNumber}

	var room models.Room
	err := collection.FindOne(context.Background(), filter).Decode(&room)
	if err != nil {
		http.Error(w, "Failed to get room", http.StatusInternalServerError)
		return
	}

	if room.ID.IsZero() {
		http.Error(w, "Room not found", http.StatusNotFound)
		return
	}

	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(room)
}

// AssignOccupantHandler maneja la asignación de un ocupante a una habitación
func (h *RoomHandler) AssignOccupantHandler(w http.ResponseWriter, r *http.Request) {
	var payload struct {
		RoomNumber string `json:"roomNumber"`
		OccupantID string `json:"occupantId"`
	}

	err := json.NewDecoder(r.Body).Decode(&payload)
	if err != nil {
		http.Error(w, "Error parsing request body", http.StatusBadRequest)
		return
	}

	roomNumber := payload.RoomNumber
	occupantID, err := primitive.ObjectIDFromHex(payload.OccupantID)
	if err != nil {
		http.Error(w, "Invalid occupant ID", http.StatusBadRequest)
		return
	}

	collection := h.Client.Database(constants.MongoDBDatabase).Collection(constants.CollectionRooms)
	filter := bson.M{"roomNumber": roomNumber}
	update := bson.M{"$set": bson.M{"occupantId": occupantID, "updatedAt": time.Now()}}

	result, err := collection.UpdateOne(context.Background(), filter, update)
	if err != nil {
		http.Error(w, "Failed to assign occupant to room", http.StatusInternalServerError)
		return
	}

	if result.MatchedCount == 0 {
		http.Error(w, "No room found with the given room number", http.StatusNotFound)
		return
	}

	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(bson.M{"message": "Occupant assigned to room successfully"})
}
