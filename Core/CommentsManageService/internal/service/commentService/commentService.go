package commentservice

import (
	"commentsManageService/internal/domain/interfaces/storage"
	"commentsManageService/internal/domain/models"
	"commentsManageService/pkg/lib/logger/sl"
	"context"
	"errors"
	"fmt"
	"log/slog"

	service_error "commentsManageService/internal/service"
	storage_error "commentsManageService/internal/storage"

	"github.com/google/uuid"
)

type CommentService struct {
	log     *slog.Logger
	storage storage.CommentStorage
}

func New(log *slog.Logger, storage storage.CommentStorage) *CommentService {
	return &CommentService{
		log:     log,
		storage: storage,
	}
}

func (c *CommentService) GetCommentById(ctx context.Context, cid uuid.UUID) (models.Comment, error) {
	const op = "service.commentService.getCommentById"
	log := c.log.With(
		slog.String("op", op),
	)

	select {
	case <-ctx.Done():
		return models.Comment{}, fmt.Errorf("%s: %w", op, ctx.Err())
	default:
	}

	comment, err := c.storage.GetCommentById(ctx, cid)
	if err != nil {
		if errors.Is(err, storage_error.ErrNotFound) {
			log.Warn("Comment not found", sl.Err(err))
			return models.Comment{}, fmt.Errorf("%s: %w", op, service_error.ErrNotFound)
		}

		log.Error("Failed to retrieve comment by id", sl.Err(err))
		return models.Comment{}, fmt.Errorf("%s: %w", op, err)
	}

	log.Info("Successfully retrieve comment")
	return comment, nil
}

func (c *CommentService) GetCommentsByArticleId(ctx context.Context, aid uuid.UUID) ([]models.Comment, error) {
	const op = "service.commentService.getCommentsByArticleId"
	log := c.log.With(
		slog.String("op", op),
	)

	select {
	case <-ctx.Done():
		return nil, fmt.Errorf("%s: %w", op, ctx.Err())
	default:
	}

	comments, err := c.storage.GetCommentsByArticleId(ctx, aid)
	if err != nil {
		log.Error("Failed to retrieve comment by article_id", sl.Err(err))
		return nil, fmt.Errorf("%s: %w", op, err)
	}

	log.Info("Successfully retrieve comments by article_id")
	return comments, nil
}

func (c *CommentService) Insert(ctx context.Context, comment models.Comment) (models.Comment, error) {
	const op = "service.commentService."
	log := c.log.With(
		slog.String("op", op),
	)

	select {
	case <-ctx.Done():
		return models.Comment{}, fmt.Errorf("%s: %w", op, ctx.Err())
	default:
	}

	comment, err := c.storage.Insert(ctx, comment)
	if err != nil {
		if errors.Is(err, storage_error.ErrAlreadyExists) {
			log.Warn("Comment already exists", sl.Err(err))
			return models.Comment{}, fmt.Errorf("%s: %w", op, service_error.ErrAlreadyExists)
		}

		log.Error("Failed to insert comment", sl.Err(err))
		return models.Comment{}, fmt.Errorf("%s: %w", op, err)
	}

	log.Info("Comment inserted successfully")
	return comment, nil
}

func (c *CommentService) Delete(ctx context.Context, cid uuid.UUID) (models.Comment, error) {
	const op = "service.commentService."
	log := c.log.With(
		slog.String("op", op),
	)

	select {
	case <-ctx.Done():
		return models.Comment{}, fmt.Errorf("%s: %w", op, ctx.Err())
	default:
	}

	comment, err := c.storage.Delete(ctx, cid)
	if err != nil {
		if errors.Is(err, storage_error.ErrNotFound) {
			log.Warn("Comment not found", sl.Err(err))
			return models.Comment{}, fmt.Errorf("%s: %w", op, service_error.ErrNotFound)
		}

		log.Error("Failed to delete comment by id", sl.Err(err))
		return models.Comment{}, fmt.Errorf("%s: %w", op, err)
	}

	log.Info("Comment deleted successfully")
	return comment, nil
}
