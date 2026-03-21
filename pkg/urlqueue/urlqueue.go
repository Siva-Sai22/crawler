// Package urlqueue implements a queue of urls.
package urlqueue

import (
	"context"

	"github.com/Siva-Sai22/crawler/pkg/fetcher"
	"github.com/Siva-Sai22/crawler/pkg/queue"
)

type URLQueue struct {
	*queue.Queue[string]
	visited map[string]bool
}

func New() *URLQueue {
	return &URLQueue{queue.New[string](), make(map[string]bool)}
}

func (q *URLQueue) InsertLinks(links []string) {
	for _, link := range links {
		if q.visited[link] {
			continue
		}
		q.Enqueue(link)
	}
}

func (q *URLQueue) Dequeue(ctx context.Context, repo *fetcher.WebsiteRepository) (bool, error) {
	url, ok := q.Queue.Dequeue()
	if ok {
		q.visited[url] = true
		outboundLinks, err := fetcher.Fetch(ctx, url, repo)
		if err != nil {
			return true, err
		}

		if len(outboundLinks) > 0 {
			q.InsertLinks(outboundLinks)
		}
	}
	return ok, nil
}
