package mockusersmanager

import (
	"apigateway/internal/domain/models"
	"apigateway/internal/storage"
	"context"

	"github.com/google/uuid"
)

type MockStorage struct {
	users map[uuid.UUID]models.User
}

func New() *MockStorage {
	return &MockStorage{
		users: make(map[uuid.UUID]models.User),
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
		return models.User{}, storage.ErrNotFound
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

	return models.User{}, storage.ErrNotFound
}

// Insert implements storage.Storage.
func (m *MockStorage) Insert(ctx context.Context, user models.User) error {
	if _, exists := m.users[user.Id]; exists {
		return storage.ErrAlreadyExists
	}
	m.users[user.Id] = user
	return nil
}

// Update implements storage.Storage.
func (m *MockStorage) Update(ctx context.Context, id uuid.UUID, user models.User) error {
	if _, exists := m.users[id]; !exists {
		return storage.ErrNotFound
	}
	m.users[id] = user
	return nil
}

// Delete implements storage.Storage.
func (m *MockStorage) Delete(ctx context.Context, id uuid.UUID) (models.User, error) {
	user, exists := m.users[id]
	if !exists {
		return models.User{}, storage.ErrNotFound
	}
	delete(m.users, id)
	return user, nil
}
