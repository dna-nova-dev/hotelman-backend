package routes

import (
	"hotelman-backend/constants"
	"hotelman-backend/handlers"
	"hotelman-backend/middleware"
	"net/http"

	"github.com/gorilla/mux"
	"go.mongodb.org/mongo-driver/mongo"
)

func RegisterRoutes(router *mux.Router, client *mongo.Client) {
	// Crear instancias de los nuevos handlers
	setupAdminHandler := &handlers.SetupAdminHandler{Client: client}
	signupHandler := &handlers.SignupHandler{Client: client}
	welcomeHandler := &handlers.WelcomeHandler{}
	addValidCURPHandler := &handlers.AddValidCURPHandler{Client: client}

	// Crear instancia de LoginHandler con jwtKey y Client
	loginHandler := handlers.NewLoginHandler(client, []byte(constants.JWTSecretKey))

	// Instancia de GetAllUsersHandler
	allUsersHandler := handlers.NewGetAllUsersHandler(client)

	// Crear instancia del middleware RequireAuth
	requireAuth := middleware.NewRequireAuth([]byte(constants.JWTSecretKey))

	// Endpoints utilizando los nuevos handlers
	router.HandleFunc("/setup", setupAdminHandler.Handle).Methods("POST")
	router.HandleFunc("/signup", signupHandler.Handle).Methods("POST")
	router.HandleFunc("/login", loginHandler.Handle).Methods("POST")
	router.HandleFunc("/welcome", welcomeHandler.Handle).Methods("GET")
	router.HandleFunc("/add-valid-curp", addValidCURPHandler.Handle).Methods("POST")

	// Endpoint protegido utilizando el middleware RequireAuth
	router.Handle("/all-users", requireAuth.Middleware(http.HandlerFunc(allUsersHandler.Handle))).Methods("GET")
}
