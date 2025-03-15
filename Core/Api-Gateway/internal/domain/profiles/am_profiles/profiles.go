package amprofiles

import (
	"apigateway/internal/domain/models"

	amv1 "github.com/chas3air/protos/gen/go/articlesManager"

	"github.com/google/uuid"
	"google.golang.org/protobuf/types/known/timestamppb"
)

func ProtoArtToArt(article *amv1.Article) (models.Article, error) {
	createdAt := article.GetCreatedAt().AsTime()

	ownerId, err := uuid.Parse(article.OwnerId)
	if err != nil {
		return models.Article{}, err
	}

	return models.Article{
		Id:        uuid.MustParse(article.Id),
		CreatedAt: createdAt,
		Title:     article.Title,
		Content:   article.Content,
		OwnerId:   ownerId,
	}, nil
}

func ArtToProtoArt(article models.Article) (*amv1.Article, error) {
	createdAt := timestamppb.New(article.CreatedAt)

	return &amv1.Article{
		Id:        article.Id.String(),
		CreatedAt: createdAt,
		Title:     article.Title,
		Content:   article.Content,
		OwnerId:   article.OwnerId.String(),
	}, nil
}
