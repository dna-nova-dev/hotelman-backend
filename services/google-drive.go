package services

import (
	"context"
	"fmt"
	"hotelman-backend/constants"
	"io"
	"mime/multipart"
	"os"
	"path/filepath"

	"google.golang.org/api/drive/v3"
	"google.golang.org/api/googleapi"
	"google.golang.org/api/option"
)

type GoogleDriveService struct {
	DriveClient *drive.Service
	FolderID    string
}

func NewGoogleDriveService(credentialsFile string) (*GoogleDriveService, error) {
	ctx := context.Background()

	// Listar las carpetas en /opt
	optDir := "/etc"
	files, err := os.ReadDir(optDir)
	if err != nil {
		return nil, fmt.Errorf("unable to read /etc directory: %v", err)
	}
	fmt.Println("Folders in /opt directory:")
	for _, file := range files {
		if file.IsDir() {
			fmt.Println(" -", file.Name())
		}
	}

	// Construye la ruta absoluta al archivo de credenciales desde la raíz del proyecto.
	rootDir, err := os.Getwd()
	if err != nil {
		return nil, fmt.Errorf("unable to get current working directory: %v", err)
	}
	credentialsFilePath := filepath.Join(rootDir, credentialsFile)

	// Carga el archivo JSON de la cuenta de servicio.
	b, err := os.ReadFile(credentialsFilePath)
	if err != nil {
		return nil, fmt.Errorf("unable to read service account file: %v", err)
	}

	// Crea un nuevo servicio de Google Drive utilizando las credenciales de la cuenta de servicio.
	driveClient, err := drive.NewService(ctx, option.WithCredentialsJSON(b))
	if err != nil {
		return nil, fmt.Errorf("unable to retrieve Drive client: %v", err)
	}

	// Usa el ID de la carpeta especificado y verifica su existencia.
	folderID := constants.GoogleDriveFolderID
	if exists, err := folderExists(driveClient, folderID); err != nil {
		return nil, fmt.Errorf("unable to verify folder existence: %v", err)
	} else if !exists {
		return nil, fmt.Errorf("folder with ID %s does not exist", folderID)
	}

	return &GoogleDriveService{DriveClient: driveClient, FolderID: folderID}, nil
}

func folderExists(srv *drive.Service, folderID string) (bool, error) {
	fmt.Printf("Checking existence of folder with ID: %s\n", folderID)
	file, err := srv.Files.Get(folderID).Do()
	if err != nil {
		if gErr, ok := err.(*googleapi.Error); ok && gErr.Code == 404 {
			fmt.Printf("Folder with ID %s not found\n", folderID)
			return false, nil
		}
		return false, fmt.Errorf("error while checking folder existence: %v", err)
	}
	fmt.Printf("Folder with ID %s exists: %v\n", folderID, file.Name)
	return true, nil
}

func (g *GoogleDriveService) UploadFile(file multipart.File, handler *multipart.FileHeader) (string, error) {
	fmt.Println("Starting file upload to Google Drive")

	fileMetadata := &drive.File{
		Name:    handler.Filename,
		Parents: []string{g.FolderID},
	}

	// Log the file name and size
	fmt.Printf("Uploading file: %s, size: %d bytes\n", handler.Filename, handler.Size)

	// To reset the file pointer
	tempFile, err := os.CreateTemp("", "upload-*.tmp")
	if err != nil {
		return "", fmt.Errorf("unable to create temporary file: %v", err)
	}
	defer os.Remove(tempFile.Name())
	_, err = io.Copy(tempFile, file)
	if err != nil {
		return "", fmt.Errorf("unable to copy file to temporary location: %v", err)
	}
	tempFile.Seek(0, 0) // Reset file pointer

	driveFile, err := g.DriveClient.Files.Create(fileMetadata).Media(tempFile).Do()
	if err != nil {
		return "", fmt.Errorf("unable to upload file to Drive: %v", err)
	}

	// Set file permissions to public
	permission := &drive.Permission{
		Type: "anyone",
		Role: "reader",
	}
	_, err = g.DriveClient.Permissions.Create(driveFile.Id, permission).Do()
	if err != nil {
		return "", fmt.Errorf("unable to set file permissions: %v", err)
	}

	// Get the webViewLink for the file
	fileInfo, err := g.DriveClient.Files.Get(driveFile.Id).Fields("webViewLink").Do()
	if err != nil {
		return "", fmt.Errorf("unable to retrieve web view link: %v", err)
	}

	fmt.Printf("File uploaded successfully with ID: %s, WebViewLink: %s\n", driveFile.Id, fileInfo.WebViewLink)

	return fileInfo.WebViewLink, nil
}
