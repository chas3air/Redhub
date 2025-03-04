package psql

import (
	"articlesManageService/internal/domain/models"
	"articlesManageService/internal/storage"
	"context"
	"database/sql"
	"errors"
	"fmt"
	"log"

	"github.com/google/uuid"

	_ "github.com/lib/pq"
)

const (
	ArticlesTableName = "Articles"
)

type PsqlStorage struct {
	DB *sql.DB
}

func New(connStr string) *PsqlStorage {
	db, err := sql.Open("postgres", connStr)
	if err != nil {
		log.Printf("Error opening database: %v", err)
		panic(err)
	}

	return &PsqlStorage{
		DB: db,
	}
}

func (s *PsqlStorage) Close() {
	err := s.DB.Close()
	if err != nil {
		log.Printf("Error closing the database: %v", err)
	}
}

func (s *PsqlStorage) GetArticles(ctx context.Context) ([]models.Article, error) {
	const op = "psql.getArticles"

	rows, err := s.DB.QueryContext(ctx, `
		SELECT * FROM `+ArticlesTableName+`;
	`)
	if err != nil {
		log.Printf("%s: %v", op, err)
		return nil, fmt.Errorf("%s: %w", op, err)
	}
	defer rows.Close()

	articles := make([]models.Article, 0, 5)
	for rows.Next() {
		var article models.Article
		err := rows.Scan(&article.Id, &article.CreatedAt, &article.Title, &article.Content, &article.OwnerId)
		if err != nil {
			log.Printf("%s: %v", op, err)
			return nil, fmt.Errorf("%s: %w", op, err)
		}
		articles = append(articles, article)
	}

	return articles, nil
}

func (s *PsqlStorage) GetArticleById(ctx context.Context, aid uuid.UUID) (models.Article, error) {
	const op = "psql.getArticleById"

	row := s.DB.QueryRowContext(ctx, `
		SELECT * FROM `+ArticlesTableName+`
		WHERE id=$1;
	`, aid)

	var article models.Article
	err := row.Scan(&article.Id, &article.CreatedAt, &article.Title, &article.Content, &article.OwnerId)
	if err != nil {
		if err == sql.ErrNoRows {
			log.Printf("%s: %v", op, storage.ErrNotFound)
			return models.Article{}, fmt.Errorf("%s: %w", op, storage.ErrNotFound)
		}
		log.Printf("%s: %v", op, err)
		return models.Article{}, fmt.Errorf("%s: %w", op, err)
	}

	return article, nil
}

func (s *PsqlStorage) GetArticlesByOwnerId(ctx context.Context, uid uuid.UUID) ([]models.Article, error) {
	const op = "psql.getArticlesByOwnerId"

	rows, err := s.DB.QueryContext(ctx, `
		SELECT * FROM `+ArticlesTableName+`
		WHERE owner_id=$1;
	`, uid)
	if err != nil {
		log.Printf("%s: %v", op, err)
		return nil, fmt.Errorf("%s: %w", op, err)
	}
	defer rows.Close()

	articles := make([]models.Article, 0, 5)
	for rows.Next() {
		var article models.Article
		err := rows.Scan(&article.Id, &article.CreatedAt, &article.Title, &article.Content, &article.OwnerId)
		if err != nil {
			log.Printf("%s: %v", op, err)
			return nil, fmt.Errorf("%s: %w", op, err)
		}
		articles = append(articles, article)
	}

	return articles, nil
}

func (s *PsqlStorage) Insert(ctx context.Context, article models.Article) error {
	const op = "psql.insert"

	_, err := s.GetArticleById(ctx, article.Id)
	if err == nil {
		log.Printf("%s: %v", op, storage.ErrAlreadyExists)
		return fmt.Errorf("%s: %w", op, storage.ErrAlreadyExists)
	} else if err != nil && !errors.Is(err, storage.ErrNotFound) {
		log.Printf("%s: %v", op, err)
		return fmt.Errorf("%s: %w", op, err)
	}

	_, err = s.DB.ExecContext(ctx, `
		INSERT INTO `+ArticlesTableName+` (id, created_at, title, content, owner_id)
		VALUES ($1, $2, $3, $4, $5);
	`, article.Id, article.CreatedAt, article.Title, article.Content, article.OwnerId)
	if err != nil {
		log.Printf("%s: %v", op, err)
		return fmt.Errorf("%s: %w", op, err)
	}

	return nil
}

func (s *PsqlStorage) Update(ctx context.Context, aid uuid.UUID, article models.Article) error {
	const op = "psql.update"

	_, err := s.GetArticleById(ctx, aid)
	if err != nil {
		if errors.Is(err, storage.ErrNotFound) {
			log.Printf("%s: %v", op, storage.ErrNotFound)
			return fmt.Errorf("%s: %w", op, storage.ErrNotFound)
		}
		log.Printf("%s: %v", op, err)
		return fmt.Errorf("%s: %w", op, err)
	}

	_, err = s.DB.ExecContext(ctx, `
		UPDATE `+ArticlesTableName+` SET 
			title = $1, 
			content = $2  
		WHERE id = $3;
	`, article.Title, article.Content, aid)
	if err != nil {
		log.Printf("%s: %v", op, err)
		return fmt.Errorf("%s: %w", op, err)
	}

	return nil
}

func (s *PsqlStorage) Delete(ctx context.Context, aid uuid.UUID) (models.Article, error) {
	const op = "psql.delete"

	article, err := s.GetArticleById(ctx, aid)
	if err != nil {
		log.Printf("%s: %v", op, err)
		return models.Article{}, fmt.Errorf("%s: %w", op, err)
	}

	_, err = s.DB.ExecContext(ctx, `
		DELETE FROM `+ArticlesTableName+`
		WHERE id=$1;
	`, aid)

	if err != nil {
		log.Printf("%s: %v", op, err)
		return models.Article{}, fmt.Errorf("%s: %w", op, err)
	}

	return article, nil
}
