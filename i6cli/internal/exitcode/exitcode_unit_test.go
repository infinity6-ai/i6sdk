package exitcode_test

import (
	"errors"
	"fmt"
	"testing"

	"github.com/infinity6-ai/i6sdk/i6cli/internal/exitcode"
	"github.com/stretchr/testify/require"
)

func TestUnitExitcode(t *testing.T) {
	type testScenario struct {
		name     string
		err      error
		wantCode int
		wantMsg  string
	}

	check := func(t *testing.T, s testScenario) {
		t.Helper()
		gotCode := exitcode.Code(s.err)
		require.Equal(t, s.wantCode, gotCode)
		if s.err != nil {
			require.Equal(t, s.wantMsg, s.err.Error())
		}
	}

	t.Run("nil error returns Success code 0", func(t *testing.T) {
		check(t, testScenario{
			name:     "nil error",
			err:      nil,
			wantCode: exitcode.Success,
		})
	})

	t.Run("generic error defaults to General code 1", func(t *testing.T) {
		check(t, testScenario{
			name:     "generic error",
			err:      errors.New("something went wrong"),
			wantCode: exitcode.General,
			wantMsg:  "something went wrong",
		})
	})

	t.Run("New creates ExitError with specific code and message", func(t *testing.T) {
		err := exitcode.New(exitcode.Usage, "missing required flag: --dataset")
		check(t, testScenario{
			name:     "usage error",
			err:      err,
			wantCode: exitcode.Usage,
			wantMsg:  "missing required flag: --dataset",
		})
	})

	t.Run("Wrap encapsulates existing error with ExitError code", func(t *testing.T) {
		baseErr := errors.New("token expired")
		wrapped := exitcode.Wrap(exitcode.Auth, baseErr)
		check(t, testScenario{
			name:     "auth error wrap",
			err:      wrapped,
			wantCode: exitcode.Auth,
			wantMsg:  "token expired",
		})
		require.Equal(t, baseErr, errors.Unwrap(wrapped))
	})

	t.Run("Wrap returns nil if error is nil", func(t *testing.T) {
		wrapped := exitcode.Wrap(exitcode.Forbidden, nil)
		require.Nil(t, wrapped)
		require.Equal(t, exitcode.Success, exitcode.Code(wrapped))
	})

	t.Run("errors.As extracts ExitError through fmt.Errorf wrapping", func(t *testing.T) {
		baseErr := exitcode.New(exitcode.NotFound, "dataset not found")
		wrappedErr := fmt.Errorf("context wrapper: %w", baseErr)
		require.Equal(t, exitcode.NotFound, exitcode.Code(wrappedErr))
	})
}
