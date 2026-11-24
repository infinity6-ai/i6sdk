package userapidatasettokenrefresh_test

import (
	"testing"

	"github.com/infinity6-ai/gox/commonz/jsonz"
	"github.com/infinity6-ai/i6sdk/apiuserz/userapidatasettokenrefresh"
	"github.com/infinity6-ai/i6sdk/sdkcorez/dsz/lakedsz"
	"github.com/infinity6-ai/i6sdk/sdkcorez/dsz/solution"
	"github.com/stretchr/testify/require"
)

func TestUnitUserDatasetTokenRefreshApi(t *testing.T) {
	api := &userapidatasettokenrefresh.Api{
		Req: userapidatasettokenrefresh.Req{
			Signature:   "assertion-sig-b64",
			Dataset:     lakedsz.MustParse("ds64l7qauv4l6512kf0c"),
			Subsolution: solution.MustParse("relevance"),
			Subgrants:   8,
			Newlabels:   map[string]string{"task": "export"},
			ExpSeconds:  1800,
			Unbound:     true,
		},
		Resp: userapidatasettokenrefresh.Resp{
			AccessToken: "fresh-access-tok-b64",
			Dataset:     lakedsz.MustParse("ds64l7qauv4l6512kf0c"),
			Solution:    solution.MustParse("relevance"),
			Grants:      8,
			Labels:      map[string]string{"task": "export"},
			ExpiresAt:   1792933262,
		},
	}

	apiCopy := jsonz.MustCopy(api, &userapidatasettokenrefresh.Api{})
	require.Equal(t, api.Req, apiCopy.Req)
	require.Equal(t, api.Resp, apiCopy.Resp)

	spec := api.ApiSpec()
	require.Equal(t, "user-dataset-token-refresh", spec.Id)
}
