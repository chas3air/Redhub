package favoritescontroller

import (
	"apigateway/internal/controllers/favorites/mock"
	"apigateway/internal/domain/models"
	"apigateway/pkg/lib/logger/sl"
	"encoding/json"
	"log/slog"
	"net/http"

	"github.com/google/uuid"
)

type FavoritesController struct {
	log *slog.Logger
	db  *mock.MockFavorites
}

func New(log *slog.Logger) *FavoritesController {
	return &FavoritesController{
		log: log,
		db:  mock.New(log),
	}
}

func (f *FavoritesController) GetByUserId(w http.ResponseWriter, r *http.Request) {
	const op = "favorites.GetByUserId"
	log := f.log.With(
		"op", op,
	)

	id, err := uuid.Parse(r.URL.Query().Get("id"))
	if err != nil {
		log.Error("id must be uuid", sl.Err(err))
		http.Error(w, "id must be uuid", http.StatusBadRequest)
		return
	}

	articles, err := f.db.GetByUserId(r.Context(), id)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	w.WriteHeader(http.StatusOK)
	w.Header().Set("Content-Type", "application/json")
	if err := json.NewEncoder(w).Encode(articles); err != nil {
		log.Error("cannot write to response", sl.Err(err))
		http.Error(w, "cannot write to response", http.StatusInternalServerError)
		return
	}
}

func (f *FavoritesController) Add(w http.ResponseWriter, r *http.Request) {
	const op = "favorites.Add"
	log := f.log.With(
		"op", op,
	)

	id, err := uuid.Parse(r.URL.Query().Get("id"))
	if err != nil {
		log.Error("id must be uuid", sl.Err(err))
		http.Error(w, "id must be uuid", http.StatusBadRequest)
		return
	}

	var article models.Article
	if err := json.NewDecoder(r.Body).Decode(&article); err != nil {
		log.Error("error reading request body", sl.Err(err))
		http.Error(w, "error reading request body", http.StatusBadRequest)
		return
	}

	f.db.Add(r.Context(), id, article)

	w.WriteHeader(http.StatusOK)
}

func (f *FavoritesController) Remove(w http.ResponseWriter, r *http.Request) {
	const op = "favorites.Remove"
	log := f.log.With(
		"op", op,
	)

	info := struct {
		UserId    uuid.UUID `json:"user_id"`
		ArticleId uuid.UUID `json:"article_id"`
	}{}

	if err := json.NewDecoder(r.Body).Decode(&info); err != nil {
		log.Error("error reading request body", sl.Err(err))
		http.Error(w, "error reading request body", http.StatusBadRequest)
		return
	}

	f.db.Remove(r.Context(), info.UserId, info.ArticleId)
	w.WriteHeader(http.StatusOK)
}
