package forgejsapi_test

import (
	"os"
	"testing"

	"github.com/infinity6-ai/gox/commonz/errorz"
	"github.com/infinity6-ai/gox/commonz/jsonz"
	"github.com/infinity6-ai/gox/routez/apispecz"
	"github.com/infinity6-ai/gox/routez/apiz"
	"github.com/infinity6-ai/gox/routez/sampleroutez/routezsamplefraction"
	"github.com/infinity6-ai/gox/schemaz/schemaz"
	"github.com/infinity6-ai/gox/schemaz/typeastz"
	"github.com/infinity6-ai/i6sdk/apicollectionz/apicollectionz"
	"github.com/infinity6-ai/i6sdk/apidatasetz/dszapimodelget"
	"github.com/infinity6-ai/i6sdk/apidatasetz/dszapistripecheckout"
	"github.com/infinity6-ai/i6sdk/forgez/forgejs/forgejsapi"
	"github.com/stretchr/testify/require"
)

func TestUnitFractionSample(t *testing.T) {
	ctx := t.Context()
	api := routezsamplefraction.NewService().Api()
	actual := forgejsapi.NewSdk().ParseApi(ctx, api)
	jsonz.MustParse(jsonz.MustFormat(actual).String(), &actual)
	println(jsonz.MustFormatPretty(actual).String())
	expected := &apispecz.Api{
		Spec: apispecz.ApiSpec{
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
		Typedefs: typeastz.NewTypedefs([]*typeastz.Typedef{
			{
				Head: &typeastz.Tuple{Typ: "object", Name: "apiSampleFractionParams", Summary: "apiSampleFraction params"},
				Properties: typeastz.NewTuples([]typeastz.Tuple{
					{Typ: "number", GoTyp: "float64", Name: "denominator", Summary: "Fraction Denominator"},
					{Typ: "number", GoTyp: "float64", Name: "numerator", Summary: "Fraction Numerator"},
				}),
			},
			{
				Head: &typeastz.Tuple{Typ: "object", Name: "apiSampleFractionQuery", Summary: "apiSampleFraction query"},
				Properties: typeastz.NewTuples([]typeastz.Tuple{
					{Typ: "number", GoTyp: "int", Name: "precision", Summary: "Fraction max precision"},
					{Typ: "apiSampleFractionQuerySmsg", Name: "smsg"},
				}),
			},
			{
				Head: &typeastz.Tuple{Typ: "object", Name: "apiSampleFractionQuerySmsg"},
				Properties: typeastz.NewTuples([]typeastz.Tuple{
					{Typ: "string", GoTyp: "string", Name: "message", Summary: "message string"},
				}),
			},
			{
				Head: &typeastz.Tuple{Typ: "object", Name: "apiSampleFractionReqHeaders", Summary: "apiSampleFraction reqHeaders"},
				Properties: typeastz.NewTuples([]typeastz.Tuple{
					{Typ: "string", GoTyp: "string", Name: "x_i6_trace_id", Summary: "Sample of Trace Id"},
				}),
			},
			{
				Head: &typeastz.Tuple{Typ: "object", Name: "apiSampleFractionReqPayload", Summary: "apiSampleFraction reqPayload"},
				Properties: typeastz.NewTuples([]typeastz.Tuple{
					{Typ: "Array<apiSampleFractionReqPayloadMessagesElement>", Name: "messages"},
					{Typ: "string", GoTyp: "string", Name: "reason", Summary: "Sample reason of this fraction"},
				}),
			},
			{
				Head: &typeastz.Tuple{Typ: "object", Name: "apiSampleFractionReqPayloadMessagesElement", Summary: "message"},
				Properties: typeastz.NewTuples([]typeastz.Tuple{
					{Typ: "string", GoTyp: "string", Name: "message", Summary: "message string"},
				}),
			},
			{
				Head: &typeastz.Tuple{Typ: "object", Name: "apiSampleFractionRespHeaders", Summary: "apiSampleFraction respHeaders"},
				Properties: typeastz.NewTuples([]typeastz.Tuple{
					{Typ: "apiSampleFractionRespHeadersXI6Smsg", Name: "x_i6_smsg"},
					{Typ: "string", GoTyp: "string", Name: "x_i6_trace_message", Summary: "Sample of any trace message"},
				}),
			},
			{
				Head: &typeastz.Tuple{Typ: "object", Name: "apiSampleFractionRespHeadersXI6Smsg"},
				Properties: typeastz.NewTuples([]typeastz.Tuple{
					{Typ: "string", GoTyp: "string", Name: "message", Summary: "message string"},
				}),
			},
			{
				Head: &typeastz.Tuple{Typ: "object", Name: "apiSampleFractionRespPayload", Summary: "apiSampleFraction respPayload"},
				Properties: typeastz.NewTuples([]typeastz.Tuple{
					{Typ: "string", GoTyp: "string", Name: "display", Summary: "Human Fraction Representation"},
					{Typ: "Array<apiSampleFractionRespPayloadMessagesElement>", Name: "messages", Summary: "Output Messages"},
					{Typ: "string", GoTyp: "string", Name: "result", Summary: "Fraction Result"},
				}),
			},
			{
				Head: &typeastz.Tuple{Typ: "string", GoTyp: "string", Name: "apiSampleFractionRespPayloadMessagesElement"},
			},
			{
				Head: &typeastz.Tuple{Typ: "object", Name: "dszFileStat"},
				Properties: typeastz.NewTuples([]typeastz.Tuple{
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
				Head: &typeastz.Tuple{Typ: "string", Name: "dszTime", Summary: "Missing Summary"},
			},
		}),
		Req: apispecz.Req{
			Params: typeastz.Typedef{
				Head:       &typeastz.Tuple{Typ: "apiSampleFractionParams", Name: "params"},
				Properties: typeastz.NewTuples(nil),
			},
			Query: typeastz.Typedef{
				Head:       &typeastz.Tuple{Typ: "apiSampleFractionQuery", Name: "query"},
				Properties: typeastz.NewTuples(nil),
			},
			Headers: typeastz.Typedef{
				Head:       &typeastz.Tuple{Typ: "apiSampleFractionReqHeaders", Name: "reqHeaders"},
				Properties: typeastz.NewTuples(nil),
			},
			Payload: typeastz.Typedef{
				Head:       &typeastz.Tuple{Typ: "apiSampleFractionReqPayload", Name: "reqPayload"},
				Properties: typeastz.NewTuples(nil),
			},
		},
		Resp: apispecz.Resp{
			Headers: typeastz.Typedef{
				Head:       &typeastz.Tuple{Typ: "apiSampleFractionRespHeaders", Name: "respHeaders"},
				Properties: typeastz.NewTuples(nil),
			},
			Payload: typeastz.Typedef{
				Head:       &typeastz.Tuple{Typ: "apiSampleFractionRespPayload", Name: "respPayload"},
				Properties: typeastz.NewTuples(nil),
			},
		},
	}
	require.Equal(t, expected, actual)
}

func TestUnitAllApis(t *testing.T) {
	ctx := t.Context()
	apis := apicollectionz.Apis()
	require.NotEmpty(t, apis)
	for _, api := range apis {
		expected := forgejsapi.NewSdk().ParseApi(ctx, api)
		require.NotNil(t, expected)
		require.Equal(t, api.ApiSpec().Id, expected.Spec.Id)
		require.NotEmpty(t, expected.Spec.CamelId)
		require.Equal(t, api.ApiSpec().Method, expected.Spec.Method)
		require.Equal(t, api.ApiSpec().Path, expected.Spec.Path)
		var actual *apispecz.Api
		jsonz.MustParse(jsonz.MustFormat(expected).String(), &actual)
		errorz.Check(os.MkdirAll("target/jsapi/api", os.ModePerm))
		errorz.Check(os.WriteFile("target/jsapi/api/"+api.ApiSpec().Id+".json", jsonz.MustFormatPretty(expected).Bytes(), os.ModePerm))
	}
}

func TestUnitWrite(t *testing.T) {
	ctx := t.Context()
	type testScenario struct {
		in   *apispecz.Apis
		want string
	}

	check := func(t *testing.T, s testScenario) {
		t.Helper()
		var apisList []*apispecz.Api
		if s.in != nil {
			apisList = s.in.Apis
		}
		got := forgejsapi.Write(ctx, apisList)
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
			in:   &apispecz.Apis{},
			want: ``,
		})
	})

	t.Run("with multiple path params and extra request fields", func(t *testing.T) {
		apis := forgejsapi.NewSdk().ParseApis(ctx, []apiz.Api{routezsamplefraction.NewService().Api()})
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
      url: ` + "`" + `${this._services.config().servUrls?.["decsuite"] || ""}/api/gox/routez/sample/fraction/${encodeURIComponent(String(req.params.numerator))}/${encodeURIComponent(String(req.params.denominator))}` + "`" + `,
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
		api := &apispecz.Api{
			Spec: apispecz.ApiSpec{
				CamelId: "apiCustomAction",
				Method:  "POST",
				Path:    "/api/custom/action",
			},
			Req: apispecz.Req{
				Payload: typeastz.Typedef{
					Head: &typeastz.Tuple{Typ: "apiCustomActionReqPayload", GoTyp: "models.CustomReq", Name: "payload"},
				},
			},
			Resp: apispecz.Resp{
				Payload: typeastz.Typedef{
					Head: &typeastz.Tuple{Typ: "apiCustomActionRespPayload", GoTyp: "models.CustomResp", Name: "payload"},
				},
			},
		}

		tds := apispecz.ExtractTypedefs([]*apispecz.Api{api})
		allTds := tds.All()
		require.Len(t, allTds, 2)
		var reqTd *typeastz.Typedef
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

		apis := &apispecz.Apis{Apis: []*apispecz.Api{api}}
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
      url: ` + "`" + `${this._services.config().servUrls?.["decsuite"] || ""}/api/custom/action` + "`" + `,
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
		api := &apispecz.Api{
			Spec: apispecz.ApiSpec{
				CamelId: "apiCustomAction",
				Method:  "POST",
				Path:    "/api/custom/{id}",
				Desc:    &schemaz.Desc{Name: "Custom", Summary: "Custom action", Labels: map[string]string{"servname": "tfinfer"}},
			},
		}
		got := forgejsapi.Write(ctx, []*apispecz.Api{api})
		require.Contains(t, got, "url: `${this._services.config().servUrls?.[\"tfinfer\"] || \"\"}/api/custom/{id}`,")
		require.NotContains(t, got, "path:")
	})
}

