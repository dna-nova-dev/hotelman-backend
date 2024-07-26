package handlers

import (
	"context"
	"encoding/json"
	"mime/multipart"
	"net/http"

	"hotelman-backend/constants"
	"hotelman-backend/models"

	"github.com/cloudinary/cloudinary-go/v2"
	"github.com/cloudinary/cloudinary-go/v2/api/uploader"
	"go.mongodb.org/mongo-driver/mongo"
	"golang.org/x/crypto/bcrypt"
)

type SignupHandler struct {
	Client     *mongo.Client
	Cloudinary *cloudinary.Cloudinary // Agrega el campo para Cloudinary
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
		profilePictureURL, err := h.uploadProfilePictureToCloudinary(file, handler)
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

func (h *SignupHandler) uploadProfilePictureToCloudinary(file multipart.File, handler *multipart.FileHeader) (string, error) {
	// Subir la imagen a Cloudinary
	resp, err := h.Cloudinary.Upload.Upload(context.Background(), file, uploader.UploadParams{Folder: "profile_pictures"})
	if err != nil {
		return "", err
	}

	// Retornar la URL de la imagen subida
	return resp.SecureURL, nil
}
