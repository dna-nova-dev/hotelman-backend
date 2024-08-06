package main

import (
	"fmt"
	"log"
	"net/http"
	"strings"
	"time"

	"github.com/gorilla/mux"
	"github.com/rs/cors"
	"go.mongodb.org/mongo-driver/mongo"

	"hotelman-backend/config"
	"hotelman-backend/constants"
	"hotelman-backend/routes"
)

var client *mongo.Client

func main() {
	client = config.ConnectDB() // Conecta a la base de datos MongoDB

	// Construye la URL de Cloudinary utilizando las constantes
	cloudinaryURL := fmt.Sprintf("cloudinary://%s:%s@%s", constants.CloudinaryAPIKey, constants.CloudinaryAPISecret, constants.CloudinaryCloudName)

	router := mux.NewRouter()
	routes.RegisterRoutes(router, client, cloudinaryURL) // Registra las rutas
	allowedOrigins := parseAllowedOrigins(constants.FrontendURL)
	// Configura CORS
	c := cors.New(cors.Options{
		AllowedOrigins:   allowedOrigins,
		AllowCredentials: true,
		AllowedMethods:   []string{"GET", "POST", "PUT", "DELETE", "OPTIONS"},
		AllowedHeaders:   []string{"Content-Type", "Authorization"},
	})

	// Aplica el middleware de CORS
	handler := c.Handler(router)

	// Configurar los certificados SSL
	certFile := "/etc/letsencrypt/live/api-v1.hotelman.pulse.lat/fullchain.pem"
	keyFile := "/etc/letsencrypt/live/api-v1.hotelman.pulse.lat/privkey.pem"

	srv := &http.Server{
		Handler:      handler,
		Addr:         constants.ServerAddress + ":" + constants.ServerPort,
		WriteTimeout: 15 * time.Second,
		ReadTimeout:  15 * time.Second,
	}

	log.Printf("Frontend URL: %s", constants.FrontendURL)
	log.Printf("Servidor iniciado en https://%s:%s", constants.ServerAddress, constants.ServerPort)
	log.Fatal(srv.ListenAndServeTLS(certFile, keyFile))
}

// parseAllowedOrigins convierte una cadena separada por punto y coma en un arreglo de URLs
func parseAllowedOrigins(origins string) []string {
	return strings.Split(origins, ",")
}
