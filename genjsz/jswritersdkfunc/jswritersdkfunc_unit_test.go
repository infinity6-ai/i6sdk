package jswritersdkfunc_test

import (
	"testing"

	"github.com/stretchr/testify/require"
	"go.code.infinity6.ai/sdkgo/genjsz/jsparsersdkfunc"
	"go.code.infinity6.ai/sdkgo/genjsz/jswritersdkfunc"
	"go.code.infinity6.ai/sdkgo/sdkgo/productz/relevance/inferz/inferzrelevancefashionfbt"
)

func TestUnitWrite(t *testing.T) {
	ctx := t.Context()
	type testScenario struct {
		in   []*jsparsersdkfunc.SdkFunc
		want string
	}

	check := func(t *testing.T, s testScenario) {
		t.Helper()
		got := jswritersdkfunc.Write(ctx, s.in)

		require.Equal(t, s.want, got)
	}

	t.Run("nil sdk funcs", func(t *testing.T) {
		check(t, testScenario{
			in:   nil,
			want: "",
		})
	})

	t.Run("empty sdk funcs", func(t *testing.T) {
		check(t, testScenario{
			in:   []*jsparsersdkfunc.SdkFunc{},
			want: "",
		})
	})

	t.Run("with client infer sdk func", func(t *testing.T) {
		parsed := jsparsersdkfunc.ParseSdkFunc(&inferzrelevancefashionfbt.Func{})
		want := `class SdkFuncs {

  /**
   * Constructs an SdkFuncs instance.
   * @param {Services} services - The services manager instance.
   */
  constructor(services) {
    /**
     * @private
     * @type {Services}
     */
    this._services = services;
  }

  /**
   * @param {sdkFuncProductRelevanceFashionFbtReq} req
   * @returns {Promise<sdkFuncProductRelevanceFashionFbtResp>}
   */
  async productRelevanceFashionFbtV1(req) {
    const infer = this._services.infer()
    const ret = await infer.run({
      dataset: req.input.dataset,
      modelName: "relevance-fashion-fbt",
      params: req.input,
      partitions: {"version": "v1", ...req.partitions},
    })
    return ret;
  }
}

export { SdkFuncs };
`
		check(t, testScenario{
			in:   []*jsparsersdkfunc.SdkFunc{parsed},
			want: want,
		})
	})
}
