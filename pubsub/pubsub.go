package pubsub

import (
	"fmt"
	"sync"
)

//type Event struct{}

type Server[T any] struct {
	subscribers map[chan T]bool
	mu          sync.Mutex
}

func NewServer[T any]() *Server[T] {
	subs := make(map[chan T]bool)
	return &Server[T]{subscribers: subs, mu: sync.Mutex{}}

}

func (s *Server[T]) Subscribe(channel chan T) {
	s.mu.Lock()
	defer s.mu.Unlock()
	if s.subscribers[channel] {
		fmt.Errorf("Already subscribed")
	}

	s.subscribers[channel] = true
}

func (s *Server[T]) Publish(value T) {
	s.mu.Lock()
	defer s.mu.Unlock()

	for ch := range s.subscribers {
		ch <- value
	}
}

func (s *Server[T]) Cancel(channel chan T) {
	s.mu.Lock()
	defer s.mu.Unlock()

	if !s.subscribers[channel] {
		fmt.Errorf("this channel is not  subscribed")
	}

	close(channel)
	delete(s.subscribers, channel)
}
