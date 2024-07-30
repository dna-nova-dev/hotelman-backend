package routes

import (
	"hotelman-backend/constants"
	"hotelman-backend/handlers"
	"hotelman-backend/middleware"
	"hotelman-backend/services"
	"net/http"
	"os"

	"github.com/gorilla/mux"
	"go.mongodb.org/mongo-driver/mongo"
)

// RegisterRoutes registra todas las rutas y handlers
func RegisterRoutes(router *mux.Router, client *mongo.Client, cloudinaryURL string) {
	var cloudinaryService *services.CloudinaryService
	var googleDriveService *services.GoogleDriveService
	var localFileSystemService *services.LocalFileSystemService
	var err error

	if constants.StorageSelector == "local" {
		localFileSystemService, err = services.NewLocalFileSystemService(constants.LocalFileSystemFolder)
		if err != nil {
			panic("Failed to initialize local file system service: " + err.Error())
		}
	} else {
		cloudinaryService, err = services.NewCloudinaryService(cloudinaryURL)
		if err != nil {
			panic("Failed to initialize Cloudinary service: " + err.Error())
		}
		googleDriveService, err = services.NewGoogleDriveService(constants.GoogleDriveCredentialsPath)
		if err != nil {
			panic("Failed to initialize Google Drive service: " + err.Error())
		}
	}

	// Crear instancias de los nuevos handlers
	setupAdminHandler := &handlers.SetupAdminHandler{Client: client}
	signupHandler := &handlers.SignupHandler{
		Client:                 client,
		CloudinaryService:      cloudinaryService,
		LocalFileSystemService: localFileSystemService,
	}
	welcomeHandler := &handlers.WelcomeHandler{}
	addValidCURPHandler := &handlers.AddValidCURPHandler{Client: client}

	// Crear instancia de LoginHandler con jwtKey y Client
	loginHandler := handlers.NewLoginHandler(client, []byte(constants.JWTSecretKey))
	logoutHandler := handlers.LogoutHandler{}

	// Crear Instancia Cliente:
	clientsHandler := &handlers.GetClientsHandler{Client: client}
	createHandler := &handlers.CreateClientHandler{
		Client:                 client,
		CloudinaryService:      cloudinaryService,
		GoogleDriveService:     googleDriveService,
		LocalFileSystemService: localFileSystemService,
	}
	// Instancia de GetAllUsersHandler
	allUsersHandler := handlers.NewGetAllUsersHandler(client)
	userDataHandler := handlers.NewUserHandler(client, []byte(constants.JWTSecretKey))

	// Instancia de Room Handler
	roomHandler := &handlers.RoomHandler{Client: client}

	// Instancia de Analytics handler
	analyticsHandler := &handlers.AnalyticsHandler{Client: client}

	// Obtener la ruta raíz del proyecto
	rootPath, err := os.Getwd()
	if err != nil {
		panic("Failed to get root directory: " + err.Error())
	}

	serveHandler := &handlers.ServeFileHandler{UploadsDir: rootPath + constants.LocalFileSystemFolder}

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

	// Endpoint rooms
	router.HandleFunc("/rooms", roomHandler.CreateRoomHandler).Methods("POST")
	router.HandleFunc("/rooms/status", roomHandler.UpdateRoomStatusHandler).Methods("PUT")
	router.HandleFunc("/rooms/occupant", roomHandler.GetRoomOccupantHandler).Methods("GET")
	router.HandleFunc("/rooms/assign", roomHandler.AssignOccupantHandler).Methods("PUT")

	// Endpoint analytics
	router.HandleFunc("/analytics", analyticsHandler.GetAnalyticsHandler).Methods("GET")

	// Content Serve
	router.HandleFunc("/serve", serveHandler.Handle).Methods("GET")
}
