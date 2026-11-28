package cmdcli

import (
	"context"
	"fmt"
	"os"

	"github.com/infinity6-ai/i6sdk/i6cli/internal/cliconfig"
	"github.com/infinity6-ai/i6sdk/i6cli/internal/exitcode"
	"github.com/infinity6-ai/i6sdk/versionz/version"
	"github.com/spf13/cobra"
)

// Prepare builds the root cobra command for i6 CLI.
func Prepare(ctx context.Context) *cobra.Command {
	var rootCmd = &cobra.Command{
		Version:       version.Version(),
		Use:           "i6",
		Short:         "Infinity6 Command Line Interface",
		Long:          "i6 is the official standalone command-line client for Infinity6 platform services, APIs, and workflows.",
		SilenceUsage:  true,
		SilenceErrors: true,
	}

	// Persistent root flags
	defaultProfile := cliconfig.ResolveProfile("")
	rootCmd.PersistentFlags().StringP("profile", "p", defaultProfile, "Configuration profile name (or I6_CLI_CONFIG_PROFILE)")
	rootCmd.PersistentFlags().String("format", "", "Output format (table|tsv|json|quiet|raw)")

	// Version command
	rootCmd.AddCommand(&cobra.Command{
		Use:   "version",
		Short: "Print the current version",
		Run: func(cmd *cobra.Command, args []string) {
			fmt.Fprintln(cmd.OutOrStdout(), version.Version())
		},
	})

	// Shell completion command
	rootCmd.AddCommand(newCompletionCmd())

	return rootCmd
}

// ExecuteWithCode runs the command and returns the deterministic exit code without calling os.Exit.
func ExecuteWithCode(ctx context.Context, rootCmd *cobra.Command) int {
	if err := rootCmd.ExecuteContext(ctx); err != nil {
		code := exitcode.Code(err)
		fmt.Fprintln(rootCmd.ErrOrStderr(), err)
		return code
	}
	return exitcode.Success
}

// Execute runs the root command with context and exits with the determined exit code.
func Execute(ctx context.Context, rootCmd *cobra.Command) {
	code := ExecuteWithCode(ctx, rootCmd)
	if code != exitcode.Success {
		os.Exit(code)
	}
}
