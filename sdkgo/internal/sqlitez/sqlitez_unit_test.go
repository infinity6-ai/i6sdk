package sqlitez_test

import (
	"context"
	"database/sql"
	"os"
	"path/filepath"
	"testing"

	"github.com/infinity6-ai/i6sdk/sdkgo/internal/sqlitez"
	"github.com/stretchr/testify/require"
)

func TestUnitSqliteOpenMemory(t *testing.T) {
	ctx := context.Background()
	conn, err := sqlitez.OpenMemory()
	require.NoError(t, err)
	defer conn.Close()

	t.Run("Exec and basic query", func(t *testing.T) {
		_, err := conn.Exec(ctx, "CREATE TABLE items (id INTEGER PRIMARY KEY, name TEXT, price REAL)")
		require.NoError(t, err)

		_, err = conn.Exec(ctx, "INSERT INTO items (name, price) VALUES (?, ?)", "Widget", 19.99)
		require.NoError(t, err)

		rows, err := conn.Query(ctx, "SELECT name, price FROM items WHERE id = ?", 1)
		require.NoError(t, err)
		require.Len(t, rows, 1)
		require.Equal(t, "Widget", rows[0]["name"])
		require.Equal(t, 19.99, rows[0]["price"])
	})
}

func TestUnitSqliteNamedQuery(t *testing.T) {
	ctx := context.Background()
	conn, err := sqlitez.OpenMemory()
	require.NoError(t, err)
	defer conn.Close()

	// Setup namedqueries table
	_, err = conn.Exec(ctx, "CREATE TABLE namedqueries (name TEXT PRIMARY KEY, query TEXT NOT NULL)")
	require.NoError(t, err)

	_, err = conn.Exec(ctx, `
		INSERT INTO namedqueries (name, query) VALUES 
		('get-user', 'SELECT json_object("id", ?1, "name", "Alice", "active", 1) AS ret'),
		('calc-metrics', 'SELECT json_object("clicks", json_extract(?1, "$.clicks") * 2, "source", json_extract(?1, "$.source")) AS ret'),
		('plain-text', 'SELECT "just plain string" AS ret'),
		('no-ret-col', 'SELECT 123 AS other_col')
	`)
	require.NoError(t, err)

	t.Run("NamedQuery executes stored query", func(t *testing.T) {
		rows, err := conn.NamedQuery(ctx, "get-user", "usr_100")
		require.NoError(t, err)
		require.Len(t, rows, 1)
		require.Contains(t, rows[0]["ret"], "Alice")
	})

	t.Run("NamedStringQuery returns ret column", func(t *testing.T) {
		str, err := conn.NamedStringQuery(ctx, "plain-text")
		require.NoError(t, err)
		require.NotNil(t, str)
		require.Equal(t, "just plain string", *str)
	})

	t.Run("NamedJsonQuery binds JSON parameters and unmarshals result", func(t *testing.T) {
		type MetricsOutput struct {
			Clicks int    `json:"clicks"`
			Source string `json:"source"`
		}

		params := map[string]any{
			"clicks": 21,
			"source": "organic",
		}

		var out MetricsOutput
		err := conn.NamedJsonQueryInto(ctx, &out, "calc-metrics", params)
		require.NoError(t, err)
		require.Equal(t, 42, out.Clicks)
		require.Equal(t, "organic", out.Source)
	})

	t.Run("Missing named query returns actionable error", func(t *testing.T) {
		_, err := conn.NamedQuery(ctx, "nonexistent-query")
		require.Error(t, err)
		require.Contains(t, err.Error(), "no named query \"nonexistent-query\"")
	})

	t.Run("Missing ret column returns actionable error", func(t *testing.T) {
		_, err := conn.NamedStringQuery(ctx, "no-ret-col")
		require.Error(t, err)
		require.Contains(t, err.Error(), "returned no 'ret' column")
	})
}

func TestUnitSqliteFileReadOnly(t *testing.T) {
	ctx := context.Background()
	tmpDir, err := os.MkdirTemp("", "sqlitez-file-test-*")
	require.NoError(t, err)
	defer os.RemoveAll(tmpDir)

	dbPath := filepath.Join(tmpDir, "test.db")

	// Create SQLite file using standard sql.Open to initialize tables and data
	setupDb, err := sql.Open("sqlite", dbPath)
	require.NoError(t, err)
	_, err = setupDb.Exec("CREATE TABLE namedqueries (name TEXT PRIMARY KEY, query TEXT);")
	require.NoError(t, err)
	_, err = setupDb.Exec("INSERT INTO namedqueries (name, query) VALUES ('test-q', 'SELECT json_object(\"status\", \"ok\") AS ret');")
	require.NoError(t, err)
	setupDb.Close()

	// Open read-only via sqlitez.Open
	conn, err := sqlitez.Open(dbPath)
	require.NoError(t, err)
	defer conn.Close()

	t.Run("Reads query successfully in read-only mode", func(t *testing.T) {
		str, err := conn.NamedStringQuery(ctx, "test-q")
		require.NoError(t, err)
		require.NotNil(t, str)
		require.Contains(t, *str, `"status":"ok"`)
	})

	t.Run("Write fails in read-only mode", func(t *testing.T) {
		_, err := conn.Exec(ctx, "INSERT INTO namedqueries (name, query) VALUES ('new-q', 'SELECT 1');")
		require.Error(t, err)
		require.Contains(t, err.Error(), "readonly")
	})
}
