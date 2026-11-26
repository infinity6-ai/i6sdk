package cmdcli_test

import (
	"bytes"
	"testing"

	"github.com/infinity6-ai/i6sdk/i6cli/cmdcli"
	"github.com/stretchr/testify/require"
)

func TestUnitCmdcli(t *testing.T) {
	t.Run("Prepare constructs root command with version", func(t *testing.T) {
		cmd := cmdcli.Prepare(t.Context())
		require.NotNil(t, cmd)
		require.Equal(t, "i6", cmd.Use)
		require.NotEmpty(t, cmd.Version)
	})

	t.Run("Version command executes without error", func(t *testing.T) {
		cmd := cmdcli.Prepare(t.Context())
		buf := new(bytes.Buffer)
		cmd.SetOut(buf)
		cmd.SetArgs([]string{"version"})

		err := cmd.ExecuteContext(t.Context())
		require.NoError(t, err)
	})
}
