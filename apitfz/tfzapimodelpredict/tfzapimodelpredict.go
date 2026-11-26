//go:generate go run github.com/infinity6-ai/i6sdk/mold@latest schemagen .
package tfzapimodelpredict

import (
	"context"

	"github.com/infinity6-ai/gox/routez/apiz"
	"github.com/infinity6-ai/gox/schemaz/schemaz"
	"github.com/infinity6-ai/i6sdk/sdkcorez/speclabels/speclabeltokentype"
)

type ReqParams struct {
	// Dataset identifier.
	Dataset string `json:"dataset"`
	// Model name.
	ModelName string `json:"model_name"`
	// Model version.
	Version string `json:"version"`
}

// Model inputs and partitions.
type ReqPayload struct {
	// Model inputs.
	Inputs string `json:"inputs"`
	// Model partition key/value pairs.
	ModelPartitions map[string]string `json:"model_partitions,omitempty"`
}

type Req struct {
	Params  *ReqParams  `json:"params,omitempty"`
	Payload *ReqPayload `json:"payload,omitempty"`
}

// Response containing the model outputs.
type RespPayload struct {
	// Model outputs.
	Outputs string `json:"outputs"`
}

type Resp struct {
	Payload *RespPayload `json:"payload,omitempty"`
}

// Model Predict
// Runs a prediction of a dataset model version on the inference service (tfinfer), sending the model inputs and the model partitions and returning the model outputs.
type Api struct {
	apiz.UnimplementedApi `json:"-"`
	Req                   *Req  `json:"req,omitempty"`
	Resp                  *Resp `json:"resp,omitempty"`
}

func (d *Api) ApiRoute() apiz.ApiSpec {
	return apiz.ApiSpec{
		Method: "POST",
		Path:   "/api/model/dataset/{dataset}/{model_name}/{version}/predict",
		Desc: func() *schemaz.Desc {
			return &schemaz.Desc{
				Labels: map[string]string{
					speclabeltokentype.Key: speclabeltokentype.Dataset,
					"servname":             "tfinfer",
				},
			}
		},
	}
}

type TfzModelPredictId struct{}

func (TfzModelPredictId) TfzModelPredict(ctx context.Context, api *Api) (int, error) {
	panic("TfzModelPredict not implemented: run 'go generate ./...'")
}
