package config

import (
	"context"
	"log"

	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"

	"hotelman-backend/constants"
)

// ConnectDB establece una conexión a MongoDB y verifica/crea la base de datos y colecciones necesarias
func ConnectDB() *mongo.Client {
	clientOptions := options.Client().ApplyURI(constants.MongoDBURI)
	client, err := mongo.Connect(context.TODO(), clientOptions)
	if err != nil {
		log.Fatal(err)
	}

	// Verificar y crear la base de datos si no existe
	err = ensureDatabase(client, constants.MongoDBDatabase)
	if err != nil {
		log.Fatal(err)
	}

	// Verificar y crear las colecciones si no existen
	err = ensureCollections(client)
	if err != nil {
		log.Fatal(err)
	}

	log.Println("Connected to MongoDB!")
	return client
}

// ensureDatabase verifica si la base de datos existe y la crea si es necesario
func ensureDatabase(client *mongo.Client, dbName string) error {
	// Obtener la base de datos
	db := client.Database(dbName)

	// Verificar si la base de datos ya existe
	var result bson.D
	err := db.RunCommand(context.Background(), bson.D{
		{"listCollections", 1},
		{"nameOnly", true},
		{"filter", bson.D{{"name", dbName}}},
	}).Decode(&result)

	if err == nil {
		// La base de datos ya existe si no hay error en la decodificación
		log.Printf("Database '%s' already exists.\n", dbName)
		return nil
	} else if err.Error() != "ns not found" {
		// Si el error no es "ns not found", retornar el error
		return err
	}

	// Si la base de datos no existe, crearla
	createResult := bson.M{}
	createErr := db.RunCommand(context.Background(), bson.D{
		{"create", dbName},
	}).Decode(&createResult)
	if createErr != nil {
		return createErr
	}

	// Verificar si la creación fue exitosa
	if createResult["ok"] == 1.0 {
		log.Printf("Database '%s' created.\n", dbName)
		return nil
	}

	return err
}

// ensureCollections verifica si las colecciones definidas en constants.AllCollections existen y las crea si es necesario
func ensureCollections(client *mongo.Client) error {
	db := client.Database(constants.MongoDBDatabase)

	// Verificar y crear colecciones según las constantes definidas
	for _, collectionName := range constants.AllCollections {
		// Verificar si la colección ya existe
		collectionExists, err := collectionExists(db, collectionName)
		if err != nil {
			return err
		}

		// Si la colección no existe, crearla
		if !collectionExists {
			err := createCollection(db, collectionName)
			if err != nil {
				return err
			}
			log.Printf("Collection '%s' created.\n", collectionName)
		}
	}

	return nil
}

// collectionExists verifica si una colección existe en la base de datos
func collectionExists(db *mongo.Database, collectionName string) (bool, error) {
	collections, err := db.ListCollectionNames(context.Background(), bson.D{})
	if err != nil {
		return false, err
	}

	for _, col := range collections {
		if col == collectionName {
			return true, nil
		}
	}

	return false, nil
}

// createCollection crea una nueva colección en la base de datos
func createCollection(db *mongo.Database, collectionName string) error {
	err := db.CreateCollection(context.Background(), collectionName)
	if err != nil {
		return err
	}
	return nil
}
