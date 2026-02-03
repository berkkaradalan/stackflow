package routes

import "github.com/gin-gonic/gin"

func setupAnalyticsRoutes(r *gin.RouterGroup) {
	// Project reports and metrics
	projects := r.Group("/projects/:id")
	{
		projects.GET("/report")
		projects.GET("/metrics")
	}

	// Agents performance
	r.GET("/agents/performance")

	// Analytics endpoints
	analytics := r.Group("/analytics")
	{
		analytics.GET("/costs")
		analytics.GET("/timeline")
	}
}
