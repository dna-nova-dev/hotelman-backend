package handlers

import (
	"context"
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"strconv"
	"strings"
	"time"

	"hotelman-backend/constants"
	"hotelman-backend/models"
	"hotelman-backend/services"

	"github.com/google/uuid"

	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
)

// CreateClientHandler maneja la creación de nuevos clientes (Rental o Guest)
type CreateClientHandler struct {
	Client                 *mongo.Client
	CloudinaryService      *services.CloudinaryService
	GoogleDriveService     *services.GoogleDriveService
	LocalFileSystemService *services.LocalFileSystemService
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
		CURP:          r.FormValue("curp"),
		RoomNumber:    r.FormValue("RoomNumber"),
		CreatedAt:     time.Now(),
		UpdatedAt:     time.Now(),
	}

	// Subir archivos según el StorageSelector
	if constants.StorageSelector == "local" {
		h.uploadFilesLocal(w, r, &rental)
	} else {
		h.uploadFilesCloud(w, r, &rental)
	}

	collection := h.Client.Database(constants.MongoDBDatabase).Collection(constants.CollectionClients)
	_, err := collection.InsertOne(context.Background(), rental)
	if err != nil {
		http.Error(w, "Failed to create rental", http.StatusInternalServerError)
		return
	}

	w.WriteHeader(http.StatusCreated)
	json.NewEncoder(w).Encode(rental)
}

func (h *CreateClientHandler) uploadFilesLocal(w http.ResponseWriter, r *http.Request, rental *models.Rental) {
	// Subir archivos a sistema de archivos local
	/*contratoFile, contratoHandler, err := r.FormFile("contratoFile")
	if err == nil {
		defer contratoFile.Close()
		contratoURL, err := h.LocalFileSystemService.UploadFilePDF(contratoFile, contratoHandler)
		if err != nil {
			http.Error(w, "Error al subir el contrato localmente", http.StatusInternalServerError)
			return
		}
		rental.ContratoURL = contratoURL
	}*/
	ineFile, ineHandler, err := r.FormFile("ineFile")
	if err == nil {
		defer ineFile.Close()
		ineURL, err := h.LocalFileSystemService.UploadFileImage(ineFile, ineHandler)
		if err != nil {
			http.Error(w, "Error al subir el INE localmente", http.StatusInternalServerError)
			return
		}
		log.Println("INE on create client handler: ", ineURL)
		rental.INEURL = ineURL
	}
}

func (h *CreateClientHandler) uploadFilesCloud(w http.ResponseWriter, r *http.Request, rental *models.Rental) {
	// Subir archivos a Google Drive y Cloudinary
	contratoFile, contratoHandler, err := r.FormFile("contratoFile")
	if err == nil {
		defer contratoFile.Close()
		contratoURL, err := h.GoogleDriveService.UploadFile(contratoFile, contratoHandler)
		if err != nil {
			http.Error(w, "Error al subir el contrato a Google Drive", http.StatusInternalServerError)
			return
		}
		rental.ContratoURL = contratoURL
	}

	ineFile, ineHandler, err := r.FormFile("ineFile")
	if err == nil {
		defer ineFile.Close()
		ineURL, err := h.CloudinaryService.UploadINEPicture(ineFile, ineHandler)
		if err != nil {
			http.Error(w, "Error al subir el INE a Cloudinary", http.StatusInternalServerError)
			return
		}
		rental.INEURL = ineURL
	}
}

func (h *CreateClientHandler) createGuest(w http.ResponseWriter, r *http.Request) {
	// Generar ID personalizado
	customID := generateCustomID(r.FormValue("hair"), r.FormValue("roomNumber"))

	guest := models.Guest{
		ID:               primitive.NewObjectID(),
		CustomID:         customID, // Asignar el ID personalizado
		ExtraDescription: r.FormValue("extraDescription"),
		Hair:             r.FormValue("hair"),
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

// generateCustomID genera un ID personalizado en formato CURP
func generateCustomID(hair string, roomNumber string) string {
	// Generar un UUID y tomar solo los primeros 8 caracteres para reducir el tamaño
	uuidPart := uuid.New().String()[:8]

	// Convertir el tipo de cabello y el número de habitación a una forma de cadena
	// Eliminar espacios del tipo de cabello y convertir a mayúsculas
	hairPart := strings.ToUpper(strings.ReplaceAll(hair, " ", ""))
	roomNumberPart := strings.ToUpper(strings.ReplaceAll(roomNumber, " ", ""))

	// Asegurar que el hairPart y roomNumberPart tengan longitud fija para simular el formato CURP
	if len(hairPart) > 2 {
		hairPart = hairPart[:2]
	}
	if len(roomNumberPart) > 2 {
		roomNumberPart = roomNumberPart[:2]
	}

	// Combinar los datos con el UUID para crear el ID personalizado
	customID := fmt.Sprintf("%s%s%s", hairPart, roomNumberPart, uuidPart)

	// Limitar la longitud del ID a 18 caracteres para aproximarse al formato CURP
	if len(customID) > 18 {
		customID = customID[:18]
	}

	return customID
}

func parseFloat(value string) float64 {
	result, _ := strconv.ParseFloat(value, 64)
	return result
}

func parseInt(value string) int {
	result, _ := strconv.Atoi(value)
	return result
}
