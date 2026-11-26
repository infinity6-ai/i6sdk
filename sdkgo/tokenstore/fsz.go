package tokenstore

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"os"
	"path/filepath"
	"sync"

	"github.com/infinity6-ai/gox/commonz/urlz"
	"github.com/infinity6-ai/gox/fsz/fsz"
)

type fszStore struct {
	mu  sync.RWMutex
	url *urlz.Url
}

// DefaultFszUrl returns the default persistent token URL: file://$HOME/.config/i6/tokens.json.
func DefaultFszUrl() (*urlz.Url, error) {
	homeDir, err := os.UserHomeDir()
	if err != nil {
		homeDir = os.TempDir()
	}
	p := filepath.Join(homeDir, ".config", "i6", "tokens.json")
	u, err := urlz.Parse("file://" + filepath.Clean(p))
	if err != nil {
		return nil, fmt.Errorf("failed to parse default tokens fsz url: %w", err)
	}
	return u, nil
}

// NewFsz constructs a persistent token store backed by gox/fsz.
// If targetUrl is nil, DefaultFszUrl() is used.
func NewFsz(targetUrl *urlz.Url) (Store, error) {
	if targetUrl == nil {
		var err error
		targetUrl, err = DefaultFszUrl()
		if err != nil {
			return nil, fmt.Errorf("failed to get default tokens fsz url: %w", err)
		}
	}
	return &fszStore{
		url: targetUrl,
	}, nil
}

func (s *fszStore) readData(ctx context.Context) (map[string]string, error) {
	var data map[string]string
	err := fsz.Download(ctx, s.url, func(found bool, _ http.Header, reader io.Reader) error {
		if !found || reader == nil {
			data = make(map[string]string)
			return nil
		}
		raw, err := io.ReadAll(reader)
		if err != nil {
			return fmt.Errorf("failed to read token store: %w", err)
		}
		if len(bytes.TrimSpace(raw)) == 0 {
			data = make(map[string]string)
			return nil
		}
		if err := json.Unmarshal(raw, &data); err != nil {
			// On corrupted content, initialize empty map to allow recovery
			data = make(map[string]string)
			return nil
		}
		return nil
	})
	if err != nil {
		return nil, fmt.Errorf("failed to download tokens from %s: %w", s.url, err)
	}
	if data == nil {
		data = make(map[string]string)
	}
	return data, nil
}

func (s *fszStore) writeData(ctx context.Context, data map[string]string) error {
	raw, err := json.MarshalIndent(data, "", "  ")
	if err != nil {
		return fmt.Errorf("failed to marshal tokens data: %w", err)
	}

	tmpUrl, err := urlz.Parse(s.url.String() + ".tmp")
	if err != nil {
		return fmt.Errorf("failed to create tmp tokens url: %w", err)
	}

	if err := fsz.Upload(ctx, tmpUrl, nil, bytes.NewReader(raw)); err != nil {
		return fmt.Errorf("failed to upload tokens to tmp file %s: %w", tmpUrl, err)
	}

	if err := fsz.Move(ctx, tmpUrl, s.url); err != nil {
		return fmt.Errorf("failed to atomically move tokens file to %s: %w", s.url, err)
	}

	// For local filesystem paths, secure file permissions to 0600
	if s.url.Scheme == "file" && s.url.Path != nil {
		_ = os.Chmod(s.url.Path.String(), 0600)
	}

	return nil
}

func (s *fszStore) Get(ctx context.Context, key string) (string, error) {
	s.mu.RLock()
	defer s.mu.RUnlock()

	data, err := s.readData(ctx)
	if err != nil {
		return "", fmt.Errorf("failed to get key %q from fsz store: %w", key, err)
	}
	return data[key], nil
}

func (s *fszStore) Set(ctx context.Context, key string, value string) error {
	s.mu.Lock()
	defer s.mu.Unlock()

	data, err := s.readData(ctx)
	if err != nil {
		return fmt.Errorf("failed to read data before set for key %q: %w", key, err)
	}
	data[key] = value
	if err := s.writeData(ctx, data); err != nil {
		return fmt.Errorf("failed to write data for key %q: %w", key, err)
	}
	return nil
}

func (s *fszStore) Delete(ctx context.Context, key string) error {
	s.mu.Lock()
	defer s.mu.Unlock()

	data, err := s.readData(ctx)
	if err != nil {
		return fmt.Errorf("failed to read data before delete for key %q: %w", key, err)
	}
	if _, exists := data[key]; !exists {
		return nil
	}
	delete(data, key)
	if err := s.writeData(ctx, data); err != nil {
		return fmt.Errorf("failed to write data after deleting key %q: %w", key, err)
	}
	return nil
}
