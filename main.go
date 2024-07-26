package main

import (
	"log"
	"net/http"
	"time"

	"github.com/cloudinary/cloudinary-go/v2"
	"github.com/gorilla/mux"
	"github.com/rs/cors" // Importa el paquete de CORS
	"go.mongodb.org/mongo-driver/mongo"

	"hotelman-backend/config"
	"hotelman-backend/constants" // Importa el paquete de constantes
	"hotelman-backend/routes"
)

var client *mongo.Client

func main() {
	client = config.ConnectDB() // Conecta a la base de datos MongoDB

	// Configura Cloudinary v2
	cloudinary, err := cloudinary.NewFromParams(constants.CloudinaryCloudName, constants.CloudinaryAPIKey, constants.CloudinaryAPISecret)
	if err != nil {
		log.Fatalf("Error al configurar Cloudinary: %v", err)
	}

	router := mux.NewRouter()
	routes.RegisterRoutes(router, client, cloudinary) // Registra las rutas

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
