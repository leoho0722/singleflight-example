package main

import (
	"log/slog"
	"sync"
	"time"

	"github.com/google/uuid"
	"leoho.io/singleflight-example/database"
)

func main() {
	wg := sync.WaitGroup{}

	db := database.New()

	wg.Add(5)
	for i := 0; i < 5; i++ {
		go func(t int) {
			defer wg.Done()
			//data := db.GetArticle(uuid.New().String())
			//data := db.GetArticleDo(uuid.New().String())
			data := db.GetArticleDoChan(uuid.NewString(), time.Duration(t*100)*time.Millisecond)
			slog.Info("data info", "data", data)
		}(i)
	}
	wg.Wait()
}
