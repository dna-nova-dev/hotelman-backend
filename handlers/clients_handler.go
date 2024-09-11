package handlers

import (
	"context"
	"encoding/json"
	"net/http"
	"strconv"

	"hotelman-backend/constants"

	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

type GetClientsHandler struct {
	Client *mongo.Client
}

func (h *GetClientsHandler) Handle(w http.ResponseWriter, r *http.Request) {
	switch r.Method {
	case http.MethodGet:
		h.getClients(w, r)
	case http.MethodPut:
		h.Update(w, r)
	default:
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
	}
}

func (h *GetClientsHandler) getClients(w http.ResponseWriter, r *http.Request) {
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

	var filter bson.M
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
	} else {
		switch clientType {
		case "rental":
			filter = bson.M{"nombres": bson.M{"$exists": true}} // Ensures we are looking at rentals
		case "guest":
			filter = bson.M{"customID": bson.M{"$exists": true}} // Ensures we are looking at guests
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
		"type":       clientType, // Agregar el tipo de cliente a la respuesta
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(response)
}

func (h *GetClientsHandler) Update(w http.ResponseWriter, r *http.Request) {
	// Parsear el ID del cliente desde los parámetros de la URL
	clientID := r.URL.Query().Get("id")
	if clientID == "" {
		http.Error(w, "Missing client ID", http.StatusBadRequest)
		return
	}

	// Convertir el ID a un ObjectID de MongoDB
	objectID, err := primitive.ObjectIDFromHex(clientID)
	if err != nil {
		http.Error(w, "Invalid client ID", http.StatusBadRequest)
		return
	}

	// Estructura para capturar los datos del cuerpo de la solicitud
	var updateData map[string]interface{}

	// Decodificar el cuerpo de la solicitud en la estructura
	if err := json.NewDecoder(r.Body).Decode(&updateData); err != nil {
		http.Error(w, "Invalid request body", http.StatusBadRequest)
		return
	}

	// Si no se proporciona ningún campo para actualizar, retornar un error
	if len(updateData) == 0 {
		http.Error(w, "No fields to update", http.StatusBadRequest)
		return
	}

	// Seleccionar la colección de clientes
	collection := h.Client.Database(constants.MongoDBDatabase).Collection(constants.CollectionClients)

	// Realizar la actualización
	filter := bson.M{"_id": objectID}
	result, err := collection.UpdateOne(context.Background(), filter, bson.M{"$set": updateData})
	if err != nil {
		http.Error(w, "Failed to update client", http.StatusInternalServerError)
		return
	}

	// Verificar si se actualizó algún documento
	if result.MatchedCount == 0 {
		http.Error(w, "Client not found", http.StatusNotFound)
		return
	}

	// Responder con un mensaje de éxito
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(map[string]string{"message": "Client updated successfully"})
}

// Implementación del algoritmo de búsqueda optimizado
func (h *GetClientsHandler) Search(w http.ResponseWriter, r *http.Request) {
	clientType := r.URL.Query().Get("type")
	search := r.URL.Query().Get("search")

	if search == "" {
		http.Error(w, "Search term is required", http.StatusBadRequest)
		return
	}

	var filter bson.M
	switch clientType {
	case "rental":
		filter = bson.M{"nombres": bson.M{"$regex": search, "$options": "i"}} // Case insensitive search for rentals
	case "guest":
		filter = bson.M{"customID": bson.M{"$regex": search, "$options": "i"}} // Case insensitive search for guests
	default:
		http.Error(w, "Invalid client type", http.StatusBadRequest)
		return
	}

	collection := h.Client.Database(constants.MongoDBDatabase).Collection(constants.CollectionClients)
	opts := options.Find().SetSort(bson.D{{Key: "_id", Value: 1}}) // Ordenar por ID para paginación

	// Realizar la búsqueda completa sin paginación, pero limitando resultados en cada solicitud
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

	response := map[string]interface{}{
		"clients": clients,
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(response)
}
