//go:generate go run github.com/infinity6-ai/i6sdk/mold schemagen .
package userapiauthpasswordreset

import (
	"context"
	"github.com/infinity6-ai/gox/routez/apiz"
	"github.com/infinity6-ai/gox/schemaz/schemaz"
)

type ReqPayload struct {
	// Account email address requesting password reset.
	Email string `json:"email"`
}

type Req struct {
	Payload *ReqPayload `json:"payload,omitempty"`
}

type RespPayload struct {
	// Status of password reset request: password_reset_email_sent.
	Status string `json:"status"`
	// Target account email address.
	Email string `json:"email"`
}

type Resp struct {
	Payload *RespPayload `json:"payload,omitempty"`
}

// Password Reset Initiation
// Initiates an out-of-band password reset email delivery.
type Api struct {
	apiz.UnimplementedApi `json:"-"`
	Req                   *Req  `json:"req,omitempty"`
	Resp                  *Resp `json:"resp,omitempty"`
}

func (d *Api) ApiRoute() apiz.ApiSpec {
	return apiz.ApiSpec{
		Method: "POST",
		Path:   "/api/userz/auth/password/reset",
		Desc: func() *schemaz.Desc {
			return &schemaz.Desc{
				Labels: map[string]string{"servname": "userz"},
			}
		},
	}
}

type UserAuthPasswordResetId struct{}

func (UserAuthPasswordResetId) UserAuthPasswordReset(ctx context.Context, api *Api) (int, error) {
	panic("UserAuthPasswordReset not implemented: run 'go generate ./...'")
}
