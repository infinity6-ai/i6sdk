//go:build !testcover

package cmdsdkgo

import (
	"context"
	"fmt"
	"os"

	"github.com/spf13/cobra"

	"go.code.infinity6.ai/sdkgo/sdkgo/version"
)

func Prepare(ctx context.Context) *cobra.Command {
	var rootCmd = &cobra.Command{
		Version: version.Version(),
		Use:     "sdkgo",
		Short:   "sdkgo",
	}

	rootCmd.AddCommand(&cobra.Command{
		Use: "version",
		Run: func(cmd *cobra.Command, args []string) {
			fmt.Printf("%s\n", version.Version())
		},
	})

	rootCmd.AddCommand(prepareGenerateCmd())

	return rootCmd
}

func Execute(rootCmd *cobra.Command) {
	if err := rootCmd.Execute(); err != nil {
		fmt.Fprintln(os.Stderr, err)
		os.Exit(1)
	}
}
