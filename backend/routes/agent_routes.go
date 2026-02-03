package routes

import "github.com/gin-gonic/gin"

func setupAgentsRoutes(r *gin.RouterGroup) {
	agents := r.Group("/agents")

	agents.Use()
	{
		agents.GET("")
		agents.GET("/:id")
		agents.GET("/:id/status")
		agents.GET("/:id/workload")
		agents.GET("/:id/performance")
		agents.GET("/:id/memory")
		agents.POST("/:id/memory")
		agents.DELETE("/:id/memory/:memory_id")
		agents.GET("/:id/notes")
		agents.PUT("/:id/notes")
	}
}