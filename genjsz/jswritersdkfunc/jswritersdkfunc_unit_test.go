package jswritersdkfunc_test

import (
	"testing"

	"github.com/stretchr/testify/require"
	"go.code.infinity6.ai/sdkgo/genjsz/jsparsersdkfunc"
	"go.code.infinity6.ai/sdkgo/genjsz/jstypez"
	"go.code.infinity6.ai/sdkgo/genjsz/jswritersdkfunc"
	"go.code.infinity6.ai/sdkgo/sdkgo/productz/relevance/dash/relavancefashiondashfunnel"
	"go.code.infinity6.ai/sdkgo/sdkgo/productz/relevance/inferz/inferzrelevancefashionfbt"
	"go.code.infinity6.ai/sdkgo/sdkgo/sdkfuncz"
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
   * @returns {Promise<sdkFuncProductRelevanceFashionFbtOutput>}
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

	t.Run("with gotyp in input and output", func(t *testing.T) {
		fn := &jsparsersdkfunc.SdkFunc{
			Spec: jsparsersdkfunc.SdkFuncSpec{
				Id:      "product-custom-task",
				CamelId: "sdkFuncProductCustomTask",
				Version: "v1",
				Typ:     sdkfuncz.TypClientInfer,
			},
			Req: jsparsersdkfunc.Req{
				Input: jstypez.Typedef{
					Head: &jstypez.Tuple{Typ: "sdkFuncProductCustomTaskInput", GoTyp: "models.CustomInput", Name: "input"},
				},
				Partitions: jstypez.Typedef{
					Head: &jstypez.Tuple{Typ: "sdkFuncProductCustomTaskPartitions", GoTyp: "models.CustomPartitions", Name: "partitions"},
				},
			},
			Resp: jsparsersdkfunc.Resp{
				Output: jstypez.Typedef{
					Head: &jstypez.Tuple{Typ: "sdkFuncProductCustomTaskOutput", GoTyp: "models.CustomOutput", Name: "output"},
				},
			},
		}

		tds := jswritersdkfunc.ExtractTypedefs([]*jsparsersdkfunc.SdkFunc{fn})
		allTds := tds.All()
		require.Len(t, allTds, 2)
		var reqTd *jstypez.Typedef
		for _, td := range allTds {
			if td.Head.Name == "sdkFuncProductCustomTaskReq" {
				reqTd = td
				break
			}
		}
		require.NotNil(t, reqTd)
		require.NotNil(t, reqTd.Properties)
		props := reqTd.Properties.All()
		require.Len(t, props, 2)
		require.Equal(t, "models.CustomInput", props[0].GoTyp)
		require.Equal(t, "models.CustomPartitions", props[1].GoTyp)

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
   * @param {sdkFuncProductCustomTaskReq} req - [Go: models.CustomInput]
   * @returns {Promise<sdkFuncProductCustomTaskOutput>} - [Go: models.CustomOutput]
   */
  async productCustomTaskV1(req) {
    const infer = this._services.infer()
    const ret = await infer.run({
      dataset: req.input.dataset,
      modelName: "custom-task",
      params: req.input,
      partitions: {"version": "v1", ...req.partitions},
    })
    return ret;
  }
}

export { SdkFuncs };
`
		check(t, testScenario{
			in:   []*jsparsersdkfunc.SdkFunc{fn},
			want: want,
		})
	})

	t.Run("with client query sdk func", func(t *testing.T) {
		parsed := jsparsersdkfunc.ParseSdkFunc(&relavancefashiondashfunnel.Func{})
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
   * @param {sdkFuncProductRelevanceFashionDashFunnelReq} req
   * @returns {Promise<sdkFuncProductRelevanceFashionDashFunnelOutput>}
   */
  async productRelevanceFashionDashFunnelV1(req) {
    const dash = this._services.dash()
    const ret = await dash.query({
      dataset: req.input.dataset,
      dashName: "product-relevance-fashion-dash-funnel",
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
