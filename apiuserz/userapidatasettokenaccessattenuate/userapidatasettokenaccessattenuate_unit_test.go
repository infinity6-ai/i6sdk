package userapidatasettokenaccessattenuate_test

import (
	"testing"

	"github.com/infinity6-ai/gox/commonz/jsonz"
	"github.com/infinity6-ai/i6sdk/apiuserz/userapidatasettokenaccessattenuate"
	"github.com/infinity6-ai/i6sdk/sdkcorez/dsz/lakedsz"
	"github.com/infinity6-ai/i6sdk/sdkcorez/dsz/solution"
	"github.com/infinity6-ai/i6sdk/sdkcorez/grants"
	"github.com/stretchr/testify/require"
)

func TestUnitUserDatasetTokenAccessAttenuateApi(t *testing.T) {
	api := &userapidatasettokenaccessattenuate.Api{
		Req: &userapidatasettokenaccessattenuate.Req{
			Headers: &userapidatasettokenaccessattenuate.ReqHeaders{
				Authorization: "Bearer token",
			},
			Payload: &userapidatasettokenaccessattenuate.ReqPayload{
				Subsolution: solution.MustParse("relevance"),
				Subgrants:   grants.Read,
				Newlabels:   map[string]string{"task": "query"},
				ExpSeconds:  1800,
			},
		},
		Resp: &userapidatasettokenaccessattenuate.Resp{
			Payload: &userapidatasettokenaccessattenuate.RespPayload{
				AccessToken: "att_token",
				Dataset:     lakedsz.MustParse("myds"),
				Solution:    solution.MustParse("relevance"),
				Grants:      grants.Read,
				Labels:      map[string]string{"task": "query"},
				ExpiresAt:   123456789,
			},
		},
	}

	apiCopy := jsonz.MustCopy(api, &userapidatasettokenaccessattenuate.Api{})
	require.Equal(t, api.Req.Headers.Authorization, apiCopy.Req.Headers.Authorization)

	spec := api.ApiSpec()
	require.Equal(t, "user-dataset-token-access-attenuate", spec.Id)
	require.Equal(t, "Attenuate Dataset Access Token", spec.Desc().Name)
}
