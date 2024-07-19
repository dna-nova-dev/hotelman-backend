package handlers

import (
	"encoding/json"
	"net/http"
	"sync"
	"time"

	"hotelman-backend/constants"
	"hotelman-backend/models"

	"github.com/dgrijalva/jwt-go"
	"github.com/google/uuid"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
	"golang.org/x/crypto/bcrypt"
)

type LoginHandler struct {
	Client   *mongo.Client
	jwtKey   []byte
	TokenMap *sync.Map // Mapa sincronizado para almacenar tokens activos
}

func NewLoginHandler(client *mongo.Client, jwtKey []byte) *LoginHandler {
	return &LoginHandler{
		Client:   client,
		jwtKey:   jwtKey,
		TokenMap: &sync.Map{},
	}
}

func (h *LoginHandler) Handle(w http.ResponseWriter, r *http.Request) {
	var creds models.Credentials
	err := json.NewDecoder(r.Body).Decode(&creds)
	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
		return
	}

	// Verificar si ya hay un token válido para el usuario
	if tokenString, ok := h.getTokenFromMap(creds.Username); ok {
		// Si hay un token válido, establecerlo en la cookie y devolverlo sin generar uno nuevo
		http.SetCookie(w, &http.Cookie{
			Name:     "token",
			Value:    tokenString,
			Expires:  time.Now().Add(7 * 24 * time.Hour),
			HttpOnly: true,
		})
		response := map[string]interface{}{
			"token":  tokenString,
			"claims": extractClaimsFromToken(tokenString),
		}
		w.Header().Set("Content-Type", "application/json")
		json.NewEncoder(w).Encode(response)
		return
	}

	collection := h.Client.Database(constants.MongoDBDatabase).Collection(constants.CollectionUsers)
	var storedUser models.User

	// Intentar buscar por correo electrónico
	err = collection.FindOne(r.Context(), bson.M{"correo": creds.Username}).Decode(&storedUser)
	if err != nil {
		// Intentar buscar por CURP si no se encontró por correo
		err = collection.FindOne(r.Context(), bson.M{"curp": creds.Username}).Decode(&storedUser)
		if err != nil {
			w.WriteHeader(http.StatusUnauthorized)
			return
		}
	}

	// Verificar la contraseña
	if bcrypt.CompareHashAndPassword([]byte(storedUser.Password), []byte(creds.Password)) != nil {
		w.WriteHeader(http.StatusUnauthorized)
		return
	}

	// Obtener el token existente o generar uno nuevo si es necesario
	tokenString, err := h.getOrCreateToken(storedUser)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		return
	}

	// Establecer el token JWT en una cookie HttpOnly
	http.SetCookie(w, &http.Cookie{
		Name:     "Autorization",
		Value:    tokenString,
		Expires:  time.Now().Add(7 * 24 * time.Hour),
		HttpOnly: true,
	})

	// Responder con el token JWT y las claims en la respuesta JSON
	response := map[string]interface{}{
		"token":  tokenString,
		"claims": extractClaimsFromToken(tokenString),
	}
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(response)
}

// Función auxiliar para obtener un token válido del mapa o crear uno nuevo
func (h *LoginHandler) getOrCreateToken(user models.User) (string, error) {
	// Verificar si ya hay un token válido para el usuario
	if tokenString, ok := h.getTokenFromMap(user.Correo); ok {
		// Verificar validez del token
		claims := extractClaimsFromToken(tokenString)
		if claims != nil && claims.ExpiresAt > time.Now().Unix() {
			return tokenString, nil // Devolver el token válido existente
		}
		// Eliminar token expirado del mapa
		h.TokenMap.Delete(user.Correo)
	}

	// Generar un nuevo TokenID aleatorio usando UUID
	tokenID := uuid.New().String()

	// Generar un nuevo token JWT
	expirationTime := time.Now().Add(7 * 24 * time.Hour)
	claims := &models.Claims{
		TokenID:  tokenID, // Asignar TokenID aleatorio
		Username: user.Correo,
		Role:     user.Rol,
		StandardClaims: jwt.StandardClaims{
			ExpiresAt: expirationTime.Unix(),
		},
	}

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	tokenString, err := token.SignedString(h.jwtKey)
	if err != nil {
		return "", err
	}

	// Almacenar el nuevo token en el mapa de tokens activos
	h.TokenMap.Store(user.Correo, tokenString)

	return tokenString, nil
}

// Función auxiliar para obtener un token válido del mapa
func (h *LoginHandler) getTokenFromMap(username string) (string, bool) {
	var tokenString string
	var found bool

	h.TokenMap.Range(func(key, value interface{}) bool {
		t := value.(string)
		claims := extractClaimsFromToken(t)

		if claims != nil && claims.Username == username {
			tokenString = t
			found = true
			return false // Detener el rango
		}
		return true // Continuar el rango
	})

	return tokenString, found
}

// Función auxiliar para extraer los claims de un token JWT
func extractClaimsFromToken(tokenString string) *models.Claims {
	token, _ := jwt.ParseWithClaims(tokenString, &models.Claims{}, func(token *jwt.Token) (interface{}, error) {
		return []byte(constants.JWTSecretKey), nil
	})

	if claims, ok := token.Claims.(*models.Claims); ok && token.Valid {
		return claims
	}

	return nil
}
