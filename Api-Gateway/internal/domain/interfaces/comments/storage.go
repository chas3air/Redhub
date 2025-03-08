package comments

import (
	"apigateway/internal/domain/models"
	"context"

	"github.com/google/uuid"
)

type CommentsStorage interface {
	GetCommentById(context.Context, uuid.UUID) (models.Comment, error)
	GetCommentsByArticleId(context.Context, uuid.UUID) ([]models.Comment, error)
	Insert(context.Context, models.Comment) (models.Comment, error)
	Delete(context.Context, uuid.UUID) (models.Comment, error)
}
