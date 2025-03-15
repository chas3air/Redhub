package articlemanager

import (
	"articlesManageService/internal/domain/interfaces/storage"
	"articlesManageService/internal/domain/models"
	"articlesManageService/internal/services"
	storage_error "articlesManageService/internal/storage"
	"articlesManageService/pkg/lib/logger/sl"
	"errors"

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
		log.Error("Error retrieving articles:", sl.Err(err))
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
		if errors.Is(err, storage_error.ErrNotFound) {
			log.Error("Article not found", sl.Err(services.ErrNotFound))
			return models.Article{}, fmt.Errorf("%s: %w", op, services.ErrNotFound)
		}

		log.Error("Error retrieving article by id", sl.Err(err))
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

	articles, err := am.storage.GetArticlesByOwnerId(ctx, uid)
	if err != nil {
		log.Error("Error retrieving articles by owner id", sl.Err(err))
		return []models.Article{}, fmt.Errorf("%s: %w", op, err)
	}

	log.Info("Successfully retrieved articles by owner id")
	return articles, nil
}

// Insert implements articlesservice.ArticlesManager.
func (am *ArticleManager) Insert(ctx context.Context, article models.Article) (models.Article, error) {
	const op = "services.articleManager.insert"
	log := am.log.With(slog.String("operation", op))

	select {
	case <-ctx.Done():
		return models.Article{}, fmt.Errorf("%s: %w", op, ctx.Err())
	default:
	}

	article, err := am.storage.Insert(ctx, article)
	if err != nil {
		if errors.Is(err, storage_error.ErrAlreadyExists) {
			log.Error("Article already exists", sl.Err(err))
			return models.Article{}, fmt.Errorf("%s: %w", op, services.ErrAlreadyExists)
		}

		log.Error("Error inserting article", sl.Err(err))
		return models.Article{}, fmt.Errorf("%s: %w", op, err)
	}

	log.Info("Successfully inserted article")
	return article, nil
}

// Update implements articlesservice.ArticlesManager.
func (am *ArticleManager) Update(ctx context.Context, aid uuid.UUID, article models.Article) (models.Article, error) {
	const op = "services.articleManager.update"
	log := am.log.With(slog.String("operation", op))

	select {
	case <-ctx.Done():
		return models.Article{}, fmt.Errorf("%s: %w", op, ctx.Err())
	default:
	}

	article, err := am.storage.Update(ctx, aid, article)
	if err != nil {
		if errors.Is(err, storage_error.ErrNotFound) {
			log.Error("Article not found for update", sl.Err(err))
			return models.Article{}, fmt.Errorf("%s: %w", op, services.ErrNotFound)
		}

		log.Error("Error updating article", sl.Err(err))
		return models.Article{}, fmt.Errorf("%s: %w", op, err)
	}
	log.Info("Successfully updated article")
	return article, nil
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
		if errors.Is(err, storage_error.ErrNotFound) {
			log.Error("Article not found for deletion", sl.Err(err))
			return models.Article{}, fmt.Errorf("%s: %w", op, services.ErrNotFound)
		}

		log.Error("Error deleting article", sl.Err(err))
		return models.Article{}, fmt.Errorf("%s: %w", op, err)
	}

	log.Info("Successfully deleted article")
	return deletedArticle, nil
}
