//go:generate go run github.com/infinity6-ai/i6sdk/mold schemagen .
package tfzapimodelmetadata

import (
	"context"

	"github.com/infinity6-ai/gox/routez/apiz"
)

// TensorDetail describes an input or output tensor of a model.
type TensorDetail struct {
	// Tensor name.
	Name string `json:"name"`
	// Tensor index.
	Index int `json:"index"`
	// Tensor dimensions.
	Shape []int `json:"shape"`
	// Tensor data type.
	DType string `json:"dtype"`
}

type ReqParams struct {
	// Model name.
	ModelName string `json:"model_name"`
	// Model version.
	Version string `json:"version"`
}

type Req struct {
	Params *ReqParams `json:"params,omitempty"`
}

// Response containing the model tensors.
type RespPayload struct {
	// Model input tensors.
	Inputs []TensorDetail `json:"inputs"`
	// Model output tensors.
	Outputs []TensorDetail `json:"outputs"`
}

type Resp struct {
	Payload *RespPayload `json:"payload,omitempty"`
}

// Get Model Metadata
// Gets the input and output tensors of a model version.
type Api struct {
	apiz.UnimplementedApi `json:"-"`
	Req                   *Req  `json:"req,omitempty"`
	Resp                  *Resp `json:"resp,omitempty"`
}

func (d *Api) ApiRoute() apiz.ApiSpec {
	return apiz.ApiSpec{
		Method: "GET",
		Path:   "/api/metadata/{model_name}/{version}",
	}
}

type TfzModelMetadataId struct{}

func (TfzModelMetadataId) TfzModelMetadata(ctx context.Context, api *Api) (int, error) {
	panic("TfzModelMetadata not implemented: run 'go generate ./...'")
}
