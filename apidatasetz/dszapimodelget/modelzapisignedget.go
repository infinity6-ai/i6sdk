//go:generate go run github.com/infinity6-ai/i6sdk/mold schemagen .
package dszapimodelget

import (
	"context"

	"github.com/infinity6-ai/gox/routez/apiz"
	"github.com/infinity6-ai/gox/schemaz/schemaz"
	"github.com/infinity6-ai/i6sdk/sdkcorez/dsz/lakedsz"
	"github.com/infinity6-ai/i6sdk/sdkcorez/speclabels/speclabeltokentype"
)

type ModelFile struct {
	// File name.
	Name string `json:"name"`
	// Signed URL.
	Url string `json:"url"`
}

type ReqParams struct {
	// Dataset identifier.
	Dataset lakedsz.Dataset `json:"dataset"`
	// Model name.
	ModelName string `json:"model_name"`
}

type ReqPayload struct {
	// Model version.
	ModelVersion string `json:"model_version,omitempty"`
	// Partition key-value pairs.
	ModelPartitions map[string]string `json:"model_partitions,omitempty"`
}

type Req struct {
	Params  *ReqParams  `json:"params,omitempty"`
	Payload *ReqPayload `json:"payload,omitempty"`
}

type RespPayload struct {
	// List of model files with signed access URLs.
	Urls []*ModelFile `json:"urls"`
}

type Resp struct {
	Payload *RespPayload `json:"payload,omitempty"`
}

// Get Model Signed URLs
// Generates presigned URLs for accessing model files in the specified dataset.
type Api struct {
	apiz.UnimplementedApi `json:"-"`
	Req                   *Req  `json:"req,omitempty"`
	Resp                  *Resp `json:"resp,omitempty"`
}

func (d *Api) ApiRoute() apiz.ApiSpec {
	return apiz.ApiSpec{
		Method: "POST",
		Path:   "/api/model/dataset/{dataset}/signed-get/{model_name}",
		Desc: func() *schemaz.Desc {
			return &schemaz.Desc{
				Labels: map[string]string{
					speclabeltokentype.Key: speclabeltokentype.Dataset,
				},
			}
		},
	}
}

type DszModelGetId struct{}

func (DszModelGetId) DszModelGet(ctx context.Context, api *Api) (int, error) {
	panic("DszModelGet not implemented: run 'go generate ./...'")
}
