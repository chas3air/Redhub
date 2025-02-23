package mockauth

import (
	"apigateway/internal/domain/models"
	"apigateway/internal/storage"
	"context"

	"github.com/google/uuid"
)

type MockAuth struct {
	users map[uuid.UUID]models.User
}

func NewMockAuth() *MockAuth {
	return &MockAuth{
		users: make(map[uuid.UUID]models.User),
	}
}

func (a *MockAuth) Login(ctx context.Context, email string, password string, appID uuid.UUID) (token string, err error) {
	for _, v := range a.users {
		if v.Email == email && v.Password == password {
			return uuid.New().String(), nil
		}
	}

	return uuid.Nil.String(), storage.ErrUserNotFound
}

func (a *MockAuth) Register(ctx context.Context, user models.User) (err error) {
	for _, v := range a.users {
		if v.Email == user.Email && v.Password == user.Password {
			return storage.ErrUserExists
		}
	}

	a.users[user.Id] = user
	return nil
}

func (a *MockAuth) IsAdmin(ctx context.Context, userID uuid.UUID) (isAdmin bool, err error) {
	for _, user := range a.users {
		if user.Id == userID {
			return user.Role == "admin", nil
		}
	}

	return false, storage.ErrUserNotFound
}
