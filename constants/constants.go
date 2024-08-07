package constants

import (
	"log"
	"os"
	"strconv"

	"github.com/pelletier/go-toml"
)

var (
	// Roles de usuario
	RoleAdmin        string
	RoleReceptionist string

	// Estados HTTP
	StatusCreated             int
	StatusBadRequest          int
	StatusUnauthorized        int
	StatusForbidden           int
	StatusInternalServerError int

	// MongoDB
	MongoDBURI      string
	MongoDBDatabase string
	FrontendURL     string

	// Collections
	CollectionUsers      string
	CollectionValidCURPs string
	CollectionClients    string
	CollectionRooms      string // Nueva colección agregada

	// JWT
	JWTSecretKey string

	// Configuración de red
	ServerAddress string
	ServerPort    string

	// Cloudinary
	CloudinaryCloudName string
	CloudinaryAPIKey    string
	CloudinaryAPISecret string

	// Google Drive
	GoogleDriveFolderID        string
	GoogleDriveCredentialsPath string

	// Local File System
	LocalFileSystemFolder string

	// Storage Selector
	StorageSelector string

	// AllCollections contiene todos los nombres de colecciones definidos
	AllCollections []string
)

func init() {
	loadConfig()
}

func loadConfig() {
	// Variables por defecto
	defaultConfig := map[string]string{
		"RoleAdmin":                  "Administrator",
		"RoleReceptionist":           "Receptionist",
		"StatusCreated":              "201",
		"StatusBadRequest":           "400",
		"StatusUnauthorized":         "401",
		"StatusForbidden":            "403",
		"StatusInternalServerError":  "500",
		"MongoDBURI":                 "mongodb://localhost:27017",
		"MongoDBDatabase":            "testdb",
		"FrontendURL":                "http://localhost:3000",
		"CollectionUsers":            "users",
		"CollectionValidCURPs":       "valid_curps",
		"CollectionClients":          "clients",
		"CollectionRooms":            "rooms", // Nueva colección agregada
		"JWTSecretKey":               "my_secret_key",
		"ServerAddress":              "0.0.0.0",
		"ServerPort":                 "8000",
		"CloudinaryCloudName":        "your-cloudinary-cloud-name",
		"CloudinaryAPIKey":           "your-cloudinary-api-key",
		"CloudinaryAPISecret":        "your-cloudinary-api-secret",
		"GoogleDriveFolderID":        "your-google-drive-folder-id",
		"GoogleDriveCredentialsPath": "credentials.json",
		"LocalFileSystemFolder":      "/uploads",
		"StorageSelector":            "local",
	}

	// Intentar cargar desde variables de entorno
	for key := range defaultConfig {
		if val, exists := os.LookupEnv(key); exists {
			defaultConfig[key] = val
		}
	}

	// Si alguna variable no está en el entorno, cargar desde config.toml
	if !allEnvVariablesSet(defaultConfig) {
		loadFromToml(defaultConfig)
	}

	assignConfigValues(defaultConfig)
}

func allEnvVariablesSet(config map[string]string) bool {
	requiredKeys := []string{
		"RoleAdmin", "RoleReceptionist", "StatusCreated", "StatusBadRequest",
		"StatusUnauthorized", "StatusForbidden", "StatusInternalServerError",
		"MongoDBURI", "MongoDBDatabase", "FrontendURL", "CollectionUsers", "CollectionValidCURPs", "CollectionClients", "CollectionRooms",
		"JWTSecretKey", "ServerAddress", "ServerPort",
		"CloudinaryCloudName", "CloudinaryAPIKey", "CloudinaryAPISecret",
		"GoogleDriveFolderID", "GoogleDriveCredentialsPath", "LocalFileSystemFolder", "StorageSelector",
	}

	for _, key := range requiredKeys {
		if config[key] == "" {
			return false
		}
	}
	return true
}

