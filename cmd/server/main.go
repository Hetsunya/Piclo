package main

import (
	"log"
	"net/http"

	"github.com/gin-gonic/gin"

)

func main() {
	router := gin.Default()

	router.GET("/health", func(c *gin.Context) {
		c.JSON(http.StatusOK, gin.H{
			"status": "healthy",
		})
	})

	router.POST("/upload", func(c *gin.Context) {
		c.JSON(http.StatusOK, gin.H{
			"message": "Not implemented",
		})
	})

	router.GET("/image/:id", func(c *gin.Context) {
		c.JSON(http.StatusOK, gin.H{
			"message": "Not implemented",
		})
	})

	log.Fatal(router.Run(":9090"))
}
