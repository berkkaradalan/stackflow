package handler

import (
	"net/http"

	"github.com/berkkaradalan/stackflow/models"
	"github.com/berkkaradalan/stackflow/service"
	"github.com/gin-gonic/gin"
)

type AuthHandler struct {
	authService *service.AuthService
}

func NewAuthHandler(authService *service.AuthService) *AuthHandler {
	return &AuthHandler{
		authService: authService,
	}
}

func (h *AuthHandler) Login(c *gin.Context) {
	var req models.LoginRequest

	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"error": err.Error(),
		})
		return
	}

	authResponse, err := h.authService.Login(c, req)

	if err != nil { 
		c.JSON(http.StatusInternalServerError, gin.H{
			"error": err.Error(),
		})
		return		
	}

	c.JSON(http.StatusOK, gin.H{
		"user": authResponse.User,
		"access_token": authResponse.Tokens.AccessToken,
		"refresh_token":authResponse.Tokens.RefreshToken, 
		//todo - add actual refresh
	})
}