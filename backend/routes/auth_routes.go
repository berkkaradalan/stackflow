package routes

import (
	"github.com/berkkaradalan/stackflow/handler"
	"github.com/berkkaradalan/stackflow/middleware"
	"github.com/berkkaradalan/stackflow/utils"
	"github.com/gin-gonic/gin"
)

func setupAuthRoutes(
	r *gin.RouterGroup, authHandler *handler.AuthHandler, jwtManager *utils.JWTManager) {
	auth := r.Group("/auth")

	// Public
	auth.POST("/login", authHandler.Login)

	// Protected
	auth.Use(middleware.AuthMiddleware(jwtManager))
	{
		auth.POST("/logout")
		auth.POST("/refresh")
		auth.GET("/profile")
		auth.PUT("/profile")
	}
}