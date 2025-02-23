package mockusers

import (
	"auth/internal/domain/models"
	"auth/internal/storage"
	"context"
	"time"

	"github.com/google/uuid"
)

type MockStorage struct {
	users map[uuid.UUID]models.User
}

func New() *MockStorage {
	generated_id := uuid.New()
	return &MockStorage{
		users: map[uuid.UUID]models.User{
			generated_id: models.User{
				Id:       generated_id,
				Email:    "testuser@example.com",
				Password: "securepassword",
				Role:     "admin",
				Nick:     "test_nick",
				Birthday: time.Date(1990, 1, 1, 0, 0, 0, 0, time.UTC),
			},
		},
	}
}

// GetUsers implements storage.Storage.
func (m *MockStorage) GetUsers(ctx context.Context) ([]models.User, error) {
	var userList []models.User
	for _, user := range m.users {
		userList = append(userList, user)
	}
	return userList, nil
}

// GetUserById implements storage.Storage.
func (m *MockStorage) GetUserById(ctx context.Context, id uuid.UUID) (models.User, error) {
	user, exists := m.users[id]
	if !exists {
		return models.User{}, storage.ErrUserNotFound
	}
	return user, nil
}

// GetUserByEmail implements storage.Storage.
func (m *MockStorage) GetUserByEmail(ctx context.Context, email string) (models.User, error) {
	for _, v := range m.users {
		if v.Email == email {
			return v, nil
		}
	}

	return models.User{}, storage.ErrUserNotFound
}

// Insert implements storage.Storage.
func (m *MockStorage) Insert(ctx context.Context, user models.User) error {
	if _, exists := m.users[user.Id]; exists {
		return storage.ErrUserExists
	}
	m.users[user.Id] = user
	return nil
}

// Update implements storage.Storage.
func (m *MockStorage) Update(ctx context.Context, id uuid.UUID, user models.User) error {
	if _, exists := m.users[id]; !exists {
		return storage.ErrUserNotFound
	}
	m.users[id] = user
	return nil
}

// Delete implements storage.Storage.
func (m *MockStorage) Delete(ctx context.Context, id uuid.UUID) (models.User, error) {
	user, exists := m.users[id]
	if !exists {
		return models.User{}, storage.ErrUserNotFound
	}
	delete(m.users, id)
	return user, nil
}
