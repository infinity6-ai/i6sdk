package resultz

import (
	"github.com/infinity6-ai/gox/schemaz/schemaz"
	"go.code.infinity6.ai/sdkgo/sdkgo/dsz/laketbpathz"
)

type Result interface {
	Schema() *schemaz.Schema
	ResultName() string

	NewQuery() any
	QuerySchema(q any) *schemaz.Schema
	Path(q any) laketbpathz.Path
}
