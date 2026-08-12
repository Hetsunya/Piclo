package handler

import (
	"errors"
	"fmt"
	"log"
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/jackc/pgx/v5"

	"Piclo/internal/config"
	"Piclo/internal/service"
	"Piclo/internal/storage"
)

func NewRouter(cfg *config.Config, imgService *service.ImageService, store storage.Storage) *gin.Engine {
	router := gin.Default()
	router.SetTrustedProxies(nil)

	// 1. API Endpoints (для React и внешних клиентов)
	api := router.Group("/api/v1")
	{
		api.GET("/health", healthHandler)

		api.POST("/upload", func(c *gin.Context) {
			uploadHandler(c, imgService, cfg.PublicURL)
		})

		// Сырой файл (для тега <img> в React и превью в мессенджерах)
		api.GET("/image/:id/raw", func(c *gin.Context) {
			imageRawHandler(c, imgService, store)
		})
	}

	// 2. UI роутинг полностью делегирован React (Vite проксирует /api/* на Go)
	// Мы больше не делаем редиректы здесь, чтобы не ломать SPA-навигацию.

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
		if err.Error() == "http: request body too large" {
			c.JSON(http.StatusRequestEntityTooLarge, gin.H{"error": "file size exceeds 10MB limit"})
			return
		}
		c.JSON(http.StatusBadRequest, gin.H{"error": "file is required"})
		return
	}

	ctx := c.Request.Context()

	// 🔴 FIX 1: ЖЕСТКАЯ ПРОВЕРКА ОШИБКИ
	id, err := imgService.ProcessAndUpload(ctx, file)
	if err != nil {
		if err.Error() == "unsupported file type" {
			c.JSON(http.StatusUnsupportedMediaType, gin.H{"error": "only jpg, png, gif, webp are allowed"})
			return
		}
		log.Printf("Upload error: %v", err)
		c.JSON(http.StatusInternalServerError, gin.H{"error": "internal processing error"})
		return
	}

	// 🔴 FIX 3: URL теперь будет корректным, если в .env указано PUBLIC_URL=http://localhost:5173
	c.JSON(http.StatusCreated, gin.H{
		"id":  id,
		"url": fmt.Sprintf("%s/image/%s", publicURL, id),
	})
}

func imageRawHandler(c *gin.Context, imgService *service.ImageService, store storage.Storage) {
	id := c.Param("id")
	ctx := c.Request.Context()

	img, err := imgService.GetImage(ctx, id)
	if err != nil || img == nil {
		if errors.Is(err, pgx.ErrNoRows) {
			c.JSON(http.StatusNotFound, gin.H{"error": "Image not found or expired"})
			return
		}
		c.JSON(http.StatusInternalServerError, gin.H{"error": "database error"})
		return
	}

	obj, err := store.GetObject(ctx, img.StorageKey)
	if err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "File not found in storage"})
		return
	}
	defer obj.Close()

	info, err := obj.Stat()
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "failed to get file info"})
		return
	}

	c.DataFromReader(http.StatusOK, info.Size, img.MIMEType, obj, nil)
}
