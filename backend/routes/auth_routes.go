package routes

import (
	"github.com/gin-gonic/gin"
)

func setupUserRoutes(r *gin.RouterGroup) {
	users := r.Group("/users")

	// Public
	users.POST("/login")

	// Protected
	users.Use()
	{
		users.POST("/logout")
		users.GET("/profile")
		users.PUT("/profile")
	}
}