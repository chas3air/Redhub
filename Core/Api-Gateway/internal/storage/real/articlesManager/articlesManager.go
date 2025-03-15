package articlesmanagerstorage

import (
	"apigateway/internal/domain/models"
	amprofiles "apigateway/internal/domain/profiles/am_profiles"
	"apigateway/pkg/lib/logger/sl"
	"context"
	"fmt"
	"log/slog"

	amv1 "github.com/chas3air/protos/gen/go/articlesManager"
	"github.com/google/uuid"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
)

type ArticlesManageStorage struct {
	log         *slog.Logger
	ServiceHost string
	ServicePort int
}

func New(log *slog.Logger, serviceHost string, servicePort int) *ArticlesManageStorage {
	return &ArticlesManageStorage{
		log:         log,
		ServiceHost: serviceHost,
		ServicePort: servicePort,
	}
}

// GetArticles implements articles.IArticlesStorage.
func (a *ArticlesManageStorage) GetArticles(ctx context.Context) ([]models.Article, error) {
	const op = "articlesmanagestorage.getArticles"
	log := a.log.With(slog.String("op", op))

	select {
	case <-ctx.Done():
		return nil, fmt.Errorf("%s: %w", op, ctx.Err())
	default:
	}

	conn, err := grpc.NewClient(
		fmt.Sprintf("%s:%d", a.ServiceHost, a.ServicePort),
		grpc.WithTransportCredentials(insecure.NewCredentials()),
	)
	if err != nil {
		log.Error("Failed to connect to gRPC server", sl.Err(err))
		return nil, fmt.Errorf("%s: %w", op, err)
	}
	defer conn.Close()

	c := amv1.NewArticlesManagerClient(conn)
	res, err := c.GetArticles(ctx, nil)
	if err != nil {
		log.Warn("Failed to get articles", sl.Err(err))
		return nil, fmt.Errorf("%s: %w", op, err)
	}

	resp_articles := make([]models.Article, 0, len(res.GetArticles()))
	for _, pbArticle := range res.GetArticles() {
		article, err := amprofiles.ProtoArtToArt(pbArticle)
		if err != nil {
			log.Error("Wrong structure, failed to customize", sl.Err(err))
			continue
		}
		resp_articles = append(resp_articles, article)
	}

	return resp_articles, nil
}

// GetArticleById implements articles.IArticlesStorage.
func (a *ArticlesManageStorage) GetArticleById(ctx context.Context, aid uuid.UUID) (models.Article, error) {
	const op = "articlesmanagestorage.getArticleById"
	log := a.log.With(slog.String("op", op))

	select {
	case <-ctx.Done():
		return models.Article{}, fmt.Errorf("%s: %w", op, ctx.Err())
	default:
	}

	conn, err := grpc.NewClient(
		fmt.Sprintf("%s:%d", a.ServiceHost, a.ServicePort),
		grpc.WithTransportCredentials(insecure.NewCredentials()),
	)
	if err != nil {
		log.Error("Failed to connect to gRPC server", sl.Err(err))
		return models.Article{}, fmt.Errorf("%s: %w", op, err)
	}
	defer conn.Close()

	c := amv1.NewArticlesManagerClient(conn)
	res, err := c.GetArticleById(ctx, &amv1.GetArticleByIdRequest{
		ArticleId: aid.String(),
	})
	if err != nil {
		log.Warn("Failed to get articles", sl.Err(err))
		return models.Article{}, fmt.Errorf("%s: %w", op, err)
	}

	resp_article, err := amprofiles.ProtoArtToArt(res.GetArticle())
	if err != nil {
		log.Warn("failed to get article by id", sl.Err(err))
		return models.Article{}, fmt.Errorf("%s: %w", op, err)
	}

	return resp_article, nil
}

// GetArticleByOwnerId implements articles.IArticlesStorage.
func (a *ArticlesManageStorage) GetArticleByOwnerId(ctx context.Context, uid uuid.UUID) ([]models.Article, error) {
	const op = "articlesmanagestorage.getArticlesByOwnerID"
	log := a.log.With(slog.String("op", op))

	select {
	case <-ctx.Done():
		return nil, fmt.Errorf("%s: %w", op, ctx.Err())
	default:
	}

	conn, err := grpc.NewClient(
		fmt.Sprintf("%s:%d", a.ServiceHost, a.ServicePort),
		grpc.WithTransportCredentials(insecure.NewCredentials()),
	)
	if err != nil {
		log.Error("Failed to connect to gRPC server", sl.Err(err))
		return nil, fmt.Errorf("%s: %w", op, err)
	}
	defer conn.Close()

	c := amv1.NewArticlesManagerClient(conn)
	res, err := c.GetArticlesByOwnerId(ctx, &amv1.GetArticlesByOwnerIdRequest{
		OwnerId: uid.String(),
	})
	if err != nil {
		log.Warn("Failed to get articles", sl.Err(err))
		return nil, fmt.Errorf("%s: %w", op, err)
	}

	resp_articles := make([]models.Article, 0, len(res.GetArticles()))
	for _, pbArticle := range res.GetArticles() {
		article, err := amprofiles.ProtoArtToArt(pbArticle)
		if err != nil {
			log.Error("Wrong structure, failed to customize", sl.Err(err))
			continue
		}
		resp_articles = append(resp_articles, article)
	}

	return resp_articles, nil
}

