package lowbot

import (
	"sync"
)

type Broadcast[T any] struct {
	mu        sync.RWMutex
	listeners []chan T
}

func NewBroadcast[T any]() *Broadcast[T] {
	return &Broadcast[T]{
		listeners: []chan T{},
	}
}

func (b *Broadcast[T]) Send(v T) {
	b.mu.RLock()
	defer b.mu.RUnlock()

	for _, listener := range b.listeners {
		select {
		case listener <- v:
		default:
		}
	}
}

func (b *Broadcast[T]) Listen() chan T {
	listener := make(chan T)
	b.listeners = append(b.listeners, listener)
	return listener
}

func (b *Broadcast[T]) Close() error {
	b.mu.RLock()
	defer b.mu.RUnlock()

	for _, listener := range b.listeners {
		close(listener)
	}

	clear(b.listeners)

	return nil
}
