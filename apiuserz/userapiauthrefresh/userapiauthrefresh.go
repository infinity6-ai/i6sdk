//go:generate go run github.com/infinity6-ai/i6sdk/mold@latest schemagen .
package userapiauthrefresh

import (
	"context"
	"github.com/infinity6-ai/gox/routez/apiz"
	"github.com/infinity6-ai/gox/schemaz/schemaz"
)

type ReqPayload struct {
	// Opaque user refresh token in format <jti>.<secret>.
	RefreshToken string `json:"refresh_token"`
}

type Req struct {
	Payload *ReqPayload `json:"payload,omitempty"`
}

type RespPayload struct {
	// Fresh signed user access token JWT.
	UserAccessToken string `json:"user_access_token"`
	// Rotated opaque user refresh token with new secret.
	UserRefreshToken string `json:"user_refresh_token"`
	// Authenticated user UUID.
	UserId string `json:"user_id"`
	// Access token expiration epoch millis (1 hour).
	ExpiresAt int64 `json:"expires_at"`
	// User refresh token expiration epoch millis (30 days).
	RefreshTokenExpiresAt int64 `json:"refresh_token_expires_at"`
}

type Resp struct {
	Payload *RespPayload `json:"payload,omitempty"`
}

// Human User Session Refresh
// Exchanges a human user refresh token for fresh access and rotated refresh tokens.
type Api struct {
	apiz.UnimplementedApi `json:"-"`
	Req                   *Req  `json:"req,omitempty"`
	Resp                  *Resp `json:"resp,omitempty"`
}

func (d *Api) ApiRoute() apiz.ApiSpec {
	return apiz.ApiSpec{
		Method: "POST",
		Path:   "/api/userz/refresh-user",
		Desc: func() *schemaz.Desc {
			return &schemaz.Desc{
				Labels: map[string]string{"servname": "userz"},
			}
		},
	}
}

type UserAuthRefreshId struct{}

func (UserAuthRefreshId) UserAuthRefresh(ctx context.Context, api *Api) (int, error) {
	panic("UserAuthRefresh not implemented: run 'go generate ./...'")
}
