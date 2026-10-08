package jsparser_test

import (
	"testing"

	"github.com/infinity6-ai/gox/commonz/jsonz"
	"github.com/infinity6-ai/gox/routez/sampleroutez/routezsamplefraction"
	"github.com/infinity6-ai/gox/schemaz/schemaz"
	"github.com/stretchr/testify/require"
	"go.code.infinity6.ai/sdkgo/genjsz/jsparser"
	"go.code.infinity6.ai/sdkgo/genjsz/jstypez"
	"go.code.infinity6.ai/sdkgo/sdkgo/sdkz"
)

func TestUnitFractionSample(t *testing.T) {
	api := routezsamplefraction.NewService().Api()
	actual := jsparser.ParseApi(api)
	jsonz.MustParse(jsonz.MustFormat(actual).String(), &actual)
	println(jsonz.MustFormatPretty(actual).String())
	expected := &jsparser.Api{
		Spec: jsparser.ApiSpec{
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
					{Typ: "number", Name: "denominator", Summary: "Fraction Denominator"},
					{Typ: "number", Name: "numerator", Summary: "Fraction Numerator"},
				}),
			},
			{
				Head: &jstypez.Tuple{Typ: "object", Name: "apiSampleFractionQuery", Summary: "apiSampleFraction query"},
				Properties: jstypez.NewTuples([]jstypez.Tuple{
					{Typ: "number", Name: "precision", Summary: "Fraction max precision"},
				}),
			},
			{
				Head: &jstypez.Tuple{Typ: "object", Name: "apiSampleFractionReqHeaders", Summary: "apiSampleFraction reqHeaders"},
				Properties: jstypez.NewTuples([]jstypez.Tuple{
					{Typ: "string", Name: "x_i6_trace_id", Summary: "Sample of Trace Id"},
				}),
			},
			{
				Head: &jstypez.Tuple{Typ: "object", Name: "apiSampleFractionReqPayload", Summary: "apiSampleFraction reqPayload"},
				Properties: jstypez.NewTuples([]jstypez.Tuple{
					{Typ: "Array<apiSampleFractionReqPayloadMessagesElement>", Name: "messages"},
					{Typ: "string", Name: "reason", Summary: "Sample reason of this fraction"},
				}),
			},
			{
				Head: &jstypez.Tuple{Typ: "object", Name: "apiSampleFractionReqPayloadMessagesElement", Summary: "message"},
				Properties: jstypez.NewTuples([]jstypez.Tuple{
					{Typ: "string", Name: "message", Summary: "message string"},
				}),
			},
			{
				Head: &jstypez.Tuple{Typ: "object", Name: "apiSampleFractionRespHeaders", Summary: "apiSampleFraction respHeaders"},
				Properties: jstypez.NewTuples([]jstypez.Tuple{
					{Typ: "string", Name: "x_i6_trace_message", Summary: "Sample of any trace message"},
				}),
			},
			{
				Head: &jstypez.Tuple{Typ: "object", Name: "apiSampleFractionRespPayload", Summary: "apiSampleFraction respPayload"},
				Properties: jstypez.NewTuples([]jstypez.Tuple{
					{Typ: "string", Name: "display", Summary: "Human Fraction Representation"},
					{Typ: "Array<apiSampleFractionRespPayloadMessagesElement>", Name: "messages", Summary: "Output Messages"},
					{Typ: "string", Name: "result", Summary: "Fraction Result"},
				}),
			},
			{
				Head: &jstypez.Tuple{Typ: "string", Name: "apiSampleFractionRespPayloadMessagesElement"},
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
		Req: jsparser.Req{
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
		Resp: jsparser.Resp{
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
	// apis := []apiz.Api{
	// &ingestzapigeturl.GetUrlApi{},
	// &dszapisignedurl.Api{},
	// 	&pipezapi.PipezStartApi{},
	// 	&dszapi.DszSolutionDomainSelectApi{},
	// }
	apis := sdkz.Apis()
	for _, api := range apis {
		expected := jsparser.ParseApi(api)
		var actual *jsparser.Api
		jsonz.MustParse(jsonz.MustFormat(expected).String(), &actual)
		// require.Equal(t, expected, actual)
		// println(jsonz.MustFormatPretty(actual).String())
	}
	// t.Fail()
}

func TestUnitAllSdkFuncs(t *testing.T) {
	// sdkFuncs := []apiz.Api{
	// &ingestzapigeturl.GetUrlApi{},
	// &dszapisignedurl.Api{},
	// 	&pipezapi.PipezStartApi{},
	// 	&dszapi.DszSolutionDomainSelectApi{},
	// }
	// sdkFuncs := sdkz.SdkFuncs()
	// for _, sdkFunc := range sdkFuncs {
	// expected := jsparser.ParseSdkFunc(sdkFunc)
	// var actual *jsparser.Api
	// jsonz.MustParse(jsonz.MustFormat(expected).String(), &actual)
	// require.Equal(t, expected, actual)
	// println(jsonz.MustFormatPretty(actual).String())
	// }
	// t.Fail()
}
