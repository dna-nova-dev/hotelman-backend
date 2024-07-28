package handlers

import (
	"context"
	"encoding/json"
	"net/http"
	"strconv"
	"time"

	"hotelman-backend/constants"
	"hotelman-backend/models"
	"hotelman-backend/services"

	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
)

// CreateClientHandler maneja la creación de nuevos clientes (Rental o Guest)
type CreateClientHandler struct {
	Client             *mongo.Client
	CloudinaryService  *services.CloudinaryService
	GoogleDriveService *services.GoogleDriveService
}

// Handle procesa la solicitud de creación de un nuevo cliente
func (h *CreateClientHandler) Handle(w http.ResponseWriter, r *http.Request) {
	err := r.ParseMultipartForm(10 << 20) // Limit to 10 MB
	if err != nil {
		http.Error(w, "Error parsing form", http.StatusBadRequest)
		return
	}

	clientType := r.FormValue("type")
	switch clientType {
	case "rental":
		h.createRental(w, r)
	case "guest":
		h.createGuest(w, r)
	default:
		http.Error(w, "Unknown client type", http.StatusBadRequest)
	}
}

func (h *CreateClientHandler) createRental(w http.ResponseWriter, r *http.Request) {
	rental := models.Rental{
		ID:            primitive.NewObjectID(),
		Nombres:       r.FormValue("nombres"),
		Apellidos:     r.FormValue("apellidos"),
		Correo:        r.FormValue("correo"),
		NumeroCelular: r.FormValue("numeroCelular"),
		INEString:     r.FormValue("INEString"),
		RoomNumber:    r.FormValue("RoomNumber"),
		CreatedAt:     time.Now(),
		UpdatedAt:     time.Now(),
	}

	// Subir archivos a Cloudinary
	contratoFile, contratoHandler, err := r.FormFile("contratoFile")
	if err == nil {
		defer contratoFile.Close()
		contratoURL, err := h.GoogleDriveService.UploadFile(contratoFile, contratoHandler)
		if err != nil {
			http.Error(w, "Error al subir el contrato", http.StatusInternalServerError)
			return
		}
		rental.ContratoURL = contratoURL
	}

	/*ineFile, ineHandler, err := r.FormFile("ineFile")
	if err == nil {
		defer ineFile.Close()
		ineURL, err := h.CloudinaryService.UploadINEPicture(ineFile, ineHandler)
		if err != nil {
			http.Error(w, "Error al subir el INE", http.StatusInternalServerError)
			return
		}
		rental.INEURL = ineURL
	}*/

	collection := h.Client.Database(constants.MongoDBDatabase).Collection(constants.CollectionClients)
	_, err = collection.InsertOne(context.Background(), rental)
	if err != nil {
		http.Error(w, "Failed to create rental", http.StatusInternalServerError)
		return
	}

	w.WriteHeader(http.StatusCreated)
	json.NewEncoder(w).Encode(rental)
}

func (h *CreateClientHandler) createGuest(w http.ResponseWriter, r *http.Request) {
	guest := models.Guest{
		ID:               primitive.NewObjectID(),
		Email:            r.FormValue("email"),
		Phone:            r.FormValue("phone"),
		ExtraDescription: r.FormValue("extraDescription"),
		Name:             r.FormValue("name"),
		Height:           r.FormValue("height"),
		RoomNumber:       r.FormValue("roomNumber"),
		Price:            parseFloat(r.FormValue("price")),
		Duration:         parseInt(r.FormValue("duration")),
		CreatedAt:        time.Now(),
		UpdatedAt:        time.Now(),
	}

	collection := h.Client.Database(constants.MongoDBDatabase).Collection(constants.CollectionClients)
	_, err := collection.InsertOne(context.Background(), guest)
	if err != nil {
		http.Error(w, "Failed to create guest", http.StatusInternalServerError)
		return
	}

	w.WriteHeader(http.StatusCreated)
	json.NewEncoder(w).Encode(guest)
}

func parseFloat(value string) float64 {
	result, _ := strconv.ParseFloat(value, 64)
	return result
}

func parseInt(value string) int {
	result, _ := strconv.Atoi(value)
	return result
}
