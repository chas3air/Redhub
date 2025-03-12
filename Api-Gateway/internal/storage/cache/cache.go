package cache

import (
	"apigateway/internal/storage"
	"sync"

	"github.com/google/uuid"
)

var (
	g_singletonCache *Cache
	once             sync.Once
)

type Cache struct {
	items map[uuid.UUID]string
}

func New() *Cache {
	once.Do(func() {
		g_singletonCache = &Cache{
			items: make(map[uuid.UUID]string),
		}
	})
	return g_singletonCache
}

func (c *Cache) Add(user_id uuid.UUID, refreshToken string) {
	c.items[user_id] = refreshToken
}

func (c *Cache) Get(user_id uuid.UUID) (string, error) {
	refresh_token, ok := c.items[user_id]
	if !ok {
		return "", storage.ErrNotFound
	}
	return refresh_token, nil
}

func (c *Cache) Delete(user_id uuid.UUID) {
	delete(c.items, user_id)
}
