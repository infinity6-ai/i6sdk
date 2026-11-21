package jsapi_test

import (
	"os"
	"testing"

	"github.com/infinity6-ai/gox/commonz/errorz"
	"github.com/infinity6-ai/gox/commonz/jsonz"
	"github.com/infinity6-ai/gox/routez/apiz"
	"github.com/infinity6-ai/gox/routez/jsapiz"
	"github.com/infinity6-ai/gox/routez/sampleroutez/routezsamplefraction"
	"github.com/infinity6-ai/gox/schemaz/jstypez"
	"github.com/infinity6-ai/gox/schemaz/schemaz"
	"github.com/stretchr/testify/require"
	"go.code.infinity6.ai/sdkgo/genjsz/jsapi"
	"go.code.infinity6.ai/sdkgo/sdkgo/sdkz"
)

func TestUnitFractionSample(t *testing.T) {
	ctx := t.Context()
	api := routezsamplefraction.NewService().Api()
	actual := jsapi.NewSdk().ParseApi(ctx, api)
	jsonz.MustParse(jsonz.MustFormat(actual).String(), &actual)
	println(jsonz.MustFormatPretty(actual).String())
	expected := &jsapiz.Api{
		Spec: jsapiz.ApiSpec{
			Id:      "sample-fraction",
			CamelId: "apiSampleFraction",
			Desc: &schemaz.Desc{
				Name:     "Sample Fraction",
				Summary:  "Sample API that performs a fraction operation",
				Markdown: "",
			},
			Method: "POST",
			Path:   "/api/gox/routez/sample/fraction/{numerator}/{denominator}",
		},
		Typedefs: jstypez.NewTypedefs([]*jstypez.Typedef{
			{
				Head: &jstypez.Tuple{Typ: "object", Name: "apiSampleFractionParams", Summary: "apiSampleFraction params"},
				Properties: jstypez.NewTuples([]jstypez.Tuple{
					{Typ: "number", GoTyp: "float64", Name: "denominator", Summary: "Fraction Denominator"},
					{Typ: "number", GoTyp: "float64", Name: "numerator", Summary: "Fraction Numerator"},
				}),
			},
			{
				Head: &jstypez.Tuple{Typ: "object", Name: "apiSampleFractionQuery", Summary: "apiSampleFraction query"},
				Properties: jstypez.NewTuples([]jstypez.Tuple{
					{Typ: "number", GoTyp: "int", Name: "precision", Summary: "Fraction max precision"},
					{Typ: "apiSampleFractionQuerySmsg", Name: "smsg"},
				}),
			},
			{
				Head: &jstypez.Tuple{Typ: "object", Name: "apiSampleFractionQuerySmsg"},
				Properties: jstypez.NewTuples([]jstypez.Tuple{
					{Typ: "string", GoTyp: "string", Name: "message", Summary: "message string"},
				}),
			},
			{
				Head: &jstypez.Tuple{Typ: "object", Name: "apiSampleFractionReqHeaders", Summary: "apiSampleFraction reqHeaders"},
				Properties: jstypez.NewTuples([]jstypez.Tuple{
					{Typ: "string", GoTyp: "string", Name: "x_i6_trace_id", Summary: "Sample of Trace Id"},
				}),
			},
			{
				Head: &jstypez.Tuple{Typ: "object", Name: "apiSampleFractionReqPayload", Summary: "apiSampleFraction reqPayload"},
				Properties: jstypez.NewTuples([]jstypez.Tuple{
					{Typ: "Array<apiSampleFractionReqPayloadMessagesElement>", Name: "messages"},
					{Typ: "string", GoTyp: "string", Name: "reason", Summary: "Sample reason of this fraction"},
				}),
			},
			{
				Head: &jstypez.Tuple{Typ: "object", Name: "apiSampleFractionReqPayloadMessagesElement", Summary: "message"},
				Properties: jstypez.NewTuples([]jstypez.Tuple{
					{Typ: "string", GoTyp: "string", Name: "message", Summary: "message string"},
				}),
			},
			{
				Head: &jstypez.Tuple{Typ: "object", Name: "apiSampleFractionRespHeaders", Summary: "apiSampleFraction respHeaders"},
				Properties: jstypez.NewTuples([]jstypez.Tuple{
					{Typ: "apiSampleFractionRespHeadersXI6Smsg", Name: "x_i6_smsg"},
					{Typ: "string", GoTyp: "string", Name: "x_i6_trace_message", Summary: "Sample of any trace message"},
				}),
			},
			{
				Head: &jstypez.Tuple{Typ: "object", Name: "apiSampleFractionRespHeadersXI6Smsg"},
				Properties: jstypez.NewTuples([]jstypez.Tuple{
					{Typ: "string", GoTyp: "string", Name: "message", Summary: "message string"},
				}),
			},
			{
				Head: &jstypez.Tuple{Typ: "object", Name: "apiSampleFractionRespPayload", Summary: "apiSampleFraction respPayload"},
				Properties: jstypez.NewTuples([]jstypez.Tuple{
					{Typ: "string", GoTyp: "string", Name: "display", Summary: "Human Fraction Representation"},
					{Typ: "Array<apiSampleFractionRespPayloadMessagesElement>", Name: "messages", Summary: "Output Messages"},
					{Typ: "string", GoTyp: "string", Name: "result", Summary: "Fraction Result"},
				}),
			},
			{
				Head: &jstypez.Tuple{Typ: "string", GoTyp: "string", Name: "apiSampleFractionRespPayloadMessagesElement"},
			},
			{
				Head: &jstypez.Tuple{Typ: "object", Name: "dszFileStat"},
				Properties: jstypez.NewTuples([]jstypez.Tuple{
					{Typ: "string", Name: "url", Summary: "missing summary"},
					{Typ: "string", Name: "content_type", Summary: "missing summary"},
					{Typ: "string", Name: "md5", Summary: "missing summary"},
					{Typ: "string", Name: "size", Summary: "missing summary"},
					{Typ: "string", Name: "etag", Summary: "missing summary"},
					{Typ: "dszTime", Name: "created_at", Summary: "missing summary"},
					{Typ: "dszTime", Name: "updated_at", Summary: "missing summary"},
				}),
			},
			{
				Head: &jstypez.Tuple{Typ: "string", Name: "dszTime", Summary: "Missing Summary"},
			},
		}),
		Req: jsapiz.Req{
			Params: jstypez.Typedef{
				Head:       &jstypez.Tuple{Typ: "apiSampleFractionParams", Name: "params"},
				Properties: jstypez.NewTuples(nil),
			},
			Query: jstypez.Typedef{
				Head:       &jstypez.Tuple{Typ: "apiSampleFractionQuery", Name: "query"},
				Properties: jstypez.NewTuples(nil),
			},
			Headers: jstypez.Typedef{
				Head:       &jstypez.Tuple{Typ: "apiSampleFractionReqHeaders", Name: "reqHeaders"},
				Properties: jstypez.NewTuples(nil),
			},
			Payload: jstypez.Typedef{
				Head:       &jstypez.Tuple{Typ: "apiSampleFractionReqPayload", Name: "reqPayload"},
				Properties: jstypez.NewTuples(nil),
			},
		},
		Resp: jsapiz.Resp{
			Headers: jstypez.Typedef{
				Head:       &jstypez.Tuple{Typ: "apiSampleFractionRespHeaders", Name: "respHeaders"},
				Properties: jstypez.NewTuples(nil),
			},
			Payload: jstypez.Typedef{
				Head:       &jstypez.Tuple{Typ: "apiSampleFractionRespPayload", Name: "respPayload"},
				Properties: jstypez.NewTuples(nil),
			},
		},
	}
	require.Equal(t, expected, actual)
}

