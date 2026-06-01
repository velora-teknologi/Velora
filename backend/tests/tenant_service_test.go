package services_test

import (
	"context"
	"encoding/json"
	"testing"

	"github.com/google/uuid"
	"github.com/velora-teknologi/velora/internal/application/dtos"
	"github.com/velora-teknologi/velora/internal/application/services"
	"github.com/velora-teknologi/velora/internal/domain/models"
	"go.uber.org/zap"
)

type mockTenantRepository struct {
	tenants map[string]*models.Tenant
}

type mockPublisher struct {
	subject string
	payload []byte
}

func (m *mockPublisher) Publish(subject string, data []byte) error {
	m.subject = subject
	m.payload = append([]byte(nil), data...)
	return nil
}

func newMockTenantRepository() *mockTenantRepository {
	return &mockTenantRepository{tenants: make(map[string]*models.Tenant)}
}

func (m *mockTenantRepository) Create(ctx context.Context, tenant *models.Tenant) error {
	if tenant.ID == "" {
		tenant.ID = uuid.NewString()
	}
	m.tenants[tenant.ID] = tenant
	return nil
}

func (m *mockTenantRepository) GetByID(ctx context.Context, id string) (*models.Tenant, error) {
	return m.tenants[id], nil
}

func (m *mockTenantRepository) Update(ctx context.Context, tenant *models.Tenant) error {
	m.tenants[tenant.ID] = tenant
	return nil
}

func (m *mockTenantRepository) Delete(ctx context.Context, id string) error {
	delete(m.tenants, id)
	return nil
}

func (m *mockTenantRepository) List(ctx context.Context, ownerID string, skip, limit int) ([]*models.Tenant, error) {
	var tenants []*models.Tenant
	for _, tenant := range m.tenants {
		if tenant.OwnerID == ownerID {
			tenants = append(tenants, tenant)
		}
	}
	return tenants, nil
}

func TestCreateTenant(t *testing.T) {
	repo := newMockTenantRepository()
	publisher := &mockPublisher{}
	logger := zap.NewNop().Sugar()
	service := services.NewTenantService(repo, publisher, logger)

	req := &dtos.CreateTenantRequest{
		Name:        "Acme Corp",
		Description: "Tenant for Acme",
	}

	resp, err := service.CreateTenant(context.Background(), "owner-1", req)
	if err != nil {
		t.Fatalf("expected no error, got %v", err)
	}

	if resp.Name != req.Name || resp.Description != req.Description || resp.OwnerID != "owner-1" {
		t.Fatalf("unexpected tenant response: %+v", resp)
	}

	if publisher.subject != "tenant.created" {
		t.Fatalf("expected tenant.created event, got %s", publisher.subject)
	}

	if len(publisher.payload) == 0 {
		t.Fatal("expected tenant.created payload to be published")
	}
}

func TestCreateTenantPublishesEvent(t *testing.T) {
	repo := newMockTenantRepository()
	publisher := &mockPublisher{}
	logger := zap.NewNop().Sugar()
	service := services.NewTenantService(repo, publisher, logger)

	req := &dtos.CreateTenantRequest{
		Name:        "Orbit Labs",
		Description: "Tenant for Orbit",
	}

	_, err := service.CreateTenant(context.Background(), "owner-2", req)
	if err != nil {
		t.Fatalf("expected no error, got %v", err)
	}

	if publisher.subject != "tenant.created" {
		t.Fatalf("expected tenant.created event, got %s", publisher.subject)
	}

	var payload map[string]interface{}
	if err := json.Unmarshal(publisher.payload, &payload); err != nil {
		t.Fatalf("expected valid json payload, got %v", err)
	}

	if payload["owner_id"] != "owner-2" {
		t.Fatalf("expected owner_id owner-2, got %v", payload["owner_id"])
	}
}

func TestGetTenant(t *testing.T) {
	repo := newMockTenantRepository()
	publisher := &mockPublisher{}
	logger := zap.NewNop().Sugar()
	service := services.NewTenantService(repo, publisher, logger)

	tenant := &models.Tenant{
		ID:          uuid.NewString(),
		Name:        "Acme Corp",
		Description: "Tenant for Acme",
		OwnerID:     "owner-1",
		IsActive:    true,
	}
	repo.Create(context.Background(), tenant)

	resp, err := service.GetTenant(context.Background(), "owner-1", tenant.ID)
	if err != nil {
		t.Fatalf("expected no error, got %v", err)
	}
	if resp.ID != tenant.ID {
		t.Fatalf("expected tenant id %s, got %s", tenant.ID, resp.ID)
	}
}

func TestListTenants(t *testing.T) {
	repo := newMockTenantRepository()
	publisher := &mockPublisher{}
	logger := zap.NewNop().Sugar()
	service := services.NewTenantService(repo, publisher, logger)

	repo.Create(context.Background(), &models.Tenant{ID: uuid.NewString(), Name: "Acme", OwnerID: "owner-1", IsActive: true})
	repo.Create(context.Background(), &models.Tenant{ID: uuid.NewString(), Name: "Beta", OwnerID: "owner-2", IsActive: true})

	list, err := service.ListTenants(context.Background(), "owner-1", 0, 10)
	if err != nil {
		t.Fatalf("expected no error, got %v", err)
	}
	if len(list) != 1 {
		t.Fatalf("expected 1 tenant, got %d", len(list))
	}
}
