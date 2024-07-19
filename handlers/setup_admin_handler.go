package handlers

import (
	"context"
	"encoding/json"
	"net/http"
	"regexp"

	"hotelman-backend/constants"
	"hotelman-backend/models"

	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
	"golang.org/x/crypto/bcrypt"
)

type SetupAdminHandler struct {
	Client *mongo.Client
}

func (h *SetupAdminHandler) Handle(w http.ResponseWriter, r *http.Request) {
	collection := h.Client.Database(constants.MongoDBDatabase).Collection(constants.CollectionUsers)
	var adminCount int64
	adminCount, err := collection.CountDocuments(context.TODO(), bson.M{"rol": "Administracion"})
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		return
	}

	if adminCount > 0 {
		w.WriteHeader(http.StatusForbidden)
		w.Write([]byte("Ya existe un administrador configurado. Use /signup para registrar nuevos usuarios."))
		return
	}

	var newUser models.User
	err = json.NewDecoder(r.Body).Decode(&newUser)
	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
		return
	}

	// Si el campo CURP está presente, validar el formato de CURP mexicano
	if newUser.CURP != "" {
		if !isValidCURP(newUser.CURP) {
			w.WriteHeader(http.StatusBadRequest)
			w.Write([]byte("El CURP proporcionado no es válido según el estándar mexicano."))
			return
		}
		// Asignar el rol como "Administracion" si se proporciona CURP válido
		newUser.Rol = "Administracion"
	}

	hashedPassword, err := bcrypt.GenerateFromPassword([]byte(newUser.Password), bcrypt.DefaultCost)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		return
	}

	newUser.Password = string(hashedPassword)
	_, err = collection.InsertOne(context.TODO(), newUser)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		return
	}

	w.WriteHeader(http.StatusCreated)
}

// isValidCURP valida si un CURP dado cumple con el formato mexicano estándar
func isValidCURP(curp string) bool {
	// Expresión regular para validar CURP mexicano
	regex := `^[A-Z]{4}[0-9]{6}[HM][A-Z]{5}[0-9]{2}$`
	match, _ := regexp.MatchString(regex, curp)
	return match
}
