package streams

import (
	"sync"

	"github.com/imiller31/servicebus-fanout/protos"
)

type ResponseManager struct {
	mu    sync.RWMutex
	store map[string]chan *protos.ClientRequest
}

func NewResponseManager() *ResponseManager {
	return &ResponseManager{
		store: make(map[string]chan *protos.ClientRequest),
	}
}

func (s *ResponseManager) Get(key string) (chan *protos.ClientRequest, bool) {
	s.mu.RLock()
	defer s.mu.RUnlock()
	ch, exists := s.store[key]
	return ch, exists
}

func (s *ResponseManager) Set(key string, ch chan *protos.ClientRequest) {
	s.mu.Lock()
	defer s.mu.Unlock()
	s.store[key] = ch
}

func (s *ResponseManager) Delete(key string) {
	s.mu.Lock()
	defer s.mu.Unlock()
	delete(s.store, key)
}
