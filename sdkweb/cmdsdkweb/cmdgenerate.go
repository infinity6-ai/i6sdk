package cmdsdkweb

import (
	"os"

	"github.com/infinity6-ai/gox/commonz/errorz"
	"github.com/infinity6-ai/gox/commonz/filez"
	"github.com/infinity6-ai/gox/commonz/logz"
	"github.com/infinity6-ai/gox/routez/apispecz"
	"github.com/infinity6-ai/i6sdk/apicollectionz/apicollectionz"
	"github.com/infinity6-ai/i6sdk/forgez/forgejs/forgejsbundle"
	"github.com/infinity6-ai/i6sdk/forgez/forgejs/forgejsschema"
	"github.com/infinity6-ai/i6sdk/mold/moldsdkfunc"
	"github.com/spf13/cobra"
)

type tlogger logz.Type

var logger = logz.Create(tlogger(true))

func prepareGenerateCmd() *cobra.Command {
	var out string
	cmd := &cobra.Command{
		Use:   "generate",
		Short: "generate the JS client of every API into the web SDK",
		Args:  cobra.NoArgs,
		RunE: func(cmd *cobra.Command, args []string) error {
			ctx := cmd.Context()
			out := errorz.Check2(cmd.Flags().GetString("out"))

			schemaParser := forgejsschema.NewSdk()
			apiParser := apispecz.New(schemaParser)
			sdkFuncParser := moldsdkfunc.New(schemaParser)

			apis := apicollectionz.Apis()
			parsedApis := apiParser.ParseApis(ctx, apis)

			funcs := apicollectionz.Funcs()
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
	cmd.Flags().StringVar(&out, "out", "src/main/sdkapis/sdkapis.js", "output path for generated sdkapis.js")
	return cmd
}
