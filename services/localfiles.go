package services

import (
	"fmt"
	"io"
	"mime/multipart"
	"os"
	"path/filepath"
)

// LocalFileSystemService es un servicio para manejar archivos en el sistema de archivos local
type LocalFileSystemService struct {
	BasePath string
}

// / NewLocalFileSystemService crea una nueva instancia de LocalFileSystemService
func NewLocalFileSystemService(relativePath string) (*LocalFileSystemService, error) {
	// Obtener la ruta raíz del programa
	rootPath, err := os.Getwd()
	if err != nil {
		return nil, fmt.Errorf("unable to get root directory: %v", err)
	}

	// Construir la ruta absoluta a partir de la raíz del proyecto
	basePath := filepath.Join(rootPath, relativePath)

	// Asegúrate de que la carpeta base exista
	err = os.MkdirAll(basePath, os.ModePerm)
	if err != nil {
		return nil, fmt.Errorf("unable to create base directory: %v", err)
	}

	return &LocalFileSystemService{BasePath: basePath}, nil
}

// UploadFilePDF maneja la carga de archivos PDF al sistema de archivos local
func (l *LocalFileSystemService) UploadFilePDF(file multipart.File, handler *multipart.FileHeader) (string, error) {
	fmt.Println("Starting PDF file upload to local file system")

	// Verifica que el archivo sea un PDF
	if filepath.Ext(handler.Filename) != ".pdf" {
		return "", fmt.Errorf("file is not a PDF")
	}

	// Crea la ruta completa del archivo
	filePath := filepath.Join(l.BasePath+"/documents", handler.Filename)

	// Crea el archivo en el sistema de archivos local
	dst, err := os.Create(filePath)
	if err != nil {
		return "", fmt.Errorf("unable to create file: %v", err)
	}
	defer dst.Close()

	// Copia el contenido del archivo cargado al nuevo archivo en el sistema de archivos local
	_, err = io.Copy(dst, file)
	if err != nil {
		return "", fmt.Errorf("unable to copy file content: %v", err)
	}

	fmt.Printf("PDF file uploaded successfully: %s\n", filePath)
	return filePath, nil
}

// UploadFileImage maneja la carga de archivos de imagen al sistema de archivos local
func (l *LocalFileSystemService) UploadFileImage(file multipart.File, handler *multipart.FileHeader) (string, error) {
	fmt.Println("Starting image file upload to local file system")

	// Verifica que el archivo sea una imagen (opcional, puedes agregar más tipos de archivos de imagen si lo deseas)
	ext := filepath.Ext(handler.Filename + "/images")
	if ext != ".jpg" && ext != ".jpeg" && ext != ".png" && ext != ".gif" {
		return "", fmt.Errorf("file is not a supported image format")
	}

	// Crea la ruta completa del archivo
	filePath := filepath.Join(l.BasePath, handler.Filename)

	// Crea el archivo en el sistema de archivos local
	dst, err := os.Create(filePath)
	if err != nil {
		return "", fmt.Errorf("unable to create file: %v", err)
	}
	defer dst.Close()

	// Copia el contenido del archivo cargado al nuevo archivo en el sistema de archivos local
	_, err = io.Copy(dst, file)
	if err != nil {
		return "", fmt.Errorf("unable to copy file content: %v", err)
	}

	fmt.Printf("Image file uploaded successfully: %s\n", filePath)
	return filePath, nil
}
