package cmprofiles

import (
	"errors"

	"apigateway/internal/domain/models"

	cmv1 "github.com/chas3air/protos/gen/go/commentsManager"
	"github.com/google/uuid"
	"google.golang.org/protobuf/types/known/timestamppb"
)

func ComToProtoCom(comment models.Comment) (*cmv1.Comment, error) {
	if comment.Id == uuid.Nil {
		return &cmv1.Comment{}, errors.New("invalid ID")
	}

	return &cmv1.Comment{
		Id:        comment.Id.String(),
		ArticleId: comment.ArticleId.String(),
		OwnerId:   comment.OwnerId.String(),
		CreatedAt: timestamppb.New(comment.CreatedAt),
		Content:   comment.Content,
	}, nil
}

func ProtoComToCom(comment *cmv1.Comment) (models.Comment, error) {
	if comment == nil || comment.Id == "" {
		return models.Comment{}, errors.New("invalid comment or ID")
	}

	createdAt := comment.GetCreatedAt().AsTime()

	id, err := uuid.Parse(comment.Id)
	if err != nil {
		return models.Comment{}, err
	}

	articleId, err := uuid.Parse(comment.ArticleId)
	if err != nil {
		return models.Comment{}, err
	}

	ownerId, err := uuid.Parse(comment.OwnerId)
	if err != nil {
		return models.Comment{}, err
	}

	return models.Comment{
		Id:        id,
		ArticleId: articleId,
		OwnerId:   ownerId,
		CreatedAt: createdAt,
		Content:   comment.Content,
	}, nil
}
