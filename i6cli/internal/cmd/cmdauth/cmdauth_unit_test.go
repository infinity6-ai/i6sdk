package cmdauth_test

import (
	"bytes"
	"testing"

	"github.com/infinity6-ai/i6sdk/i6cli/internal/cmd/cmdauth"
	"github.com/stretchr/testify/require"
)

func TestUnitCmdAuth(t *testing.T) {
	ctx := t.Context()

	t.Run("Auth command registers all subcommands", func(t *testing.T) {
		cmd := cmdauth.New(ctx)
		require.NotNil(t, cmd)
		require.Equal(t, "auth", cmd.Use)

		subcommands := make(map[string]bool)
		for _, sub := range cmd.Commands() {
			subcommands[sub.Name()] = true
		}

		require.True(t, subcommands["login"])
		require.True(t, subcommands["token"])
		require.True(t, subcommands["status"])
		require.True(t, subcommands["inspect"])
		require.True(t, subcommands["logout"])
	})

	t.Run("Auth help output displays available commands", func(t *testing.T) {
		cmd := cmdauth.New(ctx)
		outBuf := new(bytes.Buffer)
		cmd.SetOut(outBuf)
		cmd.SetArgs([]string{"--help"})

		err := cmd.ExecuteContext(ctx)
		require.NoError(t, err)

		out := outBuf.String()
		require.Contains(t, out, "login")
		require.Contains(t, out, "token")
		require.Contains(t, out, "status")
		require.Contains(t, out, "inspect")
		require.Contains(t, out, "logout")
	})
}
