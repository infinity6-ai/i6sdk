package tokenstore

import (
	"context"
	"sync"
)

// Store defines the persistence interface for authentication tokens.
type Store interface {
	Get(ctx context.Context, key string) (string, error)
	Set(ctx context.Context, key string, value string) error
	Delete(ctx context.Context, key string) error
}

type memoryStore struct {
	mu   sync.RWMutex
	data map[string]string
}

// NewMemory returns a thread-safe in-memory token store.
func NewMemory() Store {
	return &memoryStore{
		data: make(map[string]string),
	}
}

func (m *memoryStore) Get(_ context.Context, key string) (string, error) {
	m.mu.RLock()
	defer m.mu.RUnlock()
	return m.data[key], nil
}

func (m *memoryStore) Set(_ context.Context, key string, value string) error {
	m.mu.Lock()
	defer m.mu.Unlock()
	m.data[key] = value
	return nil
}

func (m *memoryStore) Delete(_ context.Context, key string) error {
	m.mu.Lock()
	defer m.mu.Unlock()
	delete(m.data, key)
	return nil
}
