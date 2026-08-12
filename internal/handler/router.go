package handler

import (
	"fmt"
	"net/http"
	"os"
	"path/filepath"

	"github.com/gin-gonic/gin"

	"Piclo/internal/config"
	"Piclo/internal/service"
)

func NewRouter(cfg *config.Config) *gin.Engine {
	router := gin.Default()
	router.SetTrustedProxies(nil)

	imgService := service.NewImageService()

	api := router.Group("/api/v1")
	{
		api.GET("/health", healthHandler)
		api.POST("/upload", func(c *gin.Context) {
			uploadHandler(c, imgService, cfg.Port)
		})
	}

	router.GET("/image/:id", func(c *gin.Context) {
		imageHandler(c)
	})

	return router
}

func healthHandler(c *gin.Context) {
	c.JSON(http.StatusOK, gin.H{"status": "healthy"})
}

func uploadHandler(c *gin.Context, imgService *service.ImageService, port string) {
	const maxSize = 10 << 20 // 10MB
	c.Request.Body = http.MaxBytesReader(c.Writer, c.Request.Body, maxSize)

	file, err := c.FormFile("file")
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "file is required"})
		return
	}

	id, ext, err := imgService.Process(file)
	if err != nil {
		c.JSON(http.StatusUnsupportedMediaType, gin.H{"error": "only images are allowed"})
		return
	}

	// Создаем папку storage, если её нет
	os.MkdirAll("storage", os.ModePerm)

	filename := id + ext
	destPath := filepath.Join("storage", filename)

	if err := c.SaveUploadedFile(file, destPath); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "failed to save file"})
		return
	}

	c.JSON(http.StatusCreated, gin.H{
		"id":  id,
		"url": fmt.Sprintf("http://localhost:%s/image/%s", port, id),
	})
}

func imageHandler(c *gin.Context) {
	id := c.Param("id")

	// Ищем файл по маске storage/id.* (так как в URL только ID без расширения)
	matches, err := filepath.Glob(filepath.Join("storage", id+".*"))
	if err != nil || len(matches) == 0 {
		c.JSON(http.StatusNotFound, gin.H{"error": "Image not found"})
		return
	}

	// Отдаем первый найденный файл (c.File принимает строку пути)
	c.File(matches[0])
}
