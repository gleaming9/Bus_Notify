package routes

import (
	"github.com/gleaming9/Bus_Notify/handlers"
	"net/http"

	"github.com/gin-gonic/gin"
)

// InitRoutes initializes the API routes
func InitRoutes() *gin.Engine {
	router := gin.Default()

	// Health check endpoint
	router.GET("/health", func(c *gin.Context) {
		c.JSON(http.StatusOK, gin.H{
			"status": "Server is running",
		})
	})

	// Endpoint to get bus information
	router.GET("/bus-info", handlers.GetBusInfoHandler)
	router.POST("/alert", handlers.AlertHandler)

	return router
}
