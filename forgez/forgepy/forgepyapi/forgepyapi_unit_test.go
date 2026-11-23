package forgepyapi_test

import (
	"testing"

	"github.com/infinity6-ai/gox/routez/apispecz"
	"github.com/infinity6-ai/gox/routez/apiz"
	"github.com/infinity6-ai/gox/routez/sampleroutez/routezsamplefraction"
	"github.com/infinity6-ai/i6sdk/apicollectionz/apicollectionz"
	"github.com/infinity6-ai/i6sdk/forgez/forgepy/forgepyapi"
	"github.com/stretchr/testify/require"
)

func TestUnitFractionSample(t *testing.T) {
	ctx := t.Context()
	api := routezsamplefraction.NewService().Api()
	actual := forgepyapi.NewSdk().ParseApi(ctx, api)
	require.NotNil(t, actual)
	require.Equal(t, "sample-fraction", actual.Spec.Id)
	require.Equal(t, "apiSampleFraction", actual.Spec.CamelId)
	require.Equal(t, "POST", actual.Spec.Method)
	require.Equal(t, "/api/gox/routez/sample/fraction/{numerator}/{denominator}", actual.Spec.Path)
}

func TestUnitAllApis(t *testing.T) {
	ctx := t.Context()
	apis := apicollectionz.Apis()
	require.NotEmpty(t, apis)
	for _, api := range apis {
		expected := forgepyapi.NewSdk().ParseApi(ctx, api)
		require.NotNil(t, expected)
		require.Equal(t, api.ApiSpec().Id, expected.Spec.Id)
		require.NotEmpty(t, expected.Spec.CamelId)
		require.Equal(t, api.ApiSpec().Method, expected.Spec.Method)
		require.Equal(t, api.ApiSpec().Path, expected.Spec.Path)
	}
}

func TestUnitWrite(t *testing.T) {
	ctx := t.Context()
	type testScenario struct {
		in          *apispecz.Apis
		mustContain []string
		wantEmpty   bool
	}

	check := func(t *testing.T, s testScenario) {
		t.Helper()
		var apisList []*apispecz.Api
		if s.in != nil {
			apisList = s.in.Apis
		}
		got := forgepyapi.Write(ctx, apisList)
		if s.wantEmpty {
			require.Empty(t, got)
			return
		}
		for _, substr := range s.mustContain {
			require.Contains(t, got, substr)
		}
	}

	t.Run("nil apis", func(t *testing.T) {
		check(t, testScenario{
			in:        nil,
			wantEmpty: true,
		})
	})

	t.Run("empty apis", func(t *testing.T) {
		check(t, testScenario{
			in:        &apispecz.Apis{},
			wantEmpty: true,
		})
	})

	t.Run("with multiple path params and extra request fields", func(t *testing.T) {
		apis := forgepyapi.NewSdk().ParseApis(ctx, []apiz.Api{routezsamplefraction.NewService().Api()})
		check(t, testScenario{
			in: apis,
			mustContain: []string{
				"class Apis:",
				"def sample_fraction(",
				"req: apiSampleFractionReq,",
				"opts: Optional[ServiceOptions] = None,",
				") -> Resp[apiSampleFractionResp]:",
				"Sample Fraction",
				"Sample API that performs a fraction operation",
				":param req: apiSampleFractionReq request payload.",
				":return: Resp containing parsed apiSampleFractionResp.",
				"quote(str((req.get('params') or {}).get('numerator', '')))",
				"quote(str((req.get('params') or {}).get('denominator', '')))",
				"\"method\": \"POST\",",
				"\"query\": req.get(\"query\"),",
				"\"headers\": req.get(\"headers\"),",
				"\"json\": req.get(\"payload\"),",
				"sampleFraction = sample_fraction",
			},
		})
	})
}
