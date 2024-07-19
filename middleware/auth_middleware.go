package middleware

import (
	"context"
	"fmt"
	"net/http"
	"strings"

	"hotelman-backend/constants"

	"github.com/dgrijalva/jwt-go"
)

func AuthMiddleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		authHeader := r.Header.Get("Authorization")
		if authHeader == "" {
			w.WriteHeader(http.StatusUnauthorized)
			return
		}

		tokenString := strings.Split(authHeader, "Bearer ")[1]
		token, err := jwt.ParseWithClaims(tokenString, &constants.Claims{}, func(token *jwt.Token) (interface{}, error) {
			return []byte(constants.JWTSecretKey), nil
		})

		if err != nil {
			w.WriteHeader(http.StatusUnauthorized)
			fmt.Fprintf(w, "Error parsing JWT: %v", err)
			return
		}

		claims, ok := token.Claims.(*constants.Claims)
		if !ok || !token.Valid {
			w.WriteHeader(http.StatusUnauthorized)
			fmt.Fprint(w, "Token no válido")
			return
		}

		ctx := context.WithValue(r.Context(), "claims", claims)
		next.ServeHTTP(w, r.WithContext(ctx))
	})
}
