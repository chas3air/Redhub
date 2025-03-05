package articlesmanager

import (
	"articlesManageService/internal/domain/interfaces/articlesservice"
	"articlesManageService/internal/domain/models"
	"articlesManageService/internal/domain/profiles"
	"articlesManageService/internal/services"
	"articlesManageService/pkg/lib/logger/sl"
	"context"
	"errors"
	"log/slog"

	amv1 "github.com/chas3air/protos/gen/go/articlesManager"
	"github.com/google/uuid"
	"google.golang.org/grpc"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

type serverAPI struct {
	amv1.UnimplementedArticlesManagerServer
	articlesManager articlesservice.ArticlesManager
	log             *slog.Logger
}

func Register(grpc *grpc.Server, articleManager articlesservice.ArticlesManager, log *slog.Logger) {
	amv1.RegisterArticlesManagerServer(grpc, &serverAPI{articlesManager: articleManager, log: log})
}

// GetArticles implements amv1.ArticlesManagerServer.
func (s *serverAPI) GetArticles(ctx context.Context, req *amv1.GetArticlesRequest) (*amv1.GetArticlesResponse, error) {
	const op = "grpc.articles.getArticles"
	log := s.log.With(
		slog.String("op", op),
	)

	select {
	case <-ctx.Done():
		return nil, status.Error(codes.DeadlineExceeded, "request time out")
	default:
	}

	app_articles, err := s.articlesManager.GetArticles(ctx)
	if err != nil {
		log.Error("Failed retrieving articles", sl.Err(err))
		return nil, status.Error(codes.Internal, "failed to retrieve app articles")
	}

	resp_articles := make([]*amv1.Article, 0, len(app_articles))
	for _, article := range app_articles {
		profiles_article, err := profiles.ArtToProtoArt(article)
		if err != nil {
			log.Error("Wrong structure, failed to customize", sl.Err(err))
			continue
		}

		resp_articles = append(resp_articles, profiles_article)
	}

	return &amv1.GetArticlesResponse{
		Articles: resp_articles,
	}, nil
}

// GetArticleById implements amv1.ArticlesManagerServer.
func (s *serverAPI) GetArticleById(ctx context.Context, req *amv1.GetArticleByIdRequest) (*amv1.GetArticleByIdResponse, error) {
	const op = "grpc.articles.getArticleById"
	log := s.log.With(
		slog.String("op", op),
	)

	select {
	case <-ctx.Done():
		return nil, status.Error(codes.DeadlineExceeded, "request time out")
	default:
	}

	id_s := req.GetArticleId()
	if id_s == "" {
		log.Error("Failed to get id", sl.Err(errors.New("required parametr id")))
		return nil, status.Error(codes.InvalidArgument, "required parameter id")
	}

	parsedUUID, err := uuid.Parse(id_s)
	if err != nil {
		log.Error("Invalid id, must be uuid", sl.Err(err))
		return nil, status.Error(codes.InvalidArgument, "invalid id, must be uuid")
	}

	app_article, err := s.articlesManager.GetArticleById(ctx, parsedUUID)
	if err != nil {
		if errors.Is(err, services.ErrNotFound) {
			log.Warn("Article with current id not found", sl.Err(err))
			return nil, status.Error(codes.NotFound, "Article with current id not found")
		}

		log.Error("failed to retrieve article by id", sl.Err(err))
		return nil, status.Error(codes.Internal, "failed to retrieve article by id")
	}

	profiles_article, err := profiles.ArtToProtoArt(app_article)
	if err != nil {
		log.Error("Wrong structure, failed to customize", sl.Err(err))
		return nil, status.Error(codes.Internal, "failed to customize")
	}

	return &amv1.GetArticleByIdResponse{
		Article: profiles_article,
	}, nil
}

// GetArticlesByOwnerId implements amv1.ArticlesManagerServer.
func (s *serverAPI) GetArticlesByOwnerId(ctx context.Context, req *amv1.GetArticlesByOwnerIdRequest) (*amv1.GetArticlesByOwnerIdResponse, error) {
	const op = "grpc.articles.getArticlesByOwnerId"
	log := s.log.With(
		slog.String("op", op),
	)

	select {
	case <-ctx.Done():
		return nil, status.Error(codes.DeadlineExceeded, "request time out")
	default:
	}

	owner_id_s := req.GetOwnerId()
	if owner_id_s == "" {
		log.Error("Owner_id is required", sl.Err(errors.New("owner_id is required")))
		return nil, status.Error(codes.InvalidArgument, "owner_id is required")
	}

	parseOwnerId, err := uuid.Parse(owner_id_s)
	if err != nil {
		log.Error("Invalid owner_id, must be uuid", sl.Err(err))
		return nil, status.Error(codes.InvalidArgument, "invalid owner_id, must be uuid")
	}

	app_articles, err := s.articlesManager.GetArticleByOwnerId(ctx, parseOwnerId)
	if err != nil {
		log.Error("Failed to retrieve articles by owner_id", sl.Err(err))
		return nil, status.Error(codes.Internal, "failed to retrieve articles by owner_id")
	}

	resp_articles := make([]*amv1.Article, 0, len(app_articles))
	for _, article := range app_articles {
		profiled_article, err := profiles.ArtToProtoArt(article)
		if err != nil {
			log.Error("Wrong structure, failed to customize", sl.Err(err))
			continue
		}

		resp_articles = append(resp_articles, profiled_article)
	}

	return &amv1.GetArticlesByOwnerIdResponse{
		Articles: resp_articles,
	}, nil
}

// InsertArticle implements amv1.ArticlesManagerServer.
func (s *serverAPI) InsertArticle(ctx context.Context, req *amv1.InsertArticleRequest) (*amv1.InsertArticleResponse, error) {
	const op = "grpc.articles.insertArticle"
	log := s.log.With(
		slog.String("op", op),
	)

	select {
	case <-ctx.Done():
		return nil, status.Error(codes.DeadlineExceeded, "request time out")
	default:
	}

	app_article, err := profiles.ProtoArtToArt(req.GetArticle())
	if err != nil {
		log.Error("Wrong structure, failed to customize", sl.Err(err))
		return nil, status.Error(codes.InvalidArgument, "failed to customize")
	}

	err = s.articlesManager.Insert(ctx, app_article)
	if err != nil {
		if errors.Is(err, services.ErrAlreadyExists) {
			log.Warn("Article already exists", sl.Err(err))
			return nil, status.Error(codes.AlreadyExists, "article already exists")
		}

		log.Error("Failed to insert article", sl.Err(err))
		return nil, status.Error(codes.Internal, "failed to insert article")
	}

	return &amv1.InsertArticleResponse{}, nil
}

// UpdateArticle implements amv1.ArticlesManagerServer.
func (s *serverAPI) UpdateArticle(ctx context.Context, req *amv1.UpdateArticleRequest) (*amv1.UpdateArticleResponse, error) {
	const op = "grpc.articles.updateArticle"
	log := s.log.With(
		slog.String("op", op),
	)

	select {
	case <-ctx.Done():
		return nil, status.Error(codes.DeadlineExceeded, "request time out")
	default:
	}

	id_s := req.GetId()
	if id_s == "" {
		log.Error("Id is required", sl.Err(errors.New("id is required")))
		return nil, status.Error(codes.InvalidArgument, "id is required")
	}

	parseUUID, err := uuid.Parse(id_s)
	if err != nil {
		log.Error("Invalid owner_id, must be uuid", sl.Err(err))
		return nil, status.Error(codes.InvalidArgument, "invalid owner_id, must be uuid")
	}

	if req.GetArticle() == nil {
		log.Error("Article is required", sl.Err(err))
		return nil, status.Error(codes.InvalidArgument, "article is required")
	}

	app_article := models.Article{Title: req.Article.Title, Content: req.Article.Content}

	err = s.articlesManager.Update(ctx, parseUUID, app_article)
	if err != nil {
		if errors.Is(err, services.ErrNotFound) {
			log.Warn("Article with current id not found", sl.Err(err))
			return nil, status.Error(codes.NotFound, "Article with current id not found")
		}

		log.Error("Failed to update article", sl.Err(err))
		return nil, status.Error(codes.Internal, "failed to update article")
	}

	return &amv1.UpdateArticleResponse{}, nil
}

// DeleteArticle implements amv1.ArticlesManagerServer.
func (s *serverAPI) DeleteArticle(ctx context.Context, req *amv1.DeleteArticleRequest) (*amv1.DeleteArticleResponse, error) {
	const op = "grpc.articles.deleteArticle"
	log := s.log.With(
		slog.String("op", op),
	)

	select {
	case <-ctx.Done():
		return nil, status.Error(codes.DeadlineExceeded, "request time out")
	default:
	}

	id_s := req.GetId()
	if id_s == "" {
		log.Error("Id is required", sl.Err(errors.New("id is required")))
		return nil, status.Error(codes.InvalidArgument, "required parameter id")
	}

	parseUUID, err := uuid.Parse(id_s)
	if err != nil {
		return nil, status.Error(codes.InvalidArgument, "invalid id")
	}

	deleted_article, err := s.articlesManager.Delete(ctx, parseUUID)
	if err != nil {
		if errors.Is(err, services.ErrNotFound) {
			log.Warn("Article with current id not found", sl.Err(err))
			return nil, status.Error(codes.NotFound, "Article with current id not found")
		}

		log.Error("Failed to delete article", sl.Err(err))
		return nil, status.Error(codes.Internal, "failed to delete article")
	}

	resp_article, err := profiles.ArtToProtoArt(deleted_article)
	if err != nil {
		log.Error("Wrong structure, failed to customize", sl.Err(err))
		return nil, status.Error(codes.Internal, "failed to customize")
	}

	return &amv1.DeleteArticleResponse{
		Article: resp_article,
	}, nil
}
