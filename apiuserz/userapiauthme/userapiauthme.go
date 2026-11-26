//go:generate go run github.com/infinity6-ai/i6sdk/mold schemagen .
package userapiauthme

import (
	"context"
	"github.com/infinity6-ai/gox/routez/apiz"
	"github.com/infinity6-ai/gox/schemaz/schemaz"
)

type ReqHeaders struct {
	// User Bearer JWT access token.
	Authorization string `json:"authorization"`
}

type Req struct {
	Headers *ReqHeaders `json:"headers,omitempty"`
}

type RespPayload struct {
	// User unique identifier.
	UserId string `json:"user_id"`
	// Redacted email address for display.
	EmailRedacted string `json:"email_redacted"`
	// List of linked identity providers.
	Providers []string `json:"providers"`
	// Account creation timestamp in epoch milliseconds.
	CreatedAt int64 `json:"created_at"`
}

type Resp struct {
	Payload *RespPayload `json:"payload,omitempty"`
}

// Current User Profile
// Retrieves current authenticated user identity profile metadata.
type Api struct {
	apiz.UnimplementedApi `json:"-"`
	Req                   *Req  `json:"req,omitempty"`
	Resp                  *Resp `json:"resp,omitempty"`
}

func (d *Api) ApiRoute() apiz.ApiSpec {
	return apiz.ApiSpec{
		Method: "GET",
		Path:   "/api/userz/auth/me",
		Desc: func() *schemaz.Desc {
			return &schemaz.Desc{
				Labels: map[string]string{"servname": "userz"},
			}
		},
	}
}

type UserAuthMeId struct{}

func (UserAuthMeId) UserAuthMe(ctx context.Context, api *Api) (int, error) {
	panic("UserAuthMe not implemented: run 'go generate ./...'")
}
