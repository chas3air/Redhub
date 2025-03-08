package mock

import (
	"commentsManageService/internal/domain/models"
	storage_error "commentsManageService/internal/storage"
	"commentsManageService/pkg/lib/logger/sl"
	"context"
	"errors"
	"log/slog"
	"sync"

	"github.com/google/uuid"
)

type MemoryStorage struct {
	log      *slog.Logger
	comments map[uuid.UUID]models.Comment
	mu       sync.RWMutex
}

func New(log *slog.Logger) *MemoryStorage {
	return &MemoryStorage{
		log:      log,
		comments: make(map[uuid.UUID]models.Comment),
	}
}

func (m *MemoryStorage) GetCommentById(ctx context.Context, cid uuid.UUID) (models.Comment, error) {
	const op = "storage.memory.getCommentById"
	log := m.log.With(slog.String("op", op))

	select {
	case <-ctx.Done():
		return models.Comment{}, ctx.Err()
	default:
	}

	m.mu.RLock()
	defer m.mu.RUnlock()

	comment, exists := m.comments[cid]
	if !exists {
		log.Error("Comment with current id not found", sl.Err(errors.New("not found")))
		return models.Comment{}, storage_error.ErrNotFound
	}

	return comment, nil
}

func (m *MemoryStorage) GetCommentsByArticleId(ctx context.Context, aid uuid.UUID) ([]models.Comment, error) {
	const op = "storage.memory.getCommentsByArticleId"

	select {
	case <-ctx.Done():
		return nil, ctx.Err()
	default:
	}

	m.mu.RLock()
	defer m.mu.RUnlock()

	var comments []models.Comment
	for _, comment := range m.comments {
		if comment.ArticleId == aid {
			comments = append(comments, comment)
		}
	}

	return comments, nil
}

func (m *MemoryStorage) Insert(ctx context.Context, comment models.Comment) (models.Comment, error) {
	const op = "storage.memory.insert"
	log := m.log.With(slog.String("op", op))

	select {
	case <-ctx.Done():
		return models.Comment{}, ctx.Err()
	default:
	}

	m.mu.Lock()
	defer m.mu.Unlock()

	if _, exists := m.comments[comment.Id]; exists {
		log.Error("Comment with this ID already exists", sl.Err(errors.New("already exists")))
		return comment, storage_error.ErrAlreadyExists
	}

	m.comments[comment.Id] = comment
	return comment, nil
}

func (m *MemoryStorage) Delete(ctx context.Context, cid uuid.UUID) (models.Comment, error) {
	const op = "storage.memory.delete"
	log := m.log.With(slog.String("op", op))

	select {
	case <-ctx.Done():
		return models.Comment{}, ctx.Err()
	default:
	}

	m.mu.Lock()
	defer m.mu.Unlock()

	comment, exists := m.comments[cid]
	if !exists {
		log.Warn("Comment not found", sl.Err(errors.New("not found")))
		return models.Comment{}, storage_error.ErrNotFound
	}

	delete(m.comments, cid)
	return comment, nil
}
