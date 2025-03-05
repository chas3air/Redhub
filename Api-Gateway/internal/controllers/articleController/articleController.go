package articlecontroller

import (
	"apigateway/internal/domain/interfaces/articles"
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

type ArticleController struct {
	log            *slog.Logger
	articleService articles.IArticlesService
}

func New(log *slog.Logger, articleService articles.IArticlesService) *ArticleController {
	return &ArticleController{
		log:            log,
		articleService: articleService,
	}
}

func (ac *ArticleController) handleError(w http.ResponseWriter, err error, log *slog.Logger) {
	if errors.Is(err, context.Canceled) {
		log.Error("Request was canceled by the user")
		http.Error(w, "Request canceled", http.StatusRequestTimeout)
	} else if errors.Is(err, context.DeadlineExceeded) || status.Code(err) == codes.DeadlineExceeded {
		log.Error("Request timed out")
		http.Error(w, "Request timeout", http.StatusRequestTimeout)
	} else {
		log.Error("Operation failed", sl.Err(err))
		http.Error(w, err.Error(), http.StatusInternalServerError)
	}
}

func (ac *ArticleController) GetArticles(w http.ResponseWriter, r *http.Request) {
	const op = "controllers.articlesController.getArticles"
	log := ac.log.With(slog.String("op", op))

	articles, err := ac.articleService.GetArticles(r.Context())
	if err != nil {
		ac.handleError(w, err, log)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	if err := json.NewEncoder(w).Encode(articles); err != nil {
		log.Error("Failed to encode response", sl.Err(err))
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	log.Info("Retrieved all articles successfully")
}

func (ac *ArticleController) GetArticleById(w http.ResponseWriter, r *http.Request) {
	const op = "controllers.articlesController.getArticleById"
	log := ac.log.With(slog.String("op", op))

	idStr := mux.Vars(r)["id"]
	uuidID, err := uuid.Parse(idStr)
	if err != nil {
		log.Error("Invalid UUID format", sl.Err(err))
		http.Error(w, "Invalid UUID", http.StatusBadRequest)
		return
	}

	article, err := ac.articleService.GetArticleById(r.Context(), uuidID)
	if err != nil {
		ac.handleError(w, err, log)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	if err := json.NewEncoder(w).Encode(article); err != nil {
		log.Error("Failed to encode response", sl.Err(err))
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	log.Info("Retrieved article by ID successfully")
}

func (ac *ArticleController) GetArticlesByOwnerId(w http.ResponseWriter, r *http.Request) {
	const op = "controllers.articlesController.getArticlesByOwnerId"
	log := ac.log.With(slog.String("op", op))

	owner_id_s := mux.Vars(r)["owner_id"]
	owner_id, err := uuid.Parse(owner_id_s)
	if err != nil {
		log.Error("Invalid UUID format", sl.Err(err))
		http.Error(w, "Invalid UUID", http.StatusBadRequest)
		return
	}

	article, err := ac.articleService.GetArticleByOwnerId(r.Context(), owner_id)
	if err != nil {
		ac.handleError(w, err, log)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	if err := json.NewEncoder(w).Encode(article); err != nil {
		log.Error("Failed to encode response", sl.Err(err))
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	log.Info("Retrieved article by OwnerID successfully")
}

func (ac *ArticleController) Insert(w http.ResponseWriter, r *http.Request) {
	const op = "controllers.articlesController.getArticlesByOwnerId"
	log := ac.log.With(slog.String("op", op))

	var article models.Article
	if err := json.NewDecoder(r.Body).Decode(&article); err != nil {
		ac.handleError(w, err, log)
		return
	}

	if err := ac.articleService.Insert(r.Context(), article); err != nil {
		ac.handleError(w, err, log)
		return
	}

	w.WriteHeader(http.StatusNoContent)
	log.Info("Inserted article successfully")
}

func (ac *ArticleController) Update(w http.ResponseWriter, r *http.Request) {
	const op = "controllers.articleController.update"
	log := ac.log.With(slog.String("op", op))

	article_id_s := mux.Vars(r)["article_id"]
	article_id, err := uuid.Parse(article_id_s)
	if err != nil {
		log.Error("Invalid UUID format", sl.Err(err))
		http.Error(w, "Invalid UUID", http.StatusBadRequest)
		return
	}

	var article models.Article
	if err := json.NewDecoder(r.Body).Decode(&article); err != nil {
		ac.handleError(w, err, log)
		return
	}

	if err := ac.articleService.Update(r.Context(), article_id, article); err != nil {
		ac.handleError(w, err, log)
		return
	}

	w.WriteHeader(http.StatusNoContent)
	log.Info("Updated article successfully")
}

func (ac *ArticleController) Delete(w http.ResponseWriter, r *http.Request) {
	const op = "controllers.articleController.delete"
	log := ac.log.With(slog.String("op", op))

	article_id_s := mux.Vars(r)["article_id"]
	article_id, err := uuid.Parse(article_id_s)
	if err != nil {
		log.Error("Invalid UUID format", sl.Err(err))
		http.Error(w, "Invalid UUID", http.StatusBadRequest)
		return
	}

	article, err := ac.articleService.Delete(r.Context(), article_id)
	if err != nil {
		ac.handleError(w, err, log)
		return
	}

	w.WriteHeader(http.StatusNoContent)
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	if err := json.NewEncoder(w).Encode(article); err != nil {
		log.Error("Failed to encode response", sl.Err(err))
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	log.Info("Deleted article successfully")
}
