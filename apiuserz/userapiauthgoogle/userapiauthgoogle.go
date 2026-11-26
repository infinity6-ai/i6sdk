//go:generate go run github.com/infinity6-ai/i6sdk/mold schemagen .
package userapiauthgoogle

import (
	"context"
	"github.com/infinity6-ai/gox/routez/apiz"
	"github.com/infinity6-ai/gox/schemaz/schemaz"
)

type ReqPayload struct {
	// Google OpenID Connect ID token credential string.
	IdToken string `json:"id_token"`
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

// Google One-Tap / ID Token Sign-In
// Authenticates a user via verified Google ID token / GIS One-Tap credentials.
type Api struct {
	apiz.UnimplementedApi `json:"-"`
	Req                   *Req  `json:"req,omitempty"`
	Resp                  *Resp `json:"resp,omitempty"`
}

func (d *Api) ApiRoute() apiz.ApiSpec {
	return apiz.ApiSpec{
		Method: "POST",
		Path:   "/api/userz/auth/google",
		Desc: func() *schemaz.Desc {
			return &schemaz.Desc{
				Labels: map[string]string{"servname": "userz"},
			}
		},
	}
}

type UserAuthGoogleId struct{}

func (UserAuthGoogleId) UserAuthGoogle(ctx context.Context, api *Api) (int, error) {
	panic("UserAuthGoogle not implemented: run 'go generate ./...'")
}
