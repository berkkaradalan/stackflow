package routes

import (
	"github.com/berkkaradalan/stackflow/handler"
	"github.com/berkkaradalan/stackflow/middleware"
	"github.com/berkkaradalan/stackflow/utils"
	"github.com/gin-gonic/gin"
)

func setupUserRoutes(
	r *gin.RouterGroup, userHandler *handler.UserHandler, jwtManager *utils.JWTManager) {
	users := r.Group("/users")

	// All user management endpoints require authentication and admin role
	users.Use(middleware.AuthMiddleware(jwtManager))
	users.Use(middleware.RoleMiddleware("admin"))
	{
		users.GET("", userHandler.GetAllUsers)
		users.GET("/:id", userHandler.GetUserByID)
		users.PUT("/:id", userHandler.UpdateUser)
		users.DELETE("/:id", userHandler.DeleteUser)
		users.POST("/invite", userHandler.InviteUser)
	}
}