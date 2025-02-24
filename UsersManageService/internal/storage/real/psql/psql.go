package psqlstorage

import (
	"context"
	"fmt"
	"usersManageService/internal/domain/models"

	"github.com/google/uuid"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
	"gorm.io/gorm/logger"
)

type PsqlStorage struct {
	DB *gorm.DB
}

func New(connStr string) *PsqlStorage {
	db, err := gorm.Open(postgres.Open(connStr), &gorm.Config{
		Logger: logger.Default.LogMode(logger.Info),
	})
	if err != nil {
		panic(err)
	}

	db.AutoMigrate(&models.User{})

	return &PsqlStorage{
		DB: db,
	}
}

func (ps *PsqlStorage) GetUsers(ctx context.Context) ([]models.User, error) {
	const op = "storage.psql.getUsers"

	select {
	case <-ctx.Done():
		return nil, fmt.Errorf("%s: %w", op, ctx.Err())
	default:
	}

	var users []models.User
	result := ps.DB.Find(&users)
	if result.Error != nil {
		return nil, fmt.Errorf("%s: %w", op, result.Error)
	}

	return users, nil
}

func (ps *PsqlStorage) GetUserById(ctx context.Context, uid uuid.UUID) (models.User, error) {
	const op = "storage.psql.getUserById"

	select {
	case <-ctx.Done():
		return models.User{}, fmt.Errorf("%s: %w", op, ctx.Err())
	default:
	}

	var user models.User
	result := ps.DB.First(&user, uid)
	if result.Error != nil {
		if result.Error == gorm.ErrRecordNotFound {
			return models.User{}, fmt.Errorf("%s: user not found: %w", op, result.Error)
		}
		return models.User{}, fmt.Errorf("%s: %w", op, result.Error)
	}

	return user, nil
}

func (ps *PsqlStorage) GetUserByEmail(ctx context.Context, email string) (models.User, error) {
	const op = "storage.psql.getUserByEmail"

	select {
	case <-ctx.Done():
		return models.User{}, fmt.Errorf("%s: %w", op, ctx.Err())
	default:
	}

	var user models.User
	result := ps.DB.Where("email = ?", email).First(&user)
	if result.Error != nil {
		if result.Error == gorm.ErrRecordNotFound {
			return models.User{}, fmt.Errorf("%s: user not found: %w", op, result.Error)
		}
		return models.User{}, fmt.Errorf("%s: %w", op, result.Error)
	}

	return user, nil
}

func (ps *PsqlStorage) Insert(ctx context.Context, user models.User) error {
	const op = "storage.psql.insert"

	select {
	case <-ctx.Done():
		return fmt.Errorf("%s: %w", op, ctx.Err())
	default:
	}

	result := ps.DB.Create(&user)
	if result.Error != nil {
		return fmt.Errorf("%s: %w", op, result.Error)
	}

	return nil
}

func (ps *PsqlStorage) Update(ctx context.Context, uid uuid.UUID, user models.User) error {
	const op = "storage.psql.update"

	select {
	case <-ctx.Done():
		return fmt.Errorf("%s: %w", op, ctx.Err())
	default:
	}

	result := ps.DB.Model(&models.User{}).Where("id = ?", uid).Updates(user)
	if result.Error != nil {
		return fmt.Errorf("%s: %w", op, result.Error)
	}

	return nil
}

func (ps *PsqlStorage) Delete(ctx context.Context, uid uuid.UUID) (models.User, error) {
	const op = "storage.psql.delete"

	select {
	case <-ctx.Done():
		return models.User{}, fmt.Errorf("%s: %w", op, ctx.Err())
	default:
	}

	var user models.User
	if err := ps.DB.First(&user, uid).Error; err != nil {
		if err == gorm.ErrRecordNotFound {
			return models.User{}, fmt.Errorf("%s: user not found: %w", op, err)
		}
		return models.User{}, fmt.Errorf("%s: %w", op, err)
	}

	if err := ps.DB.Delete(&user).Error; err != nil {
		return models.User{}, fmt.Errorf("%s: %w", op, err)
	}

	return user, nil
}
