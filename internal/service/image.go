package service

import (
	"context"
	"errors"
	"mime/multipart"
	"net/http"
	"time"

	"github.com/segmentio/ksuid"

	"Piclo/internal/model"
	"Piclo/internal/repository"
	"Piclo/internal/storage"
)

var mimeToExt = map[string]string{
	"image/jpeg": ".jpg",
	"image/png":  ".png",
	"image/gif":  ".gif",
	"image/webp": ".webp",
}

type ImageService struct {
	repo  *repository.ImageRepository
	store storage.Storage
}

func NewImageService(repo *repository.ImageRepository, store storage.Storage) *ImageService {
	return &ImageService{repo: repo, store: store}
}

func (s *ImageService) ProcessAndUpload(ctx context.Context, file *multipart.FileHeader) (string, error) {
	src, err := file.Open()
	if err != nil {
		return "", err
	}
	defer src.Close()

	// Read first 512 bytes to determine MIME type
	buffer := make([]byte, 512)
	_, err = src.Read(buffer)
	if err != nil {
		return "", err
	}

	contentType := http.DetectContentType(buffer)
	ext, isValid := mimeToExt[contentType]
	if !isValid {
		return "", errors.New("unsupported file type")
	}

	id := ksuid.New().String()
	storageKey := id + ext

	// Reset file pointer to beginning
	_, err = src.Seek(0, 0)
	if err != nil {
		return "", err
	}

	// Upload to MinIO
	err = s.store.Upload(ctx, storageKey, src, file.Size, contentType)
	if err != nil {
		return "", err
	}

	// Calculate expiration time
	expiresAt := time.Now().Add(1 * time.Hour)

	// Save metadata to DB
	img := &model.Image{
		ID:         id,
		Filename:   file.Filename,
		MIMEType:   contentType,
		Size:       file.Size,
		StorageKey: storageKey,
		CreatedAt:  time.Now(),
		ExpiresAt:  &expiresAt,
	}

	err = s.repo.Save(ctx, img)
	if err != nil {
		// Delete file from MinIO if DB fails
		_ = s.store.Delete(ctx, storageKey)
		return "", err
	}

	return id, nil
}

func (s *ImageService) GetImage(ctx context.Context, id string) (*model.Image, error) {
	return s.repo.GetByID(ctx, id)
}
