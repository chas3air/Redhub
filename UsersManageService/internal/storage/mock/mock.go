package mock

import (
	"context"
	"fmt"
	"log/slog"
	"usersManageService/internal/domain/models"
	"usersManageService/internal/storage"

	"github.com/google/uuid"
)

type MockStorage struct {
	users []models.User
	log   *slog.Logger
}

func New(log *slog.Logger) *MockStorage {
	return &MockStorage{
		users: make([]models.User, 0),
		log:   log,
	}
}

// GetUsers implements storage.Storage.
func (m *MockStorage) GetUsers(ctx context.Context) ([]models.User, error) {
	const op = "storage.mock.GetUsers"
	m.log.Info("Fetching users", slog.String("operation", op), slog.String("error", "nil"))

	return m.users, nil
}

// GetUserById implements storage.Storage.
func (m *MockStorage) GetUserById(ctx context.Context, id uuid.UUID) (models.User, error) {

	const op = "storage.mock.GetUserById"
	m.log.Info("Fetching user by ID", slog.String("operation", op), slog.String("userId", id.String()), slog.String("error", "nil"))

	for _, v := range m.users {
		if v.Id == id {
			m.log.Info("User found", slog.String("operation", op), slog.String("userId", id.String()), slog.String("error", "nil"), slog.Any("additional info", []map[string]interface{}{
				{"user": v},
			}))

			return v, nil
		}
	}

	err := fmt.Errorf("%s: %w", op, storage.ErrUserNotFound)
	m.log.Warn("User not found", slog.String("operation", op), slog.String("userId", id.String()), slog.String("error", err.Error()))
	return models.User{}, err
}

// GetUserByEmail implements storage.Storage.
func (m *MockStorage) GetUserByEmail(ctx context.Context, email string) (models.User, error) {
	const op = "storage.mock.GetUserByEmail"
	m.log.Info("Fetching user by email", slog.String("operation", op), slog.String("email", email), slog.String("error", "nil"))

	for _, v := range m.users {
		if v.Email == email {
			m.log.Info("User found", slog.String("operation", op), slog.String("email", email), slog.String("error", "nil"), slog.Any("additional info", []map[string]interface{}{
				{"user": v},
			}))

			return v, nil
		}
	}

	err := fmt.Errorf("%s: %w", op, storage.ErrUserNotFound)
	m.log.Warn("User not found", slog.String("operation", op), slog.String("email", email), slog.String("error", err.Error()))
	return models.User{}, err
}

// Insert implements storage.Storage.
func (m *MockStorage) Insert(ctx context.Context, user models.User) error {
	const op = "storage.mock.Insert"
	m.log.Info("Inserting user", slog.String("operation", op), slog.Any("additional info", []map[string]interface{}{
		{"user": user},
	}), slog.String("error", "nil"))

	m.users = append(m.users, user)
	m.log.Info("User inserted successfully", slog.String("operation", op), slog.Any("additional info", []map[string]interface{}{
		{"user": user},
	}), slog.String("error", "nil"))

	return nil
}

// Update implements storage.Storage.
func (m *MockStorage) Update(ctx context.Context, id uuid.UUID, user models.User) error {
	const op = "storage.mock.Update"
	m.log.Info("Updating user", slog.String("operation", op), slog.String("userId", id.String()), slog.Any("additional info", []map[string]interface{}{
		{"user": user},
	}), slog.String("error", "nil"))

	for i, v := range m.users {
		if v.Id == id {
			m.users[i] = user
			m.log.Info("User updated successfully", slog.String("operation", op), slog.String("userId", id.String()), slog.Any("additional info", []map[string]interface{}{
				{"user": user},
			}), slog.String("error", "nil"))
			return nil
		}
	}

	err := fmt.Errorf("%s: %w", op, storage.ErrUserNotFound)
	m.log.Warn("User not found for update", slog.String("operation", op), slog.String("userId", id.String()), slog.String("error", err.Error()))
	return err
}

// Delete implements storage.Storage.
func (m *MockStorage) Delete(ctx context.Context, id uuid.UUID) (models.User, error) {
	const op = "storage.mock.Delete"
	m.log.Info("Deleting user", slog.String("operation", op), slog.String("userId", id.String()), slog.String("error", "nil"))

	for i, v := range m.users {
		if v.Id == id {
			m.users = append(m.users[:i], m.users[i+1:]...)
			m.log.Info("User deleted successfully", slog.String("operation", op), slog.String("userId", id.String()), slog.Any("additional info", []map[string]interface{}{
				{"user": v},
			}), slog.String("error", "nil"))
			return v, nil
		}
	}

	err := fmt.Errorf("%s: %w", op, storage.ErrUserNotFound)
	m.log.Warn("User not found for deletion", slog.String("operation", op), slog.String("userId", id.String()), slog.String("error", err.Error()))
	return models.User{}, err
}
