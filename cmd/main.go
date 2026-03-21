package main

import (
	"context"
	"log"
	"time"

	"github.com/Siva-Sai22/crawler/pkg/fetcher"
	"github.com/Siva-Sai22/crawler/pkg/urlqueue"
	"github.com/jackc/pgx/v5/pgxpool"
)

func main() {
	postgresURL := "postgres://admin:postgres@localhost:5432/pages?sslmode=disable"
	db, err := pgxpool.New(context.Background(), postgresURL)
	if err != nil {
		log.Fatal(err)
	}
	defer db.Close()

	repo := fetcher.NewWebsiteRepository(db)

	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	startingURL := "https://en.wikipedia.org/wiki/Interstellar_(film)"
	que := urlqueue.New()
	que.Enqueue(startingURL)

	for {
		ok, err := que.Dequeue(ctx, repo)
		if !ok {
			break
		}
		if err != nil {
			if err == fetcher.ErrDuplicate {
				continue
			}
			log.Fatal(err)
		}
	}
}
