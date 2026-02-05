package models

import "github.com/berkkaradalan/stackflow/utils"

type LoginRequest struct {
	Email    string `json:"email" binding:"required,email"`
	Password string `json:"password" binding:"required,min=5"`
}

type RefreshTokenRequest struct {
	RefreshToken string `json:"refresh_token" binding:"required"`
}

type UpdateProfileRequest struct {
	Username    string `json:"username" binding:"omitempty,min=3,max=50"`
	Email       string `json:"email" binding:"omitempty,email"`
	AvatarUrl   string `json:"avatar_url" binding:"omitempty,url"`
	OldPassword string `json:"old_password" binding:"omitempty,min=6"`
	NewPassword string `json:"new_password" binding:"omitempty,min=6"`
}

type AuthResponse struct {
	User   User             `json:"user"`
	Tokens *utils.TokenPair `json:"tokens"`
}