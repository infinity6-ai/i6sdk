//go:generate go run github.com/infinity6-ai/i6sdk/mold schemagen .
package dszapidatasetlist

import (
	"github.com/infinity6-ai/gox/routez/apiz"
	"github.com/infinity6-ai/gox/schemaz/schemaz"
	"github.com/infinity6-ai/i6sdk/sdkcorez/grants"
	"github.com/infinity6-ai/i6sdk/sdkcorez/speclabels/speclabeltokentype"
)

type Solution struct {
	Id   string `json:"id"`
	Name string `json:"name"`
}

// Params contains path parameters embedded in the route URL.
type Params struct {
	// User identifier (id).
	UserId string `json:"userId"`
}

// Dataset represents an Infinity6 dataset container and permissions.
type Dataset struct {
	// Unique identifier of the dataset.
	Id string `json:"id"`
	// Human-readable name of the dataset.
	Name string `json:"name"`
	// Compact grant bitmask held by the user on this dataset.
	Grants grants.Grant `json:"grants"`
}

// Payload represents the JSON HTTP response payload.
type Payload struct {
	// List of datasets available to the user.
	Datasets []*Dataset `json:"datasets"`
}

// Req encapsulates all incoming HTTP data for the endpoint.
type Req struct {
	Params *Params `json:"params,omitempty"`
}

// Resp encapsulates all outgoing HTTP data for the endpoint.
type Resp struct {
	Payload *Payload `json:"payload"`
}

// Backwards-compatible aliases.
type DszListDatasetsReq = Req
type DszListDatasetsResp = Resp

type Api struct {
	apiz.UnimplementedDataRefs `json:"-"`
	Req                        *Req  `json:"req"`
	Resp                       *Resp `json:"resp"`
}

func (d *Api) ApiSpec() apiz.ApiSpec {
	return apiz.ApiSpec{
		Id: "dsz-list-datasets",
		Desc: func() *schemaz.Desc {
			return &schemaz.Desc{
				Name:    "List Datasets",
				Summary: "Lists all datasets available to the user.",
				Labels: map[string]string{
					speclabeltokentype.Key: speclabeltokentype.User,
				},
			}
		},
		Method: "GET",
		Path:   "/api/ds/user-id/{userId}/datasets",
	}
}
