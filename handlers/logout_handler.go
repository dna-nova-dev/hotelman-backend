package handlers

import (
	"net/http"
)

type LogoutHandler struct{}

func (h *LogoutHandler) Handle(w http.ResponseWriter, r *http.Request) {
	// Establecer la cookie con el nombre "Authorize" con un valor vacío y una edad máxima de -1 para eliminarla
	http.SetCookie(w, &http.Cookie{
		Name:     "Authorize",          // Nombre correcto de la cookie
		Value:    "",                   // Valor vacío para eliminar la cookie
		MaxAge:   -1,                   // Edad máxima -1 para eliminar la cookie
		Path:     "/",                  // Path debe coincidir con el path de la cookie original
		HttpOnly: true,                 // Asegura que la cookie solo sea accesible a través de HTTP
		Secure:   false,                // Cambia a true si estás usando HTTPS
		SameSite: http.SameSiteLaxMode, // Ajusta la política SameSite según sea necesario
	})

	// Configurar el encabezado de tipo de contenido
	w.Header().Set("Content-Type", "text/plain")
	// Configurar el código de estado
	w.WriteHeader(http.StatusOK)
	// Mensaje de confirmación
	w.Write([]byte("Logged out"))
}
