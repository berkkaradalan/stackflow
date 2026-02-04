package routes

import (
	"net/http"

	"github.com/berkkaradalan/stackflow/handler"
	"github.com/berkkaradalan/stackflow/utils"
	"github.com/gin-gonic/gin"
)

func SetupRouter(jwtManager *utils.JWTManager,authHandler *handler.AuthHandler) *gin.Engine {
	router := gin.New()
	// Add logger and recovery middleware
	router.Use(gin.Logger())
	router.Use(gin.Recovery())
	router.Use(func(c *gin.Context) {
		origin := c.Request.Header.Get("Origin")

		if origin == "http://localhost:3000" || origin == "http://localhost:3001" {
			c.Writer.Header().Set("Access-Control-Allow-Origin", origin)
		}
		c.Writer.Header().Set("Access-Control-Allow-Credentials", "true")
		c.Writer.Header().Set("Access-Control-Allow-Headers", 
		"Content-Type, Content-Length, Accept-Encoding, X-CSRF-Token, Authorization, accept, origin, Cache-Control, X-Requested-With")
		c.Writer.Header().Set("Access-Control-Allow-Methods", "POST, OPTIONS, GET, PUT, DELETE")

		if c.Request.Method == "OPTIONS" {
			c.AbortWithStatus(204)
			return
		}

		c.Next()
	})
	// router.Use(middleware.ErrorHandler())

	// Health endpoints
	router.GET("/health", func(ctx *gin.Context) {
		ctx.JSON(http.StatusOK, gin.H{
			"status":  "healthy",
			"message": "Server is running",
		})
	})
	router.GET("/health/live", func(ctx *gin.Context) {
		ctx.JSON(http.StatusOK, gin.H{
			"status": "live",
		})
	})
	router.GET("/health/ready", func(ctx *gin.Context) {
		ctx.JSON(http.StatusOK, gin.H{
			"status": "ready",
		})
	})

	// Metrics endpoint
	router.GET("/metrics", func(ctx *gin.Context) {
		ctx.JSON(http.StatusOK, gin.H{
			"metrics": "prometheus metrics endpoint",
		})
	})

	// WebSocket endpoint
	router.GET("/ws", func(ctx *gin.Context) {
		ctx.JSON(http.StatusOK, gin.H{
			"message": "WebSocket endpoint",
		})
	})

	api := router.Group("/api")
	{
		setupAuthRoutes(api, authHandler, jwtManager)
		setupUserRoutes(api)
		setupOrganizationRoutes(api)
		setupProjectRoutes(api)
		setupTaskRoutes(api)
		setupAgentsRoutes(api)
		setupWorkflowRoutes(api)
		setupCodeArtifactRoutes(api)
		setupAnalyticsRoutes(api)
	}

	return router
}