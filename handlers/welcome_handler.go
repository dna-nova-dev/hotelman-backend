package handlers

import (
	"hotelman-backend/models"
	"net/http"

	"github.com/dgrijalva/jwt-go"
)

type WelcomeHandler struct {
	jwtKey []byte // Agregar jwtKey como un campo de la estructura
}

func (h *WelcomeHandler) Handle(w http.ResponseWriter, r *http.Request) {
	cookie, err := r.Cookie("token")
	if err != nil {
		if err == http.ErrNoCookie {
			w.WriteHeader(http.StatusUnauthorized)
			return
		}
		w.WriteHeader(http.StatusBadRequest)
		return
	}

	tokenStr := cookie.Value

	claims := &models.Claims{}
	tkn, err := jwt.ParseWithClaims(tokenStr, claims, func(token *jwt.Token) (interface{}, error) {
		return h.jwtKey, nil // Acceder a jwtKey a través de la instancia h
	})
	if err != nil {
		if err == jwt.ErrSignatureInvalid {
			w.WriteHeader(http.StatusUnauthorized)
			return
		}
		w.WriteHeader(http.StatusBadRequest)
		return
	}
	if !tkn.Valid {
		w.WriteHeader(http.StatusUnauthorized)
		return
	}

	w.Write([]byte("Welcome " + claims.Username))
}
