package storage

import (
	"articlesManageService/internal/domain/models"
	"context"

	"github.com/google/uuid"
)

type Storage interface {
	GetArticles(ctx context.Context) ([]models.Article, error)
	GetArticleById(ctx context.Context, aid uuid.UUID) (models.Article, error)
	GetArticlesByOwnerId(ctx context.Context, uid uuid.UUID) ([]models.Article, error)
	Insert(ctx context.Context, article models.Article) (models.Article, error)
	Update(ctx context.Context, aid uuid.UUID, articler models.Article) (models.Article, error)
	Delete(ctx context.Context, aid uuid.UUID) (models.Article, error)
}
