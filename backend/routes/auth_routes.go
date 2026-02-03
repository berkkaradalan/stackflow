package routes

import (
	"github.com/gin-gonic/gin"
)

func setupAuthRoutes(r *gin.RouterGroup) {
	auth := r.Group("/auth")

	// Public
	auth.POST("/login")

	// Protected
	auth.Use()
	{
		auth.POST("/logout")
		auth.POST("/refresh")
		auth.GET("/profile")
		auth.PUT("/profile")
	}
}