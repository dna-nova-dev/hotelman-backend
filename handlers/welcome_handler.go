package handlers

import (
	"encoding/json"
	"net/http"
)

// WelcomeHandler estructura para el manejador de la ruta /welcome
type WelcomeHandler struct{}

// Handle maneja la solicitud para el endpoint /welcome
func (h *WelcomeHandler) Handle(w http.ResponseWriter, r *http.Request) {
	// Extraer todas las cookies de la solicitud
	cookies := r.Cookies()

	// Aquí puedes ajustar cómo obtienes el token, por ejemplo, de las cookies o de la sesión
	var token string
	for _, cookie := range cookies {
		if cookie.Name == "Authorize" { // Cambia "token" al nombre real de la cookie
			token = cookie.Value
			break
		}
	}

	// Configurar el encabezado de tipo de contenido y el código de estado
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)

	// Construir la respuesta JSON con el token
	response := map[string]string{
		"token": token,
	}

	// Codificar el mapa en JSON y enviarlo como respuesta
	if err := json.NewEncoder(w).Encode(response); err != nil {
		http.Error(w, "Failed to encode JSON response", http.StatusInternalServerError)
	}
}
