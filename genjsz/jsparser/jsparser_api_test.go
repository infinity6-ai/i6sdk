package jsparser_test

import (
	"testing"

	"github.com/infinity6-ai/gox/commonz/jsonz"
	"github.com/infinity6-ai/gox/routez/sampleroutez/routezsamplefraction"
	"github.com/infinity6-ai/gox/schemaz/schemaz"
	"github.com/stretchr/testify/require"
	"go.code.infinity6.ai/sdkgo/genjsz/jsparser"
	"go.code.infinity6.ai/sdkgo/genjsz/jstypez"
	"go.code.infinity6.ai/sdkgo/sdkgo/apiz/ingestzapi/ingestzgeturlapi"
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
					{Typ: "string", Name: "reason", Summary: "Sample reason of this fraction"},
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
					{Typ: "string", Name: "result", Summary: "Fraction Result"},
				}),
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

func TestUnitIngestz(t *testing.T) {
	api := &ingestzgeturlapi.GetUrlApi{}
	actual := jsparser.ParseApi(api)
	jsonz.MustParse(jsonz.MustFormat(actual).String(), &actual)
	expected := &jsparser.Api{
		Spec: jsparser.ApiSpec{
			Id:      "ingestz-get-url",
			CamelId: "apiIngestzGetUrl",
			Desc: &schemaz.Desc{
				Name:     "Get Ingestion Upload URLs",
				Summary:  "Generates presigned PUT URLs and an ingestion token for uploading files into a target dataset table and partitions.",
				Markdown: "",
			},
			Method: "POST",
			Path:   "/api/ingest/get-url/dataset/{dataset}",
		},
		Typedefs: jstypez.NewTypedefs([]*jstypez.Typedef{
			{
				Head: &jstypez.Tuple{Typ: "object", Name: "apiIngestzGetUrlParams", Summary: "apiIngestzGetUrl params"},
				Properties: jstypez.NewTuples([]jstypez.Tuple{
					{Typ: "string", Name: "dataset", Summary: "Name of the dataset that receives the uploaded files."},
				}),
			},
			{
				Head: &jstypez.Tuple{Typ: "object", Name: "apiIngestzGetUrlReqPayload", Summary: "apiIngestzGetUrl reqPayload"},
				Properties: jstypez.NewTuples([]jstypez.Tuple{
					{Typ: "number", Name: "amount", Summary: "Number of presigned upload URLs to generate."},
					{Typ: "Record<string, string>", Name: "partitions", Summary: "Partition key/value pairs that select the table partition the files are ingested into."},
					{Typ: "string", Name: "table", Summary: "Name of the table, inside the dataset, that receives the uploaded files."},
				}),
			},
			{
				Head: &jstypez.Tuple{Typ: "object", Name: "apiIngestzGetUrlRespPayload", Summary: "apiIngestzGetUrl respPayload"},
				Properties: jstypez.NewTuples([]jstypez.Tuple{
					{Typ: "string", Name: "ingest_id", Summary: "Ingestion token that identifies this ingestion and ties the uploaded files together."},
					{Typ: "Array<string>", Name: "uploads", Summary: "Presigned PUT URLs, one per requested file; upload each file to its URL."},
				}),
			},
		}),
		Req: jsparser.Req{
			Params: jstypez.Typedef{
				Head:       &jstypez.Tuple{Typ: "apiIngestzGetUrlParams", Name: "params"},
				Properties: jstypez.NewTuples(nil),
			},
			Payload: jstypez.Typedef{
				Head:       &jstypez.Tuple{Typ: "apiIngestzGetUrlReqPayload", Name: "reqPayload"},
				Properties: jstypez.NewTuples(nil),
			},
		},
		Resp: jsparser.Resp{
			Payload: jstypez.Typedef{
				Head:       &jstypez.Tuple{Typ: "apiIngestzGetUrlRespPayload", Name: "respPayload"},
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
