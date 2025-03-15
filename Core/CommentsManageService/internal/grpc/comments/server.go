package grpccomments

import (
	"commentsManageService/internal/domain/interfaces/service"
	cmprofiles "commentsManageService/internal/domain/profiles"
	"commentsManageService/pkg/lib/logger/sl"
	"context"
	"errors"
	"log/slog"

	service_error "commentsManageService/internal/service"

	cmv1 "github.com/chas3air/protos/gen/go/commentsManager"
	"github.com/google/uuid"
	"google.golang.org/grpc"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

type serverAPI struct {
	cmv1.UnimplementedCommentsManagerServer
	commentService service.CommentService
	log            *slog.Logger
}

func Register(grpc *grpc.Server, commentsService service.CommentService, log *slog.Logger) {
	cmv1.RegisterCommentsManagerServer(grpc, &serverAPI{commentService: commentsService, log: log})
}

func (s *serverAPI) GetCommentById(ctx context.Context, req *cmv1.GetCommentByIdRequest) (*cmv1.GetCommentByIdResponse, error) {
	const op = "grpc.commentsManager.getCommentById"
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
		log.Error("Failed to get id", sl.Err(errors.New("required parametr id")))
		return nil, status.Error(codes.InvalidArgument, "required parametr id")
	}
	parsedUUID, err := uuid.Parse(id_s)
	if err != nil {
		log.Error("Invalid id, must be uuid")
		return nil, status.Error(codes.InvalidArgument, "invalid id, must be uuid")
	}

	commentFromDB, err := s.commentService.GetCommentById(ctx, parsedUUID)
	if err != nil {
		if errors.Is(err, service_error.ErrNotFound) {
			log.Warn("Comment with current id not found", sl.Err(err))
			return nil, status.Error(codes.NotFound, "comment with current id not found")
		}

		log.Error("failed to retrieve comment by id", sl.Err(err))
		return nil, status.Error(codes.Internal, "failed to retrieve comment by id")
	}

	responsed_comment, err := cmprofiles.ComToProtoCom(commentFromDB)
	if err != nil {
		log.Error("Wrogn structure, failed to customize", sl.Err(err))
		return nil, status.Error(codes.Internal, "wrong structure")
	}

	return &cmv1.GetCommentByIdResponse{
		Comment: responsed_comment,
	}, nil
}

func (s *serverAPI) GetCommentsByArticleId(ctx context.Context, req *cmv1.GetCommentsByArticleIdRequest) (*cmv1.GetCommentsByArticleIdResponse, error) {
	const op = "grpc.commentsManager.getCommentsByArticleId"
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
		log.Error("Failed to get article_id", sl.Err(errors.New("required parametr article_id")))
		return nil, status.Error(codes.InvalidArgument, "required parametr article_id")
	}
	parsedUUID, err := uuid.Parse(id_s)
	if err != nil {
		log.Error("Invalid article_id, must be uuid")
		return nil, status.Error(codes.InvalidArgument, "invalid article_id, must be uuid")
	}

	commentsFromDB, err := s.commentService.GetCommentsByArticleId(ctx, parsedUUID)
	if err != nil {
		log.Error("Error retrieving comments", sl.Err(err))
		return nil, status.Error(codes.Internal, "error  retrieving comments")
	}

	resp_comments := make([]*cmv1.Comment, 0, 10)
	for _, comment := range commentsFromDB {
		profiled_comment, err := cmprofiles.ComToProtoCom(comment)
		if err != nil {
			log.Warn("Wrong structure", sl.Err(err))
			continue
		}

		resp_comments = append(resp_comments, profiled_comment)
	}

	return &cmv1.GetCommentsByArticleIdResponse{
		Comments: resp_comments,
	}, nil
}

func (s *serverAPI) Insert(ctx context.Context, req *cmv1.InsertRequest) (*cmv1.InsertResponse, error) {
	const op = "grpc.commentsManager.insert"
	log := s.log.With(
		slog.String("op", op),
	)

	select {
	case <-ctx.Done():
		return nil, status.Error(codes.DeadlineExceeded, "request time out")
	default:
	}

	commentFromReq := req.GetComment()
	if commentFromReq == nil {
		log.Error("Comment is required parametr", sl.Err(errors.New("required parament comment")))
		return nil, status.Error(codes.InvalidArgument, "required parament comment")
	}

	commentForInsert, err := cmprofiles.ProtoComToCom(commentFromReq)
	if err != nil {
		log.Error("Wrong structure", sl.Err(err))
		return nil, status.Error(codes.InvalidArgument, "wrong structure")
	}

	_, err = s.commentService.Insert(ctx, commentForInsert)
	if err != nil {
		if errors.Is(err, service_error.ErrAlreadyExists) {
			log.Warn("Comment already exists", sl.Err(err))
			return nil, status.Error(codes.AlreadyExists, "comment already exists")
		}

		log.Error("Failed to insert comment", sl.Err(err))
		return nil, status.Error(codes.Internal, "failed to insert comment")
	}

	return &cmv1.InsertResponse{
		Comment: commentFromReq,
	}, nil
}

func (s *serverAPI) Delete(ctx context.Context, req *cmv1.DeleteRequest) (*cmv1.DeleteResponse, error) {
	const op = "grpc.commentsManager.delete"
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
		log.Error("Failed to get id", sl.Err(errors.New("required parametr id")))
		return nil, status.Error(codes.InvalidArgument, "required parametr id")
	}
	parsedUUID, err := uuid.Parse(id_s)
	if err != nil {
		log.Error("Invalid id, must be uuid")
		return nil, status.Error(codes.InvalidArgument, "invalid id, must be uuid")
	}

	deleted_comment, err := s.commentService.Delete(ctx, parsedUUID)
	if err != nil {
		if errors.Is(err, service_error.ErrNotFound) {
			log.Warn("Comment not found", sl.Err(err))
			return nil, status.Error(codes.NotFound, "comment not found")
		}

		log.Error("Failed to delete comment", sl.Err(err))
		return nil, status.Error(codes.Internal, "failed to delete comment")
	}

	resp_comment, err := cmprofiles.ComToProtoCom(deleted_comment)
	if err != nil {
		log.Error("Wrong structure", sl.Err(err))
		return nil, status.Error(codes.Internal, "wrong structure")
	}

	return &cmv1.DeleteResponse{
		Comment: resp_comment,
	}, nil
}
