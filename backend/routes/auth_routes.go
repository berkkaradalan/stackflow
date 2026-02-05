package routes

import (
	"github.com/berkkaradalan/stackflow/handler"
	"github.com/berkkaradalan/stackflow/middleware"
	"github.com/berkkaradalan/stackflow/utils"
	"github.com/gin-gonic/gin"
)

func setupAuthRoutes(
	r *gin.RouterGroup, authHandler *handler.AuthHandler, jwtManager *utils.JWTManager) {
	auth := r.Group("/auth")

	// Public
	auth.POST("/login", authHandler.Login)
	auth.POST("/refresh", authHandler.Refresh)
	auth.POST("/register", authHandler.Register)
	auth.GET("/validate-invite-token", authHandler.ValidateInviteToken)

	// Protected
	auth.Use(middleware.AuthMiddleware(jwtManager))
	{
		auth.POST("/logout", authHandler.Logout)
		auth.GET("/profile", authHandler.GetProfile)
		auth.PUT("/profile", authHandler.UpdateProfile)
	}
}