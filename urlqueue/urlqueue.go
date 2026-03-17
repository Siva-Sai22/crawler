// Package urlquEue implements a queue of urls.
package urlqueue

import (
	"github.com/Siva-Sai22/crawler/queue"
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
		q.visited[link] = true
		q.Enqueue(link)
	}
}

func (q *URLQueue) Dequeue() (string, bool) {
	url, ok := q.Queue.Dequeue()
	return url, ok
}
