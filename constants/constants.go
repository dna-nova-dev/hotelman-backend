package constants

import (
	"log"
	"os"

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

	// Collections
	CollectionUsers      string
	CollectionValidCURPs string

	// JWT
	JWTSecretKey string

	// Configuración de red
	ServerAddress string
	ServerPort    string

	// AllCollections contiene todos los nombres de colecciones definidos
	AllCollections []string
)

func init() {
	configFile := "config.toml"

	// Verificar si el archivo config.toml ya existe
	if _, err := os.Stat(configFile); os.IsNotExist(err) {
		// Si no existe, crear un archivo config.toml con valores predeterminados
		createDefaultConfig(configFile)
	}

	// Abrir y leer el archivo TOML
	config, err := toml.LoadFile(configFile)
	if err != nil {
		log.Fatalf("Error loading TOML config file: %s", err)
	}

	// Mapear las variables del archivo TOML a las variables de constantes
	if err := config.Unmarshal(&Config); err != nil {
		log.Fatalf("Error unmarshaling TOML config: %s", err)
	}

	// Asignar valores a las variables de constantes
	RoleAdmin = Config.Constants.RoleAdmin
	RoleReceptionist = Config.Constants.RoleReceptionist

	StatusCreated = Config.Constants.StatusCreated
	StatusBadRequest = Config.Constants.StatusBadRequest
	StatusUnauthorized = Config.Constants.StatusUnauthorized
	StatusForbidden = Config.Constants.StatusForbidden
	StatusInternalServerError = Config.Constants.StatusInternalServerError

	MongoDBURI = Config.Constants.MongoDBURI
	MongoDBDatabase = Config.Constants.MongoDBDatabase

	CollectionUsers = Config.Constants.CollectionUsers
	CollectionValidCURPs = Config.Constants.CollectionValidCURPs

	JWTSecretKey = Config.Constants.JWTSecretKey

	ServerAddress = Config.Constants.ServerAddress
	ServerPort = Config.Constants.ServerPort

	// Inicializar AllCollections con las colecciones definidas individualmente
	AllCollections = []string{
		CollectionUsers,
		CollectionValidCURPs,
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

	CollectionUsers = "users"
	CollectionValidCURPs = "valid_curps"

	JWTSecretKey = "my_secret_key"

	ServerAddress = "0.0.0.0"
	ServerPort = "8000"
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

	CollectionUsers      string `toml:"CollectionUsers"`
	CollectionValidCURPs string `toml:"CollectionValidCURPs"`

	JWTSecretKey string `toml:"JWTSecretKey"`

	ServerAddress string `toml:"ServerAddress"`
	ServerPort    string `toml:"ServerPort"`
}

// Config contiene las constantes del archivo TOML
var Config ConfigFile
