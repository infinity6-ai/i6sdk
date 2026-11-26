//go:generate go run github.com/infinity6-ai/i6sdk/mold schemagen .
package dszapisignedurl

import (
	"context"

	"github.com/infinity6-ai/gox/fsz/fsz/fszsignmethod"
	"github.com/infinity6-ai/gox/routez/apiz"
	"github.com/infinity6-ai/gox/schemaz/schemaz"
	"github.com/infinity6-ai/i6sdk/sdkcorez/dsz/lakedsz"
	"github.com/infinity6-ai/i6sdk/sdkcorez/dsz/laketbpathz"
	"github.com/infinity6-ai/i6sdk/sdkcorez/speclabels/speclabeltokentype"
)

type Path struct {
	Method fszsignmethod.SignMethod `json:"method"`
	Path   laketbpathz.Path         `json:"path"`
}

type Url struct {
	Url string `json:"url"`
}

type ReqParams struct {
	// Dataset identifier.
	Dataset lakedsz.Dataset `json:"dataset"`
}

type ReqPayload struct {
	// Paths to sign.
	Paths []*Path `json:"paths"`
}

type Req struct {
	Params  *ReqParams  `json:"params,omitempty"`
	Payload *ReqPayload `json:"payload,omitempty"`
}

type RespPayload struct {
	// Signed URLs.
	Urls []*Url `json:"urls"`
}

type Resp struct {
	Payload *RespPayload `json:"payload,omitempty"`
}

// Generate Signed URLs
// Generates pre-signed storage URLs for dataset lake table paths.
type Api struct {
	apiz.UnimplementedApi `json:"-"`
	Req                   *Req  `json:"req,omitempty"`
	Resp                  *Resp `json:"resp,omitempty"`
}

func (d *Api) ApiRoute() apiz.ApiSpec {
	return apiz.ApiSpec{
		Method: "POST",
		Path:   "/api/ds/dataset/{dataset}/signed-url",
		Desc: func() *schemaz.Desc {
			return &schemaz.Desc{
				Labels: map[string]string{
					speclabeltokentype.Key: speclabeltokentype.Dataset,
				},
			}
		},
	}
}

type DszSignedUrlId struct{}

func (DszSignedUrlId) DszSignedUrl(ctx context.Context, api *Api) (int, error) {
	panic("DszSignedUrl not implemented: run 'go generate ./...'")
}
