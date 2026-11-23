package cmdforgez

import (
	"context"
	"fmt"
	"os"

	"github.com/infinity6-ai/gox/commonz/errorz"
	"github.com/infinity6-ai/gox/routez/apispecz"
	"github.com/infinity6-ai/i6sdk/apicollectionz/apicollectionz"
	"github.com/infinity6-ai/i6sdk/forgez/forgepy/forgepybundle"
	"github.com/infinity6-ai/i6sdk/forgez/forgepy/forgepyschema"
	"github.com/infinity6-ai/i6sdk/mold/moldsdkfunc"
	"github.com/infinity6-ai/i6sdk/versionz/version"
	"github.com/spf13/cobra"
)

func Prepare(ctx context.Context) *cobra.Command {
	var rootCmd = &cobra.Command{
		Version: version.Version(),
		Use:     "forgez",
		Short:   "forgez",
	}

	rootCmd.AddCommand(&cobra.Command{
		Use: "version",
		Run: func(cmd *cobra.Command, args []string) {
			fmt.Println(version.Version())
		},
	})

	rootCmd.AddCommand(preparePyCmd())

	return rootCmd
}

func preparePyCmd() *cobra.Command {
	var out string
	cmd := &cobra.Command{
		Use:   "py",
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
				return fmt.Errorf("failed to generate python SDK: %w", err)
			}
			return nil
		},
	}
	cmd.Flags().StringVar(&out, "out", "sdkpy/src/i6sdk/sdkapis", "output directory for generated python sdkapis")
	return cmd
}

func Execute(ctx context.Context, rootCmd *cobra.Command) {
	if err := rootCmd.ExecuteContext(ctx); err != nil {
		fmt.Fprintln(os.Stderr, err)
		os.Exit(1)
	}
}
