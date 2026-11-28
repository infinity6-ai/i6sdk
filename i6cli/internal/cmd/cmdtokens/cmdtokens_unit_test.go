package cmdtokens_test

import (
	"bytes"
	"testing"

	"github.com/infinity6-ai/i6sdk/i6cli/internal/cmd/cmdtokens"
	"github.com/stretchr/testify/require"
)

func TestUnitCmdTokens(t *testing.T) {
	ctx := t.Context()

	t.Run("Tokens command registers all subcommands", func(t *testing.T) {
		cmd := cmdtokens.New(ctx)
		require.NotNil(t, cmd)
		require.Equal(t, "tokens", cmd.Use)

		subcommands := make(map[string]bool)
		for _, sub := range cmd.Commands() {
			subcommands[sub.Name()] = true
		}

		require.True(t, subcommands["check"])
		require.True(t, subcommands["user"])
		require.True(t, subcommands["dataset-refresh"])
		require.True(t, subcommands["dataset-access"])
		require.True(t, subcommands["raw"])
	})

	t.Run("Tokens help output displays subcommands", func(t *testing.T) {
		cmd := cmdtokens.New(ctx)
		outBuf := new(bytes.Buffer)
		cmd.SetOut(outBuf)
		cmd.SetArgs([]string{"--help"})

		err := cmd.ExecuteContext(ctx)
		require.NoError(t, err)

		out := outBuf.String()
		require.Contains(t, out, "check")
		require.Contains(t, out, "user")
		require.Contains(t, out, "dataset-refresh")
		require.Contains(t, out, "dataset-access")
		require.Contains(t, out, "raw")
	})
}
