package database

import (
	"golang.org/x/sync/singleflight"
	"log/slog"
	"time"

	"leoho.io/singleflight-example/article"
	"leoho.io/singleflight-example/cache"
)

type DB struct {
	cache  *cache.Cache
	engine singleflight.Group
}

func New() *DB {
	return &DB{
		cache: cache.New(),
	}
}

func (db *DB) GetArticle(id string) *article.Article {
	data := db.cache.Get(id)
	if data != nil {
		slog.Info("cache hit...", "id", id, "data", data)
		return data
	}

	slog.Info("cache miss...", "id", id)
	data = &article.Article{
		ID:      id,
		Content: "Hello, world!",
	}
	db.cache.Set(data)

	time.Sleep(100 * time.Millisecond)
	return data
}

func (db *DB) GetArticleDo(id string) *article.Article {
	data := db.cache.Get(id)
	if data != nil {
		slog.Info("cache hit...", "id", id, "data", data)
		return data
	}

	row, err, _ := db.engine.Do(id, func() (interface{}, error) {
		slog.Info("cache miss...", "id", id)
		data = &article.Article{
			ID:      id,
			Content: "Hello, world!",
		}
		db.cache.Set(data)
		return data, nil
	})

	if err != nil {
		slog.Error("singleflight error", "err", err)
		return nil
	}

	return row.(*article.Article)
}

func (db *DB) GetArticleDoChan(id string, timeout time.Duration) *article.Article {
	data := db.cache.Get(id)
	if data != nil {
		slog.Info("cache hit...", "id", id, "data", data)
		return data
	}

	ch := db.engine.DoChan(id, func() (interface{}, error) {
		slog.Info("cache miss...", "id", id)
		data = &article.Article{
			ID:      id,
			Content: "Hello, world!",
		}
		db.cache.Set(data)
		time.Sleep(115 * time.Millisecond)
		return data, nil
	})

	select {
	case <-time.After(timeout):
		slog.Info("timeout", "id", id)
		return nil
	case res := <-ch:
		if res.Err != nil {
			slog.Error("singleflight error", "err", res.Err)
			return nil
		}
		return res.Val.(*article.Article)
	}
}
