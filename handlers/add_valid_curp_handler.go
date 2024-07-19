package handlers

import (
	"context"
	"encoding/json"
	"net/http"

	"hotelman-backend/constants"

	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
)

type AddValidCURPHandler struct {
	Client *mongo.Client
}

func (h *AddValidCURPHandler) Handle(w http.ResponseWriter, r *http.Request) {
	var curpData struct {
		CURP string `json:"curp"`
	}
	err := json.NewDecoder(r.Body).Decode(&curpData)
	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
		return
	}

	collection := h.Client.Database(constants.MongoDBDatabase).Collection(constants.CollectionValidCURPs)
	_, err = collection.InsertOne(context.TODO(), bson.M{"curp": curpData.CURP})
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		return
	}

	w.WriteHeader(http.StatusCreated)
}
