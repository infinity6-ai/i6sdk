package cmdbqschema

import (
	"testing"

	"github.com/infinity6-ai/gox/bqz/bqzschema"
	"github.com/stretchr/testify/require"
	"go.code.infinity6.ai/sdkgo/sdkgo/tabledataz/relevancefashiondataz"
)

func TestUnitTableDir(t *testing.T) {
	require.Equal(t, "relevancefashioncatalog", TableDir(&relevancefashiondataz.Catalog{}))
	require.Equal(t, "relevancefashionevent", TableDir(&relevancefashiondataz.Event{}))
}

func TestUnitAllTablesGenerateSchema(t *testing.T) {
	for _, td := range Tables() {
		_, err := bqzschema.NativeTableJSON(td)
		require.NoError(t, err, td.Table().String())
	}
}
