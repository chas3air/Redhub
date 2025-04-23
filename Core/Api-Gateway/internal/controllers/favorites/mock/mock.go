package mock

import (
	"apigateway/internal/domain/models"
	"context"
	"log/slog"

	"github.com/google/uuid"
)

type MockFavorites struct {
	log       *slog.Logger
	favorites map[uuid.UUID][]models.Article
}

func New(log *slog.Logger) *MockFavorites {
	return &MockFavorites{
		log:       log,
		favorites: make(map[uuid.UUID][]models.Article),
	}
}

func (m *MockFavorites) GetByUserId(ctx context.Context, id uuid.UUID) ([]models.Article, error) {
	return m.favorites[id], nil
}

func (m *MockFavorites) Add(ctx context.Context, id uuid.UUID, article models.Article) error {
	if _, exists := m.favorites[id]; !exists {
		m.favorites[id] = []models.Article{}
	}

	for _, existingArticle := range m.favorites[id] {
		if existingArticle.Id == article.Id {
			return nil
		}
	}

	m.favorites[id] = append(m.favorites[id], article)
	return nil
}

func (m *MockFavorites) Remove(ctx context.Context, id uuid.UUID, articleId uuid.UUID) error {
	if articles, exists := m.favorites[id]; exists {
		for i, article := range articles {
			if article.Id == articleId {
				m.favorites[id] = append(articles[:i], articles[i+1:]...)
				break
			}
		}
	}

	return nil
}
