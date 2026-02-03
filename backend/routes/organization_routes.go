package routes

import "github.com/gin-gonic/gin"

func setupOrganizationRoutes(r *gin.RouterGroup) {
	organizations := r.Group("/organizations")

	organizations.Use()
	{
		organizations.GET("")
		organizations.PUT("")
		organizations.GET("/members")
		organizations.DELETE("/members/:user_id")
	}
}
