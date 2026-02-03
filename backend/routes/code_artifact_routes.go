package routes

import "github.com/gin-gonic/gin"

func setupCodeArtifactRoutes(r *gin.RouterGroup) {
	tasks := r.Group("/tasks/:id/code")

	tasks.Use()
	{
		tasks.GET("")
		tasks.GET("/diff")
		tasks.GET("/files")
		tasks.GET("/files/*path")
		tasks.POST("/review")
	}
}
