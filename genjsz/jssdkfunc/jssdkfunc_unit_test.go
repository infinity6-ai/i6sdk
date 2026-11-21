package jssdkfunc_test

import (
	"testing"

	"github.com/infinity6-ai/gox/schemaz/typeastz"
	"github.com/stretchr/testify/require"
	"go.code.infinity6.ai/sdkgo/genjsz/jssdkfunc"
	"go.code.infinity6.ai/sdkgo/mold/moldsdkfunc"
	"go.code.infinity6.ai/sdkgo/sdkgo/productz/relevance/dash/relevancefashiondashfunnel"
	"go.code.infinity6.ai/sdkgo/sdkgo/productz/relevance/dash/relevancefashiondashlistshowcasecategories"
	"go.code.infinity6.ai/sdkgo/sdkgo/productz/relevance/inferz/inferzrelevancefashionfbt"
	"go.code.infinity6.ai/sdkgo/sdkgo/sdkfuncz"
)

func TestUnitWrite(t *testing.T) {
	ctx := t.Context()
	type testScenario struct {
		in   []*moldsdkfunc.SdkFunc
		want string
	}

	check := func(t *testing.T, s testScenario) {
		t.Helper()
		got := jssdkfunc.Write(ctx, s.in)

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
			in:   []*moldsdkfunc.SdkFunc{},
			want: "",
		})
	})

	t.Run("with client infer sdk func", func(t *testing.T) {
		parsed := jssdkfunc.NewSdk().ParseSdkFunc(ctx, &inferzrelevancefashionfbt.Func{})
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
   * Get Bought Together
   * 
   * Infers the SKUs frequently bought together with the given SKUs for a store.
   * 
   * @param {sdkFuncProductRelevanceFashionFbtReq} req
   * @param {(req: sdkFuncProductRelevanceFashionFbtReq, results: sdkFuncProductRelevanceFashionFbtOutput) => Promise<sdkFuncProductRelevanceFashionFbtSelected>} [selectResults] - picks which of the recommended skus are shown and may attach your own data to each in enriched; default: for each result, in order, the first sku not in req.input.skus (the cart) and not already picked, with the rating of that result; no duplicates, and results with no such sku are skipped
   * @param {ServiceOptions} [opts]
   * @returns {Promise<sdkFuncProductRelevanceFashionFbtSelected>}
   */
  async productRelevanceFashionFbtV1(req, selectResults, opts) {
    const infer = this._services.infer()
    const results = await infer.run({
      dataset: req.input.dataset,
      modelName: "relevance-fashion-fbt",
      params: req.input,
      partitions: {"version": "v1", ...req.partitions},
      sessionId: req.input.session_id,
      userId: req.input.user_id,
      serviceOptions: opts,
    })
    return infer.selectResults({
      dataset: req.input.dataset,
      modelName: "relevance-fashion-fbt",
      modelVersion: "v1",
      storeId: req.partitions.store_id,
      req,
      results,
      recId: results.rec_id,
      sessionId: req.input.session_id,
      userId: req.input.user_id,
      selectResults,
      serviceOptions: opts,
    })
  }
}

export { SdkFuncs };
`
		check(t, testScenario{
			in:   []*moldsdkfunc.SdkFunc{parsed},
			want: want,
		})
	})

	t.Run("with gotyp in input and output", func(t *testing.T) {
		fn := &moldsdkfunc.SdkFunc{
			Spec: moldsdkfunc.SdkFuncSpec{
				Id:      "product-custom-task",
				CamelId: "sdkFuncProductCustomTask",
				Version: "v1",
				Typ:     sdkfuncz.TypClientInfer,
			},
			Req: moldsdkfunc.Req{
				Input: typeastz.Typedef{
					Head: &typeastz.Tuple{Typ: "sdkFuncProductCustomTaskInput", GoTyp: "models.CustomInput", Name: "input"},
				},
				Partitions: typeastz.Typedef{
					Head: &typeastz.Tuple{Typ: "sdkFuncProductCustomTaskPartitions", GoTyp: "models.CustomPartitions", Name: "partitions"},
				},
			},
			Resp: moldsdkfunc.Resp{
				Output: typeastz.Typedef{
					Head: &typeastz.Tuple{Typ: "sdkFuncProductCustomTaskOutput", GoTyp: "models.CustomOutput", Name: "output"},
				},
			},
		}

		tds := jssdkfunc.ExtractTypedefs([]*moldsdkfunc.SdkFunc{fn})
		allTds := tds.All()
		require.Len(t, allTds, 4) // req, resp, enriched, selected
		var reqTd *typeastz.Typedef
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
   * @param {(req: sdkFuncProductCustomTaskReq, results: sdkFuncProductCustomTaskOutput) => Promise<sdkFuncProductCustomTaskSelected>} [selectResults] - picks which of the recommended skus are shown and may attach your own data to each in enriched; default: for each result, in order, the first sku not in req.input.skus (the cart) and not already picked, with the rating of that result; no duplicates, and results with no such sku are skipped
   * @param {ServiceOptions} [opts]
   * @returns {Promise<sdkFuncProductCustomTaskSelected>}
   */
  async productCustomTaskV1(req, selectResults, opts) {
    const infer = this._services.infer()
    const results = await infer.run({
      dataset: req.input.dataset,
      modelName: "custom-task",
      params: req.input,
      partitions: {"version": "v1", ...req.partitions},
      sessionId: req.input.session_id,
      userId: req.input.user_id,
      serviceOptions: opts,
    })
    return infer.selectResults({
      dataset: req.input.dataset,
      modelName: "custom-task",
      modelVersion: "v1",
      storeId: req.partitions.store_id,
      req,
      results,
      recId: results.rec_id,
      sessionId: req.input.session_id,
      userId: req.input.user_id,
      selectResults,
      serviceOptions: opts,
    })
  }
}

export { SdkFuncs };
`
		check(t, testScenario{
			in:   []*moldsdkfunc.SdkFunc{fn},
			want: want,
		})
	})

	t.Run("with client query sdk func", func(t *testing.T) {
		parsed := jssdkfunc.NewSdk().ParseSdkFunc(ctx, &relevancefashiondashfunnel.Func{})
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
   * Get Funnel
   * 
   * Gets the event counts of the recommendation funnel for a store and period.
   * 
   * @param {sdkFuncProductRelevanceFashionDashFunnelReq} req
   * @param {ServiceOptions} [opts]
   * @returns {Promise<sdkFuncProductRelevanceFashionDashFunnelOutput>}
   */
  async productRelevanceFashionDashFunnelV1(req, opts) {
    const dash = this._services.dash()
    const ret = await dash.query({
      dataset: req.input.dataset,
      dashName: "product-relevance-fashion-dash-funnel",
      params: req.input,
      partitions: {"version": "v1", ...req.partitions},
      serviceOptions: opts,
    })
    return ret;
  }
}

export { SdkFuncs };
`
		check(t, testScenario{
			in:   []*moldsdkfunc.SdkFunc{parsed},
			want: want,
		})
	})

	t.Run("with client query sdk func without params", func(t *testing.T) {
		parsed := jssdkfunc.NewSdk().ParseSdkFunc(ctx, &relevancefashiondashlistshowcasecategories.Func{})
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
   * List Showcase Categories
   * 
   * Lists the showcase categories of a store.
   * 
   * @param {sdkFuncProductRelevanceFashionDashListShowcaseCategoriesReq} req
   * @param {ServiceOptions} [opts]
   * @returns {Promise<sdkFuncProductRelevanceFashionDashListShowcaseCategoriesOutput>}
   */
  async productRelevanceFashionDashListShowcaseCategoriesV1(req, opts) {
    const dash = this._services.dash()
    const ret = await dash.query({
      dataset: req.input.dataset,
      dashName: "product-relevance-fashion-dash-list-showcase-categories",
      partitions: {"version": "v1", ...req.partitions},
      serviceOptions: opts,
    })
    return ret;
  }
}

export { SdkFuncs };
`
		check(t, testScenario{
			in:   []*moldsdkfunc.SdkFunc{parsed},
			want: want,
		})
	})
}
