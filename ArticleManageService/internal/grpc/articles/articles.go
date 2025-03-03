package articlesManager

import (
	"articlesManageService/internal/domain/interfaces/articlesservice"
	"articlesManageService/internal/domain/profiles"
	"context"
	"log"

	amv1 "github.com/chas3air/protos/gen/go/articlesManager"
	"github.com/google/uuid"
	"google.golang.org/grpc"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

type serverAPI struct {
	amv1.UnimplementedArticlesManagerServer
	articlesManager articlesservice.ArticlesManager
}

func Register(grpc *grpc.Server, articleManager articlesservice.ArticlesManager) {
	amv1.RegisterArticlesManagerServer(grpc, &serverAPI{articlesManager: articleManager})
}

// GetArticles implements amv1.ArticlesManagerServer.
func (s *serverAPI) GetArticles(ctx context.Context, req *amv1.GetArticlesRequest) (*amv1.GetArticlesResponse, error) {
	select {
	case <-ctx.Done():
		return nil, status.Error(codes.DeadlineExceeded, "request timed out")
	default:
	}

	app_aticles, err := s.articlesManager.GetAtricles(ctx)
	if err != nil {
		return nil, status.Error(codes.Internal, "failed to retrieve app articles")
	}

	resp_articles := make([]*amv1.Article, 0, len(app_aticles))
	for _, article := range app_aticles {
		profiles_article, err := profiles.ArtToProtoArt(article)
		if err != nil {
			log.Println("error article:", article)
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
	select {
	case <-ctx.Done():
		return nil, status.Error(codes.DeadlineExceeded, "request timed out")
	default:
	}

	id_s := req.GetArticleId()
	if id_s == "" {
		return nil, status.Error(codes.InvalidArgument, "requared parametr id")
	}

	parsedUUID, err := uuid.Parse(req.GetArticleId())
	if err != nil {
		return nil, status.Error(codes.InvalidArgument, "invalid id")
	}

	app_article, err := s.articlesManager.GetArticleById(ctx, parsedUUID)
	if err != nil {
		return nil, status.Error(codes.Internal, "failed to retrieve article by id")
	}

	profied_article, err := profiles.ArtToProtoArt(app_article)
	if err != nil {
		return nil, status.Error(codes.Internal, "failed to customise article")
	}

	return &amv1.GetArticleByIdResponse{
		Article: profied_article,
	}, nil
}

// GetArticleByOwnerId implements amv1.ArticlesManagerServer.
func (s *serverAPI) GetArticlesByOwnerId(ctx context.Context, req *amv1.GetArticlesByOwnerIdRequest) (*amv1.GetArticlesByOwnerIdResponse, error) {
	select {
	case <-ctx.Done():
		return nil, status.Error(codes.DeadlineExceeded, "request timed out")
	default:
	}

	owner_id_s := req.GetOwnerId()
	if owner_id_s == "" {
		return nil, status.Error(codes.InvalidArgument, "invalid owner_id")
	}
	parseOwnerId, err := uuid.Parse(owner_id_s)
	if err != nil {
		return nil, status.Error(codes.InvalidArgument, "invalid owner_id")
	}

	app_articles, err := s.articlesManager.GetArticleByOwnerId(ctx, parseOwnerId)
	if err != nil {
		return nil, status.Error(codes.Internal, "failed to retrieve articles by owner_id")
	}

	resp_articles := make([]*amv1.Article, 0, len(app_articles))
	for _, article := range app_articles {
		profiled_article, err := profiles.ArtToProtoArt(article)
		if err != nil {
			log.Println("error article:", err)
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
	select {
	case <-ctx.Done():
		return nil, status.Error(codes.DeadlineExceeded, "request timed out")
	default:
	}

	app_article, err := profiles.ProtoArtToArt(*req.GetArticle())
	if err != nil {
		return nil, status.Error(codes.InvalidArgument, "failed to customise")
	}

	err = s.articlesManager.Insert(ctx, app_article)
	if err != nil {
		return nil, status.Error(codes.Internal, "failed to insert article")
	}

	return nil, nil
}

// UpdateArticle implements amv1.ArticlesManagerServer.
func (s *serverAPI) UpdateArticle(ctx context.Context, req *amv1.UpdateArticleRequest) (*amv1.UpdateArticleResponse, error) {
	select {
	case <-ctx.Done():
		return nil, status.Error(codes.DeadlineExceeded, "request timed out")
	default:
	}

	id_s := req.GetId()
	if id_s == "" {
		return nil, status.Error(codes.InvalidArgument, "requared paramert id")
	}
	parseUUID, err := uuid.Parse(id_s)
	if err != nil {
		return nil, status.Error(codes.InvalidArgument, "invalid id")
	}

	app_article, err := profiles.ProtoArtToArt(*req.GetArticle())
	if err != nil {
		return nil, status.Error(codes.InvalidArgument, "failed to customise")
	}

	err = s.articlesManager.Update(ctx, parseUUID, app_article)
	if err != nil {
		return nil, status.Error(codes.Internal, "failed to update article")
	}

	return nil, nil
}

// DeleteArticle implements amv1.ArticlesManagerServer.
func (s *serverAPI) DeleteArticle(ctx context.Context, req *amv1.DeleteArticleRequest) (*amv1.DeleteArticleResponse, error) {
	select {
	case <-ctx.Done():
		return nil, status.Error(codes.DeadlineExceeded, "request timed out")
	default:
	}

	id_s := req.GetId()
	if id_s == "" {
		return nil, status.Error(codes.InvalidArgument, "required parametr id")
	}

	parseUUID, err := uuid.Parse(id_s)
	if err != nil {
		return nil, status.Error(codes.InvalidArgument, "invalid id")
	}

	deleted_article, err := s.articlesManager.Delete(ctx, parseUUID)
	if err != nil {
		return nil, status.Error(codes.Internal, "failed to delete article")
	}

	resp_article, err := profiles.ArtToProtoArt(deleted_article)
	if err != nil {
		return nil, status.Error(codes.Internal, "failed to customise")
	}

	return &amv1.DeleteArticleResponse{
		Article: resp_article,
	}, nil
}