func TestUnitAllApis(t *testing.T) {
	ctx := t.Context()
	apis := sdkz.Apis()
	for _, api := range apis {
		expected := jsapi.NewSdk().ParseApi(ctx, api)
		var actual *jsapiz.Api
		jsonz.MustParse(jsonz.MustFormat(expected).String(), &actual)
		errorz.Check(os.MkdirAll("target/jsapi/api", os.ModePerm))
		errorz.Check(os.WriteFile("target/jsapi/api/"+api.ApiSpec().Id+".json", jsonz.MustFormatPretty(expected).Bytes(), os.ModePerm))
	}
}

func TestUnitWrite(t *testing.T) {
	ctx := t.Context()
	type testScenario struct {
		in   *jsapiz.Apis
		want string
	}

	check := func(t *testing.T, s testScenario) {
		t.Helper()
		var apisList []*jsapiz.Api
		if s.in != nil {
			apisList = s.in.Apis
		}
		got := jsapi.Write(ctx, apisList)
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
			in:   &jsapiz.Apis{},
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
		api := &jsapiz.Api{
			Spec: jsapiz.ApiSpec{
				CamelId: "apiCustomAction",
				Method:  "POST",
				Path:    "/api/custom/action",
			},
			Req: jsapiz.Req{
				Payload: jstypez.Typedef{
					Head: &jstypez.Tuple{Typ: "apiCustomActionReqPayload", GoTyp: "models.CustomReq", Name: "payload"},
				},
			},
			Resp: jsapiz.Resp{
				Payload: jstypez.Typedef{
					Head: &jstypez.Tuple{Typ: "apiCustomActionRespPayload", GoTyp: "models.CustomResp", Name: "payload"},
				},
			},
		}

		tds := jsapiz.ExtractTypedefs([]*jsapiz.Api{api})
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

		apis := &jsapiz.Apis{Apis: []*jsapiz.Api{api}}
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
		api := &jsapiz.Api{
			Spec: jsapiz.ApiSpec{
				CamelId: "apiCustomAction",
				Method:  "POST",
				Path:    "/api/custom/{id}",
				Desc:    &schemaz.Desc{Name: "Custom", Summary: "Custom action", Labels: map[string]string{"servname": "tfinfer"}},
			},
		}
		got := jsapi.Write(ctx, []*jsapiz.Api{api})
		require.Contains(t, got, "url: `${this._services.config().servUrls?.[\"tfinfer\"] ?? this._services.config().baseUrl}/api/custom/{id}`,")
		require.NotContains(t, got, "path:")
	})
}
