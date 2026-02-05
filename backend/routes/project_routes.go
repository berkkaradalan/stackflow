package routes

import (
	"github.com/berkkaradalan/stackflow/handler"
	"github.com/berkkaradalan/stackflow/middleware"
	"github.com/berkkaradalan/stackflow/utils"
	"github.com/gin-gonic/gin"
)

func setupProjectRoutes(r *gin.RouterGroup, projectHandler *handler.ProjectHandler, agentHandler *handler.AgentHandler, jwtManager *utils.JWTManager) {
	projects := r.Group("/projects")

	projects.Use(middleware.AuthMiddleware(jwtManager))
	{
		projects.GET("", projectHandler.GetAllProjects)
		projects.POST("", projectHandler.CreateProject)
		projects.GET("/:id", projectHandler.GetProjectByID)
		projects.PUT("/:id", projectHandler.UpdateProject)
		projects.DELETE("/:id", projectHandler.DeleteProject)
		projects.GET("/:id/stats", projectHandler.GetProjectStats)
		projects.GET("/:id/agents", agentHandler.GetAgentsByProjectID)
	}
}