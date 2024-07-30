package handlers

import (
	"net/http"
	"os"
	"path/filepath"

	"github.com/gorilla/mux"
)

type ServeFileHandler struct {
	UploadsDir string
}

func (h *ServeFileHandler) Handle(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	folder := vars["folder"]
	filename := vars["filename"]

	// Validar que el folder sea "images" o "documents"
	if folder != "images" && folder != "documents" {
		http.Error(w, "Invalid folder", http.StatusBadRequest)
		return
	}

	// Construir la ruta completa al archivo
	filePath := filepath.Join(h.UploadsDir, folder, filename)

	// Verificar que el archivo existe y es accesible
	if _, err := os.Stat(filePath); os.IsNotExist(err) {
		http.Error(w, "File not found", http.StatusNotFound)
		return
	}

	http.ServeFile(w, r, filePath)
}
