package commentcontroller

import (
	"apigateway/internal/domain/interfaces/comments"
	"apigateway/internal/domain/models"
	"apigateway/pkg/lib/logger/sl"
	"context"
	"encoding/json"
	"errors"
	"log/slog"
	"net/http"

	"github.com/google/uuid"
	"github.com/gorilla/mux"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

type CommentController struct {
	log            *slog.Logger
	commentService comments.CommentsService
}

func New(log *slog.Logger, commentsService comments.CommentsService) *CommentController {
	return &CommentController{
		log:            log,
		commentService: commentsService,
	}
}

func (cs *CommentController) handleError(w http.ResponseWriter, err error, log *slog.Logger) {
	if errors.Is(err, context.Canceled) {
		log.Error("Request was canceled by the user")
		http.Error(w, "Request canceled", http.StatusRequestTimeout)
	} else if errors.Is(err, context.DeadlineExceeded) || status.Code(err) == codes.DeadlineExceeded {
		log.Error("Request time out")
		http.Error(w, "Request timeout", http.StatusRequestTimeout)
	} else {
		log.Error("Operation failed", sl.Err(err))
		http.Error(w, err.Error(), http.StatusInternalServerError)
	}
}

func (cs *CommentController) GetCommentById(w http.ResponseWriter, r *http.Request) {
	const op = "controllers.commentsController.getCommentById"
	log := cs.log.With(
		slog.String("op", op),
	)

	id_s := mux.Vars(r)["id"]
	if id_s == "" {
		log.Error("Failed to get id")
		http.Error(w, "failed to get id", http.StatusBadRequest)
		return
	}

	parsedUUID, err := uuid.Parse(id_s)
	if err != nil {
		log.Error("Invalid id, must be uuid", sl.Err(err))
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	comment, err := cs.commentService.GetCommentById(r.Context(), parsedUUID)
	if err != nil {
		cs.handleError(w, err, log)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	if err := json.NewEncoder(w).Encode(comment); err != nil {
		log.Error("Failed to encode response", sl.Err(err))
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	log.Info("Retrieved comment by ID successfully")
}

func (cs *CommentController) GetCommentsByArticleId(w http.ResponseWriter, r *http.Request) {
	const op = "controllers.commentsController.getCommentByArticleId"
	log := cs.log.With(
		slog.String("op", op),
	)

	id_s := mux.Vars(r)["article_id"]
	if id_s == "" {
		log.Error("Failed to get article_id")
		http.Error(w, "failed to get article_id", http.StatusBadRequest)
		return
	}

	parsedUUID, err := uuid.Parse(id_s)
	if err != nil {
		log.Error("Invalid article_id, must be uuid", sl.Err(err))
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	comments, err := cs.commentService.GetCommentsByArticleId(r.Context(), parsedUUID)
	if err != nil {
		cs.handleError(w, err, log)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	if err := json.NewEncoder(w).Encode(comments); err != nil {
		log.Error("Failed to encode response", sl.Err(err))
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	log.Info("Retrieved comments by article_id successfully")
}

func (cs *CommentController) Insert(w http.ResponseWriter, r *http.Request) {
	const op = "controllers.commentsController.insert"
	log := cs.log.With(
		slog.String("op", op),
	)

	var comment models.Comment
	if err := json.NewDecoder(r.Body).Decode(&comment); err != nil {
		log.Error("Cannot parse request body", sl.Err(err))
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	comment, err := cs.commentService.Insert(r.Context(), comment)
	if err != nil {
		cs.handleError(w, err, log)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	if err := json.NewEncoder(w).Encode(comment); err != nil {
		log.Error("Failed to encode response", sl.Err(err))
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	log.Info("Inserting comment successfully")
}

func (cs *CommentController) Delete(w http.ResponseWriter, r *http.Request) {
	const op = "controllers.commentsController.delete"
	log := cs.log.With(
		slog.String("op", op),
	)

	id_s := mux.Vars(r)["id"]
	if id_s == "" {
		log.Error("Failed to get id")
		http.Error(w, "failed to get id", http.StatusBadRequest)
		return
	}

	parsedUUID, err := uuid.Parse(id_s)
	if err != nil {
		log.Error("Invalid id, must be uuid", sl.Err(err))
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	comment, err := cs.commentService.Delete(r.Context(), parsedUUID)
	if err != nil {
		cs.handleError(w, err, log)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	if err := json.NewEncoder(w).Encode(comment); err != nil {
		log.Error("Failed to encode response", sl.Err(err))
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	log.Info("Deleting comment successfully")
}
