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

	// Читаем первые 512 байт для определения MIME
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

	// Сбрасываем указатель в начало файла
	_, err = src.Seek(0, 0)
	if err != nil {
		return "", err
	}

	// Загружаем в MinIO
	err = s.store.Upload(ctx, storageKey, src, file.Size, contentType)
	if err != nil {
		return "", err
	}

	// Рассчитываем время жизни
	expiresAt := time.Now().Add(1 * time.Second) // 1 секунда для теста, можно изменить на 24*time.Hour

	// Сохраняем метаданные в БД
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
		// Удаляем файл из MinIO, если БД упала
		_ = s.store.Delete(ctx, storageKey)
		return "", err
	}

	return id, nil
}

func (s *ImageService) GetImage(ctx context.Context, id string) (*model.Image, error) {
	return s.repo.GetByID(ctx, id)
}
