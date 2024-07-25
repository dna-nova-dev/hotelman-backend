package handlers

import (
	"encoding/json"
	"net/http"

	"hotelman-backend/constants"
	"hotelman-backend/models"

	"github.com/dgrijalva/jwt-go"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
)

type UserHandler struct {
	Client *mongo.Client
	jwtKey []byte
}

func NewUserHandler(client *mongo.Client, jwtKey []byte) *UserHandler {
	return &UserHandler{
		Client: client,
		jwtKey: jwtKey,
	}
}

func (h *UserHandler) Handle(w http.ResponseWriter, r *http.Request) {
	// Obtener el token de las cookies
	cookie, err := r.Cookie("Authorize")
	if err != nil {
		http.Error(w, "Unauthorized", http.StatusUnauthorized)
		return
	}

	tokenString := cookie.Value

	// Parsear el token y extraer los claims
	token, err := jwt.ParseWithClaims(tokenString, &models.Claims{}, func(token *jwt.Token) (interface{}, error) {
		return h.jwtKey, nil
	})
	if err != nil || !token.Valid {
		http.Error(w, "Unauthorized", http.StatusUnauthorized)
		return
	}

	claims, ok := token.Claims.(*models.Claims)
	if !ok {
		http.Error(w, "Unauthorized", http.StatusUnauthorized)
		return
	}

	// Buscar el usuario en la base de datos
	collection := h.Client.Database(constants.MongoDBDatabase).Collection(constants.CollectionUsers)
	var user models.User

	err = collection.FindOne(r.Context(), bson.M{"correo": claims.Username}).Decode(&user)
	if err != nil {
		http.Error(w, "User not found", http.StatusNotFound)
		return
	}

	// Devolver los datos del usuario en formato JSON
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(user)
}
