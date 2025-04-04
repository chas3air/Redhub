package statscontroller

import (
	"apigateway/internal/domain/interfaces/articles"
	"apigateway/internal/domain/interfaces/comments"
	"apigateway/internal/domain/interfaces/users"
	"apigateway/pkg/lib/logger/sl"
	"context"
	"encoding/json"
	"errors"
	"log/slog"
	"net/http"
	"sort"
	"time"

	"github.com/google/uuid"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
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

func (sc *StatsController) handleError(w http.ResponseWriter, err error, log *slog.Logger) {
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

// количество статей
// количество статей пользователя отсортировано по убыванию
func (sc *StatsController) GetArticlesStats(w http.ResponseWriter, r *http.Request) {
	const op = "controller.statsController.getArticlesStats"
	log := sc.log.With(
		"op", op,
	)

	type OwnerArticle struct {
		OwnerId         uuid.UUID `json:"owner_id,omitempty"`
		CountOfArticles int       `json:"count_of_articles,omitempty"`
	}

	res_struct := struct {
		CountOfArticles int            `json:"count_of_articles,omitempty"`
		OwnerArticles   []OwnerArticle `json:"owner_articles,omitempty"`
	}{}

	articles, err := sc.articleService.GetArticles(r.Context())
	if err != nil {
		sc.handleError(w, err, log)
		return
	}

	ownerStats := make(map[uuid.UUID]int)
	for _, article := range articles {
		ownerStats[article.OwnerId]++
	}

	for owner_id, countOfArticles := range ownerStats {
		res_struct.OwnerArticles = append(
			res_struct.OwnerArticles,
			OwnerArticle{
				OwnerId:         owner_id,
				CountOfArticles: countOfArticles,
			})
	}

	sort.Slice(res_struct.OwnerArticles, func(i, j int) bool {
		return res_struct.OwnerArticles[i].CountOfArticles > res_struct.OwnerArticles[j].CountOfArticles
	})

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
		sc.handleError(w, err, log)
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
