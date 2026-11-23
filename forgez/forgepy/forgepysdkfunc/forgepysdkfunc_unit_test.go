package forgepysdkfunc_test

import (
	"testing"

	"github.com/infinity6-ai/i6sdk/apicollectionz/apicollectionz"
	"github.com/infinity6-ai/i6sdk/forgez/forgepy/forgepysdkfunc"
	"github.com/infinity6-ai/i6sdk/mold/moldsdkfunc"
	"github.com/infinity6-ai/i6sdk/sdkcorez/productz/relevance/dash/relevancefashiondashfunnel"
	"github.com/infinity6-ai/i6sdk/sdkcorez/productz/relevance/inferz/inferzrelevancefashionfbt"
	"github.com/stretchr/testify/require"
)

func TestUnitWrite(t *testing.T) {
	ctx := t.Context()
	type testScenario struct {
		in          []*moldsdkfunc.SdkFunc
		mustContain []string
		wantEmpty   bool
	}

	check := func(t *testing.T, s testScenario) {
		t.Helper()
		got := forgepysdkfunc.Write(ctx, s.in)
		if s.wantEmpty {
			require.Empty(t, got)
			return
		}
		for _, substr := range s.mustContain {
			require.Contains(t, got, substr)
		}
	}

	t.Run("nil sdk funcs", func(t *testing.T) {
		check(t, testScenario{
			in:        nil,
			wantEmpty: true,
		})
	})

	t.Run("empty sdk funcs", func(t *testing.T) {
		check(t, testScenario{
			in:        []*moldsdkfunc.SdkFunc{},
			wantEmpty: true,
		})
	})

	t.Run("with client infer sdk func", func(t *testing.T) {
		parsed := forgepysdkfunc.NewSdk().ParseSdkFunc(ctx, &inferzrelevancefashionfbt.Func{})
		check(t, testScenario{
			in: []*moldsdkfunc.SdkFunc{parsed},
			mustContain: []string{
				"class SdkFuncs:",
				"def product_relevance_fashion_fbt_v1(",
				"req: sdkFuncProductRelevanceFashionFbtReq,",
				"select_results: Optional[Callable[[sdkFuncProductRelevanceFashionFbtReq, sdkFuncProductRelevanceFashionFbtOutput], sdkFuncProductRelevanceFashionFbtSelected]] = None,",
				") -> sdkFuncProductRelevanceFashionFbtSelected:",
				"Get Bought Together",
				"\"model_name\": \"relevance-fashion-fbt\",",
				"partitions[\"version\"] = \"v1\"",
				"infer.select_results({",
				"productRelevanceFashionFbtV1 = product_relevance_fashion_fbt_v1",
			},
		})
	})

	t.Run("with client query sdk func", func(t *testing.T) {
		parsed := forgepysdkfunc.NewSdk().ParseSdkFunc(ctx, &relevancefashiondashfunnel.Func{})
		check(t, testScenario{
			in: []*moldsdkfunc.SdkFunc{parsed},
			mustContain: []string{
				"class SdkFuncs:",
				"def product_relevance_fashion_dash_funnel_v1(",
				"req: sdkFuncProductRelevanceFashionDashFunnelReq,",
				") -> sdkFuncProductRelevanceFashionDashFunnelOutput:",
				"dash = self._services.dash()",
				"\"dash_name\": \"product-relevance-fashion-dash-funnel\",",
				"productRelevanceFashionDashFunnelV1 = product_relevance_fashion_dash_funnel_v1",
			},
		})
	})
}

func TestUnitExtractTypedefs(t *testing.T) {
	ctx := t.Context()
	parsed := forgepysdkfunc.NewSdk().ParseSdkFunc(ctx, &inferzrelevancefashionfbt.Func{})
	tds := forgepysdkfunc.ExtractTypedefs([]*moldsdkfunc.SdkFunc{parsed})
	allTds := tds.All()
	require.NotEmpty(t, allTds)

	var hasReq, hasResp, hasSelected, hasEnriched bool
	for _, td := range allTds {
		if td.Head.Name == "sdkFuncProductRelevanceFashionFbtReq" {
			hasReq = true
		}
		if td.Head.Name == "sdkFuncProductRelevanceFashionFbtResp" {
			hasResp = true
		}
		if td.Head.Name == "sdkFuncProductRelevanceFashionFbtSelected" {
			hasSelected = true
		}
		if td.Head.Name == "sdkFuncProductRelevanceFashionFbtOutputResultsElementSkusElementEnriched" {
			hasEnriched = true
		}
	}
	require.True(t, hasReq, "must contain req typedef")
	require.True(t, hasResp, "must contain resp typedef")
	require.True(t, hasSelected, "must contain selected typedef")
	require.True(t, hasEnriched, "must contain enriched typedef")
}

func TestUnitAllSdkFuncs(t *testing.T) {
	ctx := t.Context()
	funcs := apicollectionz.SdkFuncs()
	require.NotEmpty(t, funcs)
	for _, fn := range funcs {
		parsed := forgepysdkfunc.NewSdk().ParseSdkFunc(ctx, fn)
		require.NotNil(t, parsed)
		require.Equal(t, fn.SdkFuncSpec().Id, parsed.Spec.Id)
		require.NotEmpty(t, parsed.Spec.CamelId)
	}
}
