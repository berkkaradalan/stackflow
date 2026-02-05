package models

import "time"

type InviteToken struct {
	ID        int       `json:"id"`
	Email     string    `json:"email"`
	Username  string    `json:"username"`
	Role      string    `json:"role"`
	Token     string    `json:"token"`
	ExpiresAt time.Time `json:"expires_at"`
	UsedAt    *time.Time `json:"used_at,omitempty"`
	CreatedAt time.Time `json:"created_at"`
}

// RegisterWithTokenRequest is the request model for registering with an invite token
type RegisterWithTokenRequest struct {
	Token    string `json:"token" binding:"required"`
	Password string `json:"password" binding:"required,min=6"`
}

// InviteResponse is the response model for invite creation
type InviteResponse struct {
	InviteLink string    `json:"invite_link"`
	Token      string    `json:"token"`
	ExpiresAt  time.Time `json:"expires_at"`
	Email      string    `json:"email"`
}