func loadFromToml(config map[string]string) {
	configFile := "config.toml"

	// Verificar si el archivo config.toml ya existe
	if _, err := os.Stat(configFile); os.IsNotExist(err) {
		createDefaultConfig(configFile)
	}

	// Abrir y leer el archivo TOML
	tomlConfig, err := toml.LoadFile(configFile)
	if err != nil {
		log.Fatalf("Error loading TOML config file: %s", err)
	}

	// Mapear las variables del archivo TOML a las variables de constantes
	if err := tomlConfig.Unmarshal(&Config); err != nil {
		log.Fatalf("Error unmarshaling TOML config: %s", err)
	}

	config["RoleAdmin"] = Config.Constants.RoleAdmin
	config["RoleReceptionist"] = Config.Constants.RoleReceptionist
	config["StatusCreated"] = strconv.Itoa(Config.Constants.StatusCreated)
	config["StatusBadRequest"] = strconv.Itoa(Config.Constants.StatusBadRequest)
	config["StatusUnauthorized"] = strconv.Itoa(Config.Constants.StatusUnauthorized)
	config["StatusForbidden"] = strconv.Itoa(Config.Constants.StatusForbidden)
	config["StatusInternalServerError"] = strconv.Itoa(Config.Constants.StatusInternalServerError)
	config["MongoDBURI"] = Config.Constants.MongoDBURI
	config["MongoDBDatabase"] = Config.Constants.MongoDBDatabase
	config["FrontendURL"] = Config.Constants.FrontendURL
	config["CollectionUsers"] = Config.Constants.CollectionUsers
	config["CollectionValidCURPs"] = Config.Constants.CollectionValidCURPs
	config["CollectionClients"] = Config.Constants.CollectionClients
	config["CollectionRooms"] = Config.Constants.CollectionRooms // Nueva colección agregada
	config["JWTSecretKey"] = Config.Constants.JWTSecretKey
	config["ServerAddress"] = Config.Constants.ServerAddress
	config["ServerPort"] = Config.Constants.ServerPort
	config["CloudinaryCloudName"] = Config.Constants.CloudinaryCloudName
	config["CloudinaryAPIKey"] = Config.Constants.CloudinaryAPIKey
	config["CloudinaryAPISecret"] = Config.Constants.CloudinaryAPISecret
	config["GoogleDriveFolderID"] = Config.Constants.GoogleDriveFolderID
	config["GoogleDriveCredentialsPath"] = Config.Constants.GoogleDriveCredentialsPath
	config["LocalFileSystemFolder"] = Config.Constants.LocalFileSystemFolder
	config["StorageSelector"] = Config.Constants.StorageSelector
}

func assignConfigValues(config map[string]string) {
	RoleAdmin = config["RoleAdmin"]
	RoleReceptionist = config["RoleReceptionist"]

	StatusCreated, _ = strconv.Atoi(config["StatusCreated"])
	StatusBadRequest, _ = strconv.Atoi(config["StatusBadRequest"])
	StatusUnauthorized, _ = strconv.Atoi(config["StatusUnauthorized"])
	StatusForbidden, _ = strconv.Atoi(config["StatusForbidden"])
	StatusInternalServerError, _ = strconv.Atoi(config["StatusInternalServerError"])

	MongoDBURI = config["MongoDBURI"]
	MongoDBDatabase = config["MongoDBDatabase"]
	FrontendURL = config["FrontendURL"]

	CollectionUsers = config["CollectionUsers"]
	CollectionValidCURPs = config["CollectionValidCURPs"]
	CollectionClients = config["CollectionClients"]
	CollectionRooms = config["CollectionRooms"] // Nueva colección agregada

	JWTSecretKey = config["JWTSecretKey"]

	ServerAddress = config["ServerAddress"]
	ServerPort = config["ServerPort"]

	// Cloudinary
	CloudinaryCloudName = config["CloudinaryCloudName"]
	CloudinaryAPIKey = config["CloudinaryAPIKey"]
	CloudinaryAPISecret = config["CloudinaryAPISecret"]

	// Google Drive
	GoogleDriveFolderID = config["GoogleDriveFolderID"]
	GoogleDriveCredentialsPath = config["GoogleDriveCredentialsPath"]

	LocalFileSystemFolder = config["LocalFileSystemFolder"]

	StorageSelector = config["StorageSelector"]

	// Inicializar AllCollections con las colecciones definidas individualmente
	AllCollections = []string{
		CollectionUsers,
		CollectionValidCURPs,
		CollectionClients,
		CollectionRooms, // Nueva colección agregada
	}
}

