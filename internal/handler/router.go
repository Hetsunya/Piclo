package handler

import (
	"fmt"
	"net/http"

	"github.com/gin-gonic/gin"

	"Piclo/internal/config"
	"Piclo/internal/service"
	"Piclo/internal/storage"
)

func NewRouter(cfg *config.Config, imgService *service.ImageService, store storage.Storage) *gin.Engine {
	router := gin.Default()
	router.SetTrustedProxies(nil)

	api := router.Group("/api/v1")
	{
		api.GET("/health", healthHandler)
		api.POST("/upload", func(c *gin.Context) {
			uploadHandler(c, imgService, cfg.PublicURL)
		})
	}

	router.GET("/image/:id", func(c *gin.Context) {
		imageHandler(c, imgService, store)
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
		if err.Error() == "http: request body too large" {
			c.JSON(http.StatusRequestEntityTooLarge, gin.H{"error": "file size exceeds 10MB limit"})
			return
		}
		c.JSON(http.StatusBadRequest, gin.H{"error": "file is required"})
		return
	}

	ctx := c.Request.Context()
	id, err := imgService.ProcessAndUpload(ctx, file)

	c.JSON(http.StatusCreated, gin.H{
		"id":  id,
		"url": fmt.Sprintf("%s/image/%s", publicURL, id),
	})
}

func imageHandler(c *gin.Context, imgService *service.ImageService, store storage.Storage) {
	id := c.Param("id")

	ctx := c.Request.Context()
	img, err := imgService.GetImage(ctx, id)

	// Получаем объект из MinIO
	obj, err := store.GetObject(ctx, img.StorageKey)
	if err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "File not found in storage"})
		return
	}
	defer obj.Close()

	// Получаем размер объекта через Stat()
	info, err := obj.Stat()
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "failed to get file info"})
		return
	}

	// Отдаем поток напрямую из MinIO клиенту
	c.DataFromReader(http.StatusOK, info.Size, img.MIMEType, obj, nil)
}
