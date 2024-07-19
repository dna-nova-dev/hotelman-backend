package models

import (
	"github.com/dgrijalva/jwt-go"
)

type User struct {
	Nombres   string `json:"nombres"`
	Apellidos string `json:"apellidos"`
	Correo    string `json:"correo"`
	Celular   string `json:"celular"`
	Password  string `json:"password"`
	Rol       string `json:"rol"` // "Administracion" o "Recepcionista"
	CURP      string `json:"curp,omitempty"`
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
