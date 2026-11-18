package tabledataz_test

import (
	"errors"
	"testing"
	"time"

	"cloud.google.com/go/bigquery"
	"github.com/stretchr/testify/require"
	"go.code.infinity6.ai/sdkgo/sdkgo/tabledataz"
)

func TestUnitValidateHelpers(t *testing.T) {
	t.Run("FirstError", func(t *testing.T) {
		a, b := errors.New("a"), errors.New("b")
		require.NoError(t, tabledataz.FirstError())
		require.NoError(t, tabledataz.FirstError(nil, nil))
		require.Equal(t, a, tabledataz.FirstError(nil, a, b))
	})

	t.Run("Required", func(t *testing.T) {
		require.NoError(t, tabledataz.Required("id", "x"))
		require.ErrorContains(t, tabledataz.Required("id", ""), "failed to validate id")
		require.NoError(t, tabledataz.RequiredTime("created_at", time.Now()))
		require.ErrorContains(t, tabledataz.RequiredTime("created_at", time.Time{}), "failed to validate created_at")
	})

	t.Run("Str and MaxLen", func(t *testing.T) {
		require.NoError(t, tabledataz.Str("f", "abc", 3))
		require.ErrorContains(t, tabledataz.Str("f", "", 3), "f is required")
		require.ErrorContains(t, tabledataz.Str("f", "abcd", 3), "at most 3 characters")
		require.NoError(t, tabledataz.NullMaxLen("f", bigquery.NullString{}, 1))
		require.ErrorContains(t, tabledataz.NullMaxLen("f", tabledataz.NewBqNullString("ab"), 1), "failed to validate f")
	})

	t.Run("OneOf", func(t *testing.T) {
		require.NoError(t, tabledataz.OneOf("f", "pf", "pf", "pj"))
		require.ErrorContains(t, tabledataz.OneOf("f", "xx", "pf", "pj"), "failed to validate f")
	})

	t.Run("Between and flags", func(t *testing.T) {
		require.NoError(t, tabledataz.Between("f", 2.5, 0.0, 3.0))
		require.ErrorContains(t, tabledataz.Between("f", 4, 0, 3), "f must be <= 3")
		require.ErrorContains(t, tabledataz.Between("f", -1, 0, 3), "f must be >= 0")
		require.NoError(t, tabledataz.Flag("f", 1))
		require.ErrorContains(t, tabledataz.Flag("f", 2), "failed to validate f")
		require.NoError(t, tabledataz.NullFlag("f", bigquery.NullInt64{}))
		require.ErrorContains(t, tabledataz.NullFlag("f", tabledataz.NewBqNullInt64(2)), "failed to validate f")
	})
}
