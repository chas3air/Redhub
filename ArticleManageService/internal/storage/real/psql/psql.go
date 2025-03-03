package psql

import (
	"articlesManageService/internal/domain/models"
	"articlesManageService/internal/storage"
	"context"
	"database/sql"
	"errors"
	"fmt"

	"github.com/google/uuid"
)

const (
	ArticlesTableName         = "Articles"
	CommentsTableName         = "Comments"
	ArticleToCommentTableName = "article_to_comment"
)

type PsqlStorage struct {
	DB *sql.DB
}

func New(connStr string) *PsqlStorage {
	db, err := sql.Open("postgres", connStr)
	if err != nil {
		panic(err)
	}

	return &PsqlStorage{
		DB: db,
	}
}

func (s *PsqlStorage) Close() {
	err := s.DB.Close()
	if err != nil {
		fmt.Println("Error closing the database:", err)
	}
}

func (s *PsqlStorage) GetArticles(ctx context.Context) ([]models.Article, error) {
	const op = "psql.getArticles"

	rows, err := s.DB.QueryContext(ctx, `
		SELECT * FROM `+ArticlesTableName+`;
	`)
	if err != nil {
		return nil, fmt.Errorf("%s: %w", op, err)
	}
	defer rows.Close()

	articlesMap := make(map[uuid.UUID]*models.Article)
	for rows.Next() {
		var article models.Article
		err := rows.Scan(&article.Id, &article.CreatedAt, &article.Title, &article.Content, &article.OwnerId)
		if err != nil {
			return nil, fmt.Errorf("%s: %w", op, err)
		}
		article.Comments = []models.Comment{}
		articlesMap[article.Id] = &article
	}

	rowsForComms, err := s.DB.QueryContext(ctx, `
		SELECT article_id, comment_id FROM `+ArticleToCommentTableName+`;
	`)
	if err != nil {
		return nil, fmt.Errorf("%s: %w", op, err)
	}
	defer rowsForComms.Close()

	for rowsForComms.Next() {
		var articleID, commentID uuid.UUID
		err := rowsForComms.Scan(&articleID, &commentID)
		if err != nil {
			return nil, fmt.Errorf("%s: %w", op, err)
		}
		if article, exists := articlesMap[articleID]; exists {
			article.Comments = append(article.Comments, models.Comment{Id: commentID})
		}
	}

	articles := make([]models.Article, 0, len(articlesMap))
	for _, article := range articlesMap {
		articles = append(articles, *article)
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
			return models.Article{}, fmt.Errorf("%s: %w", op, storage.ErrNotFound)
		}
		return models.Article{}, fmt.Errorf("%s: %w", op, err)
	}

	rowsForComms, err := s.DB.QueryContext(ctx, `
		SELECT comment_id FROM `+ArticleToCommentTableName+`
		WHERE article_id=$1;
	`, aid)
	if err != nil {
		return models.Article{}, fmt.Errorf("%s: %w", op, err)
	}
	defer rowsForComms.Close()

	var cid uuid.UUID
	comments := make([]models.Comment, 0)
	for rowsForComms.Next() {
		err := rowsForComms.Scan(&cid)
		if err != nil {
			return models.Article{}, fmt.Errorf("%s: %w", op, err)
		}
		comments = append(comments, models.Comment{Id: cid})
	}
	article.Comments = comments

	return article, nil
}

func (s *PsqlStorage) GetArticlesByOwnerId(ctx context.Context, uid uuid.UUID) ([]models.Article, error) {
	const op = "psql.getArticlesByOwnerId"

	rows, err := s.DB.QueryContext(ctx, `
		SELECT * FROM `+ArticlesTableName+`
		WHERE owner_id=$1;
	`, uid)
	if err != nil {
		return nil, fmt.Errorf("%s: %w", op, err)
	}
	defer rows.Close()

	articles := make([]models.Article, 0, 5)
	for rows.Next() {
		var article models.Article
		err := rows.Scan(&article.Id, &article.CreatedAt, &article.Title, &article.Content, &article.OwnerId)
		if err != nil {
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
		return fmt.Errorf("%s: %w", op, storage.ErrAlreadyExists)
	} else if err != nil && !errors.Is(err, storage.ErrNotFound) {
		return fmt.Errorf("%s: %w", op, err)
	}

	_, err = s.DB.ExecContext(ctx, `
		INSERT INTO `+ArticlesTableName+` (id, created_at, title, content, owner_id)
		VALUES ($1, $2, $3, $4, $5);
	`, article.Id, article.CreatedAt, article.Title, article.Content, article.OwnerId)
	if err != nil {
		return fmt.Errorf("%s: %w", op, err)
	}

	return nil
}

func (s *PsqlStorage) Update(ctx context.Context, aid uuid.UUID, article models.Article) error {
	const op = "psql.update"

	_, err := s.GetArticleById(ctx, aid)
	if err != nil {
		if errors.Is(err, storage.ErrNotFound) {
			return fmt.Errorf("%s: %w", op, storage.ErrNotFound)
		}
		return fmt.Errorf("%s: %w", op, err)
	}

	_, err = s.DB.ExecContext(ctx, `
		UPDATE `+ArticlesTableName+` SET 
			title = $1, 
			content = $2  
		WHERE id = $3;
	`, article.Title, article.Content, aid)
	if err != nil {
		return fmt.Errorf("%s: %w", op, err)
	}

	return nil
}

func (s *PsqlStorage) Delete(ctx context.Context, aid uuid.UUID) (models.Article, error) {
	const op = "psql.delete"

	article, err := s.GetArticleById(ctx, aid)
	if err != nil {
		return models.Article{}, fmt.Errorf("%s: %w", op, err)
	}

	_, err = s.DB.ExecContext(ctx, `
		DELETE FROM `+ArticlesTableName+`
		WHERE id=$1;
	`, aid)

	if err != nil {
		return models.Article{}, fmt.Errorf("%s: %w", op, err)
	}

	return article, nil
}
