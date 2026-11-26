package cmdcli

import (
	"context"
	"fmt"
	"os"

	"github.com/infinity6-ai/i6sdk/versionz/version"
	"github.com/spf13/cobra"
)

// Prepare builds the root cobra command for i6 CLI.
func Prepare(ctx context.Context) *cobra.Command {
	var rootCmd = &cobra.Command{
		Version: version.Version(),
		Use:     "i6",
		Short:   "Infinity6 Command Line Interface",
		Long:    "i6 is the official standalone command-line client for Infinity6 platform services, APIs, and workflows.",
	}

	rootCmd.AddCommand(&cobra.Command{
		Use:   "version",
		Short: "Print the current version",
		Run: func(cmd *cobra.Command, args []string) {
			fmt.Println(version.Version())
		},
	})

	return rootCmd
}

// Execute runs the root command with context.
func Execute(ctx context.Context, rootCmd *cobra.Command) {
	if err := rootCmd.ExecuteContext(ctx); err != nil {
		fmt.Fprintln(os.Stderr, err)
		os.Exit(1)
	}
}
