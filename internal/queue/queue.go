// Пакет для очередей удаления ссылок
package queue

import (
	"fmt"
	"github.com/rainset/shortener/internal/storage"
	"sync"
)

type Task struct {
	UserID string
	Hashes []string
}

type DeleterQueue struct {
	mx    sync.Mutex
	s     storage.InterfaceStorage
	ch    chan *Task
	tasks []Task
}

func NewDeleterQueue(storage storage.InterfaceStorage) *DeleterQueue {
	return &DeleterQueue{
		s:  storage,
		ch: make(chan *Task),
	}
}

func (d *DeleterQueue) Init() error {
	// пометка удаляемых ссылок deleted=1
	for {
		t := <-d.ch
		//log.Println(t)
		err := d.s.DeleteUserBatchURL(t.UserID, t.Hashes)
		if err != nil {
			fmt.Printf("DeleteUserBatchURL Loop() error: %v\n", err)
			continue
		}
	}

	// удаление по времени
	//go func() {
	//	for now := range time.Tick(time.Second * 5) {
	//		log.Println("time.Tick", now)
	//		d.Exec()
	//
	//	}
	//}()
}

func (d *DeleterQueue) PopWait() *Task {
	return <-d.ch
}

func (d *DeleterQueue) Push(t *Task) {
	d.ch <- t
}
