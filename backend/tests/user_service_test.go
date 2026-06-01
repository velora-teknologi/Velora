package services_test

import (
	"context"
	"testing"

	"github.com/golang-jwt/jwt/v5"
	"github.com/google/uuid"
	"github.com/velora-teknologi/velora/internal/application/dtos"
	"github.com/velora-teknologi/velora/internal/application/services"
	"github.com/velora-teknologi/velora/internal/domain/models"
	"go.uber.org/zap"
	"golang.org/x/crypto/bcrypt"
)

// MockUserRepository is a mock implementation of UserRepository
// used for service unit tests.
type MockUserRepository struct {
	users map[string]*models.User
}

func NewMockUserRepository() *MockUserRepository {
	return &MockUserRepository{
		users: make(map[string]*models.User),
	}
}

func (m *MockUserRepository) Create(ctx context.Context, user *models.User) error {
	if user.ID == "" {
		user.ID = uuid.NewString()
	}
	m.users[user.ID] = user
	return nil
}

func (m *MockUserRepository) GetByID(ctx context.Context, id string) (*models.User, error) {
	return m.users[id], nil
}

func (m *MockUserRepository) GetByEmail(ctx context.Context, email string) (*models.User, error) {
	for _, user := range m.users {
		if user.Email == email {
			return user, nil
		}
	}
	return nil, nil
}

func (m *MockUserRepository) Update(ctx context.Context, user *models.User) error {
	m.users[user.ID] = user
	return nil
}

func (m *MockUserRepository) Delete(ctx context.Context, id string) error {
	delete(m.users, id)
	return nil
}

func (m *MockUserRepository) List(ctx context.Context, skip, limit int) ([]*models.User, error) {
	var users []*models.User
	for _, user := range m.users {
		users = append(users, user)
	}
	return users, nil
}
func (m *MockUserRepository) ListByTenantID(ctx context.Context, tenantID string, skip, limit int) ([]*models.User, error) {
	var users []*models.User
	for _, user := range m.users {
		if user.TenantID == tenantID {
			users = append(users, user)
		}
	}
	return users, nil
}

func TestCreateUser(t *testing.T) {
	repo := NewMockUserRepository()
	logger := zap.NewNop().Sugar()
	service := services.NewUserService(repo, "test-secret", logger)

	req := &dtos.CreateUserRequest{
		Email:    "test@example.com",
		Password: "password123",
		FullName: "Test User",
	}

	resp, err := service.CreateUser(context.Background(), req)
	if err != nil {
		t.Fatalf("expected no error, got %v", err)
	}

	if resp.Email != req.Email || resp.FullName != req.FullName {
		t.Fatalf("unexpected response: got %+v", resp)
	}

	saved, err := repo.GetByEmail(context.Background(), req.Email)
	if err != nil {
		t.Fatalf("unexpected repository error: %v", err)
	}
	if saved == nil {
		t.Fatal("expected user to be saved")
	}
	if bcrypt.CompareHashAndPassword([]byte(saved.Password), []byte(req.Password)) != nil {
		t.Fatal("saved password was not hashed correctly")
	}
}

func TestGetUser(t *testing.T) {
	repo := NewMockUserRepository()
	logger := zap.NewNop().Sugar()
	service := services.NewUserService(repo, "test-secret", logger)

	user := &models.User{
		ID:       uuid.NewString(),
		Email:    "get@example.com",
		FullName: "Get User",
		Password: "hashed-password",
	}
	repo.Create(context.Background(), user)

	resp, err := service.GetUser(context.Background(), user.ID)
	if err != nil {
		t.Fatalf("expected no error, got %v", err)
	}
	if resp.Email != user.Email || resp.FullName != user.FullName {
		t.Fatalf("unexpected response: got %+v", resp)
	}
}

func TestUpdateUser(t *testing.T) {
	repo := NewMockUserRepository()
	logger := zap.NewNop().Sugar()
	service := services.NewUserService(repo, "test-secret", logger)

	user := &models.User{
		ID:       uuid.NewString(),
		Email:    "update@example.com",
		FullName: "Old Name",
		Password: "hashed-password",
	}
	repo.Create(context.Background(), user)

	req := &dtos.UpdateUserRequest{
		FullName: "Updated Name",
	}

	resp, err := service.UpdateUser(context.Background(), user.ID, req)
	if err != nil {
		t.Fatalf("expected no error, got %v", err)
	}
	if resp.FullName != req.FullName {
		t.Fatalf("expected full name to be updated, got %s", resp.FullName)
	}
}

func TestDeleteUser(t *testing.T) {
	repo := NewMockUserRepository()
	logger := zap.NewNop().Sugar()
	service := services.NewUserService(repo, "test-secret", logger)

	user := &models.User{
		ID:       uuid.NewString(),
		Email:    "delete@example.com",
		FullName: "Delete User",
		Password: "hashed-password",
	}
	repo.Create(context.Background(), user)

	if err := service.DeleteUser(context.Background(), user.ID); err != nil {
		t.Fatalf("expected no error, got %v", err)
	}

	deleted, _ := repo.GetByID(context.Background(), user.ID)
	if deleted != nil {
		t.Fatal("expected user to be deleted")
	}
}

func TestLoginUser(t *testing.T) {
	repo := NewMockUserRepository()
	logger := zap.NewNop().Sugar()
	service := services.NewUserService(repo, "test-secret", logger)

	password := "login-password"
	hashed, err := bcrypt.GenerateFromPassword([]byte(password), bcrypt.DefaultCost)
	if err != nil {
		t.Fatalf("failed to hash password: %v", err)
	}

	user := &models.User{
		ID:       uuid.NewString(),
		Email:    "login@example.com",
		FullName: "Login User",
		Password: string(hashed),
	}
	repo.Create(context.Background(), user)

	resp, err := service.LoginUser(context.Background(), &dtos.LoginRequest{
		Email:    user.Email,
		Password: password,
	})
	if err != nil {
		t.Fatalf("expected no error, got %v", err)
	}

	if resp.Token == "" {
		t.Fatal("expected JWT token, got empty string")
	}
	if resp.User.Email != user.Email {
		t.Fatalf("expected user email %s, got %s", user.Email, resp.User.Email)
	}

	parsed, err := jwt.Parse(resp.Token, func(token *jwt.Token) (interface{}, error) {
		return []byte("test-secret"), nil
	})
	if err != nil || !parsed.Valid {
		t.Fatalf("expected valid JWT token, got error: %v", err)
	}
}
