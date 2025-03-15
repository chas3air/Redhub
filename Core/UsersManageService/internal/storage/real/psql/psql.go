package psqlstorage

import (
	"context"
	"database/sql"
	"errors"
	"fmt"
	"log/slog"
	"usersManageService/internal/domain/models"
	"usersManageService/internal/storage"
	storage_error "usersManageService/internal/storage"
	"usersManageService/pkg/lib/logger/sl"

	"github.com/google/uuid"
	"github.com/lib/pq"
)

type PsqlStorage struct {
	log *slog.Logger
	DB  *sql.DB
}

const UsersTableName = "users"

func New(log *slog.Logger, connStr string) *PsqlStorage {
	const op = "psql.New"
	db, err := sql.Open("postgres", connStr)
	if err != nil {
		log.With(slog.String("op", op)).Error("Error connecting to DB", sl.Err(err))
		panic(err)
	}

	return &PsqlStorage{
		log: log,
		DB:  db,
	}
}

func (ps *PsqlStorage) GetUsers(ctx context.Context) ([]models.User, error) {
	const op = "storage.psql.getUsers"
	log := ps.log.With(
		slog.String("op", op),
	)

	select {
	case <-ctx.Done():
		return nil, fmt.Errorf("%s: %w", op, ctx.Err())
	default:
	}

	rows, err := ps.DB.QueryContext(ctx, `SELECT id, email, password, role, nick, description, birthday FROM `+UsersTableName+`;`)
	if err != nil {
		log.Error("Error retrieving all users", sl.Err(err))
		return nil, fmt.Errorf("%s: %w", op, err)
	}
	defer rows.Close()

	var users []models.User
	for rows.Next() {
		var user models.User
		if err := rows.Scan(&user.Id, &user.Email, &user.Password, &user.Role, &user.Nick, &user.Description, &user.Birthday); err != nil {
			log.Error("-Error scanning row", sl.Err(err))
			continue
		}

		users = append(users, user)
	}

	return users, nil
}

func (ps *PsqlStorage) GetUserById(ctx context.Context, uid uuid.UUID) (models.User, error) {
	const op = "storage.psql.getUserById"
	log := ps.log.With(
		slog.String("op", op),
	)

	select {
	case <-ctx.Done():
		return models.User{}, fmt.Errorf("%s: %w", op, ctx.Err())
	default:
	}

	var user models.User
	err := ps.DB.QueryRowContext(ctx, `SELECT id, email, password, role, nick, description, birthday FROM `+UsersTableName+` WHERE id = $1;`, uid).
		Scan(&user.Id, &user.Email, &user.Password, &user.Role, &user.Nick, &user.Description, &user.Birthday)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			log.Error("User with current id not found", sl.Err(err))
			return models.User{}, fmt.Errorf("%s: %w", op, storage_error.ErrNotFound)
		}

		log.Error("Error scanning row", sl.Err(err))
		return models.User{}, fmt.Errorf("%s: %w", op, err)
	}

	return user, nil
}

func (ps *PsqlStorage) GetUserByEmail(ctx context.Context, email string) (models.User, error) {
	const op = "storage.psql.getUserByEmail"
	log := ps.log.With(
		slog.String("op", op),
	)

	select {
	case <-ctx.Done():
		return models.User{}, fmt.Errorf("%s: %w", op, ctx.Err())
	default:
	}

	var user models.User
	err := ps.DB.QueryRowContext(ctx, `SELECT id, email, password, role, nick, description, birthday FROM `+UsersTableName+` WHERE email = $1;`, email).
		Scan(&user.Id, &user.Email, &user.Password, &user.Role, &user.Nick, &user.Description, &user.Birthday)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			log.Warn("User with current email not found", sl.Err(err))
			return models.User{}, fmt.Errorf("%s: %w", op, storage_error.ErrNotFound)
		}

		log.Error("Error scanning row", sl.Err(err))
		return models.User{}, fmt.Errorf("%s: %w", op, err)
	}

	return user, nil
}

func (ps *PsqlStorage) Insert(ctx context.Context, user models.User) (models.User, error) {
	const op = "storage.psql.insert"
	log := ps.log.With(
		slog.String("op", op),
	)

	select {
	case <-ctx.Done():
		return models.User{}, fmt.Errorf("%s: %w", op, ctx.Err())
	default:
	}

	_, err := ps.DB.ExecContext(ctx, `
		INSERT INTO `+UsersTableName+` (id, email, password, role, nick, description, birthday) 
		VALUES ($1, $2, $3, $4, $5, $6, $7);`, user.Id, user.Email, user.Password, user.Role, user.Nick, user.Description, user.Birthday)

	if err != nil {
		if pgErr, ok := err.(*pq.Error); ok && pgErr.Code == "23505" {
			log.Error("User with this ID already exists", sl.Err(err))
			return models.User{}, fmt.Errorf("%s: %w", op, storage.ErrAlreadyExists)
		}

		log.Error("Error inserting user", sl.Err(err))
		return models.User{}, fmt.Errorf("%s: %w", op, err)
	}

	return user, nil
}

func (ps *PsqlStorage) Update(ctx context.Context, uid uuid.UUID, user models.User) (models.User, error) {
	const op = "storage.psql.update"
	log := ps.log.With(
		slog.String("op", op),
	)

	select {
	case <-ctx.Done():
		return models.User{}, fmt.Errorf("%s: %w", op, ctx.Err())
	default:
	}

	result, err := ps.DB.ExecContext(ctx, `
		UPDATE `+UsersTableName+` 
		SET email = $1, password = $2, role = $3, nick = $4, description = $5, birthday = $6 
		WHERE id = $7;`,
		user.Email, user.Password, user.Role, user.Nick, user.Description, user.Birthday, uid)
	if err != nil {
		log.Error("Error updating user", sl.Err(err))
		return models.User{}, fmt.Errorf("%s: %w", op, err)
	}

	rowsAffected, err := result.RowsAffected()
	if err != nil {
		log.Error("Error get rows affected", sl.Err(err))
		return models.User{}, fmt.Errorf("%s: %w", op, err)
	}

	if rowsAffected == 0 {
		log.Error("Zero rows affected")
		return models.User{}, fmt.Errorf("%s: %w", op, storage.ErrNotFound)
	}

	return user, nil
}

func (ps *PsqlStorage) Delete(ctx context.Context, uid uuid.UUID) (models.User, error) {
	const op = "storage.psql.delete"
	log := ps.log.With(
		slog.String("op", op),
	)

	select {
	case <-ctx.Done():
		return models.User{}, fmt.Errorf("%s: %w", op, ctx.Err())
	default:
	}

	user, err := ps.GetUserById(ctx, uid)
	if err != nil {
		if errors.Is(err, storage.ErrNotFound) {
			log.Warn("User not found", sl.Err(err))
			return models.User{}, fmt.Errorf("%s: %w", op, storage.ErrNotFound)
		}

		log.Error("Error getting user before deliting", sl.Err(err))
		return models.User{}, fmt.Errorf("%s: %w", op, err)
	}

	_, err = ps.DB.ExecContext(ctx, `
		DELETE FROM `+UsersTableName+` 
		WHERE id = $1;
	`, uid)
	if err != nil {
		log.Error("Error deleting user", sl.Err(err))
		return models.User{}, fmt.Errorf("%s: %w", op, err)
	}

	return user, nil
}
