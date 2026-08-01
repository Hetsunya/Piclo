package handler

import (
	"net/http"

	"github.com/gin-gonic/gin"

	"Piclo/internal/config"
)

func NewRouter(cfg *config.Config) *gin.Engine {

	router := gin.Default()

	// API с версионированием (для бизнес-логики)
	api := router.Group("/api/v1")
	{
		api.GET("/health", healthHandler)
		api.POST("/upload", uploadHandler)
	}

	// Публичные эндпоинты (без префикса /api, так как это прямая отдача контента)
	router.GET("/image/:id", imageHandler)

	return router
}
func healthHandler(c *gin.Context) {
	c.JSON(http.StatusOK, gin.H{"status": "healthy"})
}

func uploadHandler(c *gin.Context) {
	c.JSON(http.StatusNotImplemented, gin.H{"error": "Upload functionality is not implemented yet"})
}

func imageHandler(c *gin.Context) {
	id := c.Param("id")
	c.JSON(http.StatusNotImplemented, gin.H{"error": "Image retrieval is not implemented yet", "requested_id": id})
}
