package storage

import (
	"commentsManageService/internal/domain/models"
	"context"

	"github.com/google/uuid"
)

type CommentStorage interface {
	GetCommentById(context.Context, uuid.UUID) (models.Comment, error)
	GetCommentsByArticleId(context.Context, uuid.UUID) ([]models.Comment, error)
	Insert(context.Context, models.Comment) (models.Comment, error)
	Delete(context.Context, uuid.UUID) (models.Comment, error)
}
