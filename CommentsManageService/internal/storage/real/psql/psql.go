package psqlstorage

import (
	"commentsManageService/internal/domain/models"
	storage_error "commentsManageService/internal/storage"
	"commentsManageService/pkg/lib/logger/sl"
	"context"
	"database/sql"
	"errors"
	"fmt"
	"log/slog"

	"github.com/google/uuid"
	"github.com/lib/pq"
)

type PsqlStorage struct {
	log *slog.Logger
	DB  *sql.DB
}

const CommentTableName = "Comments"

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

func (p *PsqlStorage) GetCommentById(ctx context.Context, cid uuid.UUID) (models.Comment, error) {
	const op = "storage.psql.getCommentById"
	log := p.log.With(
		slog.String("op", op),
	)

	select {
	case <-ctx.Done():
		return models.Comment{}, fmt.Errorf("%s: %w", op, ctx.Err())
	default:
	}

	var comment models.Comment

	row := p.DB.QueryRowContext(ctx, `
		SELECT * FROM `+CommentTableName+`
		WHERE id=$1
	`, cid)
	err := row.Scan(&comment.Id, &comment.ArticleId, &comment.OwnerId, &comment.CreatedAt, &comment.Content)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			log.Error("Comment with current id not found", sl.Err(err))
			return models.Comment{}, fmt.Errorf("%s: %w", op, storage_error.ErrNotFound)
		}

		log.Error("Error scanning row", sl.Err(err))
		return models.Comment{}, fmt.Errorf("%s: %w", op, err)
	}

	return comment, nil
}

func (p *PsqlStorage) GetCommentsByArticleId(ctx context.Context, aid uuid.UUID) ([]models.Comment, error) {
	const op = "storage.psql.getCommentsByArticleId"
	log := p.log.With(
		slog.String("op", op),
	)

	select {
	case <-ctx.Done():
		return nil, fmt.Errorf("%s: %w", op, ctx.Err())
	default:
	}

	rows, err := p.DB.QueryContext(ctx, `
		SELECT * FROM `+CommentTableName+`
		WHERE article_id=$1
	`, aid)
	if err != nil {
		log.Error("Error retrieving all comments", sl.Err(err))
		return nil, fmt.Errorf("%s: %w", op, err)
	}
	defer rows.Close()

	var comments = make([]models.Comment, 0, 10)
	for rows.Next() {
		var comment models.Comment
		if err := rows.Scan(&comment.Id, &comment.ArticleId, &comment.OwnerId, &comment.CreatedAt, &comment.Content); err != nil {
			log.Error("Error scanning row", sl.Err(err))
			continue
		}

		comments = append(comments, comment)
	}

	return comments, nil
}

func (p *PsqlStorage) Insert(ctx context.Context, comment models.Comment) (models.Comment, error) {
	const op = "storage.psql.insert"
	log := p.log.With(
		slog.String("op", op),
	)

	select {
	case <-ctx.Done():
		return models.Comment{}, fmt.Errorf("%s: %w", op, ctx.Err())
	default:
	}

	_, err := p.DB.ExecContext(ctx, `
		INSERT INTO `+CommentTableName+` (id, article_id, owner_id, created_at, content)
		VALUES ($1, $2, $3, $4, $5);
	`, comment.Id, comment.ArticleId, comment.OwnerId, comment.CreatedAt, comment.Content)
	if err != nil {
		if pgErr, ok := err.(*pq.Error); ok && pgErr.Code == "23505" {
			log.Error("Comment with this ID already exists", sl.Err(err))
			return comment, fmt.Errorf("%s: %w", op, storage_error.ErrAlreadyExists)
		}

		log.Error("Error inserting comment", sl.Err(err))
		return comment, fmt.Errorf("%s: %w", op, err)
	}

	return comment, nil
}
func (p *PsqlStorage) Delete(ctx context.Context, cid uuid.UUID) (models.Comment, error) {
	const op = "storage.psql.delete"
	log := p.log.With(
		slog.String("op", op),
	)

	select {
	case <-ctx.Done():
		return models.Comment{}, fmt.Errorf("%s: %w", op, ctx.Err())
	default:
	}

	comment, err := p.GetCommentById(ctx, cid)
	if err != nil {
		if errors.Is(err, storage_error.ErrNotFound) {
			log.Warn("Comment not found", sl.Err(err))
			return models.Comment{}, fmt.Errorf("%s: %w", op, storage_error.ErrNotFound)
		}

		log.Error("Error getting comment defore deliting", sl.Err(err))
		return models.Comment{}, fmt.Errorf("%s: %w", op, err)
	}

	_, err = p.DB.ExecContext(ctx, `
		DELETE FROM `+CommentTableName+`
		WHERE id=$1
	`, cid)
	if err != nil {
		log.Error("Error deleting comment", sl.Err(err))
		return models.Comment{}, fmt.Errorf("%s: %w", op, err)
	}

	return comment, nil
}
