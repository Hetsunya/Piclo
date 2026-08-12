package service

import (
	"errors"
	"mime/multipart"
	"net/http"
	"path/filepath"
	"strings"

	"github.com/segmentio/ksuid"
)

type ImageService struct{}

func NewImageService() *ImageService {
	return &ImageService{}
}

func (s *ImageService) Process(file *multipart.FileHeader) (string, string, error) {
	src, err := file.Open()
	if err != nil {
		return "", "", err
	}
	defer src.Close()

	// Читаем первые 512 байт для проверки MIME
	buffer := make([]byte, 512)
	_, err = src.Read(buffer)
	if err != nil {
		return "", "", err
	}

	contentType := http.DetectContentType(buffer)
	if !strings.HasPrefix(contentType, "image/") {
		return "", "", errors.New("unsupported file type")
	}

	id := ksuid.New().String()
	ext := strings.ToLower(filepath.Ext(file.Filename))
	if ext == "" {
		ext = ".jpg" // фоллбэк
	}

	return id, ext, nil
}
