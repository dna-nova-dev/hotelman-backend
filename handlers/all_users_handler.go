package handlers

import (
	"context"
	"encoding/json"
	"net/http"

	"hotelman-backend/constants"
	"hotelman-backend/models"

	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
)

type GetAllUsersHandler struct {
	Client *mongo.Client
}

func NewGetAllUsersHandler(client *mongo.Client) *GetAllUsersHandler {
	return &GetAllUsersHandler{
		Client: client,
	}
}

func (h *GetAllUsersHandler) Handle(w http.ResponseWriter, r *http.Request) {
	collection := h.Client.Database(constants.MongoDBDatabase).Collection(constants.CollectionUsers)

	// Definir un slice para almacenar los usuarios recuperados
	var users []models.User

	// Obtener todos los usuarios de la colección
	cur, err := collection.Find(context.TODO(), bson.M{})
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		return
	}
	defer cur.Close(context.TODO())

	// Iterar sobre los resultados y decodificar cada documento en una estructura User
	for cur.Next(context.TODO()) {
		var user models.User
		err := cur.Decode(&user)
		if err != nil {
			w.WriteHeader(http.StatusInternalServerError)
			return
		}
		users = append(users, user)
	}

	// Verificar si hubo errores durante el cursor
	if err := cur.Err(); err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		return
	}

	// Codificar el slice de usuarios como JSON y enviar como respuesta
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(users)
}
