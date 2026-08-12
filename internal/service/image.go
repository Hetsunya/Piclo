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
)

// Маппинг разрешенных MIME-типов в правильные расширения
var mimeToExt = map[string]string{
	"image/jpeg": ".jpg",
	"image/png":  ".png",
	"image/gif":  ".gif",
	"image/webp": ".webp",
}

type ImageService struct {
	repo *repository.ImageRepository
}

func NewImageService(repo *repository.ImageRepository) *ImageService {
	return &ImageService{repo: repo}
}

// Process проверяет файл и возвращает безопасные метаданные для сохранения
func (s *ImageService) Process(file *multipart.FileHeader) (id string, storageKey string, mimeType string, size int64, err error) {
	src, err := file.Open()
	if err != nil {
		return "", "", "", 0, err
	}
	defer src.Close()

	// Читаем первые 512 байт для определения реального MIME
	buffer := make([]byte, 512)
	_, err = src.Read(buffer)
	if err != nil {
		return "", "", "", 0, err
	}

	contentType := http.DetectContentType(buffer)

	// Проверяем, есть ли такой MIME в нашем белом списке
	ext, isValid := mimeToExt[contentType]
	if !isValid {
		return "", "", "", 0, errors.New("unsupported file type")
	}

	id = ksuid.New().String()
	storageKey = id + ext // Теперь расширение гарантированно правильное

	return id, storageKey, contentType, file.Size, nil
}

// SaveMetadata записывает данные о файле в PostgreSQL
func (s *ImageService) SaveMetadata(ctx context.Context, id, originalFilename, mimeType, storageKey string, size int64) error {
	img := &model.Image{
		ID:         id,
		Filename:   originalFilename,
		MIMEType:   mimeType,
		Size:       size,
		StorageKey: storageKey,
		CreatedAt:  time.Now(),
		ExpiresAt:  nil, // Пока без TTL, можно добавить позже
	}
	return s.repo.Save(ctx, img)
}

// GetImage получает метаданные картинки из БД
func (s *ImageService) GetImage(ctx context.Context, id string) (*model.Image, error) {
	return s.repo.GetByID(ctx, id)
}
