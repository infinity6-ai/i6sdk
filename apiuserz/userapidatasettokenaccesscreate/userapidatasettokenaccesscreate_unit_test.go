package userapidatasettokenaccesscreate_test

import (
	"testing"

	"github.com/infinity6-ai/gox/commonz/jsonz"
	"github.com/infinity6-ai/i6sdk/apiuserz/userapidatasettokenaccesscreate"
	"github.com/infinity6-ai/i6sdk/sdkcorez/dsz/lakedsz"
	"github.com/infinity6-ai/i6sdk/sdkcorez/dsz/solution"
	"github.com/stretchr/testify/require"
)

func TestUnitUserDatasetTokenAccessCreateApi(t *testing.T) {
	api := &userapidatasettokenaccesscreate.Api{
		Req: userapidatasettokenaccesscreate.Req{
			Dataset:    lakedsz.MustParse("ds64l7qauv4l6512kf0c"),
			Solution:   solution.MustParse("relevance"),
			Grants:     12,
			Labels:     map[string]string{"env": "prod"},
			ExpSeconds: 3600,
			AuthHeader: "Bearer eyJ...",
		},
		Resp: userapidatasettokenaccesscreate.Resp{
			AccessToken: "tok-abc",
			Dataset:     lakedsz.MustParse("ds64l7qauv4l6512kf0c"),
			Solution:    solution.MustParse("relevance"),
			Grants:      12,
			Labels:      map[string]string{"env": "prod"},
			ExpiresAt:   1792933262,
		},
	}

	apiCopy := jsonz.MustCopy(api, &userapidatasettokenaccesscreate.Api{})
	require.Equal(t, api.Req, apiCopy.Req)
	require.Equal(t, api.Resp, apiCopy.Resp)

	spec := api.ApiSpec()
	require.Equal(t, "user-dataset-token-access-create", spec.Id)
}
