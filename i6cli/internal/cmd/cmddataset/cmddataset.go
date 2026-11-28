package cmddataset

import (
	"context"

	"github.com/infinity6-ai/i6sdk/i6cli/internal/cmd/cmddataset/cmdkeys"
	"github.com/spf13/cobra"
)

// New constructs the 'dataset' command group.
func New(ctx context.Context) *cobra.Command {
	cmd := &cobra.Command{
		Use:   "dataset",
		Short: "Manage datasets, tables, stores, and remote machine robot keys",
		Long:  `Manage datasets, tables, stores, and remote machine robot keys on the Infinity6 platform.`,
	}

	cmd.AddCommand(cmdkeys.New(ctx))

	return cmd
}
