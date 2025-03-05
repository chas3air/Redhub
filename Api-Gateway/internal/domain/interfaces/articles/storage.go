package articles

import (
	"apigateway/internal/domain/models"
	"context"

	"github.com/google/uuid"
)

type IArticlesStorage interface {
	GetArticles(ctx context.Context) ([]models.Article, error)
	GetArticleById(ctx context.Context, aid uuid.UUID) (models.Article, error)
	GetArticleByOwnerId(ctx context.Context, uid uuid.UUID) ([]models.Article, error)
	Insert(ctx context.Context, article models.Article) error
	Update(ctx context.Context, aid uuid.UUID, articler models.Article) error
	Delete(ctx context.Context, aid uuid.UUID) (models.Article, error)
}
