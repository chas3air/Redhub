package agreement

import (
	"apigateway/internal/domain/models"
	"apigateway/pkg/lib/logger/sl"
	"encoding/json"
	"log/slog"
	"net/http"

	"github.com/google/uuid"
)

type AgreementHandler struct {
	log      *slog.Logger
	articles []models.Article
}

func New(log *slog.Logger) *AgreementHandler {
	return &AgreementHandler{
		log:      log,
		articles: make([]models.Article, 0, 10),
	}
}

func (a *AgreementHandler) GetArticles(w http.ResponseWriter, r *http.Request) {
	json.NewEncoder(w).Encode(a.articles)
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
}

func (a *AgreementHandler) AddArticle(w http.ResponseWriter, r *http.Request) {
	const op = "aggrement.AddArticle"
	log := a.log.With(
		"op", op,
	)

	var article models.Article
	if err := json.NewDecoder(r.Body).Decode(&article); err != nil {
		log.Error("cannot read request body", sl.Err(err))
		http.Error(w, "cannot read request body", http.StatusBadRequest)
		return
	}

	a.articles = append(a.articles, article)
	w.WriteHeader(http.StatusOK)
}

func (a *AgreementHandler) RemoveArticle(w http.ResponseWriter, r *http.Request) {
	const op = "agreement.RemoveArticle"
	log := a.log.With(
		"op", op,
	)

	id, err := uuid.Parse(r.URL.Query().Get("id"))
	if err != nil {
		log.Error("invalid UUID format", sl.Err(err))
		http.Error(w, "invalid UUID format", http.StatusBadRequest)
		return
	}

	for i, v := range a.articles {
		if v.Id == id {
			a.articles = append(a.articles[:i], a.articles[i+1:]...)
			w.WriteHeader(http.StatusOK)
			return
		}
	}

	log.Warn("article not found", "id", id)
	http.Error(w, "article not found", http.StatusNotFound)
}
