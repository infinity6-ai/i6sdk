package cmdtokens

import (
	"context"

	"github.com/infinity6-ai/i6sdk/i6cli/internal/cmd/cmdtokens/cmdcheck"
	"github.com/infinity6-ai/i6sdk/i6cli/internal/cmd/cmdtokens/cmddatasetaccess"
	"github.com/infinity6-ai/i6sdk/i6cli/internal/cmd/cmdtokens/cmddatasetrefresh"
	"github.com/infinity6-ai/i6sdk/i6cli/internal/cmd/cmdtokens/cmdraw"
	"github.com/infinity6-ai/i6sdk/i6cli/internal/cmd/cmdtokens/cmduser"
	"github.com/spf13/cobra"
)

// New constructs the 'tokens' command group.
func New(ctx context.Context) *cobra.Command {
	cmd := &cobra.Command{
		Use:   "tokens",
		Short: "Manage, audit, and diagnose local profile tokens and keys",
		Long:  `Inspect credentials, perform health audits, manage robot keys, and run stateless crypto operations.`,
	}

	cmd.AddCommand(cmdcheck.New(ctx))
	cmd.AddCommand(cmduser.New(ctx))
	cmd.AddCommand(cmddatasetrefresh.New(ctx))
	cmd.AddCommand(cmddatasetaccess.New(ctx))
	cmd.AddCommand(cmdraw.New(ctx))

	return cmd
}
