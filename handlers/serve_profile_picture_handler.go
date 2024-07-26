package handlers

import (
	"context"
	"log"
	"net/http"

	"hotelman-backend/constants"
	"hotelman-backend/models"

	"github.com/cloudinary/cloudinary-go/v2"
	"github.com/dgrijalva/jwt-go"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
)

type ServeProfilePictureHandler struct {
	Client     *mongo.Client
	JWTKey     []byte
	Cloudinary *cloudinary.Cloudinary // Agrega el campo para Cloudinary
}

func NewServeProfilePictureHandler(client *mongo.Client, jwtKey []byte, cloudinary *cloudinary.Cloudinary) *ServeProfilePictureHandler {
	return &ServeProfilePictureHandler{
		Client:     client,
		JWTKey:     jwtKey,
		Cloudinary: cloudinary,
	}
}

func (h *ServeProfilePictureHandler) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	// Obtener el token de la cookie
	cookie, err := r.Cookie("Authorize")
	if err != nil {
		http.Error(w, "Unauthorized", http.StatusUnauthorized)
		log.Println("No token cookie found:", err)
		return
	}

	tokenString := cookie.Value
	log.Println("Token string from cookie:", tokenString)

	// Parsear y verificar el token JWT
	token, err := jwt.Parse(tokenString, func(token *jwt.Token) (interface{}, error) {
		return h.JWTKey, nil
	})

	if err != nil || !token.Valid {
		http.Error(w, "Invalid token", http.StatusUnauthorized)
		log.Println("Invalid token:", err)
		return
	}

	// Extraer los claims estándar
	claims, ok := token.Claims.(jwt.MapClaims)
	if !ok || !token.Valid {
		http.Error(w, "Invalid token", http.StatusUnauthorized)
		log.Println("Invalid token claims")
		return
	}

	// Obtener el correo del usuario desde los claims
	email, ok := claims["username"].(string)
	if !ok {
		http.Error(w, "Invalid token claims", http.StatusUnauthorized)
		log.Println("Username not found in token claims")
		return
	}
	log.Println("Email from token claims:", email)

	// Obtener la URL de la imagen de perfil del usuario desde la base de datos
	collection := h.Client.Database(constants.MongoDBDatabase).Collection(constants.CollectionUsers)
	filter := bson.M{"correo": email}
	var user models.User
	err = collection.FindOne(context.Background(), filter).Decode(&user)
	if err != nil {
		http.Error(w, "User not found", http.StatusNotFound)
		log.Println("User not found in database:", err)
		return
	}

	// Obtener la URL de la imagen de perfil desde Cloudinary
	imageURL := user.ProfilePicture
	if imageURL == "" {
		http.NotFound(w, r)
		log.Println("No profile picture URL found for user:", email)
		return
	}

	// Redirigir al cliente a la URL de la imagen de perfil en Cloudinary
	http.Redirect(w, r, imageURL, http.StatusSeeOther)
	log.Println("Redirecting to profile picture:", imageURL)
}
