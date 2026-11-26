//go:generate go run github.com/infinity6-ai/i6sdk/mold schemagen .
package userapiauthpasswordsignin

import (
	"context"
	"github.com/infinity6-ai/gox/routez/apiz"
	"github.com/infinity6-ai/gox/schemaz/schemaz"
)

type ReqPayload struct {
	// Account email address.
	Email string `json:"email"`
	// Account plaintext password.
	Password string `json:"password"`
	// PKCE code_challenge for authorization flow.
	CodeChallenge string `json:"code_challenge"`
	// Optional client application identifier.
	App string `json:"app,omitempty"`
}

type Req struct {
	Payload *ReqPayload `json:"payload,omitempty"`
}

type RespPayload struct {
	// Authorization code to be exchanged for tokens.
	Code string `json:"code"`
}

type Resp struct {
	Payload *RespPayload `json:"payload,omitempty"`
}

// Email/Password Sign-In
// Authenticates a user via email and password credentials.
type Api struct {
	apiz.UnimplementedApi `json:"-"`
	Req                   *Req  `json:"req,omitempty"`
	Resp                  *Resp `json:"resp,omitempty"`
}

func (d *Api) ApiRoute() apiz.ApiSpec {
	return apiz.ApiSpec{
		Method: "POST",
		Path:   "/api/userz/auth/signin/password",
		Desc: func() *schemaz.Desc {
			return &schemaz.Desc{
				Labels: map[string]string{"servname": "userz"},
			}
		},
	}
}

type UserAuthPasswordSigninId struct{}

func (UserAuthPasswordSigninId) UserAuthPasswordSignin(ctx context.Context, api *Api) (int, error) {
	panic("UserAuthPasswordSignin not implemented: run 'go generate ./...'")
}
