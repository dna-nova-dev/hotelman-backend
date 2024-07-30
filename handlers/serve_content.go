package handlers

import (
	"fmt"
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

	// Log para debuggear los valores de folder y filename
	fmt.Printf("Requested folder: %s\n", folder)
	fmt.Printf("Requested filename: %s\n", filename)

	// Validar que el folder sea "images" o "documents"
	if folder != "images" && folder != "documents" {
		http.Error(w, "Invalid folder", http.StatusBadRequest)
		return
	}

	// Construir la ruta completa al archivo
	filePath := filepath.Join(h.UploadsDir, folder, filename)

	// Log para debuggear la ruta completa del archivo
	fmt.Printf("Constructed file path: %s\n", filePath)

	// Verificar que el archivo existe y es accesible
	if _, err := os.Stat(filePath); os.IsNotExist(err) {
		// Log para debuggear si el archivo no existe
		fmt.Printf("File not found: %s\n", filePath)
		http.Error(w, "File not found", http.StatusNotFound)
		return
	} else if err != nil {
		// Log para debuggear cualquier otro error al acceder al archivo
		fmt.Printf("Error accessing file: %v\n", err)
		http.Error(w, "Error accessing file", http.StatusInternalServerError)
		return
	}

	// Log para debuggear que el archivo se está sirviendo
	fmt.Printf("Serving file: %s\n", filePath)
	http.ServeFile(w, r, filePath)
}
