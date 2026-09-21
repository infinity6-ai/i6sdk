//go:build !testcover

package cmdsdkgo

import (
	"fmt"
	"path/filepath"

	"github.com/spf13/cobra"

	"go.code.infinity6.ai/sdkgo/generatorjsz"
	"go.code.infinity6.ai/sdkgo/sdkgo/specz"
)

// defaultGenerateJsOut is where the legacy web SDK expects its generated API clients, relative to
// the repository root (the cwd of ./build.sh).
const defaultGenerateJsOut = "sdkweblegacy/src/apigenz"

// prepareGenerateCmd builds `generate`, the parent of every code generator.
func prepareGenerateCmd() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "generate",
		Short: "generate code from the API schemas",
		// Runnable only to print the help: it makes cobra reject unknown subcommands instead of
		// silently succeeding.
		Args: cobra.NoArgs,
		RunE: func(cmd *cobra.Command, args []string) error {
			return cmd.Help()
		},
	}
	cmd.AddCommand(prepareGenerateJsCmd())
	return cmd
}

// prepareGenerateJsCmd builds `generate js` (./build.sh run generate js): it generates the JS
// client of every API listed by specz.Sdk(), one file per API.
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
