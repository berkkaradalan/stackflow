package routes

import "github.com/gin-gonic/gin"

func setupWorkflowRoutes(r *gin.RouterGroup) {
	workflows := r.Group("/workflows")

	workflows.Use()
	{
		workflows.GET("/templates")
		workflows.GET("/templates/:name")
		workflows.POST("/custom")
	}

	// Project workflow endpoints
	projects := r.Group("/projects/:id")
	{
		projects.GET("/workflow")
		projects.PUT("/workflow")
	}

	// Task workflow statusx
	tasks := r.Group("/tasks/:id")
	{
		tasks.GET("/workflow-status")
	}
}
