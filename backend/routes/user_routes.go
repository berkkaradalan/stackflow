package routes

import "github.com/gin-gonic/gin"

func setupUserRoutes(r *gin.RouterGroup) {
	users := r.Group("/users")

	users.Use()
	{
		users.GET("")
		users.GET("/:id")
		users.PUT("/:id")
		users.DELETE("/:id")
		users.POST("/invite")
	}
}