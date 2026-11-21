package userapidatasettokenrefreshcreate_test

import (
	"testing"

	"github.com/infinity6-ai/gox/commonz/jsonz"
	"github.com/stretchr/testify/require"
	"github.com/infinity6-ai/i6sdk/apiuserz/userapidatasettokenrefreshcreate"
	"github.com/infinity6-ai/i6sdk/sdkcorez/dsz/lakedsz"
	"github.com/infinity6-ai/i6sdk/sdkcorez/dsz/solution"
)

func TestUnitUserDatasetTokenRefreshCreateApi(t *testing.T) {
	api := &userapidatasettokenrefreshcreate.Api{
		Req: userapidatasettokenrefreshcreate.Req{
			Dataset:    lakedsz.MustParse("ds64l7qauv4l6512kf0c"),
			Solution:   solution.MustParse("relevance"),
			Grants:     12,
			Labels:     map[string]string{"env": "prod"},
			Annotation: "Airflow Ingest",
			AuthHeader: "Bearer eyJ...",
		},
		Resp: userapidatasettokenrefreshcreate.Resp{
			Jti:        "jti-pubkey",
			Dataset:    lakedsz.MustParse("ds64l7qauv4l6512kf0c"),
			PrivKey:    "privkey-b64",
			Solution:   solution.MustParse("relevance"),
			Grants:     12,
			Labels:     map[string]string{"env": "prod"},
			ExpiresAt:  1792933262,
			Annotation: "Airflow Ingest",
		},
	}

	apiCopy := jsonz.MustCopy(api, &userapidatasettokenrefreshcreate.Api{})
	require.Equal(t, api.Req, apiCopy.Req)
	require.Equal(t, api.Resp, apiCopy.Resp)

	spec := api.ApiSpec()
	require.Equal(t, "user-dataset-token-refresh-create", spec.Id)
}
