package cache

import (
	"sync"

	"leoho.io/singleflight-example/article"
)

type Cache struct {
	sync.Mutex
	entries map[string]*article.Article
}

func (c *Cache) Get(id string) *article.Article {
	c.Lock()
	defer c.Unlock()
	_, isExist := c.entries[id]
	if !isExist {
		return nil
	}
	return c.entries[id]
}

func (c *Cache) Set(a *article.Article) {
	c.Lock()
	defer c.Unlock()
	c.entries[a.ID] = a
}

func New() *Cache {
	return &Cache{
		entries: make(map[string]*article.Article),
	}
}
