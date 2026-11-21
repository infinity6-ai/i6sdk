package cmdsdk

import (
	"context"
	"fmt"
	"os"

	"github.com/infinity6-ai/i6sdk/cmdz/cmdsdkjs"
	"github.com/infinity6-ai/i6sdk/versionz/version"
	"github.com/spf13/cobra"
)

func Prepare(ctx context.Context) *cobra.Command {
	var rootCmd = &cobra.Command{
		Version: version.Version(),
		Use:     "i6sdk",
		Short:   "i6sdk CLI",
	}

	rootCmd.AddCommand(&cobra.Command{
		Use: "version",
		Run: func(cmd *cobra.Command, args []string) {
			fmt.Printf("%s\n", version.Version())
		},
	})

	cmdsdkjs.Prepare(rootCmd)

	return rootCmd
}

func Execute(rootCmd *cobra.Command) {
	if err := rootCmd.Execute(); err != nil {
		fmt.Fprintln(os.Stderr, err)
		os.Exit(1)
	}
}
