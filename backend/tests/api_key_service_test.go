package services_test

import (
	"context"
	"strings"
	"testing"

	"github.com/google/uuid"
	"github.com/velora-teknologi/velora/internal/application/dtos"
	"github.com/velora-teknologi/velora/internal/application/services"
	"github.com/velora-teknologi/velora/internal/domain/models"
	"go.uber.org/zap"
)

type mockAPIKeyRepository struct {
	apiKeys map[string]*models.APIKey
}

func newMockAPIKeyRepository() *mockAPIKeyRepository {
	return &mockAPIKeyRepository{apiKeys: make(map[string]*models.APIKey)}
}

func (m *mockAPIKeyRepository) Create(ctx context.Context, apiKey *models.APIKey) error {
	if apiKey.ID == "" {
		apiKey.ID = uuid.NewString()
	}
	m.apiKeys[apiKey.ID] = apiKey
	return nil
}

func (m *mockAPIKeyRepository) GetByID(ctx context.Context, id string) (*models.APIKey, error) {
	return m.apiKeys[id], nil
}

func (m *mockAPIKeyRepository) GetByKeyID(ctx context.Context, keyID string) (*models.APIKey, error) {
	for _, apiKey := range m.apiKeys {
		if apiKey.KeyID == keyID {
			return apiKey, nil
		}
	}
	return nil, nil
}

func (m *mockAPIKeyRepository) Update(ctx context.Context, apiKey *models.APIKey) error {
	m.apiKeys[apiKey.ID] = apiKey
	return nil
}

func (m *mockAPIKeyRepository) Delete(ctx context.Context, id string) error {
	delete(m.apiKeys, id)
	return nil
}

func (m *mockAPIKeyRepository) List(ctx context.Context, userID string, skip, limit int) ([]*models.APIKey, error) {
	var list []*models.APIKey
	for _, apiKey := range m.apiKeys {
		if apiKey.UserID == userID {
			list = append(list, apiKey)
		}
	}
	return list, nil
}

func TestCreateAPIKey(t *testing.T) {
	repo := newMockAPIKeyRepository()
	userRepo := NewMockUserRepository()
	logger := zap.NewNop().Sugar()
	service := services.NewAPIKeyService(repo, userRepo, logger)

	user := &models.User{ID: uuid.NewString(), Email: "apiuser@example.com", FullName: "API User", Password: "hashed-password"}
	userRepo.Create(context.Background(), user)

	req := &dtos.CreateAPIKeyRequest{Name: "Default Key", Scopes: []string{"read", "write"}, ExpiresIn: 3600}
	resp, err := service.CreateAPIKey(context.Background(), user.ID, req)
	if err != nil {
		t.Fatalf("expected no error, got %v", err)
	}
	if resp.Key == "" {
		t.Fatal("expected raw API key value in response")
	}
	if !strings.Contains(resp.Key, ".") {
		t.Fatal("expected api key to include key id prefix")
	}
	if resp.Name != req.Name || !resp.IsActive {
		t.Fatalf("unexpected response: %+v", resp)
	}
}

func TestGetAndListAPIKeys(t *testing.T) {
	repo := newMockAPIKeyRepository()
	userRepo := NewMockUserRepository()
	logger := zap.NewNop().Sugar()
	service := services.NewAPIKeyService(repo, userRepo, logger)

	user := &models.User{ID: uuid.NewString(), Email: "apiuser2@example.com", FullName: "API User 2", Password: "hashed-password"}
	userRepo.Create(context.Background(), user)

	apiKey := &models.APIKey{ID: uuid.NewString(), UserID: user.ID, Name: "Read Key", KeyHash: "hashed", Scopes: "read", IsActive: true}
	repo.Create(context.Background(), apiKey)

	resp, err := service.GetAPIKey(context.Background(), user.ID, apiKey.ID)
	if err != nil {
		t.Fatalf("expected no error, got %v", err)
	}
	if resp.ID != apiKey.ID || resp.Key != "" {
		t.Fatalf("unexpected API key response: %+v", resp)
	}

	list, err := service.ListAPIKeys(context.Background(), user.ID, 0, 10)
	if err != nil {
		t.Fatalf("expected no error listing API keys, got %v", err)
	}
	if len(list) != 1 {
		t.Fatalf("expected 1 api key, got %d", len(list))
	}
}

func TestRevokeAPIKey(t *testing.T) {
	repo := newMockAPIKeyRepository()
	userRepo := NewMockUserRepository()
	logger := zap.NewNop().Sugar()
	service := services.NewAPIKeyService(repo, userRepo, logger)

	user := &models.User{ID: uuid.NewString(), Email: "apiuser3@example.com", FullName: "API User 3", Password: "hashed-password"}
	userRepo.Create(context.Background(), user)

	apiKey := &models.APIKey{ID: uuid.NewString(), UserID: user.ID, Name: "Revocable Key", KeyHash: "hashed", Scopes: "read", IsActive: true}
	repo.Create(context.Background(), apiKey)

	if err := service.RevokeAPIKey(context.Background(), user.ID, apiKey.ID); err != nil {
		t.Fatalf("expected no error revoking api key, got %v", err)
	}

	updated, _ := repo.GetByID(context.Background(), apiKey.ID)
	if updated == nil || updated.IsActive {
		t.Fatal("expected api key to be deactivated")
	}
}
