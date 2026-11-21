package cmdsdkweb_test

import (
	"context"
	"testing"

	"github.com/infinity6-ai/i6sdk/sdkweb/cmdsdkweb"
	"github.com/infinity6-ai/i6sdk/sdkweb/version"
	"github.com/stretchr/testify/require"
)

func TestUnitCmdSdkWebPrepare(t *testing.T) {
	ctx := context.Background()
	cmd := cmdsdkweb.Prepare(ctx)
	require.NotNil(t, cmd)
	require.Equal(t, "sdkweb", cmd.Use)
	require.Equal(t, version.Version(), cmd.Version)

	subcommands := cmd.Commands()
	require.NotEmpty(t, subcommands)

	names := make(map[string]bool)
	for _, c := range subcommands {
		names[c.Name()] = true
	}
	require.True(t, names["version"])
	require.True(t, names["server"])
	require.True(t, names["generate"])
}
