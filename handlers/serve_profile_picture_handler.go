package handlers

import (
	"context"
	"log"
	"net/http"
	"os"
	"path/filepath"

	"hotelman-backend/constants"
	"hotelman-backend/models"

	"github.com/dgrijalva/jwt-go"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
)

type ServeProfilePictureHandler struct {
	Client  *mongo.Client
	JWTKey  []byte
	PicPath string // Directorio donde se almacenan las imágenes de perfil
}

func NewServeProfilePictureHandler(client *mongo.Client, jwtKey []byte, picPath string) *ServeProfilePictureHandler {
	return &ServeProfilePictureHandler{
		Client:  client,
		JWTKey:  jwtKey,
		PicPath: picPath,
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
	log.Println("User profile picture from database:", user.ProfilePicture)

	// Construir la ruta del archivo de la imagen
	filePath := filepath.Join(h.PicPath, user.ProfilePicture)
	log.Println("Constructed file path:", filePath)

	// Verificar si el archivo existe
	if _, err := os.Stat(filePath); os.IsNotExist(err) {
		http.NotFound(w, r)
		log.Println("Profile picture not found at path:", filePath)
		return
	}

	// Establecer el tipo de contenido
	w.Header().Set("Content-Type", "image/jpeg")

	// Leer y escribir el archivo de la imagen
	http.ServeFile(w, r, filePath)
	log.Println("Serving profile picture:", filePath)
}
