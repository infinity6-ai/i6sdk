//go:generate go run github.com/infinity6-ai/i6sdk/mold schemagen .
package userapiauthlogout

import (
	"context"
	"github.com/infinity6-ai/gox/routez/apiz"
	"github.com/infinity6-ai/gox/schemaz/schemaz"
)

type ReqPayload struct {
	// Human user refresh token to revoke.
	RefreshToken string `json:"refresh_token"`
}

type Req struct {
	Payload *ReqPayload `json:"payload,omitempty"`
}

type RespPayload struct {
	// Confirmation flag that refresh token was successfully revoked.
	Success bool `json:"success"`
}

type Resp struct {
	Payload *RespPayload `json:"payload,omitempty"`
}

// User Logout
// Revokes an active human user refresh token and terminates session.
type Api struct {
	apiz.UnimplementedApi `json:"-"`
	Req                   *Req  `json:"req,omitempty"`
	Resp                  *Resp `json:"resp,omitempty"`
}

func (d *Api) ApiRoute() apiz.ApiSpec {
	return apiz.ApiSpec{
		Method: "POST",
		Path:   "/api/userz/auth/logout",
		Desc: func() *schemaz.Desc {
			return &schemaz.Desc{
				Labels: map[string]string{"servname": "userz"},
			}
		},
	}
}

type UserAuthLogoutId struct{}

func (UserAuthLogoutId) UserAuthLogout(ctx context.Context, api *Api) (int, error) {
	panic("UserAuthLogout not implemented: run 'go generate ./...'")
}
