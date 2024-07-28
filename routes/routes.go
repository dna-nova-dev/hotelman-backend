package routes

import (
	"hotelman-backend/constants"
	"hotelman-backend/handlers"
	"hotelman-backend/middleware"
	"hotelman-backend/services" // Importa el paquete de servicios
	"net/http"

	"github.com/gorilla/mux"
	"go.mongodb.org/mongo-driver/mongo"
)

// RegisterRoutes registra todas las rutas y handlers
func RegisterRoutes(router *mux.Router, client *mongo.Client, cloudinaryURL string) {
	// Crear instancia de CloudinaryService
	cloudinaryService, err := services.NewCloudinaryService(cloudinaryURL)
	if err != nil {
		panic("Failed to initialize Cloudinary service: " + err.Error())
	}
	googleDriveService, err := services.NewGoogleDriveService(constants.GoogleDriveCredentialsPath)
	if err != nil {
		panic("Failed to initialize Google Drive service: " + err.Error())
	}

	// Crear instancias de los nuevos handlers
	setupAdminHandler := &handlers.SetupAdminHandler{Client: client}
	signupHandler := &handlers.SignupHandler{Client: client, CloudinaryService: cloudinaryService}
	welcomeHandler := &handlers.WelcomeHandler{}
	addValidCURPHandler := &handlers.AddValidCURPHandler{Client: client}

	// Crear instancia de LoginHandler con jwtKey y Client
	loginHandler := handlers.NewLoginHandler(client, []byte(constants.JWTSecretKey))
	logoutHandler := handlers.LogoutHandler{}

	// Crear Instancia Cliente:
	clientsHandler := &handlers.GetClientsHandler{Client: client}
	createHandler := &handlers.CreateClientHandler{Client: client, CloudinaryService: cloudinaryService, GoogleDriveService: googleDriveService}
	// Instancia de GetAllUsersHandler
	allUsersHandler := handlers.NewGetAllUsersHandler(client)
	userDataHandler := handlers.NewUserHandler(client, []byte(constants.JWTSecretKey))

	// Crear instancia del middleware RequireAuth para roles específicos
	requireAuthAdmin := middleware.NewRequireAuth([]byte(constants.JWTSecretKey), []string{"Administracion"})
	requireAuthReceptionist := middleware.NewRequireAuth([]byte(constants.JWTSecretKey), []string{"Recepcionista", "Administracion"})

	// Endpoints utilizando los nuevos handlers
	router.HandleFunc("/setup", setupAdminHandler.Handle).Methods("POST")
	router.HandleFunc("/signup", signupHandler.Handle).Methods("POST")
	router.HandleFunc("/login", loginHandler.Handle).Methods("POST")
	router.HandleFunc("/logout", logoutHandler.Handle).Methods("POST")
	router.HandleFunc("/add-valid-curp", addValidCURPHandler.Handle).Methods("POST")

	// Endpoint protegido utilizando el middleware RequireAuth para administradores
	router.Handle("/welcome", requireAuthAdmin.Middleware(http.HandlerFunc(welcomeHandler.Handle))).Methods("GET")

	// Endpoint protegido utilizando el middleware RequireAuth para recepcionistas y administradores
	router.Handle("/all-users", requireAuthReceptionist.Middleware(http.HandlerFunc(allUsersHandler.Handle))).Methods("GET")
	router.Handle("/user", requireAuthReceptionist.Middleware(http.HandlerFunc(userDataHandler.Handle))).Methods("GET")
	router.HandleFunc("/create-client", createHandler.Handle).Methods("POST")
	router.HandleFunc("/clients", clientsHandler.Handle).Methods("GET")
}
