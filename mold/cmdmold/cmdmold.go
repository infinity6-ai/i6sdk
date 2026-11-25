package cmdmold

import (
	"context"
	"fmt"
	"os"

	"github.com/infinity6-ai/i6sdk/mold/schemagen"
	"github.com/infinity6-ai/i6sdk/versionz/version"
	"github.com/spf13/cobra"
)

func Prepare(ctx context.Context) *cobra.Command {
	var rootCmd = &cobra.Command{
		Version: version.Version(),
		Use:     "mold",
		Short:   "mold",
	}

	rootCmd.AddCommand(&cobra.Command{
		Use: "version",
		Run: func(cmd *cobra.Command, args []string) {
			fmt.Println(version.Version())
		},
	})

	rootCmd.AddCommand(prepareSchemagenCmd())

	return rootCmd
}

func prepareSchemagenCmd() *cobra.Command {
	var dir string
	cmd := &cobra.Command{
		Use:   "schemagen [dir]",
		Short: "generate ApiDataRefs schema code for Go API packages",
		Args:  cobra.MaximumNArgs(1),
		RunE: func(cmd *cobra.Command, args []string) error {
			targetDir := dir
			if len(args) > 0 {
				targetDir = args[0]
			}
			if targetDir == "" {
				targetDir = "."
			}
			outFile, err := schemagen.GenerateDir(targetDir)
			if err != nil {
				return err
			}
			fmt.Printf("Generated %s\n", outFile)
			return nil
		},
	}
	cmd.Flags().StringVar(&dir, "dir", ".", "target package directory")
	return cmd
}

func Execute(ctx context.Context, rootCmd *cobra.Command) {
	if err := rootCmd.ExecuteContext(ctx); err != nil {
		fmt.Fprintln(os.Stderr, err)
		os.Exit(1)
	}
}
