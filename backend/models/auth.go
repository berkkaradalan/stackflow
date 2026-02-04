package models

import "github.com/berkkaradalan/stackflow/utils"

type LoginRequest struct {
	Email    string `json:"email" binding:"required,email"`
	Password string `json:"password" binding:"required,min=6"`
}

type AuthResponse struct {
	User         User    `json:"user"`
	Tokens       *utils.TokenPair `json:"tokens"`
}