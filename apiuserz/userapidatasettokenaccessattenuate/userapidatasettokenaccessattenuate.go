//go:generate go run github.com/infinity6-ai/i6sdk/mold@latest schemagen .
package userapidatasettokenaccessattenuate

import (
	"context"
	"github.com/infinity6-ai/gox/routez/apiz"
	"github.com/infinity6-ai/gox/schemaz/schemaz"
	"github.com/infinity6-ai/i6sdk/sdkcorez/dsz/lakedsz"
	"github.com/infinity6-ai/i6sdk/sdkcorez/dsz/solution"
	"github.com/infinity6-ai/i6sdk/sdkcorez/grants"
)

type ReqHeaders struct {
	// Parent Bearer dataset access token.
	Authorization string `json:"authorization"`
}

type ReqPayload struct {
	// Narrowed solution domain.
	Subsolution solution.Solution `json:"subsolution,omitempty"`
	// Reduced uint8 grant bitmask.
	Subgrants grants.Grant `json:"subgrants,omitempty"`
	// Additional restrictive label constraints.
	Newlabels map[string]string `json:"newlabels,omitempty"`
	// Requested lifetime in seconds (capped at parent exp).
	ExpSeconds int64 `json:"exp_seconds,omitempty"`
}

type Req struct {
	Headers *ReqHeaders `json:"headers,omitempty"`
	Payload *ReqPayload `json:"payload,omitempty"`
}

type RespPayload struct {
	// Attenuated Base64URL dataset access token.
	AccessToken string `json:"access_token"`
	// Target dataset ID.
	Dataset lakedsz.Dataset `json:"dataset"`
	// Effective solution domain.
	Solution solution.Solution `json:"solution"`
	// Effective uint8 grant bitmask.
	Grants grants.Grant `json:"grants"`
	// Effective metadata labels.
	Labels map[string]string `json:"labels,omitempty"`
	// Expiration epoch milliseconds.
	ExpiresAt int64 `json:"expires_at"`
}

type Resp struct {
	Payload *RespPayload `json:"payload,omitempty"`
}

// Attenuate Dataset Access Token
// Derives an in-flight attenuated child access token with dropped privileges from an active access token.
type Api struct {
	apiz.UnimplementedApi `json:"-"`
	Req                   *Req  `json:"req,omitempty"`
	Resp                  *Resp `json:"resp,omitempty"`
}

func (d *Api) ApiRoute() apiz.ApiSpec {
	return apiz.ApiSpec{
		Method: "POST",
		Path:   "/api/userz/attenuate",
		Desc: func() *schemaz.Desc {
			return &schemaz.Desc{
				Labels: map[string]string{"servname": "userz"},
			}
		},
	}
}

type UserDatasetTokenAccessAttenuateId struct{}

func (UserDatasetTokenAccessAttenuateId) UserDatasetTokenAccessAttenuate(ctx context.Context, api *Api) (int, error) {
	panic("UserDatasetTokenAccessAttenuate not implemented: run 'go generate ./...'")
}
