package psql

import (
	"articlesManageService/internal/domain/models"
	"articlesManageService/internal/storage"
	"articlesManageService/pkg/lib/logger/sl"
	"context"
	"database/sql"
	"errors"
	"fmt"
	"log/slog"
	"os"
	"path/filepath"

	"github.com/google/uuid"
	"github.com/pressly/goose/v3"

	"github.com/lib/pq"
)

const (
	ArticlesTableName = "Articles"
)

type PsqlStorage struct {
	log *slog.Logger
	DB  *sql.DB
}

func New(log *slog.Logger, connStr string) *PsqlStorage {
	db, err := sql.Open("postgres", connStr)
	if err != nil {
		log.Error("error opening database", sl.Err(err))
		panic(err)
	}

	wd, _ := os.Getwd()
	migrationPath := filepath.Join(wd, "app", "migrations")
	if err := applyMigrations(db, migrationPath); err != nil {
		panic(err)
	}

	return &PsqlStorage{
		log: log,
		DB:  db,
	}
}

func applyMigrations(db *sql.DB, migrationsPath string) error {
	return goose.Up(db, migrationsPath)
}

func (s *PsqlStorage) Close() {
	err := s.DB.Close()
	if err != nil {
		s.log.Info("Error closing the database: %v", sl.Err(err))
	}
}

func (s *PsqlStorage) GetArticles(ctx context.Context) ([]models.Article, error) {
	const op = "psql.getArticles"
	log := s.log.With(
		slog.String("op", op),
	)

	rows, err := s.DB.QueryContext(ctx, `
		SELECT * FROM `+ArticlesTableName+`;
	`)
	if err != nil {
		log.Error("error retrieving all articles", sl.Err(err))
		return nil, fmt.Errorf("%s: %w", op, err)
	}
	defer rows.Close()

	articles := make([]models.Article, 0, 5)
	for rows.Next() {
		var article models.Article
		err := rows.Scan(&article.Id, &article.CreatedAt, &article.Title, &article.Content, &article.OwnerId)
		if err != nil {
			log.Warn("Error scaning row", sl.Err(err))
			continue
		}
		articles = append(articles, article)
	}

	return articles, nil
}

func (s *PsqlStorage) GetArticleById(ctx context.Context, aid uuid.UUID) (models.Article, error) {
	const op = "psql.getArticleById"
	log := s.log.With(
		slog.String("op", op),
	)

	row := s.DB.QueryRowContext(ctx, `
		SELECT * FROM `+ArticlesTableName+`
		WHERE id=$1;
	`, aid)

	var article models.Article
	err := row.Scan(&article.Id, &article.CreatedAt, &article.Title, &article.Content, &article.OwnerId)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			log.Warn("Article with current id not found", sl.Err(err))
			return models.Article{}, fmt.Errorf("%s: %w", op, storage.ErrNotFound)
		}

		log.Error("Error scaning row", sl.Err(err))
		return models.Article{}, fmt.Errorf("%s: %w", op, err)
	}

	return article, nil
}

func (s *PsqlStorage) GetArticlesByOwnerId(ctx context.Context, uid uuid.UUID) ([]models.Article, error) {
	const op = "psql.getArticlesByOwnerId"
	log := s.log.With(
		slog.String("op", op),
	)

	rows, err := s.DB.QueryContext(ctx, `
		SELECT * FROM `+ArticlesTableName+`
		WHERE owner_id=$1;
	`, uid)
	if err != nil {
		log.Info("Error retriening articles by owner_id")
		return nil, fmt.Errorf("%s: %w", op, err)
	}
	defer rows.Close()

	articles := make([]models.Article, 0, 5)
	for rows.Next() {
		var article models.Article
		err := rows.Scan(&article.Id, &article.CreatedAt, &article.Title, &article.Content, &article.OwnerId)
		if err != nil {
			log.Warn("Error scaning row", sl.Err(err))
			continue
		}

		articles = append(articles, article)
	}

	return articles, nil
}

func (s *PsqlStorage) Insert(ctx context.Context, article models.Article) (models.Article, error) {
	const op = "psql.insert"
	log := s.log.With(
		slog.String("op", op),
	)

	_, err := s.DB.ExecContext(ctx, `
		INSERT INTO `+ArticlesTableName+` (id, created_at, title, content, owner_id)
		VALUES ($1, $2, $3, $4, $5);
	`, article.Id, article.CreatedAt, article.Title, article.Content, article.OwnerId)

	if err != nil {
		if pgErr, ok := err.(*pq.Error); ok && pgErr.Code == "23505" {
			log.Error("Article with this ID already exists", sl.Err(err))
			return models.Article{}, fmt.Errorf("%s: %w", op, storage.ErrAlreadyExists)
		}

		log.Error("Error inserting article", sl.Err(err))
		return models.Article{}, fmt.Errorf("%s: %w", op, err)
	}

	return article, nil
}

func (s *PsqlStorage) Update(ctx context.Context, aid uuid.UUID, article models.Article) (models.Article, error) {
	const op = "psql.update"
	log := s.log.With(
		slog.String("op", op),
	)

	result, err := s.DB.ExecContext(ctx, `
		UPDATE `+ArticlesTableName+` SET 
			title = $1, 
			content = $2  
		WHERE id = $3;
	`, article.Title, article.Content, aid)
	if err != nil {
		log.Error("Error updating article", sl.Err(err))
		return models.Article{}, fmt.Errorf("%s: %w", op, err)
	}

	rowsAffected, err := result.RowsAffected()
	if err != nil {
		log.Error("Error get rows affected", sl.Err(err))
		return models.Article{}, fmt.Errorf("%s: %w", op, err)
	}

	if rowsAffected == 0 {
		log.Error("Zero rows affected")
		return models.Article{}, fmt.Errorf("%s: %w", op, storage.ErrNotFound)
	}

	return article, nil
}

func (s *PsqlStorage) Delete(ctx context.Context, aid uuid.UUID) (models.Article, error) {
	const op = "psql.delete"
	log := s.log.With(
		slog.String("op", op),
	)

	article, err := s.GetArticleById(ctx, aid)
	if err != nil {
		if errors.Is(err, storage.ErrNotFound) {
			log.Warn("Article not found", sl.Err(err))
			return models.Article{}, fmt.Errorf("%s: %w", op, storage.ErrNotFound)
		}
		log.Error("Error getting article before deliting", sl.Err(err))
		return models.Article{}, fmt.Errorf("%s: %w", op, err)
	}

	_, err = s.DB.ExecContext(ctx, `
		DELETE FROM `+ArticlesTableName+`
		WHERE id=$1;
	`, aid)

	if err != nil {
		log.Error("Error deleting article", sl.Err(err))
		return models.Article{}, fmt.Errorf("%s: %w", op, err)
	}

	return article, nil
}
