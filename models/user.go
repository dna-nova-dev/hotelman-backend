package models

import (
	"github.com/dgrijalva/jwt-go"
)

type User struct {
	Nombres        string `json:"nombres" bson:"nombres"`
	Apellidos      string `json:"apellidos" bson:"apellidos"`
	Correo         string `json:"correo" bson:"correo"`
	Celular        string `json:"celular" bson:"celular"`
	Password       string `json:"password" bson:"password"`
	Rol            string `json:"rol" bson:"rol"` // "Administracion" o "Recepcionista"
	CURP           string `json:"curp,omitempty" bson:"curp,omitempty"`
	ProfilePicture string `json:"profilePicture,omitempty" bson:"profilePicture,omitempty"` // URL de la imagen de perfil
}

type Credentials struct {
	Username string `json:"username"`
	Password string `json:"password"`
}

type Claims struct {
	Username string `json:"username"`
	Role     string `json:"rol"`
	TokenID  string `json:"token_id"`
	jwt.StandardClaims
}
