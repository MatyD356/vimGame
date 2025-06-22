package cache

import (
	"sync"
	"time"
)

type PageCache struct {
	ChildDatabaseId string               `json:"child_database_id"`
	Title           string               `json:"title"`
	CreatedTime     time.Time            `json:"created_time"`
	ChildDatabase   []ChildDatabaseCache `json:"child_database,omitempty"`
}

type ChildDatabaseCache struct {
	ID    string `json:"id"`
	Title string `json:"title"`
	Value string `json:"value"`
}

type Cache struct {
	page                map[string]PageCache
	childDatabase       map[string]ChildDatabaseCache
	parsedChildDatabase map[string][]ChildDatabaseCache

	mu sync.RWMutex
}

func NewCache() *Cache {
	return &Cache{
		page:                make(map[string]PageCache),
		childDatabase:       make(map[string]ChildDatabaseCache),
		parsedChildDatabase: make(map[string][]ChildDatabaseCache),
	}
}

func (c *Cache) SetPage(key string, value PageCache) {
	if value.ChildDatabaseId == "" {
		return
	}
	c.mu.Lock()
	defer c.mu.Unlock()
	c.page[key] = value
}

func (c *Cache) GetPage(key string) (PageCache, bool) {
	c.mu.RLock()
	defer c.mu.RUnlock()

	value, exists := c.page[key]
	if !exists || value.ChildDatabaseId == "" {
		return PageCache{}, false
	}
	return value, true
}

func (c *Cache) SetChildDatabase(key string, value ChildDatabaseCache) {
	if value.ID == "" {
		return
	}
	c.mu.Lock()
	defer c.mu.Unlock()
	c.childDatabase[key] = value
}

func (c *Cache) GetChildDatabase(key string) (ChildDatabaseCache, bool) {
	c.mu.RLock()
	defer c.mu.RUnlock()

	value, exists := c.childDatabase[key]
	if !exists || value.ID == "" {
		return ChildDatabaseCache{}, false
	}
	return value, true
}

func (c *Cache) SetParsedChildDatabase(key string, value []ChildDatabaseCache) {
	if len(value) == 0 {
		return
	}
	c.mu.Lock()
	defer c.mu.Unlock()
	c.parsedChildDatabase[key] = value
}

func (c *Cache) GetParsedChildDatabase(key string) ([]ChildDatabaseCache, bool) {
	c.mu.RLock()
	defer c.mu.RUnlock()

	value, exists := c.parsedChildDatabase[key]
	if !exists || len(value) == 0 {
		return nil, false
	}
	return value, true
}