func TestUnitStripeCheckoutParseApi(t *testing.T) {
	ctx := t.Context()
	parsed := forgejsapi.NewSdk().ParseApi(ctx, &dszapistripecheckout.Api{})
	require.NotNil(t, parsed)
	require.Equal(t, "dsz-stripe-checkout", parsed.Spec.Id)
	require.Equal(t, "apiDszStripeCheckout", parsed.Spec.CamelId)
	require.Equal(t, "POST", parsed.Spec.Method)
	require.Equal(t, "/api/stripe/checkout", parsed.Spec.Path)

	tdfMap := make(map[string]*typeastz.Typedef)
	for _, tdf := range parsed.Typedefs.All() {
		tdfMap[tdf.Head.Name] = tdf
	}

	props := func(name string) map[string]typeastz.Tuple {
		tdf, ok := tdfMap[name]
		require.True(t, ok, name)
		ret := make(map[string]typeastz.Tuple)
		for _, p := range tdf.Properties.All() {
			ret[p.Name] = p
		}
		return ret
	}

	req := props("apiDszStripeCheckoutReqPayload")
	for _, name := range []string{"dataset_name", "plan", "success_url", "cancel_url", "campaign_state", "payment_phone", "payment_email", "document"} {
		require.Contains(t, req, name)
	}
	for _, name := range []string{
		"user", "mode", "locale", "allow_promotion_codes", "line_items", "customer_email",
		"subscription_data_trial_period_days", "cnpj_hub", "responsible_name_hub", "responsible_email_hub",
		"id_solution_hub", "id_campanha_hub", "terms_version", "contract_version", "terms_accepted_at",
	} {
		require.NotContains(t, req, name)
	}
	require.Equal(t, "string", req["campaign_state"].Typ)

	resp := props("apiDszStripeCheckoutRespPayload")
	require.Equal(t, "string", resp["id"].Typ)
	require.Equal(t, "string", resp["url"].Typ)
	require.Contains(t, resp, "dataset")
}

