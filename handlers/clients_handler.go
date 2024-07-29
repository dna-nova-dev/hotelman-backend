// src/handlers/GetClientsHandler.go

package handlers

import (
	"context"
	"encoding/json"
	"net/http"
	"strconv"

	"hotelman-backend/constants"

	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

type GetClientsHandler struct {
	Client *mongo.Client
}

func (h *GetClientsHandler) Handle(w http.ResponseWriter, r *http.Request) {
	clientType := r.URL.Query().Get("type")
	search := r.URL.Query().Get("search")
	pageStr := r.URL.Query().Get("page")
	pageSizeStr := r.URL.Query().Get("pageSize")

	page, err := strconv.Atoi(pageStr)
	if err != nil || page <= 0 {
		page = 1 // Default page
	}

	pageSize, err := strconv.Atoi(pageSizeStr)
	if err != nil || pageSize <= 0 {
		pageSize = 10 // Default page size
	}

	skip := (page - 1) * pageSize
	limit := int64(pageSize)

	filter := bson.M{}
	if search != "" {
		switch clientType {
		case "rental":
			filter = bson.M{"nombres": bson.M{"$regex": search, "$options": "i"}} // Case insensitive search for rentals
		case "guest":
			filter = bson.M{"customID": bson.M{"$regex": search, "$options": "i"}} // Case insensitive search for guests
		default:
			http.Error(w, "Invalid client type", http.StatusBadRequest)
			return
		}
	}

	collection := h.Client.Database(constants.MongoDBDatabase).Collection(constants.CollectionClients)
	opts := options.Find().SetSkip(int64(skip)).SetLimit(limit)
	cursor, err := collection.Find(context.Background(), filter, opts)
	if err != nil {
		http.Error(w, "Failed to retrieve clients", http.StatusInternalServerError)
		return
	}
	defer cursor.Close(context.Background())

	var clients []bson.M
	if err := cursor.All(context.Background(), &clients); err != nil {
		http.Error(w, "Failed to decode clients", http.StatusInternalServerError)
		return
	}

	// Obtener el total de documentos para la paginación
	totalDocs, err := collection.CountDocuments(context.Background(), filter)
	if err != nil {
		http.Error(w, "Failed to count documents", http.StatusInternalServerError)
		return
	}

	totalPages := (totalDocs + int64(pageSize) - 1) / int64(pageSize) // Calcular el número total de páginas

	response := map[string]interface{}{
		"clients":    clients,
		"totalPages": totalPages,
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(response)
}
