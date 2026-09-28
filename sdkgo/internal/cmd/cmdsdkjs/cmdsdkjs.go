package cmdsdkjs

import (
	"os"

	"github.com/infinity6-ai/gox/commonz/errorz"
	"github.com/spf13/cobra"
	"go.code.infinity6.ai/platform/logz"
	"go.code.infinity6.ai/sdkgo/genjsz/jsapiwriter"
	"go.code.infinity6.ai/sdkgo/genjsz/jsparser"
	"go.code.infinity6.ai/sdkgo/sdkgo/sdkz"
)

type tlogger logz.Type

var logger = logz.Create(tlogger(true))

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
			ctx := cmd.Context()
			out := errorz.Check2(cmd.Flags().GetString("out"))
			apis := sdkz.All()
			parsedApis := jsparser.ParseApis(ctx, apis)
			code := []byte(jsapiwriter.Write(ctx, parsedApis))
			err := os.WriteFile(out, []byte(code), os.ModePerm)
			errorz.Check(err)
			logger.Info(ctx, "sdkapis generated", map[string]any{"out": out, "len": len(code)})
			return nil
		},
	}
	cmd.Flags().StringVar(&out, "out", "sdkweblegacy/src/sdkapis/sdkapis.js", "output directory, one .js file per api")
	return cmd
}
