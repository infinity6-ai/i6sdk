package resultz

import (
	"context"

	"github.com/infinity6-ai/gox/commonz/urlz"
	"github.com/infinity6-ai/gox/schemaz/schemaz"
	"github.com/infinity6-ai/i6sdk/sdkcorez/dsz/laketbpathz"
)

type Result interface {
	Schema() *schemaz.Schema
	ResultName() string

	QuerySchema() *schemaz.Schema
	Path() laketbpathz.Path
	Download(ctx context.Context, u *urlz.Url) bool
}
