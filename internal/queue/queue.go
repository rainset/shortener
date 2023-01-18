// Пакет для очередей удаления ссылок
package queue

import (
	"fmt"
	"github.com/rainset/shortener/internal/storage"
	"sync"
)

type Task struct {
	CookieID string
	Hashes   []string
}

type DeleteURLQueue struct {
	mu   sync.Mutex
	ch   chan *Task
	urls []string
	s    storage.InterfaceStorage
}

func NewDeleteURLQueue(storage storage.InterfaceStorage) *DeleteURLQueue {
	return &DeleteURLQueue{
		ch:   make(chan *Task, 1),
		urls: make([]string, 0),
		s:    storage,
	}
}

func (q *DeleteURLQueue) Push(t *Task) {
	q.ch <- t
}

func (q *DeleteURLQueue) PopWait() *Task {
	return <-q.ch
}

func (q *DeleteURLQueue) PeriodicURLDelete() {
	var err error
	for {
		//time.Sleep(5 * time.Second)

		if len(q.urls) == 0 {
			continue
		}
		q.mu.Lock()
		urls := q.urls
		q.urls = nil
		q.mu.Unlock()

		err = q.s.DeleteBatchURL(urls)
		if err != nil {
			fmt.Printf("PeriodicURLDelete Loop() error: %v\n", err)
			continue
		}

	}
}

type DeleteURLWorker struct {
	id    int
	queue *DeleteURLQueue
	s     storage.InterfaceStorage
}

func NewDeleteURLWorker(id int, queue *DeleteURLQueue, s storage.InterfaceStorage) *DeleteURLWorker {
	w := DeleteURLWorker{
		id:    id,
		queue: queue,
		s:     s,
	}
	return &w
}

func (w *DeleteURLWorker) Loop() {
	var err error
	for {
		t := w.queue.PopWait()
		w.queue.mu.Lock()
		w.queue.urls = append(w.queue.urls, t.Hashes...)
		if len(w.queue.urls) > 1 {
			err = w.s.DeleteUserBatchURL(t.CookieID, w.queue.urls)
			if err != nil {
				fmt.Printf("DeleteURLWorker Loop() error: %v\n", err)
				continue
			}
			fmt.Printf("worker #%d delete %s\n", w.id, w.queue.urls)
			w.queue.urls = nil
		} else {
			fmt.Printf("worker #%d add %s\n", w.id, w.queue.urls)
		}

		w.queue.mu.Unlock()
	}
}
