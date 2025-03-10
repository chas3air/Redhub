package app

import (
	articlecontroller "apigateway/internal/controllers/articleController"
	authcontroller "apigateway/internal/controllers/auth"
	commentsmanagercontroller "apigateway/internal/controllers/commentController"
	userscontroller "apigateway/internal/controllers/usersManager"
	articlemanageservice "apigateway/internal/services/articleManager"
	authservice "apigateway/internal/services/auth"
	commentsmanagerservice "apigateway/internal/services/comments"
	usersmanagerservice "apigateway/internal/services/usersManager"
	articlesmanagerstorage "apigateway/internal/storage/real/articlesManager"
	authstorage "apigateway/internal/storage/real/auth"
	commentsmanagerstorage "apigateway/internal/storage/real/comments"
	usersmanagerstorage "apigateway/internal/storage/real/usersManager"
	"apigateway/pkg/config"
	"fmt"
	"log/slog"
	"net/http"

	"github.com/gorilla/mux"
)

type App struct {
	log *slog.Logger
	cfg *config.Config
}

func New(log *slog.Logger, cfg *config.Config) *App {
	return &App{
		log: log,
		cfg: cfg,
	}
}

func (a *App) Start() {
	authstorage := authstorage.New(a.log, a.cfg.AuthHost, a.cfg.AuthPort)
	auth_service := authservice.New(a.log, authstorage)
	authcontroller := authcontroller.New(a.log, auth_service)

	usersmanagerstorage := usersmanagerstorage.New(a.log, a.cfg.UsersStorageHost, a.cfg.UsersStoragePort)
	usersmanager_service := usersmanagerservice.New(a.log, usersmanagerstorage)
	userscontroller := userscontroller.New(a.log, usersmanager_service)

	articlemanagerstrorage := articlesmanagerstorage.New(a.log, a.cfg.ArticlesStorageHost, a.cfg.ArticlesStoragePort)
	articlemanager_service := articlemanageservice.New(a.log, articlemanagerstrorage)
	articlecontroller := articlecontroller.New(a.log, articlemanager_service)

	commentmanagerstorage := commentsmanagerstorage.New(a.log, a.cfg.CommentsStorageHost, a.cfg.CommentsStoragePort)
	commentsmanagerservice := commentsmanagerservice.New(a.log, commentmanagerstorage)
	commentsmanagercontroller := commentsmanagercontroller.New(a.log, commentsmanagerservice)

	r := mux.NewRouter()
	r.HandleFunc("/api/v1/health-check", func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusOK)
	})

	r.HandleFunc("/api/v1/login", authcontroller.Login).Methods(http.MethodPost)
	r.HandleFunc("/api/v1/register", authcontroller.Register).Methods(http.MethodPost)

	r.HandleFunc("/api/v1/users", userscontroller.GetUsers).Methods(http.MethodGet)
	r.HandleFunc("/api/v1/users/{id}", userscontroller.GetUserById).Methods(http.MethodGet)
	r.HandleFunc("/api/v1/users", userscontroller.Insert).Methods(http.MethodPost)
	r.HandleFunc("/api/v1/users/{id}", userscontroller.Update).Methods(http.MethodPut)
	r.HandleFunc("/api/v1/users/{id}", userscontroller.Delete).Methods(http.MethodDelete)

	r.HandleFunc("/api/v1/articles", articlecontroller.GetArticles).Methods(http.MethodGet)
	r.HandleFunc("/api/v1/articles/{article_id}/", articlecontroller.GetArticleById).Methods(http.MethodGet)
	r.HandleFunc("/api/v1/articles/{owner_id}", articlecontroller.GetArticlesByOwnerId).Methods(http.MethodGet)
	r.HandleFunc("/api/v1/articles", articlecontroller.Insert).Methods(http.MethodPost)
	r.HandleFunc("/api/v1/articles/{article_id}", articlecontroller.Update).Methods(http.MethodPut)
	r.HandleFunc("/api/v1/articles/{article_id}", articlecontroller.Delete).Methods(http.MethodDelete)

	r.HandleFunc("/api/v1/comments/{id}", commentsmanagercontroller.GetCommentById).Methods(http.MethodGet)
	r.HandleFunc("/api/v1/{article_id}/comments", commentsmanagercontroller.GetCommentsByArticleId).Methods(http.MethodGet)
	r.HandleFunc("/api/v1/comments", commentsmanagercontroller.Insert).Methods(http.MethodPost)
	r.HandleFunc("/api/v1/comments/{id}", commentsmanagercontroller.Delete).Methods(http.MethodDelete)

	if err := http.ListenAndServe(fmt.Sprintf(":%d", a.cfg.API.Port), r); err != nil {
		panic(err)
	}

}
