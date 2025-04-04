package statscontroller

import (
	"apigateway/internal/domain/interfaces/articles"
	"apigateway/internal/domain/interfaces/comments"
	"apigateway/internal/domain/interfaces/users"
	"apigateway/pkg/lib/logger/sl"
	"encoding/json"
	"log/slog"
	"net/http"
	"time"
)

type StatsController struct {
	log            *slog.Logger
	articleService articles.IArticlesService
	usersService   users.UsersService
	commentService comments.CommentsService
}

func New(log *slog.Logger,
	articleService articles.IArticlesService,
	usersService users.UsersService,
	commentService comments.CommentsService,
) *StatsController {
	return &StatsController{
		log:            log,
		articleService: articleService,
		usersService:   usersService,
		commentService: commentService,
	}
}

// количество статей
// количество статей пользователя отсортировано по убыванию
func (sc *StatsController) GetArticlesStats(w http.ResponseWriter, r *http.Request) {
	const op = "controller.statsController.getArticlesStats"
	log := sc.log.With(
		"op", op,
	)

	res_struct := struct {
		CountOfArticles int `json:"count_of_articles,omitempty"`
	}{}

	articles, err := sc.articleService.GetArticles(r.Context())
	if err != nil {
		// дописать
		return
	}

	res_struct.CountOfArticles = len(articles)

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	if err := json.NewEncoder(w).Encode(res_struct); err != nil {
		log.Error("Failed to encode response", sl.Err(err))
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
}

// количество пользователей
// вывод возрастных групп (0-18, 19-25, 26-45, 46-...)
func (sc *StatsController) GetUsersStats(w http.ResponseWriter, r *http.Request) {
	const op = "controller.statsController.getUsersStats"
	log := sc.log.With(
		"op", op,
	)

	users, err := sc.usersService.GetUsers(r.Context())
	if err != nil {
		// дописать
		return
	}

	res_struct := struct {
		CountOfUsers int    `json:"count_of_users,omitempty"`
		ArrayOfAges  [4]int `json:"array_of_ages,omitempty"`
	}{}

	res_struct.CountOfUsers = len(users)

	for _, user := range users {
		buf_age := int(time.Since(user.Birthday).Hours() / 24 / 365.25)
		if buf_age > 46 {
			res_struct.ArrayOfAges[3]++
		} else if buf_age > 26 {
			res_struct.ArrayOfAges[2]++
		} else if buf_age > 19 {
			res_struct.ArrayOfAges[1]++
		} else {
			res_struct.ArrayOfAges[0]++
		}
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	if err := json.NewEncoder(w).Encode(res_struct); err != nil {
		log.Error("Failed to encode response", sl.Err(err))
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	log.Info("Retrieved usersStats")
}
