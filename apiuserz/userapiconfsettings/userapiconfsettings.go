//go:generate go run github.com/infinity6-ai/i6sdk/mold schemagen .
package userapiconfsettings

import (
	"context"
	"github.com/infinity6-ai/gox/routez/apiz"
	"github.com/infinity6-ai/gox/schemaz/schemaz"
)

type Req struct{}

type RespPayload struct {
	// Maximum idle duration for refresh tokens in seconds before expiration.
	RefreshTokenMaxIdleSeconds int64 `json:"refresh_token_max_idle_seconds"`
	// Time-to-live for human user access tokens in seconds.
	UserAccessTokenTTLSeconds int64 `json:"user_access_token_ttl_seconds"`
	// Time-to-live for dataset access tokens in seconds.
	DatasetAccessTokenTTLSeconds int64 `json:"dataset_access_token_ttl_seconds"`
	// Time-to-live for RFC 8628 device codes in seconds.
	DeviceCodeTTLSeconds int64 `json:"device_code_ttl_seconds"`
	// Maximum allowed future expiration for Ed25519 request signatures in seconds.
	SignedRequestMaxTTLSeconds int64 `json:"signed_request_max_ttl_seconds"`
}

type Resp struct {
	Payload *RespPayload `json:"payload,omitempty"`
}

// Public Settings Discovery
// Retrieves public operational policies and token expiration parameters.
type Api struct {
	apiz.UnimplementedApi `json:"-"`
	Req                   *Req  `json:"req,omitempty"`
	Resp                  *Resp `json:"resp,omitempty"`
}

func (d *Api) ApiRoute() apiz.ApiSpec {
	return apiz.ApiSpec{
		Method: "GET",
		Path:   "/api/userz/conf/settings.json",
		Desc: func() *schemaz.Desc {
			return &schemaz.Desc{
				Labels: map[string]string{"servname": "userz"},
			}
		},
	}
}

type UserConfSettingsId struct{}

func (UserConfSettingsId) UserConfSettings(ctx context.Context, api *Api) (int, error) {
	panic("UserConfSettings not implemented: run 'go generate ./...'")
}
