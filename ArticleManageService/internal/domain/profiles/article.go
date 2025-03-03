package profiles

import (
	"articlesManageService/internal/domain/models"

	amv1 "github.com/chas3air/protos/gen/go/articlesManager"
)

func ProtoArtToArt(article amv1.Article) (models.Article, error) {

	return models.Article{}, nil
}

func ArtToProtoArt(article models.Article) (*amv1.Article, error) {
	return &amv1.Article{}, nil
}
