package main

import (
	"fmt"
	"log"
	"net/http"
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

	// Configura CORS
	c := cors.New(cors.Options{
		AllowedOrigins:   []string{"%s", constants.FrontendURL},
		AllowCredentials: true,
		AllowedMethods:   []string{"GET", "POST", "PUT", "DELETE", "OPTIONS"},
		AllowedHeaders:   []string{"Content-Type", "Authorization"},
	})

	// Aplica el middleware de CORS
	handler := c.Handler(router)

	srv := &http.Server{
		Handler:      handler,
		Addr:         constants.ServerAddress + ":" + constants.ServerPort,
		WriteTimeout: 15 * time.Second,
		ReadTimeout:  15 * time.Second,
	}
	log.Printf("Frontend URL: %s", constants.FrontendURL)
	log.Printf("Servidor iniciado en http://%s:%s", constants.ServerAddress, constants.ServerPort)
	log.Fatal(srv.ListenAndServe())
}
