package tokenstore

import (
	"context"
	"database/sql"
	"errors"
	"fmt"
	"os"
	"path/filepath"
	"strings"
	"time"

	"github.com/infinity6-ai/gox/commonz/urlz"
	_ "modernc.org/sqlite"
)

type sqliteStore struct {
	db *sql.DB
}

// DefaultSqliteUrl returns the default persistent token URL: file://$HOME/.config/i6/tokens.db.
func DefaultSqliteUrl() (*urlz.Url, error) {
	homeDir, err := os.UserHomeDir()
	if err != nil {
		homeDir = os.TempDir()
	}
	p := filepath.Join(homeDir, ".config", "i6", "tokens.db")
	u, err := urlz.Parse("file://" + filepath.Clean(p))
	if err != nil {
		return nil, fmt.Errorf("failed to parse default tokens sqlite url: %w", err)
	}
	return u, nil
}

// NewSqlite constructs a persistent SQLite token store with WAL mode and 0600 file permissions.
// If targetUrl is nil, DefaultSqliteUrl() is used.
func NewSqlite(targetUrl *urlz.Url) (Store, error) {
	if targetUrl == nil {
		var err error
		targetUrl, err = DefaultSqliteUrl()
		if err != nil {
			return nil, fmt.Errorf("failed to get default tokens sqlite url: %w", err)
		}
	}

	var dsn string
	var filePath string
	isMemory := false

	if targetUrl.Path != nil && (targetUrl.Path.String() == ":memory:" || targetUrl.Path.String() == "/:memory:" || strings.Contains(targetUrl.Path.String(), "mode=memory")) {
		isMemory = true
		dsn = "file::memory:?mode=memory&cache=shared"
	} else if targetUrl.Path != nil {
		filePath = targetUrl.Path.String()
		dir := filepath.Dir(filePath)
		if err := os.MkdirAll(dir, 0700); err != nil {
			return nil, fmt.Errorf("failed to create sqlite tokens directory %q: %w", dir, err)
		}
		dsn = fmt.Sprintf("file:%s?_pragma=busy_timeout(5000)&_pragma=journal_mode(WAL)&_pragma=synchronous(NORMAL)", filepath.ToSlash(filePath))
	} else {
		return nil, fmt.Errorf("invalid sqlite tokenstore target url: %s", targetUrl)
	}

	db, err := sql.Open("sqlite", dsn)
	if err != nil {
		return nil, fmt.Errorf("failed to open sqlite database at %q: %w", dsn, err)
	}

	if isMemory {
		db.SetMaxOpenConns(1)
	} else {
		db.SetMaxOpenConns(10)
		db.SetMaxIdleConns(5)
	}

	if err := db.Ping(); err != nil {
		_ = db.Close()
		return nil, fmt.Errorf("failed to connect to sqlite token store: %w", err)
	}

	initSql := `
		PRAGMA journal_mode = WAL;
		PRAGMA busy_timeout = 5000;
		PRAGMA synchronous = NORMAL;
		CREATE TABLE IF NOT EXISTS tokens (
			key TEXT PRIMARY KEY,
			value TEXT NOT NULL,
			updated_at INTEGER NOT NULL
		);
	`
	if _, err := db.Exec(initSql); err != nil {
		_ = db.Close()
		return nil, fmt.Errorf("failed to initialize sqlite tokens schema: %w", err)
	}

	if !isMemory && filePath != "" {
		_ = os.Chmod(filePath, 0600)
	}

	return &sqliteStore{db: db}, nil
}

func (s *sqliteStore) Get(ctx context.Context, key string) (string, error) {
	var val string
	err := s.db.QueryRowContext(ctx, "SELECT value FROM tokens WHERE key = ?", key).Scan(&val)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return "", nil
		}
		return "", fmt.Errorf("failed to get key %q from sqlite token store: %w", key, err)
	}
	return val, nil
}

func (s *sqliteStore) Set(ctx context.Context, key string, value string) error {
	now := time.Now().UnixMilli()
	_, err := s.db.ExecContext(ctx, `
		INSERT INTO tokens (key, value, updated_at) VALUES (?, ?, ?)
		ON CONFLICT(key) DO UPDATE SET value = excluded.value, updated_at = excluded.updated_at
	`, key, value, now)
	if err != nil {
		return fmt.Errorf("failed to set key %q in sqlite token store: %w", key, err)
	}
	return nil
}

func (s *sqliteStore) Delete(ctx context.Context, key string) error {
	_, err := s.db.ExecContext(ctx, "DELETE FROM tokens WHERE key = ?", key)
	if err != nil {
		return fmt.Errorf("failed to delete key %q from sqlite token store: %w", key, err)
	}
	return nil
}

func (s *sqliteStore) DeletePrefix(ctx context.Context, prefix string) error {
	_, err := s.db.ExecContext(ctx, "DELETE FROM tokens WHERE substr(key, 1, length(?)) = ?", prefix, prefix)
	if err != nil {
		return fmt.Errorf("failed to delete prefix %q from sqlite token store: %w", prefix, err)
	}
	return nil
}

func (s *sqliteStore) ListPrefix(ctx context.Context, prefix string) (map[string]string, error) {
	rows, err := s.db.QueryContext(ctx, "SELECT key, value FROM tokens WHERE substr(key, 1, length(?)) = ?", prefix, prefix)
	if err != nil {
		return nil, fmt.Errorf("failed to list prefix %q from sqlite token store: %w", prefix, err)
	}
	defer rows.Close()

	res := make(map[string]string)
	for rows.Next() {
		var k, v string
		if err := rows.Scan(&k, &v); err != nil {
			return nil, fmt.Errorf("failed to scan token row for prefix %q: %w", prefix, err)
		}
		res[k] = v
	}
	if err := rows.Err(); err != nil {
		return nil, fmt.Errorf("error reading token rows for prefix %q: %w", prefix, err)
	}
	return res, nil
}

func (s *sqliteStore) Close() error {
	if s.db == nil {
		return nil
	}
	err := s.db.Close()
	s.db = nil
	if err != nil {
		return fmt.Errorf("failed to close sqlite token store: %w", err)
	}
	return nil
}
