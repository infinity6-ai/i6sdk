package cmdcli

import (
	"context"
	"fmt"
	"os"

	"github.com/infinity6-ai/i6sdk/i6cli/internal/cliconfig"
	"github.com/infinity6-ai/i6sdk/i6cli/internal/cmd/cmdauth"
	"github.com/infinity6-ai/i6sdk/i6cli/internal/cmd/cmdconfig"
	"github.com/infinity6-ai/i6sdk/i6cli/internal/exitcode"
	"github.com/infinity6-ai/i6sdk/versionz/version"
	"github.com/spf13/cobra"
)

// Prepare builds the root cobra command for i6 CLI.
func Prepare(ctx context.Context) *cobra.Command {
	var rootCmd = &cobra.Command{
		Version: version.Version(),
		Use:     "i6",
		Short:   "Infinity6 Command Line Interface",
		Long: `i6 is the official standalone command-line client for Infinity6 platform services, APIs, and workflows.

ENVIRONMENT VARIABLES:
  I6_CLI_CONFIG_DIR        Base directory for CLI configuration and profiles (default: $HOME/.config/i6)
  I6_CLI_CONFIG_PROFILE    Active profile name when --profile is not provided (default: "main")
  I6_CLI_CONFIG_JSON       Inline JSON configuration for zero-disk execution (CI/CD runners)
  I6_DATASET               Default dataset identifier when --dataset flag is omitted
  I6_LOGZ_NO_COLOR         Disable ANSI terminal colors in logs (true|false, auto-detected)
  I6_LOGZ_STRUCTURED       Output structured JSON logs to stderr (true|false)`,
		SilenceUsage:  true,
		SilenceErrors: true,
	}

	// Persistent root flags
	defaultProfile := cliconfig.ResolveProfile("", ctx)
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

	// Profile configuration command group
	rootCmd.AddCommand(cmdconfig.New(ctx))

	// Authentication and session command group
	rootCmd.AddCommand(cmdauth.New(ctx))

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
