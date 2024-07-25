package middleware

import (
	"context"
	"net/http"
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
		// Obtener el token de la cookie
		cookie, err := r.Cookie("Authorize")
		if err != nil {
			w.WriteHeader(http.StatusUnauthorized)
			w.Write([]byte("Invalid token"))
			return
		}

		tokenString := cookie.Value

		// Parsear y verificar el token JWT
		token, err := jwt.Parse(tokenString, func(token *jwt.Token) (interface{}, error) {
			return ra.jwtKey, nil
		})

		if err != nil {
			if err == jwt.ErrSignatureInvalid {
				w.WriteHeader(http.StatusUnauthorized)
				w.Write([]byte("Token has expired"))
				return
			}
			w.WriteHeader(http.StatusBadRequest)
			return
		}
		if !token.Valid {
			w.WriteHeader(http.StatusUnauthorized)
			return
		}

		// Extraer los claims estándar
		if claims, ok := token.Claims.(jwt.MapClaims); ok && token.Valid {
			// Verificar si el token ha expirado
			if exp, ok := claims["exp"].(float64); ok {
				if time.Unix(int64(exp), 0).Before(time.Now()) {
					w.WriteHeader(http.StatusUnauthorized)
					w.Write([]byte("Token has expired"))
					return
				}
			} else {
				w.WriteHeader(http.StatusUnauthorized)
				w.Write([]byte("Invalid token"))
				return
			}

			// Verificar el rol del usuario
			if role, ok := claims["rol"].(string); ok {
				if !contains(ra.roles, role) {
					w.WriteHeader(http.StatusForbidden)
					w.Write([]byte("Access denied"))
					return
				}
			} else {
				w.WriteHeader(http.StatusUnauthorized)
				w.Write([]byte("Invalid token - Role"))
				return
			}

			// Añadir los claims al contexto de la solicitud
			ctx := context.WithValue(r.Context(), "claims", claims)
			r = r.WithContext(ctx)

			// Continuar con el siguiente handler
			next.ServeHTTP(w, r)
		} else {
			w.WriteHeader(http.StatusUnauthorized)
			w.Write([]byte("Invalid token"))
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
