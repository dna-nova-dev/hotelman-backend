package handlers

import (
	"bytes"
	"context"
	"encoding/base64"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"os"
	"path/filepath"
	"strconv"
	"strings"
	"time"

	"hotelman-backend/constants"
	"hotelman-backend/models"
	"hotelman-backend/services"
	"hotelman-backend/utils"

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
	var clientData map[string]interface{}
	err := json.NewDecoder(r.Body).Decode(&clientData)
	if err != nil {
		http.Error(w, "Error parsing JSON", http.StatusBadRequest)
		return
	}

	clientType, ok := clientData["type"].(string)
	if !ok {
		http.Error(w, "Client type is required", http.StatusBadRequest)
		return
	}

	switch clientType {
	case "rental":
		h.createRental(w, clientData)
	case "guest":
		h.createGuest(w, clientData)
	default:
		http.Error(w, "Unknown client type", http.StatusBadRequest)
	}
}

func (h *CreateClientHandler) createRental(w http.ResponseWriter, clientData map[string]interface{}) {
	rental := models.Rental{
		ID:            primitive.NewObjectID(),
		Nombres:       clientData["nombres"].(string),
		Apellidos:     clientData["apellidos"].(string),
		Correo:        clientData["correo"].(string),
		NumeroCelular: clientData["numeroCelular"].(string),
		CURP:          clientData["curp"].(string),
		RoomNumber:    clientData["RoomNumber"].(string),
		CreatedAt:     time.Now(),
		UpdatedAt:     time.Now(),
	}

	if constants.StorageSelector == "local" {
		h.uploadFilesLocal(w, clientData, &rental)
	} else {
		// Si usas almacenamiento en la nube, implementa aquí la lógica
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

func (h *CreateClientHandler) uploadFilesLocal(w http.ResponseWriter, clientData map[string]interface{}, rental *models.Rental) {
	// Subir archivos desde Base64 a sistema de archivos local

	// Subir INE File
	if ineFileBase64, ok := clientData["ineFile"].(map[string]interface{})["fileData"].(string); ok {
		ineURL, err := h.saveBase64File(ineFileBase64, "images")
		if err != nil {
			http.Error(w, "Error al subir el INE localmente", http.StatusInternalServerError)
			return
		}
		rental.INEURL = ineURL
	}

	// Subir Contrato File
	if contratoFileBase64, ok := clientData["contratoFile"].(map[string]interface{})["fileData"].(string); ok {
		contratoURL, err := h.saveBase64File(contratoFileBase64, "documents")
		if err != nil {
			http.Error(w, "Error al subir el contrato localmente", http.StatusInternalServerError)
			return
		}
		rental.ContratoURL = contratoURL
	}
}

func (h *CreateClientHandler) saveBase64File(base64Data string, folder string) (string, error) {
	// Decodificar Base64
	decodedData, err := base64.StdEncoding.DecodeString(base64Data)
	if err != nil {
		return "", fmt.Errorf("failed to decode base64 data: %w", err)
	}

	// Generar un nombre de archivo único
	fileName := fmt.Sprintf("%s_%s", uuid.New().String(), folder)
	filePath := filepath.Join(h.LocalFileSystemService.BasePath, folder, fileName)

	// Crear el archivo en el sistema de archivos local
	err = h.createFile(filePath, bytes.NewReader(decodedData))
	if err != nil {
		return "", fmt.Errorf("failed to save file: %w", err)
	}

	ip, err := utils.GetPublicIP()
	if err != nil {
		return "", fmt.Errorf("unable to get public IP: %v", err)
	}

	url := fmt.Sprintf("http://%s:8000/serve?folder=%s&filename=%s", ip, folder, fileName)
	return url, nil
}

func (h *CreateClientHandler) createFile(filePath string, data io.Reader) error {
	// Crea el archivo en el sistema de archivos local
	dst, err := os.Create(filePath)
	if err != nil {
		return fmt.Errorf("unable to create file: %v", err)
	}
	defer dst.Close()

	// Copia el contenido del archivo cargado al nuevo archivo en el sistema de archivos local
	_, err = io.Copy(dst, data)
	if err != nil {
		return fmt.Errorf("unable to copy file content: %v", err)
	}

	return nil
}

func (h *CreateClientHandler) createGuest(w http.ResponseWriter, clientData map[string]interface{}) {
	// Generar ID personalizado
	customID := generateCustomID(clientData["hair"].(string), clientData["roomNumber"].(string))

	guest := models.Guest{
		ID:               primitive.NewObjectID(),
		CustomID:         customID, // Asignar el ID personalizado
		ExtraDescription: clientData["extraDescription"].(string),
		Hair:             clientData["hair"].(string),
		Height:           clientData["height"].(string),
		RoomNumber:       clientData["roomNumber"].(string),
		Price:            parseFloat(clientData["price"].(string)),
		Duration:         parseInt(clientData["duration"].(string)),
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
