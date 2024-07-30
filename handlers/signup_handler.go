package handlers

import (
	"context"
	"encoding/json"
	"net/http"

	"hotelman-backend/constants"
	"hotelman-backend/models"
	"hotelman-backend/services"

	"go.mongodb.org/mongo-driver/mongo"
	"golang.org/x/crypto/bcrypt"
)

type SignupHandler struct {
	Client                 *mongo.Client
	CloudinaryService      *services.CloudinaryService
	LocalFileSystemService *services.LocalFileSystemService
}

func (h *SignupHandler) Handle(w http.ResponseWriter, r *http.Request) {
	err := r.ParseMultipartForm(10 << 20) // Limit to 10 MB
	if err != nil {
		http.Error(w, "Error parsing form", http.StatusBadRequest)
		return
	}

	collection := h.Client.Database(constants.MongoDBDatabase).Collection(constants.CollectionUsers)
	var newUser models.User

	newUser.Nombres = r.FormValue("nombres")
	newUser.Apellidos = r.FormValue("apellidos")
	newUser.Correo = r.FormValue("correo")
	newUser.Celular = r.FormValue("numeroCelular")
	newUser.Password = r.FormValue("contrasena")
	newUser.Rol = "Recepcionista"
	confirmarContrasena := r.FormValue("confirmarContrasena")
	newUser.CURP = r.FormValue("curp")

	if newUser.CURP != "" {
		if !isValidCURP(newUser.CURP) {
			http.Error(w, "CURP inválido", http.StatusBadRequest)
			return
		} else {
			newUser.Rol = "Administracion"
		}
	}

	if newUser.Password != confirmarContrasena {
		http.Error(w, "Las contraseñas no coinciden", http.StatusBadRequest)
		return
	}

	hashedPassword, err := bcrypt.GenerateFromPassword([]byte(newUser.Password), bcrypt.DefaultCost)
	if err != nil {
		http.Error(w, "Error al hashear la contraseña", http.StatusInternalServerError)
		return
	}
	newUser.Password = string(hashedPassword)

	file, handler, err := r.FormFile("profilePicture")
	if err == nil {
		defer file.Close()
		if constants.StorageSelector == "local" {
			profilePictureURL, err := h.LocalFileSystemService.UploadFileImage(file, handler)
			if err != nil {
				http.Error(w, "Error al guardar la imagen de perfil localmente", http.StatusInternalServerError)
				return
			}
			newUser.ProfilePicture = profilePictureURL
		} else {
			profilePictureURL, err := h.CloudinaryService.UploadProfilePicture(file, handler)
			if err != nil {
				http.Error(w, "Error al guardar la imagen de perfil en la nube", http.StatusInternalServerError)
				return
			}
			newUser.ProfilePicture = profilePictureURL
		}
	} else {
		newUser.ProfilePicture = ""
	}

	_, err = collection.InsertOne(context.TODO(), newUser)
	if err != nil {
		http.Error(w, "Error al registrar el usuario", http.StatusInternalServerError)
		return
	}

	w.WriteHeader(http.StatusCreated)
	json.NewEncoder(w).Encode(map[string]string{"message": "Usuario registrado con éxito"})
}
