package main

import (
	"log"
	"net/http"
	"os"
	"time"

	"github.com/gorilla/mux"
	"github.com/rs/cors" // Importa el paquete de CORS
	"go.mongodb.org/mongo-driver/mongo"

	"hotelman-backend/config"
	"hotelman-backend/constants" // Importa el paquete de constantes
	"hotelman-backend/routes"
)

var client *mongo.Client

func main() {
	// Crear la carpeta uploads si no existe
	err := os.MkdirAll("uploads", os.ModePerm)
	if err != nil {
		log.Fatalf("Error al crear la carpeta 'uploads': %v", err)
	}

	client = config.ConnectDB() // Conecta a la base de datos MongoDB

	router := mux.NewRouter()
	routes.RegisterRoutes(router, client) // Registra las rutas usando las constantes

	// Configura CORS
	c := cors.New(cors.Options{
		AllowedOrigins:   []string{"%s", constants.FrontendURL}, // Ajusta esto a la URL de tu frontend
		AllowCredentials: true,
		AllowedMethods:   []string{"GET", "POST", "PUT", "DELETE", "OPTIONS"},
		AllowedHeaders:   []string{"Content-Type", "Authorization"},
	})

	// Aplica el middleware de CORS
	handler := c.Handler(router)

	srv := &http.Server{
		Handler:      handler,
		Addr:         constants.ServerAddress + ":" + constants.ServerPort, // Usa las constantes para la dirección y puerto
		WriteTimeout: 15 * time.Second,
		ReadTimeout:  15 * time.Second,
	}
	log.Printf("Frontend URL: %s", constants.FrontendURL)
	log.Printf("Servidor iniciado en http://%s:%s", constants.ServerAddress, constants.ServerPort)
	log.Fatal(srv.ListenAndServe())
}
