package service

import (
	"context"
	"crypto/rand"
	"encoding/base64"
	"errors"
	"fmt"
	"time"

	"github.com/berkkaradalan/stackflow/models"
	"github.com/berkkaradalan/stackflow/repository/postgres"
	"github.com/berkkaradalan/stackflow/utils"
)

type UserService struct {
	userRepo        *repository.UserRepository
	inviteTokenRepo *repository.InviteTokenRepository
}

func NewUserService(userRepo *repository.UserRepository, inviteTokenRepo *repository.InviteTokenRepository) *UserService {
	return &UserService{
		userRepo:        userRepo,
		inviteTokenRepo: inviteTokenRepo,
	}
}

// GetAllUsers returns all users from the database
func (s *UserService) GetAllUsers(ctx context.Context) (*models.UserListResponse, error) {
	users, err := s.userRepo.GetAll(ctx)
	if err != nil {
		return nil, fmt.Errorf("failed to get users: %w", err)
	}

	return &models.UserListResponse{
		Users:      users,
		TotalCount: len(users),
	}, nil
}

// GetUserByID returns a single user by ID
func (s *UserService) GetUserByID(ctx context.Context, id int) (*models.User, error) {
	user, err := s.userRepo.GetByID(ctx, id)
	if err != nil {
		return nil, fmt.Errorf("failed to get user: %w", err)
	}

	return user, nil
}

// UpdateUser updates a user's information
func (s *UserService) UpdateUser(ctx context.Context, id int, req *models.UpdateUserRequest) (*models.User, error) {
	// First check if user exists
	_, err := s.userRepo.GetByID(ctx, id)
	if err != nil {
		return nil, fmt.Errorf("user not found: %w", err)
	}

	// Build updates map
	updates := make(map[string]interface{})

	if req.Username != nil {
		updates["username"] = *req.Username
	}

	if req.Email != nil {
		// Check if email is already taken by another user
		existingUser, err := s.userRepo.GetByEmail(ctx, *req.Email)
		if err == nil && existingUser.ID != id {
			return nil, errors.New("email already in use")
		}
		updates["email"] = *req.Email
	}

	if req.AvatarUrl != nil {
		updates["avatar_url"] = *req.AvatarUrl
	}

	if req.Role != nil {
		updates["role"] = *req.Role
	}

	if req.IsActive != nil {
		updates["is_active"] = *req.IsActive
	}

	// Perform partial update
	updatedUser, err := s.userRepo.UpdatePartial(ctx, id, updates)
	if err != nil {
		return nil, fmt.Errorf("failed to update user: %w", err)
	}

	return updatedUser, nil
}

// DeleteUser deletes a user by ID
func (s *UserService) DeleteUser(ctx context.Context, id int) error {
	// Check if user exists
	_, err := s.userRepo.GetByID(ctx, id)
	if err != nil {
		return fmt.Errorf("user not found: %w", err)
	}

	err = s.userRepo.Delete(ctx, id)
	if err != nil {
		return fmt.Errorf("failed to delete user: %w", err)
	}

	return nil
}

// InviteUser creates an invite token for a new user
func (s *UserService) InviteUser(ctx context.Context, req *models.InviteUserRequest, baseURL string) (*models.InviteResponse, error) {
	// Check if email already exists
	existingUser, err := s.userRepo.GetByEmail(ctx, req.Email)
	if err == nil && existingUser != nil {
		return nil, errors.New("email already in use")
	}

	// Generate a secure random token
	token, err := generateSecureToken(32)
	if err != nil {
		return nil, fmt.Errorf("failed to generate token: %w", err)
	}

	// Token expires in 7 days
	expiresAt := time.Now().Add(7 * 24 * time.Hour)

	// Create invite token
	inviteToken := &models.InviteToken{
		Email:     req.Email,
		Username:  req.Username,
		Role:      req.Role,
		Token:     token,
		ExpiresAt: expiresAt,
	}

	err = s.inviteTokenRepo.Create(ctx, inviteToken)
	if err != nil {
		return nil, fmt.Errorf("failed to create invite token: %w", err)
	}

	// Generate invite link
	inviteLink := fmt.Sprintf("%s/register?token=%s", baseURL, token)

	return &models.InviteResponse{
		InviteLink: inviteLink,
		Token:      token,
		ExpiresAt:  expiresAt,
		Email:      req.Email,
	}, nil
}

// RegisterWithToken registers a user using an invite token
func (s *UserService) RegisterWithToken(ctx context.Context, req *models.RegisterWithTokenRequest) (*models.User, error) {
	// Get invite token
	inviteToken, err := s.inviteTokenRepo.GetByToken(ctx, req.Token)
	if err != nil {
		return nil, errors.New("invalid or expired invite token")
	}

	// Check if already used
	if inviteToken.UsedAt != nil {
		return nil, errors.New("invite token already used")
	}

	// Check if expired
	if time.Now().After(inviteToken.ExpiresAt) {
		return nil, errors.New("invite token expired")
	}

	// Check if email already exists (race condition check)
	existingUser, err := s.userRepo.GetByEmail(ctx, inviteToken.Email)
	if err == nil && existingUser != nil {
		return nil, errors.New("email already in use")
	}

	// Hash the password
	hashedPassword, err := utils.HashPassword(req.Password)
	if err != nil {
		return nil, fmt.Errorf("failed to hash password: %w", err)
	}

	// Create the user
	user := &models.User{
		Username:     inviteToken.Username,
		Email:        inviteToken.Email,
		PasswordHash: hashedPassword,
		AvatarUrl:    "",
		Role:         inviteToken.Role,
		IsActive:     true,
	}

	err = s.userRepo.Create(ctx, user)
	if err != nil {
		return nil, fmt.Errorf("failed to create user: %w", err)
	}

	// Mark token as used
	err = s.inviteTokenRepo.MarkAsUsed(ctx, req.Token)
	if err != nil {
		return nil, fmt.Errorf("failed to mark token as used: %w", err)
	}

	return user, nil
}

// generateSecureToken generates a cryptographically secure random token
func generateSecureToken(length int) (string, error) {
	bytes := make([]byte, length)
	_, err := rand.Read(bytes)
	if err != nil {
		return "", err
	}
	return base64.URLEncoding.EncodeToString(bytes), nil
}
