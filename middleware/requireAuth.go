package middleware

import (
	"context"
	"net/http"
	"time"

	"github.com/dgrijalva/jwt-go"
)

type RequireAuth struct {
	jwtKey []byte
}

func NewRequireAuth(jwtKey []byte) *RequireAuth {
	return &RequireAuth{jwtKey: jwtKey}
}

func (ra *RequireAuth) Middleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		// Obtener el token de la cookie
		cookie, err := r.Cookie("Autorization")
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
			// Verificar que el rol del usuario sea "Administrador"
			if role, ok := claims["rol"].(string); ok {
				if role != "Administracion" {
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
