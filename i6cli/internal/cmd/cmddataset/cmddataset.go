package cmddataset

import (
	"context"

	"github.com/infinity6-ai/i6sdk/i6cli/internal/cmd/cmddataset/cmdget"
	"github.com/infinity6-ai/i6sdk/i6cli/internal/cmd/cmddataset/cmdkeys"
	"github.com/infinity6-ai/i6sdk/i6cli/internal/cmd/cmddataset/cmdlist"
	"github.com/infinity6-ai/i6sdk/i6cli/internal/cmd/cmddataset/cmdstore"
	"github.com/infinity6-ai/i6sdk/i6cli/internal/cmd/cmddataset/cmdtable"
	"github.com/spf13/cobra"
)

// New constructs the 'dataset' command group.
func New(ctx context.Context) *cobra.Command {
	cmd := &cobra.Command{
		Use:   "dataset",
		Short: "Manage datasets, tables, stores, and remote machine robot keys",
		Long:  `Manage datasets, tables, stores, and remote machine robot keys on the Infinity6 platform.`,
	}

	cmd.AddCommand(cmdlist.New(ctx))
	cmd.AddCommand(cmdget.New(ctx))
	cmd.AddCommand(cmdtable.New(ctx))
	cmd.AddCommand(cmdstore.New(ctx))
	cmd.AddCommand(cmdkeys.New(ctx))

	return cmd
}
