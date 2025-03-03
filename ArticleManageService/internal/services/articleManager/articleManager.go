package articlemanager

import (
	"articlesManageService/internal/domain/interfaces/storage"
	"articlesManageService/internal/domain/models"
	"articlesManageService/internal/services"

	"context"
	"fmt"
	"log/slog"

	"github.com/google/uuid"
)

type ArticleManager struct {
	log     *slog.Logger
	storage storage.Storage
}

func New(log *slog.Logger, storage storage.Storage) *ArticleManager {
	return &ArticleManager{
		log:     log,
		storage: storage,
	}
}

// GetArticles implements articlesservice.ArticlesManager.
func (am *ArticleManager) GetArticles(ctx context.Context) ([]models.Article, error) {
	const op = "services.articleManager.getArticles"
	log := am.log.With(slog.String("operation", op))

	select {
	case <-ctx.Done():
		return nil, fmt.Errorf("%s: %w", op, ctx.Err())
	default:
	}

	articles, err := am.storage.GetArticles(ctx)
	if err != nil {
		return nil, fmt.Errorf("%s: %w", op, err)
	}

	log.Info("Successfully retrieved articles")
	return articles, nil
}

// GetArticleById implements articlesservice.ArticlesManager.
func (am *ArticleManager) GetArticleById(ctx context.Context, aid uuid.UUID) (models.Article, error) {
	const op = "services.articleManager.getArticleById"
	log := am.log.With(slog.String("operation", op))

	select {
	case <-ctx.Done():
		return models.Article{}, fmt.Errorf("%s: %w", op, ctx.Err())
	default:
	}

	article, err := am.storage.GetArticleById(ctx, aid)
	if err != nil {
		if err == services.ErrNotFound {
			return models.Article{}, fmt.Errorf("%s: %w", op, err)
		}
		return models.Article{}, fmt.Errorf("%s: %w", op, err)
	}

	log.Info("Successfully retrieved article by id")
	return article, nil
}

// GetArticleByOwnerId implements articlesservice.ArticlesManager.
func (am *ArticleManager) GetArticleByOwnerId(ctx context.Context, uid uuid.UUID) ([]models.Article, error) {
	const op = "services.articleManager.getArticleByOwnerId"
	log := am.log.With(slog.String("operation", op))

	select {
	case <-ctx.Done():
		return nil, fmt.Errorf("%s: %w", op, ctx.Err())
	default:
	}

	articles, err := am.storage.GetArticleByOwnerId(ctx, uid)
	if err != nil {
		return []models.Article{}, fmt.Errorf("%s: %w", op, err)
	}

	log.Info("Successfully retrieved articles by owner id")
	return articles, nil
}

// Insert implements articlesservice.ArticlesManager.
func (am *ArticleManager) Insert(ctx context.Context, article models.Article) error {
	const op = "services.articleManager.insert"
	log := am.log.With(slog.String("operation", op))

	select {
	case <-ctx.Done():
		return fmt.Errorf("%s: %w", op, ctx.Err())
	default:
	}

	err := am.storage.Insert(ctx, article)
	if err != nil {
		if err == services.ErrAlreadyExists {
			return fmt.Errorf("%s: %w", op, err)
		}
		return fmt.Errorf("%s: %w", op, err)
	}

	log.Info("Successfully inserted article")
	return nil
}

// Update implements articlesservice.ArticlesManager.
func (am *ArticleManager) Update(ctx context.Context, aid uuid.UUID, article models.Article) error {
	const op = "services.articleManager.update"
	log := am.log.With(slog.String("operation", op))

	select {
	case <-ctx.Done():
		return fmt.Errorf("%s: %w", op, ctx.Err())
	default:
	}

	err := am.storage.Update(ctx, aid, article)
	if err != nil {
		if err == services.ErrNotFound {
			return fmt.Errorf("%s: %w", op, err)
		}
		return fmt.Errorf("%s: %w", op, err)
	}
	log.Info("Successfully updated article")
	return nil
}

// Delete implements articlesservice.ArticlesManager.
func (am *ArticleManager) Delete(ctx context.Context, aid uuid.UUID) (models.Article, error) {
	const op = "services.articleManager.delete"
	log := am.log.With(slog.String("operation", op))

	select {
	case <-ctx.Done():
		return models.Article{}, fmt.Errorf("%s: %w", op, ctx.Err())
	default:
	}

	deletedArticle, err := am.storage.Delete(ctx, aid)
	if err != nil {
		if err == services.ErrNotFound {
			return models.Article{}, fmt.Errorf("%s: %w", op, err)
		}
		return models.Article{}, fmt.Errorf("%s: %w", op, err)
	}

	log.Info("Successfully deleted article")
	return deletedArticle, nil
}
