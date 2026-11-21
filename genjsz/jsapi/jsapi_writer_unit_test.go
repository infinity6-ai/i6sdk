package jsapi_test

import (
	"testing"

	"github.com/infinity6-ai/gox/routez/apiz"
	"github.com/infinity6-ai/gox/routez/sampleroutez/routezsamplefraction"
	"github.com/infinity6-ai/gox/schemaz/schemaz"
	"github.com/stretchr/testify/require"
	"go.code.infinity6.ai/sdkgo/genjsz/jsapi"
	"go.code.infinity6.ai/sdkgo/genjsz/jstypez"
)

func TestUnitWrite(t *testing.T) {
	ctx := t.Context()
	type testScenario struct {
		in   *jsapi.Apis
		want string
	}

	check := func(t *testing.T, s testScenario) {
		t.Helper()
		var apisList []*jsapi.Api
		if s.in != nil {
			apisList = s.in.Apis
		}
		got := jsapi.Write(ctx, apisList)
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
			in:   &jsapi.Apis{},
			want: ``,
		})
	})

	t.Run("with multiple path params and extra request fields", func(t *testing.T) {
		apis := jsapi.NewSdk().ParseApis(ctx, []apiz.Api{routezsamplefraction.NewService().Api()})
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
   * @returns {Promise<import("../invoker/invoker.js").Resp<apiSampleFractionResp>>}
   */
  async sampleFraction(req, opts) {
    return this._services.invoker().invokeResult({
      method: "POST",
      path: ` + "`" + `api/gox/routez/sample/fraction/${encodeURIComponent(String(req.params.numerator))}/${encodeURIComponent(String(req.params.denominator))}` + "`" + `,
      query: req.query,
      headers: req.headers,
      json: req.payload,
      opts: opts,
    });
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
		api := &jsapi.Api{
			Spec: jsapi.ApiSpec{
				CamelId: "apiCustomAction",
				Method:  "POST",
				Path:    "/api/custom/action",
			},
			Req: jsapi.Req{
				Payload: jstypez.Typedef{
					Head: &jstypez.Tuple{Typ: "apiCustomActionReqPayload", GoTyp: "models.CustomReq", Name: "payload"},
				},
			},
			Resp: jsapi.Resp{
				Payload: jstypez.Typedef{
					Head: &jstypez.Tuple{Typ: "apiCustomActionRespPayload", GoTyp: "models.CustomResp", Name: "payload"},
				},
			},
		}

		tds := jsapi.ExtractTypedefs([]*jsapi.Api{api})
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

		apis := &jsapi.Apis{Apis: []*jsapi.Api{api}}
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
   * @returns {Promise<import("../invoker/invoker.js").Resp<apiCustomActionResp>>} - [Go: models.CustomResp]
   */
  async customAction(req, opts) {
    return this._services.invoker().invokeResult({
      method: "POST",
      path: ` + "`" + `api/custom/action` + "`" + `,
      json: req.payload,
      opts: opts,
    });
  }
}

export { Apis };
`
		check(t, testScenario{
			in:   apis,
			want: want,
		})
	})

	t.Run("servname label uses the service url", func(t *testing.T) {
		api := &jsapi.Api{
			Spec: jsapi.ApiSpec{
				CamelId: "apiCustomAction",
				Method:  "POST",
				Path:    "/api/custom/{id}",
				Desc:    &schemaz.Desc{Name: "Custom", Summary: "Custom action", Labels: map[string]string{"servname": "tfinfer"}},
			},
		}
		got := jsapi.Write(ctx, []*jsapi.Api{api})
		require.Contains(t, got, "url: `${this._services.config().servUrls?.[\"tfinfer\"] ?? this._services.config().baseUrl}/api/custom/{id}`,")
		require.NotContains(t, got, "path:")
	})
}
