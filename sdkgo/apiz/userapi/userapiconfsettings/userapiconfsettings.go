package userapiconfsettings

import (
	"github.com/infinity6-ai/gox/routez/apiz"
	"github.com/infinity6-ai/gox/schemaz/schemaz"
)

type Req struct{}

type Resp struct {
	RefreshTokenMaxIdleSeconds  int64 `json:"refresh_token_max_idle_seconds"`
	UserAccessTokenTTLSeconds    int64 `json:"user_access_token_ttl_seconds"`
	DatasetAccessTokenTTLSeconds int64 `json:"dataset_access_token_ttl_seconds"`
	DeviceCodeTTLSeconds        int64 `json:"device_code_ttl_seconds"`
}

type Api struct {
	Req  Req  `json:"req"`
	Resp Resp `json:"resp"`
}

func (d *Api) ApiSpec() apiz.ApiSpec {
	return apiz.ApiSpec{
		Id: "user-conf-settings",
		Desc: func() *schemaz.Desc {
			return &schemaz.Desc{
				Name:    "Public Settings Discovery",
				Summary: "Retrieves public operational policies and token expiration parameters.",
				Labels:  map[string]string{"servname": "userz"},
			}
		},
		Method: "GET",
		Path:   "/api/userz/conf/settings.json",
	}
}

func (d *Api) ApiDataRefs() *apiz.DataRefs {
	return &apiz.DataRefs{
		PathParams:  nil,
		QueryParams: nil,
		ReqHeaders:  nil,
		ReqBody:     nil,
		RespHeaders: nil,
		RespBody: &schemaz.Schema{
			Desc: &schemaz.Desc{Summary: "Public configuration and settings."},
			Object: func(read bool) map[string]*schemaz.Schema {
				return map[string]*schemaz.Schema{
					"refresh_token_max_idle_seconds": {
						Desc:  &schemaz.Desc{Summary: "Maximum idle duration for refresh tokens in seconds before expiration."},
						Field: func() any { return &d.Resp.RefreshTokenMaxIdleSeconds },
					},
					"user_access_token_ttl_seconds": {
						Desc:  &schemaz.Desc{Summary: "Time-to-live for human user access tokens in seconds."},
						Field: func() any { return &d.Resp.UserAccessTokenTTLSeconds },
					},
					"dataset_access_token_ttl_seconds": {
						Desc:  &schemaz.Desc{Summary: "Time-to-live for dataset access tokens in seconds."},
						Field: func() any { return &d.Resp.DatasetAccessTokenTTLSeconds },
					},
					"device_code_ttl_seconds": {
						Desc:  &schemaz.Desc{Summary: "Time-to-live for RFC 8628 device codes in seconds."},
						Field: func() any { return &d.Resp.DeviceCodeTTLSeconds },
					},
				}
			},
		},
	}
}
