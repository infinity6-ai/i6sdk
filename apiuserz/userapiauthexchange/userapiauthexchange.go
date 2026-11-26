//go:generate go run github.com/infinity6-ai/i6sdk/mold schemagen .
package userapiauthexchange

import (
	"context"
	"github.com/infinity6-ai/gox/routez/apiz"
	"github.com/infinity6-ai/gox/schemaz/schemaz"
)

type ReqPayload struct {
	// Authorization code issued by OAuth identity provider callback.
	Code string `json:"code"`
	// PKCE code_verifier corresponding to the original code_challenge.
	CodeVerifier string `json:"code_verifier"`
	// Optional extension parameter.
	Ext string `json:"ext,omitempty"`
}

type Req struct {
	Payload *ReqPayload `json:"payload,omitempty"`
}

type RespPayload struct {
	// Signed user access token JWT.
	UserAccessToken string `json:"user_access_token"`
	// Opaque user refresh token with new secret.
	UserRefreshToken string `json:"user_refresh_token"`
	// Authenticated user UUID.
	UserId string `json:"user_id"`
	// Access token expiration epoch millis.
	ExpiresAt int64 `json:"expires_at"`
	// User refresh token expiration epoch millis.
	RefreshTokenExpiresAt int64 `json:"refresh_token_expires_at"`
}

type Resp struct {
	Payload *RespPayload `json:"payload,omitempty"`
}

// OAuth 2.0 PKCE Code Exchange
// Exchanges an authorization code and PKCE code_verifier for authenticated user tokens.
type Api struct {
	apiz.UnimplementedApi `json:"-"`
	Req                   *Req  `json:"req,omitempty"`
	Resp                  *Resp `json:"resp,omitempty"`
}

func (d *Api) ApiRoute() apiz.ApiSpec {
	return apiz.ApiSpec{
		Method: "POST",
		Path:   "/api/userz/auth/exchange",
		Desc: func() *schemaz.Desc {
			return &schemaz.Desc{
				Labels: map[string]string{"servname": "userz"},
			}
		},
	}
}

type UserAuthExchangeId struct{}

func (UserAuthExchangeId) UserAuthExchange(ctx context.Context, api *Api) (int, error) {
	panic("UserAuthExchange not implemented: run 'go generate ./...'")
}
