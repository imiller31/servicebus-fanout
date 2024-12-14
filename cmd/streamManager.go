package cmd

import (
	"sync"

	"github.com/imiller31/servicebus-fanout/protos"
	"google.golang.org/grpc"
)

type streamType string
type clientName string

type StreamManager struct {
	mu sync.Mutex

	store map[streamType]map[clientName]grpc.ServerStreamingServer[protos.Response]
}

func NewStreamManager() *StreamManager {
	return &StreamManager{
		store: make(map[streamType]map[clientName]grpc.ServerStreamingServer[protos.Response]),
	}
}

func (s *StreamManager) GetRandomStream(strType streamType) (grpc.ServerStreamingServer[protos.Response], bool) {
	s.mu.Lock()
	defer s.mu.Unlock()

	if _, ok := s.store[strType]; !ok {
		return nil, false
	}

	for client := range s.store[strType] {
		return s.store[strType][client], true
	}

	return nil, false
}

func (s *StreamManager) GetStream(strType streamType, client clientName) (grpc.ServerStreamingServer[protos.Response], bool) {
	s.mu.Lock()
	defer s.mu.Unlock()

	if _, ok := s.store[strType]; !ok {
		return nil, false
	}

	if _, ok := s.store[strType][client]; !ok {
		return nil, false
	}

	return s.store[strType][client], true
}

func (s *StreamManager) SetStream(strType streamType, client clientName, stream grpc.ServerStreamingServer[protos.Response]) {
	s.mu.Lock()
	defer s.mu.Unlock()

	if _, ok := s.store[strType]; !ok {
		s.store[strType] = make(map[clientName]grpc.ServerStreamingServer[protos.Response])
	}

	s.store[strType][client] = stream
}

func (s *StreamManager) DeleteStream(strType streamType, client clientName) {
	s.mu.Lock()
	defer s.mu.Unlock()

	if _, ok := s.store[strType]; !ok {
		return
	}

	delete(s.store[strType], client)
	if len(s.store[strType]) == 0 {
		delete(s.store, strType)
	}
}
