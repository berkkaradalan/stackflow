package models

import "time"

type User struct {
	ID           int       `json:"id"`
	Username     string    `json:"username"`
	Email        string    `json:"email"`
	PasswordHash string    `json:"-"`
	AvatarUrl    string    `json:"avatar_url"`
	Role         string    `json:"role"`
	IsActive     bool      `json:"is_active"`
	CreatedAt    time.Time `json:"created_at"`
	UpdatedAt    time.Time `json:"updated_at"`
}

// UpdateUserRequest is the request model for updating a user
type UpdateUserRequest struct {
	Username  *string `json:"username" binding:"omitempty,min=3,max=50"`
	Email     *string `json:"email" binding:"omitempty,email"`
	AvatarUrl *string `json:"avatar_url" binding:"omitempty,url"`
	Role      *string `json:"role" binding:"omitempty,oneof=admin user"`
	IsActive  *bool   `json:"is_active"`
}

// InviteUserRequest is the request model for inviting a new user
type InviteUserRequest struct {
	Username string `json:"username" binding:"required,min=3,max=50"`
	Email    string `json:"email" binding:"required,email"`
	Role     string `json:"role" binding:"required,oneof=admin user"`
}

// UserListResponse is the response model for listing users
type UserListResponse struct {
	Users      []User `json:"users"`
	TotalCount int    `json:"total_count"`
}