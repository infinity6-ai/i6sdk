package userapiauthme

import (
	"github.com/infinity6-ai/gox/routez/apiz"
	"github.com/infinity6-ai/gox/schemaz/schemaz"
	"go.code.infinity6.ai/sdkgo/sdkgo/dsz/lakedsz"
)

type UserDatasetGrant struct {
	Dataset lakedsz.Dataset `json:"dataset"`
	Grants  uint8           `json:"grants"`
}

type Req struct {
	AuthHeader string `json:"authorization"`
}

type Resp struct {
	UserId        string             `json:"user_id"`
	EmailRedacted string             `json:"email_redacted"`
	Providers     []string           `json:"providers"`
	Datasets      []UserDatasetGrant `json:"datasets"`
}

type Api struct {
	Req  Req  `json:"req"`
	Resp Resp `json:"resp"`
}

func (d *Api) ApiSpec() apiz.ApiSpec {
	return apiz.ApiSpec{
		Id: "user-auth-me",
		Desc: func() *schemaz.Desc {
			return &schemaz.Desc{
				Name:    "Current User & Workspaces",
				Summary: "Retrieves current authenticated user identity and accessible dataset workspaces.",
				Labels:  map[string]string{"servname": "userz"},
			}
		},
		Method: "GET",
		Path:   "/api/userz/auth/me",
	}
}

func (d *Api) ApiDataRefs() *apiz.DataRefs {
	return &apiz.DataRefs{
		PathParams:  nil,
		QueryParams: nil,
		ReqHeaders: &schemaz.Schema{
			Object: func(read bool) map[string]*schemaz.Schema {
				return map[string]*schemaz.Schema{
					"authorization": {
						Desc:     &schemaz.Desc{Summary: "User Bearer JWT access token."},
						UrlValue: schemaz.ParseStr(&d.Req.AuthHeader),
					},
				}
			},
		},
		ReqBody:     nil,
		RespHeaders: nil,
		RespBody: &schemaz.Schema{
			Desc: &schemaz.Desc{Summary: "Current user profile and accessible workspaces."},
			Object: func(read bool) map[string]*schemaz.Schema {
				return map[string]*schemaz.Schema{
					"user_id": {
						Desc:  &schemaz.Desc{Summary: "User unique identifier."},
						Field: func() any { return &d.Resp.UserId },
					},
					"email_redacted": {
						Desc:  &schemaz.Desc{Summary: "Redacted email address for display."},
						Field: func() any { return &d.Resp.EmailRedacted },
					},
					"providers": {
						Desc: &schemaz.Desc{Summary: "List of linked identity providers."},
						Array: schemaz.NewArray(schemaz.ArrayOptions[[]string, string]{
							Slice: func() *[]string { return &d.Resp.Providers },
							Get: func(idx int, element *string) *schemaz.Schema {
								return &schemaz.Schema{
									Field: func() any { return element },
								}
							},
						}),
					},
					"datasets": {
						Desc: &schemaz.Desc{Summary: "List of accessible dataset workspaces."},
						Array: schemaz.NewArray(schemaz.ArrayOptions[[]UserDatasetGrant, UserDatasetGrant]{
							Slice: func() *[]UserDatasetGrant { return &d.Resp.Datasets },
							Get: func(idx int, element *UserDatasetGrant) *schemaz.Schema {
								return &schemaz.Schema{
									Object: func(_ bool) map[string]*schemaz.Schema {
										return map[string]*schemaz.Schema{
											"dataset": {
												Desc:  &schemaz.Desc{Summary: "Dataset identifier."},
												Field: func() any { return &element.Dataset },
											},
											"grants": {
												Desc:  &schemaz.Desc{Summary: "Grants bitmask."},
												Field: func() any { return &element.Grants },
											},
										}
									},
								}
							},
						}),
					},
				}
			},
		},
	}
}
