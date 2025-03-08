package commentsstorage

import (
	"apigateway/internal/domain/models"
	cmprofiles "apigateway/internal/domain/profiles/cm_profiles"
	"apigateway/pkg/lib/logger/sl"
	"context"
	"errors"
	"fmt"
	"log/slog"

	"github.com/google/uuid"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"

	cmv1 "github.com/chas3air/protos/gen/go/commentsManager"
)

type CommentsManageStorage struct {
	log         *slog.Logger
	ServiceHost string
	ServicePost int
}

func New(log *slog.Logger, serviceHost string, servicePort int) *CommentsManageStorage {
	return &CommentsManageStorage{
		log:         log,
		ServiceHost: serviceHost,
		ServicePost: servicePort,
	}
}

func (cms *CommentsManageStorage) GetCommentById(ctx context.Context, cid uuid.UUID) (models.Comment, error) {
	const op = "commentsManageStorage.getCommentById"
	log := cms.log.With(
		slog.String("op", op),
	)

	select {
	case <-ctx.Done():
		return models.Comment{}, fmt.Errorf("%s: %w", op, ctx.Err())
	default:
	}

	conn, err := grpc.NewClient(
		fmt.Sprintf("%s:%d", cms.ServiceHost, cms.ServicePost),
		grpc.WithTransportCredentials(insecure.NewCredentials()),
	)
	if err != nil {
		log.Error("failed to connect to gRPC server", sl.Err(err))
		return models.Comment{}, fmt.Errorf("%s: %w", op, err)
	}
	defer conn.Close()

	c := cmv1.NewCommentsManagerClient(conn)
	res, err := c.GetCommentById(ctx, &cmv1.GetCommentByIdRequest{
		Id: cid.String(),
	})
	if err != nil {
		log.Error("Failed to get comment by id", sl.Err(err))
		return models.Comment{}, fmt.Errorf("%s: %w", op, err)
	}

	if res.GetComment() == nil {
		log.Warn("Comment is empty", sl.Err(errors.New("empty comment")))
		return models.Comment{}, fmt.Errorf("%s: %w", op, errors.New("empty comment"))
	}

	resComment, err := cmprofiles.ProtoComToCom(res.GetComment())
	if err != nil {
		log.Error("Wrong structure", sl.Err(err))
		return models.Comment{}, fmt.Errorf("%s: %w", op, errors.New("wrong structure"))
	}

	return resComment, nil
}

func (cms *CommentsManageStorage) GetCommentsByArticleId(ctx context.Context, aid uuid.UUID) ([]models.Comment, error) {
	const op = "commentsManageStorage.getCommentsByArticleId"
	log := cms.log.With(
		slog.String("op", op),
	)

	select {
	case <-ctx.Done():
		return nil, fmt.Errorf("%s: %w", op, ctx.Err())
	default:
	}

	conn, err := grpc.NewClient(
		fmt.Sprintf("%s:%d", cms.ServiceHost, cms.ServicePost),
		grpc.WithTransportCredentials(insecure.NewCredentials()),
	)
	if err != nil {
		log.Error("failed to connect to gRPC server", sl.Err(err))
		return nil, fmt.Errorf("%s: %w", op, err)
	}
	defer conn.Close()

	c := cmv1.NewCommentsManagerClient(conn)
	res, err := c.GetCommentsByArticleId(ctx, &cmv1.GetCommentsByArticleIdRequest{
		ArticleId: aid.String(),
	})
	if err != nil {
		log.Error("failed to get comments by article_id", sl.Err(err))
		return nil, fmt.Errorf("%s: %w", op, err)
	}

	if res.GetComments() == nil {
		log.Warn("No one comment")
		return nil, nil
	}

	var resComment = make([]models.Comment, 0, 10)
	for _, pbComment := range res.GetComments() {
		comment, err := cmprofiles.ProtoComToCom(pbComment)
		if err != nil {
			log.Warn("Wrong structure", sl.Err(err))
			continue
		}

		resComment = append(resComment, comment)
	}

	return resComment, nil
}

func (cms *CommentsManageStorage) Insert(ctx context.Context, comment models.Comment) (models.Comment, error) {
	const op = "commentsManageStorage.insert"
	log := cms.log.With(
		slog.String("op", op),
	)

	select {
	case <-ctx.Done():
		return models.Comment{}, fmt.Errorf("%s: %w", op, ctx.Err())
	default:
	}

	conn, err := grpc.NewClient(
		fmt.Sprintf("%s:%d", cms.ServiceHost, cms.ServicePost),
		grpc.WithTransportCredentials(insecure.NewCredentials()),
	)
	if err != nil {
		log.Error("failed to connect to gRPC server", sl.Err(err))
		return models.Comment{}, fmt.Errorf("%s: %w", op, err)
	}
	defer conn.Close()

	commentForInsert, err := cmprofiles.ComToProtoCom(comment)
	if err != nil {
		log.Error("Wrong structure", sl.Err(err))
		return models.Comment{}, fmt.Errorf("%s: %w", op, err)
	}

	c := cmv1.NewCommentsManagerClient(conn)
	_, err = c.Insert(ctx, &cmv1.InsertRequest{
		Comment: commentForInsert,
	})
	if err != nil {
		log.Error("Failed to insert comment", sl.Err(err))
		return models.Comment{}, fmt.Errorf("%s: %w", op, err)
	}

	return comment, nil
}

func (cms *CommentsManageStorage) Delete(ctx context.Context, cid uuid.UUID) (models.Comment, error) {
	const op = "commentsManageStorage.delete"
	log := cms.log.With(
		slog.String("op", op),
	)

	select {
	case <-ctx.Done():
		return models.Comment{}, fmt.Errorf("%s: %w", op, ctx.Err())
	default:
	}

	conn, err := grpc.NewClient(
		fmt.Sprintf("%s:%d", cms.ServiceHost, cms.ServicePost),
		grpc.WithTransportCredentials(insecure.NewCredentials()),
	)
	if err != nil {
		log.Error("failed to connect to gRPC server", sl.Err(err))
		return models.Comment{}, fmt.Errorf("%s: %w", op, err)
	}
	defer conn.Close()

	c := cmv1.NewCommentsManagerClient(conn)
	res, err := c.Delete(ctx, &cmv1.DeleteRequest{
		Id: cid.String(),
	})
	if err != nil {
		log.Error("Failed to delete comment", sl.Err(err))
		return models.Comment{}, fmt.Errorf("%s: %w", op, err)
	}

	deleted_comment, err := cmprofiles.ProtoComToCom(res.GetComment())
	if err != nil {
		log.Error("Wrong structure", sl.Err(err))
		return models.Comment{}, fmt.Errorf("%s: %w", op, err)
	}

	return deleted_comment, nil
}
