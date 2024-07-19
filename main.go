package main

import (
	"log"
	"net/http"
	"time"

	"github.com/gorilla/mux"
	"go.mongodb.org/mongo-driver/mongo"

	"hotelman-backend/config"
	"hotelman-backend/constants" // Importa el paquete de constantes
	"hotelman-backend/routes"
)

var client *mongo.Client

func main() {
	client = config.ConnectDB() // Conecta a la base de datos MongoDB

	router := mux.NewRouter()
	routes.RegisterRoutes(router, client) // Registra las rutas usando las constantes

	srv := &http.Server{
		Handler:      router,
		Addr:         constants.ServerAddress + ":" + constants.ServerPort, // Usa las constantes para la dirección y puerto
		WriteTimeout: 15 * time.Second,
		ReadTimeout:  15 * time.Second,
	}

	log.Printf("Servidor iniciado en http://%s:%s", constants.ServerAddress, constants.ServerPort)
	log.Fatal(srv.ListenAndServe())
}
