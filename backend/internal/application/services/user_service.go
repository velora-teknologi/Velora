package services

import (
	"context"
	"time"

	"github.com/golang-jwt/jwt/v5"
	"go.uber.org/zap"
	"golang.org/x/crypto/bcrypt"

	"github.com/velora-teknologi/velora/internal/application/dtos"
	"github.com/velora-teknologi/velora/internal/domain/models"
	"github.com/velora-teknologi/velora/internal/domain/repositories"
	customErrors "github.com/velora-teknologi/velora/pkg/errors"
)

type UserService interface {
	CreateUser(ctx context.Context, req *dtos.CreateUserRequest) (*dtos.UserResponse, error)
	GetUser(ctx context.Context, id string) (*dtos.UserResponse, error)
	UpdateUser(ctx context.Context, id string, req *dtos.UpdateUserRequest) (*dtos.UserResponse, error)
	DeleteUser(ctx context.Context, id string) error
	ListUsers(ctx context.Context, skip, limit int) ([]*dtos.UserResponse, error)
	LoginUser(ctx context.Context, req *dtos.LoginRequest) (*dtos.LoginResponse, error)
}

type userService struct {
	repo      repositories.UserRepository
	logger    *zap.SugaredLogger
	jwtSecret string
}

func NewUserService(repo repositories.UserRepository, jwtSecret string, logger *zap.SugaredLogger) UserService {
	return &userService{
		repo:      repo,
		jwtSecret: jwtSecret,
		logger:    logger,
	}
}

func (s *userService) CreateUser(ctx context.Context, req *dtos.CreateUserRequest) (*dtos.UserResponse, error) {
	// Check if user already exists
	existing, err := s.repo.GetByEmail(ctx, req.Email)
	if err != nil {
		s.logger.Errorf("Error checking existing user: %v", err)
		return nil, customErrors.NewInternalError("Failed to create user")
	}

	if existing != nil {
		return nil, customErrors.NewConflictError("User with this email already exists")
	}

	// Hash password
	hashedPassword, err := bcrypt.GenerateFromPassword([]byte(req.Password), bcrypt.DefaultCost)
	if err != nil {
		s.logger.Errorf("Error hashing password: %v", err)
		return nil, customErrors.NewInternalError("Failed to create user")
	}

	user := &models.User{
		Email:    req.Email,
		Password: string(hashedPassword),
		FullName: req.FullName,
		Role:     "user",
		IsActive: true,
	}

	if err := s.repo.Create(ctx, user); err != nil {
		s.logger.Errorf("Error creating user: %v", err)
		return nil, customErrors.NewInternalError("Failed to create user")
	}

	return userToResponse(user), nil
}

func (s *userService) GetUser(ctx context.Context, id string) (*dtos.UserResponse, error) {
	user, err := s.repo.GetByID(ctx, id)
	if err != nil {
		s.logger.Errorf("Error getting user: %v", err)
		return nil, customErrors.NewInternalError("Failed to get user")
	}

	if user == nil {
		return nil, customErrors.NewNotFoundError("User")
	}

	return userToResponse(user), nil
}

func (s *userService) UpdateUser(ctx context.Context, id string, req *dtos.UpdateUserRequest) (*dtos.UserResponse, error) {
	user, err := s.repo.GetByID(ctx, id)
	if err != nil {
		s.logger.Errorf("Error getting user: %v", err)
		return nil, customErrors.NewInternalError("Failed to update user")
	}

	if user == nil {
		return nil, customErrors.NewNotFoundError("User")
	}

	// Update fields if provided
	if req.Email != "" {
		user.Email = req.Email
	}
	if req.FullName != "" {
		user.FullName = req.FullName
	}
	if req.Avatar != "" {
		user.Avatar = req.Avatar
	}
	if req.Role != "" {
		user.Role = req.Role
	}

	if err := s.repo.Update(ctx, user); err != nil {
		s.logger.Errorf("Error updating user: %v", err)
		return nil, customErrors.NewInternalError("Failed to update user")
	}

	return userToResponse(user), nil
}

func (s *userService) DeleteUser(ctx context.Context, id string) error {
	user, err := s.repo.GetByID(ctx, id)
	if err != nil {
		s.logger.Errorf("Error getting user: %v", err)
		return customErrors.NewInternalError("Failed to delete user")
	}

	if user == nil {
		return customErrors.NewNotFoundError("User")
	}

	if err := s.repo.Delete(ctx, id); err != nil {
		s.logger.Errorf("Error deleting user: %v", err)
		return customErrors.NewInternalError("Failed to delete user")
	}

	return nil
}

func (s *userService) ListUsers(ctx context.Context, skip, limit int) ([]*dtos.UserResponse, error) {
	users, err := s.repo.List(ctx, skip, limit)
	if err != nil {
		s.logger.Errorf("Error listing users: %v", err)
		return nil, customErrors.NewInternalError("Failed to list users")
	}

	var responses []*dtos.UserResponse
	for _, user := range users {
		responses = append(responses, userToResponse(user))
	}

	return responses, nil
}

func (s *userService) LoginUser(ctx context.Context, req *dtos.LoginRequest) (*dtos.LoginResponse, error) {
	user, err := s.repo.GetByEmail(ctx, req.Email)
	if err != nil {
		s.logger.Errorf("Error getting user: %v", err)
		return nil, customErrors.NewInternalError("Login failed")
	}

	if user == nil {
		return nil, customErrors.NewUnauthorizedError("Invalid email or password")
	}

	// Verify password
	if err := bcrypt.CompareHashAndPassword([]byte(user.Password), []byte(req.Password)); err != nil {
		s.logger.Warnf("Invalid password for user %s", req.Email)
		return nil, customErrors.NewUnauthorizedError("Invalid email or password")
	}

	token, err := s.generateToken(user)
	if err != nil {
		s.logger.Errorf("Error generating JWT token: %v", err)
		return nil, customErrors.NewInternalError("Login failed")
	}

	return &dtos.LoginResponse{
		Token: token,
		User:  *userToResponse(user),
	}, nil
}

func (s *userService) generateToken(user *models.User) (string, error) {
	if s.jwtSecret == "" {
		return "", customErrors.NewInternalError("JWT secret is not configured")
	}

	claims := jwt.MapClaims{
		"sub":   user.ID,
		"email": user.Email,
		"role":  user.Role,
		"exp":   time.Now().Add(24 * time.Hour).Unix(),
	}

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	return token.SignedString([]byte(s.jwtSecret))
}

func userToResponse(user *models.User) *dtos.UserResponse {
	return &dtos.UserResponse{
		ID:       user.ID,
		Email:    user.Email,
		FullName: user.FullName,
		Avatar:   user.Avatar,
		Role:     user.Role,
		IsActive: user.IsActive,
	}
}
