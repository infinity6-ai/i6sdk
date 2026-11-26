//go:generate go run github.com/infinity6-ai/i6sdk/mold schemagen .
package userapiauthdevicecode

import (
	"context"
	"github.com/infinity6-ai/gox/routez/apiz"
	"github.com/infinity6-ai/gox/schemaz/schemaz"
)

type ReqPayload struct {
	// CLI client identifier requesting device authorization.
	App string `json:"app,omitempty"`
}

type Req struct {
	Payload *ReqPayload `json:"payload,omitempty"`
}

type RespPayload struct {
	// High-entropy internal device verification code for polling.
	DeviceCode string `json:"device_code"`
	// Human-friendly verification code displayed to the user.
	UserCode string `json:"user_code"`
	// Base browser verification URL.
	VerificationUri string `json:"verification_uri"`
	// Pre-filled browser verification URL including the user_code.
	VerificationUriComplete string `json:"verification_uri_complete"`
	// Lifetime of device and user codes in seconds.
	ExpiresInSec int64 `json:"expires_in_sec"`
	// Minimum polling interval in seconds.
	IntervalSec int `json:"interval_sec"`
}

type Resp struct {
	Payload *RespPayload `json:"payload,omitempty"`
}

// Device Authorization Code Initiation
// Initiates an RFC 8628 device authorization flow for CLI and remote headless SSH sessions.
type Api struct {
	apiz.UnimplementedApi `json:"-"`
	Req                   *Req  `json:"req,omitempty"`
	Resp                  *Resp `json:"resp,omitempty"`
}

func (d *Api) ApiRoute() apiz.ApiSpec {
	return apiz.ApiSpec{
		Method: "POST",
		Path:   "/api/userz/auth/device/code",
		Desc: func() *schemaz.Desc {
			return &schemaz.Desc{
				Labels: map[string]string{"servname": "userz"},
			}
		},
	}
}

type UserAuthDeviceCodeId struct{}

func (UserAuthDeviceCodeId) UserAuthDeviceCode(ctx context.Context, api *Api) (int, error) {
	panic("UserAuthDeviceCode not implemented: run 'go generate ./...'")
}
