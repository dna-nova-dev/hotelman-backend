// services/cloudinary.go
package services

import (
	"context"
	"mime/multipart"

	"github.com/cloudinary/cloudinary-go/v2"
	"github.com/cloudinary/cloudinary-go/v2/api/uploader"
)

type CloudinaryService struct {
	Cloudinary *cloudinary.Cloudinary
}

func NewCloudinaryService(cloudinaryURL string) (*CloudinaryService, error) {
	cld, err := cloudinary.NewFromURL(cloudinaryURL)
	if err != nil {
		return nil, err
	}
	return &CloudinaryService{Cloudinary: cld}, nil
}

func (s *CloudinaryService) UploadProfilePicture(file multipart.File, handler *multipart.FileHeader) (string, error) {
	resp, err := s.Cloudinary.Upload.Upload(context.Background(), file, uploader.UploadParams{Folder: "profile_pictures"})
	if err != nil {
		return "", err
	}
	return resp.SecureURL, nil
}

func (s *CloudinaryService) UploadINEPicture(file multipart.File, handler *multipart.FileHeader) (string, error) {
	resp, err := s.Cloudinary.Upload.Upload(context.Background(), file, uploader.UploadParams{Folder: "ine_pictures"})
	if err != nil {
		return "", err
	}
	return resp.SecureURL, nil
}
