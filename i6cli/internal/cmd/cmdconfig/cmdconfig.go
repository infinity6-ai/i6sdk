package cmdconfig

import (
	"context"

	"github.com/infinity6-ai/i6sdk/i6cli/internal/cmd/cmdconfig/cmdinit"
	"github.com/infinity6-ai/i6sdk/i6cli/internal/cmd/cmdconfig/cmdset"
	"github.com/infinity6-ai/i6sdk/i6cli/internal/cmd/cmdconfig/cmdshow"
	"github.com/spf13/cobra"
)

// New constructs the 'config' command group.
func New(ctx context.Context) *cobra.Command {
	cmd := &cobra.Command{
		Use:   "config",
		Short: "Manage local profile configuration, directories, and defaults",
		Long:  `Manage Infinity6 CLI configuration profiles, target box environments, and defaults.`,
	}

	cmd.AddCommand(cmdinit.New(ctx))
	cmd.AddCommand(cmdshow.New(ctx))
	cmd.AddCommand(cmdset.New(ctx))

	return cmd
}
