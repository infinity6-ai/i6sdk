//go:generate go run github.com/infinity6-ai/i6sdk/mold schemagen .
package dszapitableget

import (
	"context"

	"github.com/infinity6-ai/gox/routez/apiz"
	"github.com/infinity6-ai/gox/schemaz/schemaz"
	"github.com/infinity6-ai/i6sdk/sdkcorez/dsz/lakedsz"
	"github.com/infinity6-ai/i6sdk/sdkcorez/dsz/laketbpathz"
	"github.com/infinity6-ai/i6sdk/sdkcorez/speclabels/speclabeltokentype"
)

type ReqParams struct {
	// Dataset identifier.
	Dataset lakedsz.Dataset `json:"dataset"`
	// Table file path without scheme.
	Path laketbpathz.Path `json:"path"`
}

type Req struct {
	Params *ReqParams `json:"params,omitempty"`
}

type RespHeaders struct {
	// Redirect location URL.
	Location string `json:"location"`
}

type Resp struct {
	Headers *RespHeaders `json:"headers,omitempty"`
}

// Get Dataset Table File
// Get table file by redirecting to its signed URL.
type Api struct {
	apiz.UnimplementedApi `json:"-"`
	Req                   *Req  `json:"req,omitempty"`
	Resp                  *Resp `json:"resp,omitempty"`
}

func (d *Api) ApiRoute() apiz.ApiSpec {
	return apiz.ApiSpec{
		Method: "GET",
		Path:   "/api/ds/dataset/{dataset}/table-file/{path...}",
		Desc: func() *schemaz.Desc {
			return &schemaz.Desc{
				Labels: map[string]string{
					speclabeltokentype.Key: speclabeltokentype.Dataset,
				},
			}
		},
	}
}

type DszTableGetId struct{}

func (DszTableGetId) DszTableGet(ctx context.Context, api *Api) (int, error) {
	panic("DszTableGet not implemented: run 'go generate ./...'")
}
