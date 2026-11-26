package sqlitez

import (
	"context"
	"database/sql"
	"encoding/json"
	"fmt"
	"path/filepath"
	"strings"

	_ "modernc.org/sqlite"
)

// Conn represents an open SQLite database connection with named query capabilities.
type Conn struct {
	db *sql.DB
}

// Open opens a SQLite connection in read-only mode to the file at path.
func Open(path string) (*Conn, error) {
	absPath, err := filepath.Abs(path)
	if err != nil {
		return nil, fmt.Errorf("failed to resolve absolute sqlite path %q: %w", path, err)
	}

	// URI format for read-only access
	dsn := fmt.Sprintf("file:%s?mode=ro", filepath.ToSlash(absPath))
	db, err := sql.Open("sqlite", dsn)
	if err != nil {
		return nil, fmt.Errorf("failed to open sqlite database at %q: %w", path, err)
	}

	if err := db.Ping(); err != nil {
		_ = db.Close()
		return nil, fmt.Errorf("failed to connect to sqlite database at %q: %w", path, err)
	}

	return &Conn{db: db}, nil
}

// OpenMemory opens an in-memory SQLite connection (useful for tests or ephemeral runs).
func OpenMemory() (*Conn, error) {
	db, err := sql.Open("sqlite", ":memory:")
	if err != nil {
		return nil, fmt.Errorf("failed to open in-memory sqlite database: %w", err)
	}
	if err := db.Ping(); err != nil {
		_ = db.Close()
		return nil, fmt.Errorf("failed to connect to in-memory sqlite database: %w", err)
	}
	return &Conn{db: db}, nil
}

// DB returns the underlying standard library *sql.DB.
func (c *Conn) DB() *sql.DB {
	return c.db
}

// Close closes the underlying SQLite database connection.
func (c *Conn) Close() error {
	if c.db == nil {
		return nil
	}
	err := c.db.Close()
	c.db = nil
	if err != nil {
		return fmt.Errorf("failed to close sqlite connection: %w", err)
	}
	return nil
}

// Exec executes a query without returning any rows (useful for setup and testing).
func (c *Conn) Exec(ctx context.Context, query string, args ...any) (sql.Result, error) {
	if c.db == nil {
		return nil, fmt.Errorf("sqlitez: connection is closed")
	}
	res, err := c.db.ExecContext(ctx, query, args...)
	if err != nil {
		return nil, fmt.Errorf("sql execution failed: %w", err)
	}
	return res, nil
}

// Query executes a query and returns each row as a map of column name to value.
func (c *Conn) Query(ctx context.Context, query string, params ...any) ([]map[string]any, error) {
	if c.db == nil {
		return nil, fmt.Errorf("sqlitez: connection is closed")
	}

	rows, err := c.db.QueryContext(ctx, query, params...)
	if err != nil {
		return nil, fmt.Errorf("sql query failed: %w", err)
	}
	defer rows.Close()

	cols, err := rows.Columns()
	if err != nil {
		return nil, fmt.Errorf("failed to get columns: %w", err)
	}

	var results []map[string]any

	for rows.Next() {
		colVals := make([]any, len(cols))
		colPointers := make([]any, len(cols))
		for i := range colVals {
			colPointers[i] = &colVals[i]
		}

		if err := rows.Scan(colPointers...); err != nil {
			return nil, fmt.Errorf("failed to scan row: %w", err)
		}

		rowMap := make(map[string]any, len(cols))
		for i, colName := range cols {
			val := colVals[i]
			// Normalize []byte to string for cleaner text handling
			if b, ok := val.([]byte); ok {
				val = string(b)
			}
			rowMap[colName] = val
		}
		results = append(results, rowMap)
	}

	if err := rows.Err(); err != nil {
		return nil, fmt.Errorf("row iteration error: %w", err)
	}

	return results, nil
}

// NamedQuery executes a named query stored in the namedqueries table.
func (c *Conn) NamedQuery(ctx context.Context, name string, params ...any) ([]map[string]any, error) {
	findRows, err := c.Query(ctx, "SELECT query FROM namedqueries WHERE name = ?", name)
	if err != nil {
		return nil, fmt.Errorf("failed to lookup named query %q: %w", name, err)
	}
	if len(findRows) == 0 {
		return nil, fmt.Errorf("sqlitez: no named query %q found in namedqueries table", name)
	}

	rawQuery, ok := findRows[0]["query"]
	if !ok || rawQuery == nil {
		return nil, fmt.Errorf("sqlitez: named query %q has nil or missing query column", name)
	}

	sqlText, ok := rawQuery.(string)
	if !ok {
		sqlText = fmt.Sprintf("%v", rawQuery)
	}

	return c.Query(ctx, sqlText, params...)
}

// NamedStringQuery binds JSON-serialized params to positional placeholders (?1, ?2, ...)
// and returns the string value of the 'ret' column from the first row.
func (c *Conn) NamedStringQuery(ctx context.Context, name string, params ...any) (*string, error) {
	boundParams := make([]any, len(params))
	for i, p := range params {
		if p == nil {
			boundParams[i] = "null"
			continue
		}
		if s, ok := p.(string); ok {
			boundParams[i] = s
			continue
		}
		b, err := json.Marshal(p)
		if err != nil {
			return nil, fmt.Errorf("failed to serialize parameter %d for query %q to JSON: %w", i, name, err)
		}
		boundParams[i] = string(b)
	}

	rows, err := c.NamedQuery(ctx, name, boundParams...)
	if err != nil {
		return nil, err
	}
	if len(rows) == 0 {
		return nil, nil
	}

	row := rows[0]
	// Look up 'ret' column case-insensitively
	var retVal any
	found := false
	for k, v := range row {
		if strings.EqualFold(k, "ret") {
			retVal = v
			found = true
			break
		}
	}
	if !found {
		return nil, fmt.Errorf("sqlitez: named query %q returned no 'ret' column in output", name)
	}
	if retVal == nil {
		return nil, nil
	}

	retStr := fmt.Sprintf("%v", retVal)
	return &retStr, nil
}

// NamedJsonQuery executes NamedStringQuery and returns raw JSON bytes of the 'ret' column.
func (c *Conn) NamedJsonQuery(ctx context.Context, name string, params ...any) ([]byte, error) {
	str, err := c.NamedStringQuery(ctx, name, params...)
	if err != nil {
		return nil, err
	}
	if str == nil {
		return nil, nil
	}
	return []byte(*str), nil
}

// NamedJsonQueryInto executes NamedJsonQuery and unmarshals the result into target.
func (c *Conn) NamedJsonQueryInto(ctx context.Context, target any, name string, params ...any) error {
	raw, err := c.NamedJsonQuery(ctx, name, params...)
	if err != nil {
		return err
	}
	if len(raw) == 0 {
		return nil
	}
	if err := json.Unmarshal(raw, target); err != nil {
		return fmt.Errorf("failed to unmarshal JSON query result for %q: %w\nRaw:\n%s", name, err, string(raw))
	}
	return nil
}