func TestUnitModelzApiSignedGetParseApi(t *testing.T) {
	ctx := t.Context()
	api := &dszapimodelget.Api{}
	parsed := forgejsapi.NewSdk().ParseApi(ctx, api)
	require.NotNil(t, parsed)
	require.Equal(t, "dsz-model-get", parsed.Spec.Id)
	require.Equal(t, "apiDszModelGet", parsed.Spec.CamelId)
	require.Equal(t, "POST", parsed.Spec.Method)
	require.Equal(t, "/api/model/dataset/{dataset}/signed-get/{model_name}", parsed.Spec.Path)

	require.Equal(t, "apiDszModelGetParams", parsed.Req.Params.Head.Typ)
	require.Equal(t, "params", parsed.Req.Params.Head.Name)
	require.Equal(t, "apiDszModelGetReqPayload", parsed.Req.Payload.Head.Typ)
	require.Equal(t, "reqPayload", parsed.Req.Payload.Head.Name)
	require.Equal(t, "apiDszModelGetRespPayload", parsed.Resp.Payload.Head.Typ)
	require.Equal(t, "respPayload", parsed.Resp.Payload.Head.Name)

	tdfMap := make(map[string]*typeastz.Typedef)
	for _, tdf := range parsed.Typedefs.All() {
		tdfMap[tdf.Head.Name] = tdf
	}

	reqPayloadTdf, ok := tdfMap["apiDszModelGetReqPayload"]
	require.True(t, ok)
	require.NotNil(t, reqPayloadTdf)
	require.NotNil(t, reqPayloadTdf.Properties)

	propMap := make(map[string]typeastz.Tuple)
	for _, prop := range reqPayloadTdf.Properties.All() {
		propMap[prop.Name] = prop
	}
	require.Contains(t, propMap, "model_version")
	require.Equal(t, "string", propMap["model_version"].Typ)
	require.Contains(t, propMap, "model_partitions")
	require.Equal(t, "Record<apiDszModelGetReqPayloadModelPartitionsKey, apiDszModelGetReqPayloadModelPartitionsValue>", propMap["model_partitions"].Typ)

	respPayloadTdf, ok := tdfMap["apiDszModelGetRespPayload"]
	require.True(t, ok)
	require.NotNil(t, respPayloadTdf)
	require.NotNil(t, respPayloadTdf.Properties)

	respPropMap := make(map[string]typeastz.Tuple)
	for _, prop := range respPayloadTdf.Properties.All() {
		respPropMap[prop.Name] = prop
	}
	require.Contains(t, respPropMap, "urls")
}
