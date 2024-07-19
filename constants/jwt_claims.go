package constants

import (
	"github.com/dgrijalva/jwt-go"
)

type Claims struct {
	Username string `json:"username"`
	Role     string `json:"role"`
	TokenID  string `json:"token_id"`
	jwt.StandardClaims
}
