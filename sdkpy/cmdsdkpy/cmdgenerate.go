package cmdsdkpy

import (
	"fmt"

	"github.com/infinity6-ai/gox/commonz/errorz"
	"github.com/infinity6-ai/gox/commonz/logz"
	"github.com/infinity6-ai/gox/routez/apispecz"
	"github.com/infinity6-ai/i6sdk/apicollectionz/apicollectionz"
	"github.com/infinity6-ai/i6sdk/forgez/forgepy/forgepybundle"
	"github.com/infinity6-ai/i6sdk/forgez/forgepy/forgepyschema"
	"github.com/infinity6-ai/i6sdk/mold/moldsdkfunc"
	"github.com/spf13/cobra"
)

type tlogger logz.Type

var logger = logz.Create(tlogger(true))

func prepareGenerateCmd() *cobra.Command {
	var out string
	cmd := &cobra.Command{
		Use:   "generate",
		Short: "generate the Python client of every API into sdkpy",
		Args:  cobra.NoArgs,
		RunE: func(cmd *cobra.Command, args []string) error {
			ctx := cmd.Context()
			outVal := errorz.Check2(cmd.Flags().GetString("out"))

			schemaParser := forgepyschema.NewSdk()
			apiParser := apispecz.New(schemaParser)
			sdkFuncParser := moldsdkfunc.New(schemaParser)

			apis := apicollectionz.Apis()
			parsedApis := apiParser.ParseApis(ctx, apis)

			funcs := apicollectionz.SdkFuncs()
			parsedSdkFuncs := sdkFuncParser.ParseSdkFuncs(ctx, funcs)

			if err := forgepybundle.WriteDir(ctx, outVal, forgepybundle.Options{
				Apis:     parsedApis.Apis,
				SdkFuncs: parsedSdkFuncs.SdkFuncs,
			}); err != nil {
				return fmt.Errorf("failed to generate python SDK into %s: %w", outVal, err)
			}
			logger.Info(ctx, "sdkapis generated for Python", map[string]any{"out": outVal})
			return nil
		},
	}
	cmd.Flags().StringVar(&out, "out", "src/i6sdk/sdkapis", "output directory for generated python sdkapis")
	return cmd
}
