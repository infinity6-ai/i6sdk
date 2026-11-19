package cmdbqschema

import (
	"os"
	"path/filepath"
	"strings"

	"github.com/infinity6-ai/gox/bqz/bqzschema"
	"github.com/infinity6-ai/gox/commonz/errorz"
	"github.com/infinity6-ai/gox/commonz/filez"
	"github.com/infinity6-ai/gox/commonz/logz"
	"github.com/spf13/cobra"
	"go.code.infinity6.ai/sdkgo/sdkgo/apiz/ingestzapi/ingestzapistream"
	"go.code.infinity6.ai/sdkgo/sdkgo/tabledataz"
)

type tlogger logz.Type

var logger = logz.Create(tlogger(true))

func Prepare(parent *cobra.Command) *cobra.Command {
	cmd := &cobra.Command{
		Use: "bq",
	}
	parent.AddCommand(cmd)

	schemaCmd := &cobra.Command{
		Use: "schema",
	}
	cmd.AddCommand(schemaCmd)

	schemaCmd.AddCommand(prepareGenerateSchemaCmd())

	return cmd
}

func TableDir(td tabledataz.TableData) string {
	return strings.ReplaceAll(td.Table().String(), "-", "")
}

func prepareGenerateSchemaCmd() *cobra.Command {
	var out string
	cmd := &cobra.Command{
		Use:   "generate",
		Short: "generate the BigQuery schema.json of every streamed table into the i6-infra terraform modules",
		Args:  cobra.NoArgs,
		RunE: func(cmd *cobra.Command, args []string) error {
			ctx := cmd.Context()
			out := errorz.Check2(cmd.Flags().GetString("out"))

			for _, td := range ingestzapistream.AllTables() {
				code := errorz.Check2(bqzschema.NativeTableJSON(td))
				file := filepath.Join(out, TableDir(td), "schema.json")
				errorz.Check(filez.CreateParentDirs(file))
				errorz.Check(os.WriteFile(file, code, 0o644))
				logger.Info(ctx, "bq schema generated", map[string]any{"table": td.Table().String(), "out": file, "len": len(code)})
			}
			return nil
		},
	}
	cmd.Flags().StringVar(&out, "out", "../i6box/i6-infra/src/bq/bqtables", "terraform bqtables directory, one <table>/schema.json per streamed table")
	return cmd
}
