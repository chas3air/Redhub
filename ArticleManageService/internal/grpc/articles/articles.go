package articlesManager

import (
	"articlesManageService/internal/domain/interfaces/articlesservice"
	"context"

	amv1 "github.com/chas3air/protos/gen/go/articlesManager"
	"google.golang.org/grpc"
)

type serverAPI struct {
	amv1.UnimplementedArticlesManagerServer
	articlesManager articlesservice.ArticlesManager
}

func Register(grpc *grpc.Server, articleManager articlesservice.ArticlesManager) {
	amv1.RegisterArticlesManagerServer(grpc, &serverAPI{articlesManager: articleManager})
}

// GetArticles implements amv1.ArticlesManagerServer.
func (s *serverAPI) GetArticles(context.Context, *amv1.GetArticlesRequest) (*amv1.GetArticlesResponse, error) {
	panic("unimplemented")
}

// GetArticleById implements amv1.ArticlesManagerServer.
func (s *serverAPI) GetArticleById(context.Context, *amv1.GetArticleByIdRequest) (*amv1.GetArticleByIdResponse, error) {
	panic("unimplemented")
}

// GetArticleByOwnerId implements amv1.ArticlesManagerServer.
func (s *serverAPI) GetArticleByOwnerId(context.Context, *amv1.GetArticleByOwnerIdRequest) (*amv1.GetArticleByOwnerIdResponse, error) {
	panic("unimplemented")
}

// InsertArticle implements amv1.ArticlesManagerServer.
func (s *serverAPI) InsertArticle(context.Context, *amv1.InsertArticleRequest) (*amv1.InsertArticleResponse, error) {
	panic("unimplemented")
}

// UpdateArticle implements amv1.ArticlesManagerServer.
func (s *serverAPI) UpdateArticle(context.Context, *amv1.UpdateArticleRequest) (*amv1.UpdateArticleResponse, error) {
	panic("unimplemented")
}

// DeleteArticle implements amv1.ArticlesManagerServer.
func (s *serverAPI) DeleteArticle(context.Context, *amv1.DeleteArticleRequest) (*amv1.DeleteArticleResponse, error) {
	panic("unimplemented")
}
