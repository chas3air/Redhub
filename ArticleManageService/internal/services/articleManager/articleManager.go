package articlemanager

import (
	"articlesManageService/internal/domain/interfaces/storage"
	"articlesManageService/internal/domain/models"
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

// GetAtricles implements articlesservice.ArticlesManager.
func (am *ArticleManager) GetAtricles(ctx context.Context) ([]models.Article, error) {
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
		return fmt.Errorf("%s: %w", op, err)
	}

	log.Info("Successfully insert article")
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
		return fmt.Errorf("%s: %w", op, err)
	}
	log.Info("Successfully update article")
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

	deleted_article, err := am.storage.Delete(ctx, aid)
	if err != nil {
		return models.Article{}, fmt.Errorf("%s: %w", op, err)
	}

	log.Info("Successfully delete article")
	return deleted_article, nil
}
