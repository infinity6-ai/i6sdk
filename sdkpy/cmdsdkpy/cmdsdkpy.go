package cmdsdkpy

import (
	"context"
	"fmt"
	"os"
	"os/signal"
	"syscall"

	"github.com/infinity6-ai/i6sdk/sdkpy/version"
	"github.com/spf13/cobra"
)

func Prepare(ctx context.Context) *cobra.Command {
	var rootCmd = &cobra.Command{
		Version: version.Version(),
		Use:     "sdkpy",
		Short:   "i6 sdk python",
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
	ctx, stop := signal.NotifyContext(context.Background(), os.Interrupt, syscall.SIGTERM)
	defer stop()
	if err := rootCmd.ExecuteContext(ctx); err != nil {
		fmt.Fprintln(os.Stderr, err)
		os.Exit(1)
	}
}
