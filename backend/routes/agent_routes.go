package routes

import (
	"github.com/berkkaradalan/stackflow/handler"
	"github.com/berkkaradalan/stackflow/middleware"
	"github.com/berkkaradalan/stackflow/utils"
	"github.com/gin-gonic/gin"
)

func setupAgentsRoutes(r *gin.RouterGroup, agentHandler *handler.AgentHandler, jwtManager *utils.JWTManager) {
	agents := r.Group("/agents")
	agents.Use(middleware.AuthMiddleware(jwtManager))
	{
		agents.POST("", agentHandler.CreateAgent)
		agents.GET("", agentHandler.GetAllAgents)
		agents.GET("/:id", agentHandler.GetAgentByID)
		agents.PUT("/:id", agentHandler.UpdateAgent)
		agents.DELETE("/:id", agentHandler.DeleteAgent)
		agents.GET("/:id/status", agentHandler.GetAgentStatus)
		agents.GET("/:id/workload", agentHandler.GetAgentWorkload)
		agents.GET("/:id/performance", agentHandler.GetAgentPerformance)
		agents.GET("/:id/health", agentHandler.HealthCheck)
	}
}