// createDefaultConfig crea un archivo config.toml con valores predeterminados
func createDefaultConfig(filename string) {
	// Definir la estructura del archivo TOML con valores predeterminados
	defaultConfig := `
	[constants]
	RoleAdmin = "Administrator"
	RoleReceptionist = "Receptionist"

	StatusCreated = 201
	StatusBadRequest = 400
	StatusUnauthorized = 401
	StatusForbidden = 403
	StatusInternalServerError = 500

	MongoDBURI = "mongodb://localhost:27017"
	MongoDBDatabase = "testdb"
	FrontendURL = "http://localhost:3000"

	CollectionUsers = "users"
	CollectionValidCURPs = "valid_curps"
	CollectionClients = "clients"
	CollectionRooms = "rooms"  // Nueva colección agregada

	JWTSecretKey = "my_secret_key"

	ServerAddress = "0.0.0.0"
	ServerPort = "8000"

	CloudinaryCloudName = "your-cloudinary-cloud-name"
	CloudinaryAPIKey = "your-cloudinary-api-key"
	CloudinaryAPISecret = "your-cloudinary-api-secret"

	GoogleDriveFolderID = "your-google-drive-folder-id"
	GoogleDriveCredentialsPath = "credentials.json"
	LocalFileSystemFolder = "/uploads"
	StorageSelector = "local"
	`

	// Crear el archivo config.toml con los valores predeterminados
	file, err := os.Create(filename)
	if err != nil {
		log.Fatalf("Error creating default TOML config file: %s", err)
	}
	defer file.Close()

	// Escribir el contenido predeterminado en el archivo
	_, err = file.WriteString(defaultConfig)
	if err != nil {
		log.Fatalf("Error writing default TOML config: %s", err)
	}

	log.Printf("Default TOML config file created: %s\n", filename)
}

// Config representa la estructura del archivo TOML
type ConfigFile struct {
	Constants Constants `toml:"constants"`
}

// Constants contiene todas las constantes definidas en el archivo TOML
type Constants struct {
	RoleAdmin        string `toml:"RoleAdmin"`
	RoleReceptionist string `toml:"RoleReceptionist"`

	StatusCreated             int `toml:"StatusCreated"`
	StatusBadRequest          int `toml:"StatusBadRequest"`
	StatusUnauthorized        int `toml:"StatusUnauthorized"`
	StatusForbidden           int `toml:"StatusForbidden"`
	StatusInternalServerError int `toml:"StatusInternalServerError"`

	MongoDBURI      string `toml:"MongoDBURI"`
	MongoDBDatabase string `toml:"MongoDBDatabase"`
	FrontendURL     string `toml:"FrontendURL"`

	CollectionUsers      string `toml:"CollectionUsers"`
	CollectionValidCURPs string `toml:"CollectionValidCURPs"`
	CollectionClients    string `toml:"CollectionClients"`
	CollectionRooms      string `toml:"CollectionRooms"` // Nueva colección agregada

	JWTSecretKey string `toml:"JWTSecretKey"`

	ServerAddress string `toml:"ServerAddress"`
	ServerPort    string `toml:"ServerPort"`

	CloudinaryCloudName string `toml:"CloudinaryCloudName"`
	CloudinaryAPIKey    string `toml:"CloudinaryAPIKey"`
	CloudinaryAPISecret string `toml:"CloudinaryAPISecret"`

	GoogleDriveFolderID        string `toml:"GoogleDriveFolderID"`
	GoogleDriveCredentialsPath string `toml:"GoogleDriveCredentialsPath"`

	LocalFileSystemFolder string `toml:"/uploads"`
	StorageSelector       string `toml:"local"`
}

// Config es una instancia global de ConfigFile que contiene la configuración cargada
var Config ConfigFile
