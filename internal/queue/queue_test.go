package queue

import (
	"github.com/rainset/shortener/internal/storage"
	"reflect"
	"sync"
	"testing"
)

func TestDeleteURLQueue_PeriodicURLDelete(t *testing.T) {
	type fields struct {
		mu   sync.Mutex
		ch   chan *Task
		urls []string
		s    storage.InterfaceStorage
	}
	tests := []struct {
		name   string
		fields fields
	}{
		// TODO: Add test cases.
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			q := &DeleteURLQueue{
				mu:   tt.fields.mu,
				ch:   tt.fields.ch,
				urls: tt.fields.urls,
				s:    tt.fields.s,
			}
			q.PeriodicURLDelete()
		})
	}
}

func TestDeleteURLQueue_PopWait(t *testing.T) {
	type fields struct {
		mu   sync.Mutex
		ch   chan *Task
		urls []string
		s    storage.InterfaceStorage
	}
	tests := []struct {
		name   string
		fields fields
		want   *Task
	}{
		// TODO: Add test cases.
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			q := &DeleteURLQueue{
				mu:   tt.fields.mu,
				ch:   tt.fields.ch,
				urls: tt.fields.urls,
				s:    tt.fields.s,
			}
			if got := q.PopWait(); !reflect.DeepEqual(got, tt.want) {
				t.Errorf("PopWait() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestDeleteURLQueue_Push(t *testing.T) {
	type fields struct {
		mu   sync.Mutex
		ch   chan *Task
		urls []string
		s    storage.InterfaceStorage
	}
	type args struct {
		t *Task
	}
	tests := []struct {
		name   string
		fields fields
		args   args
	}{
		// TODO: Add test cases.
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			q := &DeleteURLQueue{
				mu:   tt.fields.mu,
				ch:   tt.fields.ch,
				urls: tt.fields.urls,
				s:    tt.fields.s,
			}
			q.Push(tt.args.t)
		})
	}
}

func TestDeleteURLWorker_Loop(t *testing.T) {
	type fields struct {
		id    int
		queue *DeleteURLQueue
		s     storage.InterfaceStorage
	}
	tests := []struct {
		name   string
		fields fields
	}{
		// TODO: Add test cases.
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			w := &DeleteURLWorker{
				id:    tt.fields.id,
				queue: tt.fields.queue,
				s:     tt.fields.s,
			}
			w.Loop()
		})
	}
}

func TestNewDeleteURLQueue(t *testing.T) {
	type args struct {
		storage storage.InterfaceStorage
	}
	tests := []struct {
		name string
		args args
		want *DeleteURLQueue
	}{
		// TODO: Add test cases.
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := NewDeleteURLQueue(tt.args.storage); !reflect.DeepEqual(got, tt.want) {
				t.Errorf("NewDeleteURLQueue() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestNewDeleteURLWorker(t *testing.T) {
	type args struct {
		id    int
		queue *DeleteURLQueue
		s     storage.InterfaceStorage
	}
	tests := []struct {
		name string
		args args
		want *DeleteURLWorker
	}{
		// TODO: Add test cases.
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := NewDeleteURLWorker(tt.args.id, tt.args.queue, tt.args.s); !reflect.DeepEqual(got, tt.want) {
				t.Errorf("NewDeleteURLWorker() = %v, want %v", got, tt.want)
			}
		})
	}
}
