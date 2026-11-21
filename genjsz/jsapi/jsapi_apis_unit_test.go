package jsapi_test

import (
	"testing"

	"github.com/infinity6-ai/gox/commonz/jsonz"
	"github.com/infinity6-ai/gox/routez/apiz"
	"github.com/infinity6-ai/gox/routez/sampleroutez/routezsamplefraction"
	"github.com/infinity6-ai/gox/schemaz/schemaz"
	"github.com/stretchr/testify/require"
	"go.code.infinity6.ai/sdkgo/genjsz/jsapi"
	"github.com/infinity6-ai/gox/schemaz/jstypez"
	"go.code.infinity6.ai/sdkgo/sdkgo/apiz/ingestzapi/ingestzapigeturl"
)

func expectedFractionApi() *jsapi.Api {
	return &jsapi.Api{
		Spec: jsapi.ApiSpec{
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
		Req: jsapi.Req{
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
		Resp: jsapi.Resp{
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
}

func expectedIngestzApi() *jsapi.Api {
	return &jsapi.Api{
		Spec: jsapi.ApiSpec{
			Id:      "ingestz-get-url",
			CamelId: "apiIngestzGetUrl",
			Desc: &schemaz.Desc{
				Name:     "Get Ingestion Upload URLs",
				Summary:  "Generates presigned PUT URLs and an ingestion token for uploading files into a target dataset table and partitions.",
				Markdown: "",
				Labels:   map[string]string{"token_types": "dataset"},
			},
			Method: "POST",
			Path:   "/api/ingest/dataset/{dataset}/get-url",
		},
		Typedefs: jstypez.NewTypedefs([]*jstypez.Typedef{
			{
				Head: &jstypez.Tuple{Typ: "object", Name: "apiIngestzGetUrlParams", Summary: "apiIngestzGetUrl params"},
				Properties: jstypez.NewTuples([]jstypez.Tuple{
					{Typ: "string", GoTyp: "lakedsz.Dataset", Name: "dataset", Summary: "Name of the dataset that receives the uploaded files."},
				}),
			},
			{
				Head: &jstypez.Tuple{Typ: "object", Name: "apiIngestzGetUrlReqPayload", Summary: "apiIngestzGetUrl reqPayload"},
				Properties: jstypez.NewTuples([]jstypez.Tuple{
					{Typ: "number", GoTyp: "int", Name: "amount", Summary: "Number of presigned upload URLs to generate."},
					{Typ: "Record<apiIngestzGetUrlReqPayloadPartitionsKey, apiIngestzGetUrlReqPayloadPartitionsValue>", Name: "partitions", Summary: "Partition key/value pairs that select the table partition the files are ingested into."},
					{Typ: "string", GoTyp: "laketbnamez.TableName", Name: "table", Summary: "Name of the table, inside the dataset, that receives the uploaded files."},
				}),
			},
			{
				Head: &jstypez.Tuple{Typ: "object", Name: "apiIngestzGetUrlRespPayload", Summary: "apiIngestzGetUrl respPayload"},
				Properties: jstypez.NewTuples([]jstypez.Tuple{
					{Typ: "string", GoTyp: "string", Name: "ingest_id", Summary: "Ingestion token that identifies this ingestion and ties the uploaded files together."},
					{Typ: "Array<apiIngestzGetUrlRespPayloadUploadsElement>", Name: "uploads", Summary: "Presigned PUT URLs, one per requested file; upload each file to its URL."},
				}),
			},
			{
				Head: &jstypez.Tuple{Typ: "string", GoTyp: "string", Name: "apiIngestzGetUrlReqPayloadPartitionsKey"},
			},
			{
				Head: &jstypez.Tuple{Typ: "string", GoTyp: "string", Name: "apiIngestzGetUrlReqPayloadPartitionsValue"},
			},
			{
				Head: &jstypez.Tuple{Typ: "string", GoTyp: "string", Name: "apiIngestzGetUrlRespPayloadUploadsElement"},
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
		Req: jsapi.Req{
			Params: jstypez.Typedef{
				Head:       &jstypez.Tuple{Typ: "apiIngestzGetUrlParams", Name: "params"},
				Properties: jstypez.NewTuples(nil),
			},
			Payload: jstypez.Typedef{
				Head:       &jstypez.Tuple{Typ: "apiIngestzGetUrlReqPayload", Name: "reqPayload"},
				Properties: jstypez.NewTuples(nil),
			},
		},
		Resp: jsapi.Resp{
			Payload: jstypez.Typedef{
				Head:       &jstypez.Tuple{Typ: "apiIngestzGetUrlRespPayload", Name: "respPayload"},
				Properties: jstypez.NewTuples(nil),
			},
		},
	}
}

func TestUnitParseApis(t *testing.T) {
	ctx := t.Context()
	type testScenario struct {
		in   []apiz.Api
		want *jsapi.Apis
	}

	check := func(t *testing.T, s testScenario) {
		t.Helper()
		actual := jsapi.NewSdk().ParseApis(ctx, s.in)
		jsonz.MustParse(jsonz.MustFormat(actual).String(), &actual)
		require.Equal(t, s.want, actual)
	}

	t.Run("empty apis slice", func(t *testing.T) {
		check(t, testScenario{
			in: []apiz.Api{},
			want: &jsapi.Apis{
				Apis:     []*jsapi.Api{},
				Typedefs: jstypez.NewTypedefs(nil),
			},
		})
	})

	t.Run("nil apis slice", func(t *testing.T) {
		check(t, testScenario{
			in: nil,
			want: &jsapi.Apis{
				Apis:     []*jsapi.Api{},
				Typedefs: jstypez.NewTypedefs(nil),
			},
		})
	})

	t.Run("single api", func(t *testing.T) {
		ingestz := expectedIngestzApi()
		check(t, testScenario{
			in: []apiz.Api{
				&ingestzapigeturl.Api{},
			},
			want: &jsapi.Apis{
				Apis:     []*jsapi.Api{ingestz},
				Typedefs: ingestz.Typedefs,
			},
		})
	})

	t.Run("multiple apis with merged typedefs", func(t *testing.T) {
		fraction := expectedFractionApi()
		ingestz := expectedIngestzApi()

		mergedTypedefs := jstypez.NewTypedefs(nil)
		mergedTypedefs.Merge(fraction.Typedefs)
		mergedTypedefs.Merge(ingestz.Typedefs)

		check(t, testScenario{
			in: []apiz.Api{
				routezsamplefraction.NewService().Api(),
				&ingestzapigeturl.Api{},
			},
			want: &jsapi.Apis{
				Apis:     []*jsapi.Api{fraction, ingestz},
				Typedefs: mergedTypedefs,
			},
		})
	})

	t.Run("duplicate apis merge typedefs idempotently", func(t *testing.T) {
		fraction1 := expectedFractionApi()
		fraction2 := expectedFractionApi()

		check(t, testScenario{
			in: []apiz.Api{
				routezsamplefraction.NewService().Api(),
				routezsamplefraction.NewService().Api(),
			},
			want: &jsapi.Apis{
				Apis:     []*jsapi.Api{fraction1, fraction2},
				Typedefs: fraction1.Typedefs,
			},
		})
	})
}