// Insert implements articles.IArticlesStorage.
func (a *ArticlesManageStorage) Insert(ctx context.Context, article models.Article) (models.Article, error) {
	const op = "articlesmanagestorage.insert"
	log := a.log.With(slog.String("op", op))

	select {
	case <-ctx.Done():
		return models.Article{}, fmt.Errorf("%s: %w", op, ctx.Err())
	default:
	}

	articleForInsert, err := amprofiles.ArtToProtoArt(article)
	if err != nil {
		log.Error("Failed format of article", sl.Err(err))
		return models.Article{}, fmt.Errorf("%s: %w", op, err)
	}

	conn, err := grpc.NewClient(
		fmt.Sprintf("%s:%d", a.ServiceHost, a.ServicePort),
		grpc.WithTransportCredentials(insecure.NewCredentials()),
	)
	if err != nil {
		log.Error("Failed to connect to gRPC server", sl.Err(err))
		return models.Article{}, fmt.Errorf("%s: %w", op, err)
	}
	defer conn.Close()

	c := amv1.NewArticlesManagerClient(conn)
	_, err = c.InsertArticle(ctx, &amv1.InsertArticleRequest{
		Article: articleForInsert,
	})
	if err != nil {
		log.Error("Failed to insert article:", sl.Err(err))
		return models.Article{}, fmt.Errorf("%s: %w", op, err)
	}

	return article, nil
}

// Update implements articles.IArticlesStorage.
func (a *ArticlesManageStorage) Update(ctx context.Context, aid uuid.UUID, article models.Article) (models.Article, error) {
	const op = "articlesmanagestorage.update"
	log := a.log.With(slog.String("op", op))

	select {
	case <-ctx.Done():
		return models.Article{}, fmt.Errorf("%s: %w", op, ctx.Err())
	default:
	}

	articleForInsert, err := amprofiles.ArtToProtoArt(article)
	if err != nil {
		log.Error("Failed format of article", sl.Err(err))
		return models.Article{}, fmt.Errorf("%s: %w", op, err)
	}

	conn, err := grpc.NewClient(
		fmt.Sprintf("%s:%d", a.ServiceHost, a.ServicePort),
		grpc.WithTransportCredentials(insecure.NewCredentials()),
	)
	if err != nil {
		log.Error("Failed to connect to gRPC server", sl.Err(err))
		return models.Article{}, fmt.Errorf("%s: %w", op, err)
	}
	defer conn.Close()

	c := amv1.NewArticlesManagerClient(conn)
	_, err = c.UpdateArticle(ctx, &amv1.UpdateArticleRequest{
		Id:      aid.String(),
		Article: articleForInsert,
	})
	if err != nil {
		log.Error("Failed to update article", sl.Err(err))
		return models.Article{}, fmt.Errorf("%s: %w", op, err)
	}

	return article, nil
}

// Delete implements articles.IArticlesStorage.
func (a *ArticlesManageStorage) Delete(ctx context.Context, aid uuid.UUID) (models.Article, error) {
	const op = "articlesmanagestorage.delete"
	log := a.log.With(slog.String("op", op))

	select {
	case <-ctx.Done():
		return models.Article{}, fmt.Errorf("%s: %w", op, ctx.Err())
	default:
	}

	conn, err := grpc.NewClient(
		fmt.Sprintf("%s:%d", a.ServiceHost, a.ServicePort),
		grpc.WithTransportCredentials(insecure.NewCredentials()),
	)
	if err != nil {
		log.Error("Failed to connect to gRPC server", sl.Err(err))
		return models.Article{}, fmt.Errorf("%s: %w", op, err)
	}
	defer conn.Close()

	c := amv1.NewArticlesManagerClient(conn)
	res, err := c.DeleteArticle(ctx, &amv1.DeleteArticleRequest{
		Id: aid.String(),
	})
	if err != nil {
		log.Error("Failed to delete article", sl.Err(err))
		return models.Article{}, fmt.Errorf("%s: %w", op, err)
	}

	resp_article, err := amprofiles.ProtoArtToArt(res.GetArticle())
	if err != nil {
		log.Error("Wrong structure, failed to customize", sl.Err(err))
		return models.Article{}, fmt.Errorf("%s: %w", op, err)
	}

	return resp_article, nil
}
