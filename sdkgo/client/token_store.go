package client

import (
	"context"
	"sync"
)

// TokenStore defines the persistence interface for authentication tokens.
type TokenStore interface {
	Get(ctx context.Context, key string) (string, error)
	Set(ctx context.Context, key string, value string) error
	Delete(ctx context.Context, key string) error
}

type memoryTokenStore struct {
	mu   sync.RWMutex
	data map[string]string
}

func newMemoryTokenStore() *memoryTokenStore {
	return &memoryTokenStore{
		data: make(map[string]string),
	}
}

func (m *memoryTokenStore) Get(_ context.Context, key string) (string, error) {
	m.mu.RLock()
	defer m.mu.RUnlock()
	return m.data[key], nil
}

func (m *memoryTokenStore) Set(_ context.Context, key string, value string) error {
	m.mu.Lock()
	defer m.mu.Unlock()
	m.data[key] = value
	return nil
}

func (m *memoryTokenStore) Delete(_ context.Context, key string) error {
	m.mu.Lock()
	defer m.mu.Unlock()
	delete(m.data, key)
	return nil
}
