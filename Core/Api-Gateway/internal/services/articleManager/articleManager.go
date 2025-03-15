package articlemanageservice

import (
	"apigateway/internal/domain/interfaces/articles"
	"apigateway/internal/domain/models"
	"apigateway/pkg/lib/logger/sl"
	"context"
	"fmt"
	"log/slog"

	"github.com/google/uuid"
)

type ArticleManageService struct {
	log     *slog.Logger
	storage articles.IArticlesStorage
}

func New(log *slog.Logger, storage articles.IArticlesStorage) *ArticleManageService {
	return &ArticleManageService{
		log:     log,
		storage: storage,
	}
}

// GetArticles implements articles.IArticlesService.
func (a *ArticleManageService) GetArticles(ctx context.Context) ([]models.Article, error) {
	const op = "services.articleManager.getArticles"
	log := a.log.With(
		slog.String("op", op),
	)

	select {
	case <-ctx.Done():
		return nil, fmt.Errorf("%s: %w", op, ctx.Err())
	default:
	}

	articles, err := a.storage.GetArticles(ctx)
	if err != nil {
		log.Error("error retrieving articles", sl.Err(err))
		return nil, fmt.Errorf("%s: %w", op, err)
	}

	return articles, nil
}

// GetArticleById implements articles.IArticlesService.
func (a *ArticleManageService) GetArticleById(ctx context.Context, aid uuid.UUID) (models.Article, error) {
	const op = "services.articleManager.getArticleById"
	log := a.log.With(
		slog.String("op", op),
	)

	select {
	case <-ctx.Done():
		return models.Article{}, fmt.Errorf("%s: %w", op, ctx.Err())
	default:
	}

	article, err := a.storage.GetArticleById(ctx, aid)
	if err != nil {
		log.Error("error retrieving article by id", sl.Err(err))
		return models.Article{}, fmt.Errorf("%s: %w", op, err)
	}

	return article, nil
}

// GetArticleByOwnerId implements articles.IArticlesService.
func (a *ArticleManageService) GetArticleByOwnerId(ctx context.Context, uid uuid.UUID) ([]models.Article, error) {
	const op = "services.articleManager.getArticleByOwnerId"
	log := a.log.With(
		slog.String("op", op),
	)

	select {
	case <-ctx.Done():
		return nil, fmt.Errorf("%s: %w", op, ctx.Err())
	default:
	}

	articles, err := a.storage.GetArticleByOwnerId(ctx, uid)
	if err != nil {
		log.Error("error retrieving articles by uid", sl.Err(err))
		return nil, fmt.Errorf("%s: %w", op, err)
	}

	return articles, nil
}

// Insert implements articles.IArticlesService.
func (a *ArticleManageService) Insert(ctx context.Context, article models.Article) (models.Article, error) {
	const op = "services.articleManager.insert"
	log := a.log.With(
		slog.String("op", op),
	)

	select {
	case <-ctx.Done():
		return models.Article{}, fmt.Errorf("%s: %w", op, ctx.Err())
	default:
	}

	article, err := a.storage.Insert(ctx, article)
	if err != nil {
		log.Error("failed to insert article", sl.Err(err))
		return models.Article{}, fmt.Errorf("%s: %w", op, err)
	}

	return article, nil
}

// Update implements articles.IArticlesService.
func (a ArticleManageService) Update(ctx context.Context, aid uuid.UUID, article models.Article) (models.Article, error) {
	const op = "services.articleManager.update"
	log := a.log.With(
		slog.String("op", op),
	)

	select {
	case <-ctx.Done():
		return models.Article{}, fmt.Errorf("%s: %w", op, ctx.Err())
	default:
	}

	article, err := a.storage.Update(ctx, aid, article)
	if err != nil {
		log.Error("failed to update article", sl.Err(err))
		return models.Article{}, fmt.Errorf("%s: %w", op, err)
	}

	return article, nil
}

// Delete implements articles.IArticlesService.
func (a ArticleManageService) Delete(ctx context.Context, aid uuid.UUID) (models.Article, error) {
	const op = "services.articleManager.delete"

	log := a.log.With(
		slog.String("op", op),
	)

	select {
	case <-ctx.Done():
		return models.Article{}, fmt.Errorf("%s: %w", op, ctx.Err())
	default:
	}

	article, err := a.storage.Delete(ctx, aid)
	if err != nil {
		log.Error("failed to update article", sl.Err(err))
		return models.Article{}, fmt.Errorf("%s: %w", op, err)
	}

	return article, nil
}
