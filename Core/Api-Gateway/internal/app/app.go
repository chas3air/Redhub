package app

import (
	"apigateway/internal/controllers/agreement"
	articlecontroller "apigateway/internal/controllers/articleController"
	authcontroller "apigateway/internal/controllers/auth"
	commentsmanagercontroller "apigateway/internal/controllers/commentController"
	favoritescontroller "apigateway/internal/controllers/favorites"
	"apigateway/internal/controllers/middleware"
	statscontroller "apigateway/internal/controllers/statsController"
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
	// Пачка для микросервиса авторизации
	authStorage := authstorage.New(a.log, a.cfg.AuthHost, a.cfg.AuthPort)
	authService := authservice.New(a.log, authStorage)
	authController := authcontroller.New(a.log, authService)

	// Пачка для микросервиса пользователей
	usersManagerStorage := usersmanagerstorage.New(a.log, a.cfg.UsersStorageHost, a.cfg.UsersStoragePort)
	usersManagerService := usersmanagerservice.New(a.log, usersManagerStorage)
	usersController := userscontroller.New(a.log, usersManagerService)

	// Пачка для микросервиса постов
	articleManagerStorage := articlesmanagerstorage.New(a.log, a.cfg.ArticlesStorageHost, a.cfg.ArticlesStoragePort)
	articleManagerService := articlemanageservice.New(a.log, articleManagerStorage)
	articleController := articlecontroller.New(a.log, articleManagerService)

	// Пачка для микросервиса комментариев
	commentManagerStorage := commentsmanagerstorage.New(a.log, a.cfg.CommentsStorageHost, a.cfg.CommentsStoragePort)
	commentsManagerService := commentsmanagerservice.New(a.log, commentManagerStorage)
	commentsManagerController := commentsmanagercontroller.New(a.log, commentsManagerService)

	// Контроллер для статы
	statsController := statscontroller.New(a.log, articleManagerService, usersManagerService, commentsManagerService)

	// Контроллер для модерации
	moderationController := agreement.New(a.log)

	// Избранное
	favoritesController := favoritescontroller.New(a.log)

	// Создание объекта middleware
	middleware := middleware.New()

	r := mux.NewRouter()
	r.Use(middleware.CORS)
	r.HandleFunc("/api/v1/health-check", func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusOK)
	})

	// Группа для авторизации, не пропускает если пользователь уже существует
	authRouter := r.PathPrefix("/api/v1").Subrouter()
	authRouter.Use(middleware.CORS)
	authRouter.Use(middleware.PreventAccessIfLoggedIn)
	authRouter.HandleFunc("/login", authController.Login).Methods(http.MethodPost, http.MethodOptions)
	authRouter.HandleFunc("/register", authController.Register).Methods(http.MethodPost, http.MethodOptions)

	// Группа для работы с пользователями
	route_for_user_admin := r.PathPrefix("/api/v1/users").Subrouter()
	route_for_user_admin.Use(middleware.ValidateToken)
	route_for_user_admin.Use(middleware.RequireUserAdmin)

	r.HandleFunc("/api/v1/users", usersController.GetUsers).Methods(http.MethodGet, http.MethodOptions)
	r.HandleFunc("/api/v1/users/{id}", usersController.GetUserById).Methods(http.MethodGet, http.MethodOptions)
	route_for_user_admin.HandleFunc("", usersController.Insert).Methods(http.MethodPost, http.MethodOptions)
	route_for_user_admin.HandleFunc("/{id}", usersController.Update).Methods(http.MethodPut, http.MethodOptions)
	route_for_user_admin.HandleFunc("/{id}", usersController.Delete).Methods(http.MethodDelete, http.MethodOptions)

	// Группа для работы с постами и комментариями
	route_for_article_admin := r.PathPrefix("/api/v1").Subrouter()
	route_for_article_admin.Use(middleware.ValidateToken)
	route_for_article_admin.Use(middleware.RequireArticleAdmin)

	route_for_user := r.PathPrefix("/api/v1").Subrouter()
	route_for_user.Use(middleware.ValidateToken)
	route_for_user.Use(middleware.RequireUser)

	r.HandleFunc("/api/v1/articles", articleController.GetArticles).Methods(http.MethodGet, http.MethodOptions)
	r.HandleFunc("/api/v1/articles/{article_id}/", articleController.GetArticleById).Methods(http.MethodGet, http.MethodOptions)
	r.HandleFunc("/api/v1/articles/{owner_id}", articleController.GetArticlesByOwnerId).Methods(http.MethodGet, http.MethodOptions)
	route_for_user.HandleFunc("/articles", articleController.Insert).Methods(http.MethodPost, http.MethodOptions)
	route_for_article_admin.HandleFunc("/articles/{article_id}", articleController.Update).Methods(http.MethodPut, http.MethodOptions)
	route_for_article_admin.HandleFunc("/articles/{article_id}", articleController.Delete).Methods(http.MethodDelete, http.MethodOptions)

	r.HandleFunc("/api/v1/comments/{id}", commentsManagerController.GetCommentById).Methods(http.MethodGet, http.MethodOptions)
	r.HandleFunc("/api/v1/{article_id}/comments", commentsManagerController.GetCommentsByArticleId).Methods(http.MethodGet, http.MethodOptions)
	route_for_user.HandleFunc("/comments", commentsManagerController.Insert).Methods(http.MethodPost, http.MethodOptions)
	route_for_article_admin.HandleFunc("/comments/{id}", commentsManagerController.Delete).Methods(http.MethodDelete, http.MethodOptions)

	route_for_analyst := r.PathPrefix("/api/v1/stats").Subrouter()
	route_for_analyst.Use(middleware.ValidateToken)
	route_for_analyst.HandleFunc("/articles", statsController.GetArticlesStats).Methods(http.MethodGet, http.MethodOptions)
	route_for_analyst.HandleFunc("/users", statsController.GetUsersStats).Methods(http.MethodGet, http.MethodOptions)

	// route_for_moderation := r.PathPrefix("/api/v1/moderation").Subrouter()
	// route_for_moderation.Use(middleware.RequireModerator)
	r.HandleFunc("/api/v1/moderation/get", moderationController.GetArticles).Methods(http.MethodGet, http.MethodOptions)
	r.HandleFunc("/api/v1/moderation/add", moderationController.AddArticle).Methods(http.MethodPost, http.MethodOptions)
	r.HandleFunc("/api/v1/moderation/remove", moderationController.RemoveArticle).Methods(http.MethodDelete, http.MethodOptions)

	route_for_favorites := r.PathPrefix("/api/v1/favorites").Subrouter()
	route_for_favorites.Use(middleware.ValidateToken)
	route_for_favorites.Use(middleware.RequireUser)
	route_for_favorites.HandleFunc("/get", favoritesController.GetByUserId).Methods(http.MethodGet, http.MethodOptions)
	route_for_favorites.HandleFunc("/add", favoritesController.Add).Methods(http.MethodPost, http.MethodOptions)
	route_for_favorites.HandleFunc("/delete", favoritesController.Remove).Methods(http.MethodDelete, http.MethodOptions)

	if err := http.ListenAndServe(fmt.Sprintf(":%d", a.cfg.API.Port), r); err != nil {
		panic(err)
	}
}
