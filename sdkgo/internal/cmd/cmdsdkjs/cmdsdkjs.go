package cmdsdkjs

import (
	"os"

	"github.com/infinity6-ai/gox/commonz/errorz"
	"github.com/infinity6-ai/gox/commonz/filez"
	"github.com/infinity6-ai/gox/commonz/logz"
	"github.com/infinity6-ai/gox/routez/apispecz"
	"github.com/infinity6-ai/i6sdk/apicollectionz"
	"github.com/spf13/cobra"
	"go.code.infinity6.ai/sdkgo/forgez/forgejs/forgejsbundle"
	"go.code.infinity6.ai/sdkgo/forgez/forgejs/forgejsschema"
	"go.code.infinity6.ai/sdkgo/mold/moldsdkfunc"
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

			schemaParser := forgejsschema.NewSdk()
			apiParser := apispecz.New(schemaParser)
			sdkFuncParser := moldsdkfunc.New(schemaParser)

			apis := apicollectionz.Apis()
			parsedApis := apiParser.ParseApis(ctx, apis)

			funcs := apicollectionz.SdkFuncs()
			parsedSdkFuncs := sdkFuncParser.ParseSdkFuncs(ctx, funcs)

			code := []byte(forgejsbundle.Write(ctx, forgejsbundle.Options{
				Apis:     parsedApis.Apis,
				SdkFuncs: parsedSdkFuncs.SdkFuncs,
			}))

			errorz.Check(filez.CreateParentDirs(out))
			err := os.WriteFile(out, []byte(code), os.ModePerm)
			errorz.Check(err)
			logger.Info(ctx, "sdkapis generated", map[string]any{"out": out, "len": len(code)})
			return nil
		},
	}
	cmd.Flags().StringVar(&out, "out", "sdkweblegacy/src/sdkapis/sdkapis.js", "output directory, one .js file per api")
	return cmd
}
