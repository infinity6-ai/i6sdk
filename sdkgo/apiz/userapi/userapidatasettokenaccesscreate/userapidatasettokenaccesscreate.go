package userapidatasettokenaccesscreate

import (
	"github.com/infinity6-ai/gox/routez/apiz"
	"github.com/infinity6-ai/gox/schemaz/schemaz"
	"go.code.infinity6.ai/sdkgo/sdkgo/dsz/lakedsz"
	"go.code.infinity6.ai/sdkgo/sdkgo/dsz/solution"
)

type Req struct {
	Dataset      lakedsz.Dataset   `json:"dataset"`
	Solution     solution.Solution `json:"solution"`
	Grants       uint8             `json:"grants"`
	Labels       map[string]string `json:"labels,omitempty"`
	ClientPubKey string            `json:"client_pub_key,omitempty"`
	ExpSeconds   int64             `json:"exp_seconds,omitempty"`
	AuthHeader   string            `json:"authorization"`
}

type Resp struct {
	AccessToken string            `json:"access_token"`
	Dataset     lakedsz.Dataset   `json:"dataset"`
	Solution    solution.Solution `json:"solution"`
	Grants      uint8             `json:"grants"`
	Labels      map[string]string `json:"labels,omitempty"`
	ExpiresAt   int64             `json:"expires_at"`
}

type Api struct {
	Req  Req  `json:"req"`
	Resp Resp `json:"resp"`
}

func (d *Api) ApiSpec() apiz.ApiSpec {
	return apiz.ApiSpec{
		Id: "user-dataset-token-access-create",
		Desc: func() *schemaz.Desc {
			return &schemaz.Desc{
				Name:    "Create Dataset Access Token",
				Summary: "Mints a short-lived dataset access token for an authenticated user.",
			}
		},
		Method: "POST",
		Path:   "/api/user/dataset/access-token",
	}
}

func (d *Api) ApiDataRefs() *apiz.DataRefs {
	return &apiz.DataRefs{
		PathParams:  nil,
		QueryParams: nil,
		ReqHeaders: &schemaz.Schema{
			Object: func(read bool) map[string]*schemaz.Schema {
				return map[string]*schemaz.Schema{
					"Authorization": {
						Desc:  &schemaz.Desc{Summary: "User Bearer JWT access token."},
						Field: func() any { return &d.Req.AuthHeader },
					},
				}
			},
		},
		ReqBody: &schemaz.Schema{
			Desc: &schemaz.Desc{Summary: "Dataset access token creation request."},
			Object: func(read bool) map[string]*schemaz.Schema {
				return map[string]*schemaz.Schema{
					"dataset": {
						Desc:     &schemaz.Desc{Summary: "Target dataset ID."},
						UrlValue: schemaz.ParseStrGet(&d.Req.Dataset, lakedsz.MustParse),
					},
					"solution": {
						Desc:     &schemaz.Desc{Summary: "Target solution domain restriction (empty for all)."},
						UrlValue: schemaz.ParseStrGet(&d.Req.Solution, solution.MustParse),
					},
					"grants": {
						Desc:  &schemaz.Desc{Summary: "Compact uint8 grant bitmask."},
						Field: func() any { return &d.Req.Grants },
					},
					"labels": {
						Desc:  &schemaz.Desc{Summary: "Contextual metadata labels."},
						Field: func() any { return &d.Req.Labels },
					},
					"client_pub_key": {
						Desc:  &schemaz.Desc{Summary: "Base64URL client public key for PoP binding (empty for unbound bearer)."},
						Field: func() any { return &d.Req.ClientPubKey },
					},
					"exp_seconds": {
						Desc:  &schemaz.Desc{Summary: "Requested lifetime in seconds."},
						Field: func() any { return &d.Req.ExpSeconds },
					},
				}
			},
		},
		RespHeaders: nil,
		RespBody: &schemaz.Schema{
			Desc: &schemaz.Desc{Summary: "Dataset access token response."},
			Object: func(read bool) map[string]*schemaz.Schema {
				return map[string]*schemaz.Schema{
					"access_token": {
						Desc:  &schemaz.Desc{Summary: "Binary-packed Base64URL dataset access token."},
						Field: func() any { return &d.Resp.AccessToken },
					},
					"dataset": {
						Desc:     &schemaz.Desc{Summary: "Target dataset ID."},
						UrlValue: schemaz.ParseStrGet(&d.Resp.Dataset, lakedsz.MustParse),
					},
					"solution": {
						Desc:     &schemaz.Desc{Summary: "Solution domain restriction."},
						UrlValue: schemaz.ParseStrGet(&d.Resp.Solution, solution.MustParse),
					},
					"grants": {
						Desc:  &schemaz.Desc{Summary: "Effective uint8 grant bitmask."},
						Field: func() any { return &d.Resp.Grants },
					},
					"labels": {
						Desc:  &schemaz.Desc{Summary: "Contextual metadata labels."},
						Field: func() any { return &d.Resp.Labels },
					},
					"expires_at": {
						Desc:  &schemaz.Desc{Summary: "Expiration epoch milliseconds."},
						Field: func() any { return &d.Resp.ExpiresAt },
					},
				}
			},
		},
	}
}
