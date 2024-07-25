package handlers

import (
	"context"
	"encoding/json"
	"io"
	"mime/multipart"
	"net/http"
	"os"
	"path/filepath"

	"hotelman-backend/constants"
	"hotelman-backend/models"

	"go.mongodb.org/mongo-driver/mongo"
	"golang.org/x/crypto/bcrypt"
)

type SignupHandler struct {
	Client *mongo.Client
}

func (h *SignupHandler) Handle(w http.ResponseWriter, r *http.Request) {
	// Parse the form to handle file uploads
	err := r.ParseMultipartForm(10 << 20) // Limit to 10 MB
	if err != nil {
		http.Error(w, "Error parsing form", http.StatusBadRequest)
		return
	}

	collection := h.Client.Database(constants.MongoDBDatabase).Collection(constants.CollectionUsers)
	var newUser models.User

	// Parse user data
	newUser.Nombres = r.FormValue("nombres")
	newUser.Apellidos = r.FormValue("apellidos")
	newUser.Correo = r.FormValue("correo")
	newUser.Celular = r.FormValue("numeroCelular")
	newUser.Password = r.FormValue("contrasena")
	newUser.Rol = "Recepcionista"                             // Actualizado a Password
	confirmarContrasena := r.FormValue("confirmarContrasena") // Se usa solo para validación
	newUser.CURP = r.FormValue("curp")

	// Validar el formato del CURP si está presente
	if newUser.CURP != "" {
		if !isValidCURP(newUser.CURP) {
			http.Error(w, "CURP inválido", http.StatusBadRequest)
			return
		}
	}

	// Validar las contraseñas
	if newUser.Password != confirmarContrasena {
		http.Error(w, "Las contraseñas no coinciden", http.StatusBadRequest)
		return
	}

	// Hash the password
	hashedPassword, err := bcrypt.GenerateFromPassword([]byte(newUser.Password), bcrypt.DefaultCost)
	if err != nil {
		http.Error(w, "Error al hashear la contraseña", http.StatusInternalServerError)
		return
	}
	newUser.Password = string(hashedPassword)

	// Handle profile picture upload
	file, handler, err := r.FormFile("profilePicture")
	if err == nil {
		defer file.Close()
		profilePictureURL, err := saveProfilePicture(file, handler)
		if err != nil {
			http.Error(w, "Error al guardar la imagen de perfil", http.StatusInternalServerError)
			return
		}
		newUser.ProfilePicture = profilePictureURL
	} else {
		newUser.ProfilePicture = ""
	}

	// Insert user into database
	_, err = collection.InsertOne(context.TODO(), newUser)
	if err != nil {
		http.Error(w, "Error al registrar el usuario", http.StatusInternalServerError)
		return
	}

	w.WriteHeader(http.StatusCreated)
	json.NewEncoder(w).Encode(map[string]string{"message": "Usuario registrado con éxito"})
}

func saveProfilePicture(file multipart.File, handler *multipart.FileHeader) (string, error) {
	tempFile, err := os.Create(filepath.Join("uploads", handler.Filename))
	if err != nil {
		return "", err
	}
	defer tempFile.Close()

	_, err = io.Copy(tempFile, file)
	if err != nil {
		return "", err
	}

	return tempFile.Name(), nil
}
