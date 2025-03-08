package commentsservice

import (
	"apigateway/internal/domain/interfaces/comments"
	"apigateway/internal/domain/models"
	"apigateway/pkg/lib/logger/sl"
	"context"
	"fmt"
	"log/slog"

	"github.com/google/uuid"
)

type CommentsService struct {
	log     *slog.Logger
	storage comments.CommentsStorage
}

func New(log *slog.Logger, storage comments.CommentsStorage) *CommentsService {
	return &CommentsService{
		log:     log,
		storage: storage,
	}
}

func (cs *CommentsService) GetCommentById(ctx context.Context, cid uuid.UUID) (models.Comment, error) {
	const op = "services.comments.getCommentById"
	log := cs.log.With(
		slog.String("op", op),
	)

	select {
	case <-ctx.Done():
		return models.Comment{}, fmt.Errorf("%s: %w", op, ctx.Err())
	default:
	}

	comment, err := cs.storage.GetCommentById(ctx, cid)
	if err != nil {
		log.Error("error getting comment by id", sl.Err(err))
		return models.Comment{}, fmt.Errorf("%s: %w", op, err)
	}

	return comment, nil
}

func (cs *CommentsService) GetCommentsByArticleId(ctx context.Context, aid uuid.UUID) ([]models.Comment, error) {
	const op = "services.comments.getCommentsByArticleId"
	log := cs.log.With(
		slog.String("op", op),
	)

	select {
	case <-ctx.Done():
		return nil, fmt.Errorf("%s: %w", op, ctx.Err())
	default:
	}

	comments, err := cs.storage.GetCommentsByArticleId(ctx, aid)
	if err != nil {
		log.Error("error getting comments by article_id", sl.Err(err))
		return nil, fmt.Errorf("%s: %w", op, err)
	}

	return comments, nil
}

func (cs *CommentsService) Insert(ctx context.Context, comment models.Comment) (models.Comment, error) {
	const op = "services.comments.insert"
	log := cs.log.With(
		slog.String("op", op),
	)

	select {
	case <-ctx.Done():
		return models.Comment{}, fmt.Errorf("%s: %w", op, ctx.Err())
	default:
	}

	comment, err := cs.storage.Insert(ctx, comment)
	if err != nil {
		log.Error("error inserting comment", sl.Err(err))
		return models.Comment{}, fmt.Errorf("%s: %w", op, err)
	}

	return comment, nil
}

func (cs *CommentsService) Delete(ctx context.Context, cid uuid.UUID) (models.Comment, error) {
	const op = "services.comments.delete"
	log := cs.log.With(
		slog.String("op", op),
	)

	select {
	case <-ctx.Done():
		return models.Comment{}, fmt.Errorf("%s: %w", op, ctx.Err())
	default:
	}

	comment, err := cs.storage.Delete(ctx, cid)
	if err != nil {
		log.Error("Error deleting comment", sl.Err(err))
		return models.Comment{}, fmt.Errorf("%s: %w", op, err)
	}

	return comment, nil
}
