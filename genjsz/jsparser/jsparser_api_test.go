package jsparser_test

import (
	"testing"

	"github.com/infinity6-ai/gox/commonz/jsonz"
	"github.com/infinity6-ai/gox/routez/sampleroutez/routezsamplefraction"
	"github.com/infinity6-ai/gox/schemaz/schemaz"
	"github.com/stretchr/testify/require"
	"go.code.infinity6.ai/sdkgo/genjsz/jsparser"
	"go.code.infinity6.ai/sdkgo/genjsz/jstypez"
	"go.code.infinity6.ai/sdkgo/sdkgo/apiz/ingestzgeturlapi"
	"go.code.infinity6.ai/sdkgo/sdkgo/apiz/pipezapi"
)

func TestUnitFractionSample(t *testing.T) {
	api := routezsamplefraction.NewService().Api()
	actual := jsparser.ParseApi(api)
	jsonz.MustParse(jsonz.MustFormat(actual).String(), &actual)
	expected := &jsparser.Api{
		Spec: jsparser.ApiSpec{
			Id:      "sample-fraction",
			CamelId: "apiSampleFraction",
			Desc: &schemaz.Desc{
				Name:     "Sample Fraction",
				Summary:  "Sample API that performs a fraction operation",
				Markdown: "",
			},
			Method:  "POST",
			Path:    "/api/gox/routez/sample/fraction/{numerator}/{denominator}",
		},
		Typedefs: jstypez.NewTypedefs([]*jstypez.Typedef{
			{
				Head: &jstypez.Tuple{Typ: "object", Name: "apiSampleFractionparams"},
				Properties: jstypez.NewTuples([]jstypez.Tuple{
					{Typ: "number", Name: "denominator", Summary: "Fraction Denominator"},
					{Typ: "number", Name: "numerator", Summary: "Fraction Numerator"},
				}),
			},
			{
				Head: &jstypez.Tuple{Typ: "object", Name: "apiSampleFractionquery"},
				Properties: jstypez.NewTuples([]jstypez.Tuple{
					{Typ: "number", Name: "precision", Summary: "Fraction max precision"},
				}),
			},
			{
				Head: &jstypez.Tuple{Typ: "object", Name: "apiSampleFractionreqHeaders"},
				Properties: jstypez.NewTuples([]jstypez.Tuple{
					{Typ: "string", Name: "x_i6_trace_id", Summary: "Sample of Trace Id"},
				}),
			},
			{
				Head: &jstypez.Tuple{Typ: "object", Name: "apiSampleFractionreqPayload"},
				Properties: jstypez.NewTuples([]jstypez.Tuple{
					{Typ: "string", Name: "reason", Summary: "Sample reason of this fraction"},
				}),
			},
			{
				Head: &jstypez.Tuple{Typ: "object", Name: "apiSampleFractionrespHeaders"},
				Properties: jstypez.NewTuples([]jstypez.Tuple{
					{Typ: "string", Name: "x_i6_trace_message", Summary: "Sample of any trace message"},
				}),
			},
			{
				Head: &jstypez.Tuple{Typ: "object", Name: "apiSampleFractionrespPayload"},
				Properties: jstypez.NewTuples([]jstypez.Tuple{
					{Typ: "string", Name: "display", Summary: "Human Fraction Representation"},
					{Typ: "string", Name: "result", Summary: "Fraction Result"},
				}),
			},
		}),
		Req: jsparser.Req{
			Params: jstypez.Typedef{
				Head:       &jstypez.Tuple{Typ: "apiSampleFractionparams", Name: "params"},
				Properties: jstypez.NewTuples(nil),
			},
			Query: jstypez.Typedef{
				Head:       &jstypez.Tuple{Typ: "apiSampleFractionquery", Name: "query"},
				Properties: jstypez.NewTuples(nil),
			},
			Headers: jstypez.Typedef{
				Head:       &jstypez.Tuple{Typ: "apiSampleFractionreqHeaders", Name: "reqHeaders"},
				Properties: jstypez.NewTuples(nil),
			},
			Payload: jstypez.Typedef{
				Head:       &jstypez.Tuple{Typ: "apiSampleFractionreqPayload", Name: "reqPayload"},
				Properties: jstypez.NewTuples(nil),
			},
		},
		Resp: jsparser.Resp{
			Headers: jstypez.Typedef{
				Head:       &jstypez.Tuple{Typ: "apiSampleFractionrespHeaders", Name: "respHeaders"},
				Properties: jstypez.NewTuples(nil),
			},
			Payload: jstypez.Typedef{
				Head:       &jstypez.Tuple{Typ: "apiSampleFractionrespPayload", Name: "respPayload"},
				Properties: jstypez.NewTuples(nil),
			},
		},
	}
	require.Equal(t, expected, actual)
}

func TestUnitIngestz(t *testing.T) {
	api := &ingestzgeturlapi.GetUrlApi{}
	actual := jsparser.ParseApi(api)
	jsonz.MustParse(jsonz.MustFormat(actual).String(), &actual)
	expected := &jsparser.Api{
		Spec: jsparser.ApiSpec{
			Id:      "ingestz-get-url",
			CamelId: "apiIngestzGetUrl",
			Method:  "POST",
			Path:    "/api/ingest/get-url/dataset/{dataset}",
		},
		Typedefs: jstypez.NewTypedefs([]*jstypez.Typedef{
			{
				Head: &jstypez.Tuple{Typ: "object", Name: "apiIngestzGetUrlparams"},
				Properties: jstypez.NewTuples([]jstypez.Tuple{
					{Typ: "string", Name: "dataset"},
				}),
			},
			{
				Head: &jstypez.Tuple{Typ: "object", Name: "apiIngestzGetUrlreqPayload"},
				Properties: jstypez.NewTuples([]jstypez.Tuple{
					{Typ: "number", Name: "amount"},
					{Typ: "Record<string, string>", Name: "partitions"},
					{Typ: "string", Name: "table"},
				}),
			},
			{
				Head: &jstypez.Tuple{Typ: "object", Name: "apiIngestzGetUrlrespPayload"},
			},
		}),
		Req: jsparser.Req{
			Params: jstypez.Typedef{
				Head:       &jstypez.Tuple{Typ: "apiIngestzGetUrlparams", Name: "params"},
				Properties: jstypez.NewTuples(nil),
			},
			Payload: jstypez.Typedef{
				Head:       &jstypez.Tuple{Typ: "apiIngestzGetUrlreqPayload", Name: "reqPayload"},
				Properties: jstypez.NewTuples(nil),
			},
		},
		Resp: jsparser.Resp{
			Payload: jstypez.Typedef{
				Head:       &jstypez.Tuple{Typ: "apiIngestzGetUrlrespPayload", Name: "respPayload"},
				Properties: jstypez.NewTuples(nil),
			},
		},
	}
	require.Equal(t, expected, actual)
}

func TestManualOne(t *testing.T) {
	api := &pipezapi.PipezStartApi{}
	actual := jsparser.ParseApi(api)
	println(jsonz.MustFormat(actual).String())
}
