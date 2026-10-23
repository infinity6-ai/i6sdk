package jswriterapi_test

import (
	"testing"

	"github.com/infinity6-ai/gox/routez/apiz"
	"github.com/infinity6-ai/gox/routez/sampleroutez/routezsamplefraction"
	"github.com/stretchr/testify/require"
	"go.code.infinity6.ai/sdkgo/genjsz/jsparserapi"
	"go.code.infinity6.ai/sdkgo/genjsz/jstypez"
	"go.code.infinity6.ai/sdkgo/genjsz/jswriterapi"
)

func TestUnitWrite(t *testing.T) {
	ctx := t.Context()
	type testScenario struct {
		in   *jsparserapi.Apis
		want string
	}

	check := func(t *testing.T, s testScenario) {
		t.Helper()
		var apisList []*jsparserapi.Api
		if s.in != nil {
			apisList = s.in.Apis
		}
		got := jswriterapi.WriteV2(ctx, apisList)
		println(got)

		require.Equal(t, s.want, got)
	}

	t.Run("nil apis", func(t *testing.T) {
		check(t, testScenario{
			in:   nil,
			want: "",
		})
	})

	t.Run("empty apis", func(t *testing.T) {
		check(t, testScenario{
			in:   &jsparserapi.Apis{},
			want: ``,
		})
	})

	t.Run("with multiple path params and extra request fields", func(t *testing.T) {
		apis := jsparserapi.ParseApis(ctx, []apiz.Api{routezsamplefraction.NewService().Api()})
		want := `class Apis {

  /**
   * Constructs an Ingest client instance.
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
   * Sample Fraction
   * 
   * Sample API that performs a fraction operation
   * 
   * @param {apiSampleFractionReq} req
   * @param {ServiceOptions} [opts]
   * @returns {Promise<import("../errors.js").I6Result<apiSampleFractionResp>>}
   */
  async sampleFraction(req, opts) {
    return this._services.invoker().invokeResult({
      method: "POST",
      path: ` + "`" + `api/gox/routez/sample/fraction/${encodeURIComponent(String(req.params.numerator))}/${encodeURIComponent(String(req.params.denominator))}` + "`" + `,
      query: req.query,
      headers: req.headers,
      json: req.payload,
      opts: opts,
    }, "payload");
  }
}

export { Apis };
`
		check(t, testScenario{
			in:   apis,
			want: want,
		})
	})

	t.Run("with gotyp in request and response payload", func(t *testing.T) {
		api := &jsparserapi.Api{
			Spec: jsparserapi.ApiSpec{
				CamelId: "apiCustomAction",
				Method:  "POST",
				Path:    "/api/custom/action",
			},
			Req: jsparserapi.Req{
				Payload: jstypez.Typedef{
					Head: &jstypez.Tuple{Typ: "apiCustomActionReqPayload", GoTyp: "models.CustomReq", Name: "payload"},
				},
			},
			Resp: jsparserapi.Resp{
				Payload: jstypez.Typedef{
					Head: &jstypez.Tuple{Typ: "apiCustomActionRespPayload", GoTyp: "models.CustomResp", Name: "payload"},
				},
			},
		}

		tds := jswriterapi.ExtractTypedefs([]*jsparserapi.Api{api})
		allTds := tds.All()
		require.Len(t, allTds, 2)
		var reqTd *jstypez.Typedef
		for _, td := range allTds {
			if td.Head.Name == "apiCustomActionReq" {
				reqTd = td
				break
			}
		}
		require.NotNil(t, reqTd)
		require.NotNil(t, reqTd.Properties)
		props := reqTd.Properties.All()
		require.Len(t, props, 1)
		require.Equal(t, "models.CustomReq", props[0].GoTyp)

		apis := &jsparserapi.Apis{Apis: []*jsparserapi.Api{api}}
		want := `class Apis {

  /**
   * Constructs an Ingest client instance.
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
   * 
   * 
   * 
   * 
   * @param {apiCustomActionReq} req - [Go: models.CustomReq]
   * @param {ServiceOptions} [opts]
   * @returns {Promise<import("../errors.js").I6Result<apiCustomActionResp>>} - [Go: models.CustomResp]
   */
  async customAction(req, opts) {
    return this._services.invoker().invokeResult({
      method: "POST",
      path: ` + "`" + `api/custom/action` + "`" + `,
      json: req.payload,
      opts: opts,
    }, "payload");
  }
}

export { Apis };
`
		check(t, testScenario{
			in:   apis,
			want: want,
		})
	})
}
