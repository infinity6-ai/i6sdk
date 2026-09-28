package cmdsdkjs

import "github.com/spf13/cobra"

func Prepare(parent *cobra.Command) *cobra.Command {
	cmd := &cobra.Command{
		Use: "js",
	}
	parent.AddCommand(cmd)

	cmd.AddCommand(prepareGenerateJsCmd())

	return cmd
}

func prepareGenerateJsCmd() *cobra.Command {
	var out string
	cmd := &cobra.Command{
		Use:   "generate",
		Short: "generate the JS client of every API into the legacy web SDK",
		Args:  cobra.NoArgs,
		RunE: func(cmd *cobra.Command, args []string) error {
			// names, err := generatorjsz.Generate(out, sdkz.All())
			// if err != nil {
			// 	return fmt.Errorf("failed to generate js clients into %s: %w", out, err)
			// }
			// w := cmd.OutOrStdout()
			// for _, name := range names {
			// 	fmt.Fprintln(w, filepath.Join(out, name))
			// }
			// fmt.Fprintf(w, "generated %d apis into %s\n", len(names), out)
			return nil
		},
	}
	cmd.Flags().StringVar(&out, "out", "sdkweblegacy/src/apigenz", "output directory, one .js file per api")
	return cmd
}
