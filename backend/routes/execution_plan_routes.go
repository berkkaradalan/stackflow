package routes

import (
	"github.com/berkkaradalan/stackflow/handler"
	"github.com/berkkaradalan/stackflow/middleware"
	"github.com/berkkaradalan/stackflow/utils"
	"github.com/gin-gonic/gin"
)

func setupExecutionPlanRoutes(r *gin.RouterGroup, planHandler *handler.ExecutionPlanHandler, jwtManager *utils.JWTManager) {
	// Project execution plan endpoints (requires auth)
	projects := r.Group("/projects")
	projects.Use(middleware.AuthMiddleware(jwtManager))
	{
		// PM Agent creates and manages execution plans
		projects.POST("/:id/execution-plan", planHandler.CreatePlan)
		projects.GET("/:id/execution-plan", planHandler.GetActivePlan)
		projects.GET("/:id/execution-plans", planHandler.GetAllPlans)
		projects.PUT("/:id/execution-plan", planHandler.UpdatePlan)

		// Reporting endpoints
		projects.GET("/:id/reports/daily", planHandler.GetDailyReport)
		projects.GET("/:id/reports/weekly", planHandler.GetWeeklyReport)
		projects.POST("/:id/reports/generate", planHandler.GenerateReport)
	}

	// Agent task flow endpoints (requires auth)
	agents := r.Group("/agents")
	agents.Use(middleware.AuthMiddleware(jwtManager))
	{
		// Developer/QA bots request tasks and report completion
		agents.GET("/:id/next-task", planHandler.GetNextTask)
		agents.POST("/:id/task-complete", planHandler.TaskComplete)
		agents.GET("/:id/context", planHandler.GetAgentContext)
	}
}
