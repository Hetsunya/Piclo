package handler

import (
	"context"
	"fmt"
	"log"
	"net/http"
	"os"
	"path/filepath"

	"github.com/gin-gonic/gin"
	"github.com/jackc/pgx/v5"

	"Piclo/internal/config"
	"Piclo/internal/service"
)

func NewRouter(cfg *config.Config, imgService *service.ImageService) *gin.Engine {
	router := gin.Default()
	router.SetTrustedProxies(nil)

	// Создаем папку storage при старте, если её нет
	if err := os.MkdirAll("storage", 0755); err != nil {
		panic(fmt.Sprintf("failed to create storage dir: %v", err))
	}

	api := router.Group("/api/v1")
	{
		api.GET("/health", healthHandler)
		api.POST("/upload", func(c *gin.Context) {
			uploadHandler(c, imgService, cfg.PublicURL)
		})
	}

	router.GET("/image/:id", func(c *gin.Context) {
		imageHandler(c, imgService)
	})

	return router
}

func healthHandler(c *gin.Context) {
	c.JSON(http.StatusOK, gin.H{"status": "healthy"})
}

func uploadHandler(c *gin.Context, imgService *service.ImageService, publicURL string) {
	const maxSize = 10 << 20 // 10MB
	c.Request.Body = http.MaxBytesReader(c.Writer, c.Request.Body, maxSize)

	file, err := c.FormFile("file")
	if err != nil {
		// Если ошибка из-за превышения размера, MaxBytesReader вернет специфичную ошибку
		if err.Error() == "http: request body too large" {
			c.JSON(http.StatusRequestEntityTooLarge, gin.H{"error": "file size exceeds 10MB limit"})
			return
		}
		c.JSON(http.StatusBadRequest, gin.H{"error": "file is required"})
		return
	}

	// 1. Процессим файл (проверка MIME, генерация ID и правильного расширения)
	id, storageKey, mimeType, size, err := imgService.Process(file)
	if err != nil {
		if err.Error() == "unsupported file type" {
			c.JSON(http.StatusUnsupportedMediaType, gin.H{"error": "only jpg, png, gif, webp are allowed"})
			return
		}
		c.JSON(http.StatusInternalServerError, gin.H{"error": "internal processing error"})
		return
	}

	// 2. Сохраняем физический файл на диск
	destPath := filepath.Join("storage", storageKey)
	if err := c.SaveUploadedFile(file, destPath); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "failed to save file"})
		return
	}

	// 3. Сохраняем метаданные в PostgreSQL
	ctx := context.Background()
	if err := imgService.SaveMetadata(ctx, id, file.Filename, mimeType, storageKey, size); err != nil {
		log.Printf("!!! РЕАЛЬНАЯ ОШИБКА БД: %v !!!", err) // <-- ВСТАВИТЬ РОВНО ЭТУ СТРОКУ
		c.JSON(http.StatusInternalServerError, gin.H{"error": "failed to save metadata to database"})
		return
	}

	// 4. Возвращаем успешный ответ
	c.JSON(http.StatusCreated, gin.H{
		"id":  id,
		"url": fmt.Sprintf("%s/image/%s", publicURL, id),
	})
}

func imageHandler(c *gin.Context, imgService *service.ImageService) {
	id := c.Param("id")

	// 1. Ищем запись в БД (никаких filepath.Glob!)
	ctx := context.Background()
	img, err := imgService.GetImage(ctx, id)
	if err != nil {
		if err == pgx.ErrNoRows {
			c.JSON(http.StatusNotFound, gin.H{"error": "Image not found"})
			return
		}
		c.JSON(http.StatusInternalServerError, gin.H{"error": "database error"})
		return
	}

	// 2. Формируем путь к файлу на основе storageKey из БД
	filePath := filepath.Join("storage", img.StorageKey)

	// 3. Отдаем файл. Gin сам определит Content-Type по расширению файла и отдаст его
	c.File(filePath)
}
