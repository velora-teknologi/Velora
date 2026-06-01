package services

import (
	"context"
	"crypto/rand"
	"encoding/hex"
	"strings"
	"time"

	"go.uber.org/zap"
	"golang.org/x/crypto/bcrypt"

	"github.com/velora-teknologi/velora/internal/application/dtos"
	"github.com/velora-teknologi/velora/internal/domain/models"
	"github.com/velora-teknologi/velora/internal/domain/repositories"
	customErrors "github.com/velora-teknologi/velora/pkg/errors"
)

type APIKeyService interface {
	CreateAPIKey(ctx context.Context, userID string, req *dtos.CreateAPIKeyRequest) (*dtos.APIKeyResponse, error)
	GetAPIKey(ctx context.Context, userID, id string) (*dtos.APIKeyResponse, error)
	ListAPIKeys(ctx context.Context, userID string, skip, limit int) ([]*dtos.APIKeyResponse, error)
	RevokeAPIKey(ctx context.Context, userID, id string) error
}

type apiKeyService struct {
	repo     repositories.APIKeyRepository
	userRepo repositories.UserRepository
	logger   *zap.SugaredLogger
}

func NewAPIKeyService(repo repositories.APIKeyRepository, userRepo repositories.UserRepository, logger *zap.SugaredLogger) APIKeyService {
	return &apiKeyService{repo: repo, userRepo: userRepo, logger: logger}
}

func (s *apiKeyService) CreateAPIKey(ctx context.Context, userID string, req *dtos.CreateAPIKeyRequest) (*dtos.APIKeyResponse, error) {
	if req.Name == "" {
		return nil, customErrors.NewValidationError("API key name is required")
	}

	user, err := s.userRepo.GetByID(ctx, userID)
	if err != nil {
		s.logger.Errorf("Error retrieving user: %v", err)
		return nil, customErrors.NewInternalError("Failed to create API key")
	}
	if user == nil {
		return nil, customErrors.NewNotFoundError("User")
	}

	secret, err := generateAPIKeySecret()
	if err != nil {
		s.logger.Errorf("Error generating API key: %v", err)
		return nil, customErrors.NewInternalError("Failed to create API key")
	}

	hash, err := bcrypt.GenerateFromPassword([]byte(secret), bcrypt.DefaultCost)
	if err != nil {
		s.logger.Errorf("Error hashing API key secret: %v", err)
		return nil, customErrors.NewInternalError("Failed to create API key")
	}

	apiKey := &models.APIKey{
		UserID:   userID,
		Name:     req.Name,
		KeyHash:  string(hash),
		Scopes:   strings.Join(req.Scopes, ","),
		IsActive: true,
	}

	if req.ExpiresIn > 0 {
		expiresAt := time.Now().Add(time.Duration(req.ExpiresIn) * time.Second)
		apiKey.ExpiresAt = &expiresAt
	}

	if err := s.repo.Create(ctx, apiKey); err != nil {
		s.logger.Errorf("Error creating API key: %v", err)
		return nil, customErrors.NewInternalError("Failed to create API key")
	}

	key := apiKey.KeyID + "." + secret
	return apiKeyToResponse(apiKey, &key), nil
}

func (s *apiKeyService) GetAPIKey(ctx context.Context, userID, id string) (*dtos.APIKeyResponse, error) {
	apiKey, err := s.repo.GetByID(ctx, id)
	if err != nil {
		s.logger.Errorf("Error retrieving API key: %v", err)
		return nil, customErrors.NewInternalError("Failed to get API key")
	}
	if apiKey == nil || apiKey.UserID != userID {
		return nil, customErrors.NewNotFoundError("API key")
	}

	return apiKeyToResponse(apiKey, nil), nil
}

func (s *apiKeyService) ListAPIKeys(ctx context.Context, userID string, skip, limit int) ([]*dtos.APIKeyResponse, error) {
	apiKeys, err := s.repo.List(ctx, userID, skip, limit)
	if err != nil {
		s.logger.Errorf("Error listing API keys: %v", err)
		return nil, customErrors.NewInternalError("Failed to list API keys")
	}

	var responses []*dtos.APIKeyResponse
	for _, apiKey := range apiKeys {
		responses = append(responses, apiKeyToResponse(apiKey, nil))
	}
	return responses, nil
}

func (s *apiKeyService) RevokeAPIKey(ctx context.Context, userID, id string) error {
	apiKey, err := s.repo.GetByID(ctx, id)
	if err != nil {
		s.logger.Errorf("Error retrieving API key: %v", err)
		return customErrors.NewInternalError("Failed to revoke API key")
	}
	if apiKey == nil || apiKey.UserID != userID {
		return customErrors.NewNotFoundError("API key")
	}

	apiKey.IsActive = false
	if err := s.repo.Update(ctx, apiKey); err != nil {
		s.logger.Errorf("Error revoking API key: %v", err)
		return customErrors.NewInternalError("Failed to revoke API key")
	}

	return nil
}

func generateAPIKeySecret() (string, error) {
	bytes := make([]byte, 32)
	if _, err := rand.Read(bytes); err != nil {
		return "", err
	}
	return hex.EncodeToString(bytes), nil
}

func apiKeyToResponse(apiKey *models.APIKey, secret *string) *dtos.APIKeyResponse {
	var scopes []string
	if apiKey.Scopes != "" {
		scopes = strings.Split(apiKey.Scopes, ",")
	}

	response := &dtos.APIKeyResponse{
		ID:        apiKey.ID,
		Name:      apiKey.Name,
		Scopes:    scopes,
		IsActive:  apiKey.IsActive,
		ExpiresAt: formatTime(apiKey.ExpiresAt),
	}
	if secret != nil {
		response.Key = *secret
	}

	return response
}

func formatTime(t *time.Time) *string {
	if t == nil {
		return nil
	}
	value := t.UTC().Format(time.RFC3339)
	return &value
}
