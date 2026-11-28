package cmdauth

import (
	"context"

	"github.com/infinity6-ai/i6sdk/i6cli/internal/cmd/cmdauth/cmdlogin"
	"github.com/infinity6-ai/i6sdk/i6cli/internal/cmd/cmdauth/cmdlogout"
	"github.com/infinity6-ai/i6sdk/i6cli/internal/cmd/cmdauth/cmdstatus"
	"github.com/infinity6-ai/i6sdk/i6cli/internal/cmd/cmdauth/cmdtoken"
	"github.com/spf13/cobra"
)

// New constructs the 'auth' command group.
func New(ctx context.Context) *cobra.Command {
	cmd := &cobra.Command{
		Use:   "auth",
		Short: "Manage authentication, device login, tokens, and sessions",
		Long:  `Authenticate, inspect identities, universally ingest tokens, and manage sessions.`,
	}

	cmd.AddCommand(cmdlogin.New(ctx))
	cmd.AddCommand(cmdtoken.New(ctx))
	cmd.AddCommand(cmdstatus.New(ctx))
	cmd.AddCommand(cmdlogout.New(ctx))

	return cmd
}
