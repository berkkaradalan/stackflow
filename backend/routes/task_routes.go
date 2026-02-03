package routes

import "github.com/gin-gonic/gin"

func setupTaskRoutes(r *gin.RouterGroup) {
	// Tasks endpoints under projects
	projects := r.Group("/projects/:id")
	{
		projects.GET("/tasks")
		projects.POST("/tasks")
	}

	// Individual task endpoints
	tasks := r.Group("/tasks")
	{
		tasks.GET("/:id")
		tasks.PUT("/:id")
		tasks.DELETE("/:id")
		tasks.POST("/:id/assign")
		tasks.POST("/:id/start")
		tasks.POST("/:id/complete")
		tasks.POST("/:id/review")
		tasks.POST("/:id/approve")
		tasks.POST("/:id/reject")
		tasks.GET("/:id/history")
		tasks.GET("/:id/conversation")
	}
}