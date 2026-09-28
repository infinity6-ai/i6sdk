package jsparser_test

import (
	"testing"

	"github.com/infinity6-ai/gox/commonz/jsonz"
	"github.com/infinity6-ai/gox/routez/apiz"
	"github.com/infinity6-ai/gox/routez/sampleroutez/routezsamplefraction"
	"github.com/stretchr/testify/require"
	"go.code.infinity6.ai/sdkgo/genjsz/jsparser"
	"go.code.infinity6.ai/sdkgo/genjsz/jstypez"
	"go.code.infinity6.ai/sdkgo/sdkgo/apiz/ingestzgeturlapi"
)

func expectedFractionApi() *jsparser.Api {
	return &jsparser.Api{
		Spec: jsparser.ApiSpec{
			Id:      "samplefraction",
			CamelId: "samplefraction",
			Method:  "POST",
			Path:    "/api/gox/routez/sample/fraction/{numerator}/{denominator}",
		},
		Typedefs: jstypez.NewTypedefs([]*jstypez.Typedef{
			{
				Head: &jstypez.Tuple{Typ: "object", Name: "samplefractionparams"},
				Properties: jstypez.NewTuples([]jstypez.Tuple{
					{Typ: "number", Name: "denominator"},
					{Typ: "number", Name: "numerator"},
				}),
			},
			{
				Head: &jstypez.Tuple{Typ: "object", Name: "samplefractionquery"},
				Properties: jstypez.NewTuples([]jstypez.Tuple{
					{Typ: "number", Name: "precision"},
				}),
			},
			{
				Head: &jstypez.Tuple{Typ: "object", Name: "samplefractionreqHeaders"},
				Properties: jstypez.NewTuples([]jstypez.Tuple{
					{Typ: "string", Name: "x_i6_trace_id"},
				}),
			},
			{
				Head: &jstypez.Tuple{Typ: "object", Name: "samplefractionreqPayload"},
				Properties: jstypez.NewTuples([]jstypez.Tuple{
					{Typ: "string", Name: "reason"},
				}),
			},
			{
				Head: &jstypez.Tuple{Typ: "object", Name: "samplefractionrespHeaders"},
				Properties: jstypez.NewTuples([]jstypez.Tuple{
					{Typ: "string", Name: "x_i6_trace_message"},
				}),
			},
			{
				Head: &jstypez.Tuple{Typ: "object", Name: "samplefractionrespPayload"},
				Properties: jstypez.NewTuples([]jstypez.Tuple{
					{Typ: "string", Name: "display"},
					{Typ: "string", Name: "result"},
				}),
			},
		}),
		Req: jsparser.Req{
			Params: jstypez.Typedef{
				Head:       &jstypez.Tuple{Typ: "samplefractionparams", Name: "params"},
				Properties: jstypez.NewTuples(nil),
			},
			Query: jstypez.Typedef{
				Head:       &jstypez.Tuple{Typ: "samplefractionquery", Name: "query"},
				Properties: jstypez.NewTuples(nil),
			},
			Headers: jstypez.Typedef{
				Head:       &jstypez.Tuple{Typ: "samplefractionreqHeaders", Name: "reqHeaders"},
				Properties: jstypez.NewTuples(nil),
			},
			Payload: jstypez.Typedef{
				Head:       &jstypez.Tuple{Typ: "samplefractionreqPayload", Name: "reqPayload"},
				Properties: jstypez.NewTuples(nil),
			},
		},
		Resp: jsparser.Resp{
			Headers: jstypez.Typedef{
				Head:       &jstypez.Tuple{Typ: "samplefractionrespHeaders", Name: "respHeaders"},
				Properties: jstypez.NewTuples(nil),
			},
			Payload: jstypez.Typedef{
				Head:       &jstypez.Tuple{Typ: "samplefractionrespPayload", Name: "respPayload"},
				Properties: jstypez.NewTuples(nil),
			},
		},
	}
}

func expectedIngestzApi() *jsparser.Api {
	return &jsparser.Api{
		Spec: jsparser.ApiSpec{
			Id:      "ingestz-get-url",
			CamelId: "ingestzGetUrl",
			Method:  "POST",
			Path:    "/api/ingest/get-url/dataset/{dataset}",
		},
		Typedefs: jstypez.NewTypedefs([]*jstypez.Typedef{
			{
				Head: &jstypez.Tuple{Typ: "object", Name: "ingestzGetUrlparams"},
				Properties: jstypez.NewTuples([]jstypez.Tuple{
					{Typ: "string", Name: "dataset"},
				}),
			},
			{
				Head: &jstypez.Tuple{Typ: "object", Name: "ingestzGetUrlreqPayload"},
				Properties: jstypez.NewTuples([]jstypez.Tuple{
					{Typ: "number", Name: "amount"},
					{Typ: "Record<string, string>", Name: "partitions"},
					{Typ: "string", Name: "table"},
				}),
			},
			{
				Head: &jstypez.Tuple{Typ: "object", Name: "ingestzGetUrlrespPayload"},
			},
		}),
		Req: jsparser.Req{
			Params: jstypez.Typedef{
				Head:       &jstypez.Tuple{Typ: "ingestzGetUrlparams", Name: "params"},
				Properties: jstypez.NewTuples(nil),
			},
			Payload: jstypez.Typedef{
				Head:       &jstypez.Tuple{Typ: "ingestzGetUrlreqPayload", Name: "reqPayload"},
				Properties: jstypez.NewTuples(nil),
			},
		},
		Resp: jsparser.Resp{
			Payload: jstypez.Typedef{
				Head:       &jstypez.Tuple{Typ: "ingestzGetUrlrespPayload", Name: "respPayload"},
				Properties: jstypez.NewTuples(nil),
			},
		},
	}
}

func TestUnitParseApis(t *testing.T) {
	ctx := t.Context()
	type testScenario struct {
		in   []apiz.Api
		want *jsparser.Apis
	}

	check := func(t *testing.T, s testScenario) {
		t.Helper()
		actual := jsparser.ParseApis(ctx, s.in)
		jsonz.MustParse(jsonz.MustFormat(actual).String(), &actual)
		require.Equal(t, s.want, actual)
	}

	t.Run("empty apis slice", func(t *testing.T) {
		check(t, testScenario{
			in: []apiz.Api{},
			want: &jsparser.Apis{
				Apis:     []*jsparser.Api{},
				Typedefs: jstypez.NewTypedefs(nil),
			},
		})
	})

	t.Run("nil apis slice", func(t *testing.T) {
		check(t, testScenario{
			in: nil,
			want: &jsparser.Apis{
				Apis:     []*jsparser.Api{},
				Typedefs: jstypez.NewTypedefs(nil),
			},
		})
	})

	t.Run("single api", func(t *testing.T) {
		ingestz := expectedIngestzApi()
		check(t, testScenario{
			in: []apiz.Api{
				&ingestzgeturlapi.GetUrlApi{},
			},
			want: &jsparser.Apis{
				Apis:     []*jsparser.Api{ingestz},
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
				&ingestzgeturlapi.GetUrlApi{},
			},
			want: &jsparser.Apis{
				Apis:     []*jsparser.Api{fraction, ingestz},
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
			want: &jsparser.Apis{
				Apis:     []*jsparser.Api{fraction1, fraction2},
				Typedefs: fraction1.Typedefs,
			},
		})
	})
}
