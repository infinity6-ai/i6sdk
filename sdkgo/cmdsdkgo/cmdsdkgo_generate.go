//go:build !testcover

package cmdsdkgo

import (
	"fmt"
	"path/filepath"

	"github.com/spf13/cobra"

	"go.code.infinity6.ai/sdkgo/generatorjsz"
	"go.code.infinity6.ai/sdkgo/sdkgo/specz"
)

const defaultGenerateJsOut = "sdkweblegacy/src/apigenz"

func prepareGenerateCmd() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "generate",
		Short: "generate code from the API schemas",
		Args:  cobra.NoArgs,
		RunE: func(cmd *cobra.Command, args []string) error {
			return cmd.Help()
		},
	}
	cmd.AddCommand(prepareGenerateJsCmd())
	return cmd
}

func prepareGenerateJsCmd() *cobra.Command {
	var out string
	cmd := &cobra.Command{
		Use:   "js",
		Short: "generate the JS client of every API into the legacy web SDK",
		Args:  cobra.NoArgs,
		RunE: func(cmd *cobra.Command, args []string) error {
			names, err := generatorjsz.Generate(out, specz.Sdk())
			if err != nil {
				return fmt.Errorf("failed to generate js clients into %s: %w", out, err)
			}
			w := cmd.OutOrStdout()
			for _, name := range names {
				fmt.Fprintln(w, filepath.Join(out, name))
			}
			fmt.Fprintf(w, "generated %d apis into %s\n", len(names), out)
			return nil
		},
	}
	cmd.Flags().StringVar(&out, "out", defaultGenerateJsOut, "output directory, one .js file per api")
	return cmd
}
