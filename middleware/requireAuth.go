package middleware

import (
	"context"
	"net/http"
	"strings"
	"time"

	"github.com/dgrijalva/jwt-go"
)

type RequireAuth struct {
	jwtKey []byte
	roles  []string
}

func NewRequireAuth(jwtKey []byte, roles []string) *RequireAuth {
	return &RequireAuth{jwtKey: jwtKey, roles: roles}
}

func (ra *RequireAuth) Middleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		var tokenString string

		// Obtener el token de la cookie
		if cookie, err := r.Cookie("Authorize"); err == nil {
			tokenString = cookie.Value
		} else {
			// Obtener el token del encabezado Authorization
			authHeader := r.Header.Get("Authorization")
			if strings.HasPrefix(authHeader, "Bearer ") {
				tokenString = strings.TrimPrefix(authHeader, "Bearer ")
			} else {
				w.WriteHeader(http.StatusUnauthorized)
				w.Write([]byte("Unauthorized - Token missing or malformed"))
				return
			}
		}

		if tokenString == "" {
			w.WriteHeader(http.StatusUnauthorized)
			w.Write([]byte("Unauthorized - Token missing"))
			return
		}

		// Parsear y verificar el token JWT
		token, err := jwt.Parse(tokenString, func(token *jwt.Token) (interface{}, error) {
			return ra.jwtKey, nil
		})

		if err != nil {
			if err == jwt.ErrSignatureInvalid {
				w.WriteHeader(http.StatusUnauthorized)
				w.Write([]byte("Unauthorized - Invalid token signature"))
				return
			}
			w.WriteHeader(http.StatusBadRequest)
			w.Write([]byte("Bad request - Token parsing error"))
			return
		}
		if !token.Valid {
			w.WriteHeader(http.StatusUnauthorized)
			w.Write([]byte("Unauthorized - Token is not valid"))
			return
		}

		// Extraer los claims estándar
		if claims, ok := token.Claims.(jwt.MapClaims); ok && token.Valid {
			// Verificar si el token ha expirado
			if exp, ok := claims["exp"].(float64); ok {
				if time.Unix(int64(exp), 0).Before(time.Now()) {
					w.WriteHeader(http.StatusUnauthorized)
					w.Write([]byte("Unauthorized - Token has expired"))
					return
				}
			} else {
				w.WriteHeader(http.StatusUnauthorized)
				w.Write([]byte("Unauthorized - Token missing expiration"))
				return
			}

			// Verificar el rol del usuario
			if role, ok := claims["rol"].(string); ok {
				if !contains(ra.roles, role) {
					w.WriteHeader(http.StatusForbidden)
					w.Write([]byte("Forbidden - Access denied"))
					return
				}
			} else {
				w.WriteHeader(http.StatusUnauthorized)
				w.Write([]byte("Unauthorized - Token missing role"))
				return
			}

			// Añadir los claims al contexto de la solicitud
			ctx := context.WithValue(r.Context(), "claims", claims)
			r = r.WithContext(ctx)

			// Continuar con el siguiente handler
			next.ServeHTTP(w, r)
		} else {
			w.WriteHeader(http.StatusUnauthorized)
			w.Write([]byte("Unauthorized - Token claims error"))
			return
		}
	})
}

// Helper function to check if a slice contains a specific string
func contains(slice []string, item string) bool {
	for _, v := range slice {
		if v == item {
			return true
		}
	}
	return false
}
