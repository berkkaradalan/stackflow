package routes

import "github.com/gin-gonic/gin"

func setupProjectRoutes(r *gin.RouterGroup) {
	projects := r.Group("/projects")

	projects.Use()
	{
		projects.GET("")
		projects.POST("")
		projects.GET("/:id")
		projects.PUT("/:id")
		projects.DELETE("/:id")
		projects.GET("/:id/stats")
	}
}