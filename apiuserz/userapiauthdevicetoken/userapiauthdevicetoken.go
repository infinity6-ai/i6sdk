//go:generate go run github.com/infinity6-ai/i6sdk/mold@latest schemagen .
package userapiauthdevicetoken

import (
	"context"
	"github.com/infinity6-ai/gox/routez/apiz"
	"github.com/infinity6-ai/gox/schemaz/schemaz"
)

type ReqPayload struct {
	// Secret device code returned by /auth/device/code.
	DeviceCode string `json:"device_code"`
}

type Req struct {
	Payload *ReqPayload `json:"payload,omitempty"`
}

type RespPayload struct {
	// Signed user access token JWT (when authorized).
	UserAccessToken string `json:"user_access_token,omitempty"`
	// Opaque user refresh token session (when authorized).
	UserRefreshToken string `json:"user_refresh_token,omitempty"`
	// Authenticated user UUID.
	UserId string `json:"user_id,omitempty"`
	// Access token expiration epoch millis.
	ExpiresAt int64 `json:"expires_at,omitempty"`
	// User refresh token expiration epoch millis (30 days).
	RefreshTokenExpiresAt int64 `json:"refresh_token_expires_at,omitempty"`
	// Status of the device code: pending | authorized | denied | expired.
	Status string `json:"status"`
}

type Resp struct {
	Payload *RespPayload `json:"payload,omitempty"`
}

// Device Token Polling & Exchange
// Polled by CLI clients to check status and exchange an authorized device code for active user tokens.
type Api struct {
	apiz.UnimplementedApi `json:"-"`
	Req                   *Req  `json:"req,omitempty"`
	Resp                  *Resp `json:"resp,omitempty"`
}

func (d *Api) ApiRoute() apiz.ApiSpec {
	return apiz.ApiSpec{
		Method: "POST",
		Path:   "/api/userz/auth/device/token",
		Desc: func() *schemaz.Desc {
			return &schemaz.Desc{
				Labels: map[string]string{"servname": "userz"},
			}
		},
	}
}

type UserAuthDeviceTokenId struct{}

func (UserAuthDeviceTokenId) UserAuthDeviceToken(ctx context.Context, api *Api) (int, error) {
	panic("UserAuthDeviceToken not implemented: run 'go generate ./...'")
}
