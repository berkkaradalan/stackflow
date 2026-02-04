package service

import (
	"context"
	"errors"
	"fmt"

	"github.com/berkkaradalan/stackflow/models"
	repository "github.com/berkkaradalan/stackflow/repository/postgres"
	"github.com/berkkaradalan/stackflow/utils"
)

type AuthService struct {
	userRepo    *repository.UserRepository
	// cacheRepo   *redis.CacheRepository
	jwtManager  *utils.JWTManager
}

func NewAuthService(
	userRepo *repository.UserRepository, 
	// cacheRepo *redis.CacheRepository,
	jwtManager *utils.JWTManager,
) *AuthService {
	return &AuthService{
		userRepo:   userRepo,
		// cacheRepo:  cacheRepo,
		jwtManager: jwtManager,
	}
}

type AuthResponse struct {
	User         *models.User    `json:"user"`
	Tokens       *utils.TokenPair `json:"tokens"`
}

func (s *AuthService) Login(ctx context.Context, req models.LoginRequest) (*AuthResponse, error) {
	user, err := s.userRepo.GetByEmail(ctx, req.Email)
	if err != nil {
		return nil, errors.New("invalid credentials")
	}

	if !utils.CheckPassword(req.Password, user.PasswordHash) {
		return nil, errors.New("invalid credentials")
	}

	if !user.IsActive {
		return nil, errors.New("account is deactivated")
	}

	tokens, err := s.jwtManager.GenerateTokenPair(user.ID, user.Email, user.Role)
	if err != nil {
		return nil, fmt.Errorf("failed to generate tokens: %w", err)
	}

	// // Store refresh token in Redis (for logout/blacklist functionality)
	// err = s.cacheRepo.SetRefreshToken(ctx, user.ID, tokens.RefreshToken, s.jwtManager.GetRefreshExpiry())
	// if err != nil {
	// 	return nil, fmt.Errorf("failed to store refresh token: %w", err)
	// }

	user.PasswordHash = ""

	return &AuthResponse{
		User:   user,
		Tokens: tokens,
	}, nil
}

func (s *AuthService) RefreshToken(ctx context.Context, refreshToken string) (*utils.TokenPair, error) {
	// Validate refresh token
	claims, err := s.jwtManager.ValidateRefreshToken(refreshToken)
	if err != nil {
		return nil, errors.New("invalid refresh token")
	}

	// Check if refresh token exists in Redis (not revoked)
	// storedToken, err := s.cacheRepo.GetRefreshToken(ctx, claims.UserID)
	// if err != nil || storedToken != refreshToken {
	// 	return nil, errors.New("refresh token revoked or expired")
	// }

	// Get user to ensure they still exist and are active
	user, err := s.userRepo.GetByID(ctx, claims.UserID) // You need to implement GetByID
	if err != nil || !user.IsActive {
		return nil, errors.New("user not found or inactive")
	}

	// Generate new token pair
	newTokens, err := s.jwtManager.GenerateTokenPair(user.ID, user.Email, user.Role)
	if err != nil {
		return nil, err
	}

	// Replace old refresh token with new one (rotation)
	// err = s.cacheRepo.SetRefreshToken(ctx, user.ID, newTokens.RefreshToken, s.jwtManager.GetRefreshExpiry())
	// if err != nil {
	// 	return nil, err
	// }

	return newTokens, nil
}

func (s *AuthService) Logout(ctx context.Context, userID int) error {
	// TODO: Implement token blacklisting with Redis when cacheRepo is enabled
	// For now, client-side token removal is sufficient
	return nil
}

func (s *AuthService) GetProfile(ctx context.Context, userID int) (*models.User, error) {
	user, err := s.userRepo.GetByID(ctx, userID)
	if err != nil {
		return nil, err
	}
	user.PasswordHash = ""
	return user, nil
}

func (s *AuthService) UpdateProfile(ctx context.Context, userID int, req models.UpdateProfileRequest) (*models.User, error) {
	// Get current user
	user, err := s.userRepo.GetByID(ctx, userID)
	if err != nil {
		return nil, err
	}

	// Update fields if provided
	utils.SetIfNotEmpty(&user.Username, req.Username)
	utils.SetIfNotEmpty(&user.Email, req.Email)
	utils.SetIfNotEmpty(&user.AvatarUrl, req.AvatarUrl)

	// Update password if provided
	if req.NewPassword != "" {
		if req.OldPassword == "" {
			return nil, errors.New("old password is required to change password")
		}
		if !utils.CheckPassword(req.OldPassword, user.PasswordHash) {
			return nil, errors.New("invalid old password")
		}
		hashedPassword, err := utils.HashPassword(req.NewPassword)
		if err != nil {
			return nil, fmt.Errorf("failed to hash password: %w", err)
		}
		user.PasswordHash = hashedPassword
	}

	// Save updated user
	if err := s.userRepo.Update(ctx, user); err != nil {
		return nil, err
	}

	user.PasswordHash = ""
	return user, nil
}