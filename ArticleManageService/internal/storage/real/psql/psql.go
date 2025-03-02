package psql

import (
	"articlesManageService/internal/domain/models"
	"context"
	"database/sql"
	"fmt"

	"github.com/google/uuid"
)

const (
	ArticlesTableName         = "Articles"
	CommentsTableName         = "Comments"
	ArticleToCommentTableName = "article_to_comment"
	ArticleToTagTableName     = "article_to_tag"
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
	s.DB.Close()
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
			continue
		}
		article.Comments = []models.Comment{}
		article.Tags = []string{}
		articlesMap[article.Id] = &article
	}

	// Получаем комментарии
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
			continue
		}
		if article, exists := articlesMap[articleID]; exists {
			article.Comments = append(article.Comments, models.Comment{Id: commentID})
		}
	}

	// Получаем теги
	rowsForTags, err := s.DB.QueryContext(ctx, `
		SELECT article_id, tag FROM `+ArticleToTagTableName+`;
	`)
	if err != nil {
		return nil, fmt.Errorf("%s: %w", op, err)
	}
	defer rowsForTags.Close()

	for rowsForTags.Next() {
		var articleID uuid.UUID
		var tag string
		err := rowsForTags.Scan(&articleID, &tag)
		if err != nil {
			continue
		}
		if article, exists := articlesMap[articleID]; exists {
			article.Tags = append(article.Tags, tag)
		}
	}

	// Конвертируем карту в срез
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
		return models.Article{}, fmt.Errorf("%s: %w", op, err)
	}

	// Получаем комментарии
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
			continue
		}
		comments = append(comments, models.Comment{Id: cid})
	}
	article.Comments = comments

	// Получаем теги
	rowsForTags, err := s.DB.QueryContext(ctx, `
		SELECT tag FROM `+ArticleToTagTableName+`
		WHERE article_id=$1;
	`, aid)
	if err != nil {
		return models.Article{}, fmt.Errorf("%s: %w", op, err)
	}
	defer rowsForTags.Close()

	var tag string
	tags := make([]string, 0)
	for rowsForTags.Next() {
		err := rowsForTags.Scan(&tag)
		if err != nil {
			continue
		}
		tags = append(tags, tag)
	}
	article.Tags = tags

	return article, nil
}

func (s *PsqlStorage) GetArticleByOwnerId(ctx context.Context, uid uuid.UUID) ([]models.Article, error) {
	const op = "psql.getArticleByOwnerId"

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
			continue
		}
		articles = append(articles, article)
	}

	return articles, nil
}

func (s *PsqlStorage) Insert(ctx context.Context, article models.Article) error {
	const op = "psql.insert"

	_, err := s.DB.ExecContext(ctx, `
		INSERT INTO `+ArticlesTableName+` (id, created_at, title, content, owner_id)
		VALUES ($1, $2, $3, $4, $5);
	`, article.Id, article.CreatedAt, article.Title, article.Content, article.OwnerId)
	if err != nil {
		return fmt.Errorf("%s: %w", op, err)
	}

	// Добавляем комментарии, если есть
	for _, comment := range article.Comments {
		_, err := s.DB.ExecContext(ctx, `
			INSERT INTO `+ArticleToCommentTableName+` (article_id, comment_id)
			VALUES ($1, $2);
		`, article.Id, comment.Id)
		if err != nil {
			return fmt.Errorf("%s: %w", op, err)
		}
	}

	// Добавляем теги, если есть
	for _, tag := range article.Tags {
		_, err := s.DB.ExecContext(ctx, `
			INSERT INTO `+ArticleToTagTableName+` (article_id, tag)
			VALUES ($1, $2);
		`, article.Id, tag)
		if err != nil {
			return fmt.Errorf("%s: %w", op, err)
		}
	}

	return nil
}

func (s *PsqlStorage) Update(ctx context.Context, aid uuid.UUID, article models.Article) error {
	const op = "psql.update"

	_, err := s.DB.ExecContext(ctx, `
		UPDATE `+ArticlesTableName+` SET 
			title = $1, 
			content = $2, 
			owner_id = $3 
		WHERE id = $4;
	`, article.Title, article.Content, article.OwnerId, aid)
	if err != nil {
		return fmt.Errorf("%s: %w", op, err)
	}

	// Обновляем комментарии
	_, err = s.DB.ExecContext(ctx, `
		DELETE FROM `+ArticleToCommentTableName+` 
		WHERE article_id = $1;
	`, aid)
	if err != nil {
		return fmt.Errorf("%s: %w", op, err)
	}

	for _, comment := range article.Comments {
		_, err := s.DB.ExecContext(ctx, `
			INSERT INTO `+ArticleToCommentTableName+` (article_id, comment_id)
			VALUES ($1, $2);
		`, aid, comment.Id)
		if err != nil {
			return fmt.Errorf("%s: %w", op, err)
		}
	}

	// Обновляем теги
	_, err = s.DB.ExecContext(ctx, `
		DELETE FROM `+ArticleToTagTableName+` 
		WHERE article_id = $1;
	`, aid)
	if err != nil {
		return fmt.Errorf("%s: %w", op, err)
	}

	for _, tag := range article.Tags {
		_, err := s.DB.ExecContext(ctx, `
			INSERT INTO `+ArticleToTagTableName+` (article_id, tag)
			VALUES ($1, $2);
		`, aid, tag)
		if err != nil {
			return fmt.Errorf("%s: %w", op, err)
		}
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
