package nullchecker

import (
	"auth/internal/domain/models"

	"github.com/google/uuid"
)

func IsUserNullChecker(user models.User) bool {
	return user.Id != uuid.UUID{} &&
		user.Email != "" &&
		user.Password != "" &&
		user.Role != "" &&
		user.Nick != ""
}
