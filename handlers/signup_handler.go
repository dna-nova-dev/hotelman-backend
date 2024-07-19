package handlers

import (
	"context"
	"encoding/json"
	"net/http"

	"hotelman-backend/constants"
	"hotelman-backend/models"

	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
	"golang.org/x/crypto/bcrypt"
)

type SignupHandler struct {
	Client *mongo.Client
}

func (h *SignupHandler) Handle(w http.ResponseWriter, r *http.Request) {
	collection := h.Client.Database(constants.MongoDBDatabase).Collection(constants.CollectionUsers)
	var newUser models.User
	err := json.NewDecoder(r.Body).Decode(&newUser)
	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
		return
	}

	// Validar el formato del CURP si está presente
	if newUser.CURP != "" {
		if !isValidCURP(newUser.CURP) {
			w.WriteHeader(http.StatusBadRequest)
			w.Write([]byte("El CURP proporcionado no es válido según el estándar mexicano."))
			return
		}
		// Asignar automáticamente el rol de "Administracion" si se proporcionó un CURP válido
		newUser.Rol = "Administracion"
	}

	// Si el rol es "Administracion", verificar si el CURP está en la lista de CURPs válidos
	if newUser.Rol == "Administracion" {
		validCURPs := h.Client.Database(constants.MongoDBDatabase).Collection(constants.CollectionValidCURPs)
		var result bson.M
		err := validCURPs.FindOne(context.TODO(), bson.M{"curp": newUser.CURP}).Decode(&result)
		if err != nil {
			w.WriteHeader(http.StatusForbidden)
			w.Write([]byte("CURP no válida para el registro de administradores."))
			return
		}
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
