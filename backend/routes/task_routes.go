package routes

import (
	"github.com/berkkaradalan/stackflow/handler"
	"github.com/berkkaradalan/stackflow/middleware"
	"github.com/berkkaradalan/stackflow/utils"
	"github.com/gin-gonic/gin"
)

func setupTaskRoutes(r *gin.RouterGroup, taskHandler *handler.TaskHandler, jwtManager *utils.JWTManager) {
	// Tasks under projects (requires auth)
	projects := r.Group("/projects")
	projects.Use(middleware.AuthMiddleware(jwtManager))
	{
		projects.GET("/:id/tasks", taskHandler.GetTasksByProject)
		projects.POST("/:id/tasks", taskHandler.CreateTask)
	}

	// Individual task endpoints (requires auth)
	tasks := r.Group("/tasks")
	tasks.Use(middleware.AuthMiddleware(jwtManager))
	{
		// List all tasks with filters
		tasks.GET("", taskHandler.GetAllTasks)

		// CRUD
		tasks.GET("/:id", taskHandler.GetTaskByID)
		tasks.PUT("/:id", taskHandler.UpdateTask)
		tasks.DELETE("/:id", taskHandler.DeleteTask)

		// Assignment
		tasks.POST("/:id/assign", taskHandler.AssignAgent)
		tasks.POST("/:id/reviewer", taskHandler.SetReviewer)

		// Status transitions
		tasks.POST("/:id/start", taskHandler.StartTask)
		tasks.POST("/:id/done", taskHandler.CompleteTask)
		tasks.POST("/:id/close", taskHandler.CloseTask)
		tasks.POST("/:id/wontdo", taskHandler.WontDoTask)
		tasks.POST("/:id/reopen", taskHandler.ReopenTask)

		// Activities (AI progress)
		tasks.GET("/:id/activities", taskHandler.GetTaskActivities)
		tasks.POST("/:id/activities", taskHandler.AddProgress)
	}
}
