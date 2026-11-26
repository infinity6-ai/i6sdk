//go:generate go run github.com/infinity6-ai/i6sdk/mold@latest schemagen .
package userapiauthpasswordsignup

import (
	"context"
	"github.com/infinity6-ai/gox/routez/apiz"
	"github.com/infinity6-ai/gox/schemaz/schemaz"
)

type ReqPayload struct {
	// Desired account email address.
	Email string `json:"email"`
	// Desired account plaintext password.
	Password string `json:"password"`
	// Optional PKCE code_challenge.
	CodeChallenge string `json:"code_challenge,omitempty"`
	// Optional client application identifier.
	App string `json:"app,omitempty"`
}

type Req struct {
	Payload *ReqPayload `json:"payload,omitempty"`
}

type RespPayload struct {
	// Status of user registration: verification_email_sent.
	Status string `json:"status"`
	// Target account email address.
	Email string `json:"email"`
}

type Resp struct {
	Payload *RespPayload `json:"payload,omitempty"`
}

// Email/Password Registration
// Registers a new user account with email and password credentials.
type Api struct {
	apiz.UnimplementedApi `json:"-"`
	Req                   *Req  `json:"req,omitempty"`
	Resp                  *Resp `json:"resp,omitempty"`
}

func (d *Api) ApiRoute() apiz.ApiSpec {
	return apiz.ApiSpec{
		Method: "POST",
		Path:   "/api/userz/auth/signup/password",
		Desc: func() *schemaz.Desc {
			return &schemaz.Desc{
				Labels: map[string]string{"servname": "userz"},
			}
		},
	}
}

type UserAuthPasswordSignupId struct{}

func (UserAuthPasswordSignupId) UserAuthPasswordSignup(ctx context.Context, api *Api) (int, error) {
	panic("UserAuthPasswordSignup not implemented: run 'go generate ./...'")
}
