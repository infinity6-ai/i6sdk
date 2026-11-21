package jssdkfunc_test

import (
	"os"
	"testing"

	"github.com/infinity6-ai/gox/commonz/errorz"
	"github.com/infinity6-ai/gox/commonz/jsonz"
	"github.com/infinity6-ai/gox/schemaz/typeastz"
	"github.com/infinity6-ai/gox/schemaz/schemaz"
	"github.com/stretchr/testify/require"
	"go.code.infinity6.ai/sdkgo/genjsz/jssdkfunc"
	"go.code.infinity6.ai/sdkgo/sdkgo/productz/relevance/inferz/inferzrelevancefashionfbt"
	"go.code.infinity6.ai/sdkgo/sdkgo/sdkfuncz"
	"go.code.infinity6.ai/sdkgo/sdkgo/sdkz"
)

func expectedFashionFbtSdkFunc() *jssdkfunc.SdkFunc {
	return &jssdkfunc.SdkFunc{
		Spec: jssdkfunc.SdkFuncSpec{
			Id:      "product-relevance-fashion-fbt",
			CamelId: "sdkFuncProductRelevanceFashionFbt",
			Version: "v1",
			Desc: &schemaz.Desc{
				Name:    "Get Bought Together",
				Summary: "Infers the SKUs frequently bought together with the given SKUs for a store.",
			},
			Typ: sdkfuncz.TypClientInfer,
		},
		Typedefs: typeastz.NewTypedefs([]*typeastz.Typedef{
			{
				Head: &typeastz.Tuple{Typ: "string", Name: "dszTime", Summary: "Missing Summary"},
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
				Head: &typeastz.Tuple{Typ: "object", Name: "sdkFuncProductRelevanceFashionFbtInput", Summary: "sdkFuncProductRelevanceFashionFbt input"},
				Properties: typeastz.NewTuples([]typeastz.Tuple{
					{Typ: "string", GoTyp: "lakedsz.Dataset", Name: "dataset", Summary: "Dataset identifier."},
					{Typ: "string", GoTyp: "string", Name: "session_id", Summary: "Session identifier."},
					{Typ: "Array<sdkFuncProductRelevanceFashionFbtInputSkusElement>", Name: "skus", Summary: "SKUs to find matches for."},
					{Typ: "string", GoTyp: "string", Name: "user_id", Summary: "User identifier."},
				}),
			},
			{
				Head: &typeastz.Tuple{Typ: "object", Name: "sdkFuncProductRelevanceFashionFbtInputSkusElement"},
				Properties: typeastz.NewTuples([]typeastz.Tuple{
					{Typ: "string", GoTyp: "string", Name: "sku_id", Summary: "SKU identifier."},
				}),
			},
			{
				Head: &typeastz.Tuple{Typ: "object", Name: "sdkFuncProductRelevanceFashionFbtPartitions", Summary: "sdkFuncProductRelevanceFashionFbt partitions"},
				Properties: typeastz.NewTuples([]typeastz.Tuple{
					{Typ: "string", GoTyp: "string", Name: "store_id", Summary: "Store identifier."},
				}),
			},
			{
				Head: &typeastz.Tuple{Typ: "object", Name: "sdkFuncProductRelevanceFashionFbtOutput", Summary: "sdkFuncProductRelevanceFashionFbt output"},
				Properties: typeastz.NewTuples([]typeastz.Tuple{
					{Typ: "Array<sdkFuncProductRelevanceFashionFbtOutputResultsElement>", Name: "results", Summary: "Inference results."},
				}),
			},
			{
				Head: &typeastz.Tuple{Typ: "object", Name: "sdkFuncProductRelevanceFashionFbtOutputResultsElement"},
				Properties: typeastz.NewTuples([]typeastz.Tuple{
					{Typ: "number", GoTyp: "float64", Name: "rating", Summary: "Result rating."},
					{Typ: "Array<sdkFuncProductRelevanceFashionFbtOutputResultsElementSkusElement>", Name: "skus", Summary: "Matching SKUs."},
				}),
			},
			{
				Head: &typeastz.Tuple{Typ: "object", Name: "sdkFuncProductRelevanceFashionFbtOutputResultsElementSkusElement"},
				Properties: typeastz.NewTuples([]typeastz.Tuple{
					{Typ: "string", GoTyp: "string", Name: "sku_id", Summary: "SKU identifier."},
				}),
			},
		}),
		Req: jssdkfunc.Req{
			Input: typeastz.Typedef{
				Head:       &typeastz.Tuple{Typ: "sdkFuncProductRelevanceFashionFbtInput", Name: "input"},
				Properties: typeastz.NewTuples(nil),
			},
			Partitions: typeastz.Typedef{
				Head:       &typeastz.Tuple{Typ: "sdkFuncProductRelevanceFashionFbtPartitions", Name: "partitions"},
				Properties: typeastz.NewTuples(nil),
			},
		},
		Resp: jssdkfunc.Resp{
			Output: typeastz.Typedef{
				Head:       &typeastz.Tuple{Typ: "sdkFuncProductRelevanceFashionFbtOutput", Name: "output"},
				Properties: typeastz.NewTuples(nil),
			},
		},
	}
}

func TestUnitFashionFbtSdkFunc(t *testing.T) {
	ctx := t.Context()
	fn := &inferzrelevancefashionfbt.Func{}
	actual := jssdkfunc.NewSdk().ParseSdkFunc(ctx, fn)
	jsonz.MustParse(jsonz.MustFormat(actual).String(), &actual)
	expected := expectedFashionFbtSdkFunc()
	require.Equal(t, expected, actual)
}

func TestUnitAllSdkFuncs(t *testing.T) {
	ctx := t.Context()
	funcs := sdkz.SdkFuncs()
	require.NotEmpty(t, funcs)
	for _, fn := range funcs {
		expected := jssdkfunc.NewSdk().ParseSdkFunc(ctx, fn)
		var actual *jssdkfunc.SdkFunc
		jsonz.MustParse(jsonz.MustFormat(expected).String(), &actual)
		errorz.Check(os.MkdirAll("target/jssdkfunc/sdkfunc", os.ModePerm))
		errorz.Check(os.WriteFile("target/jssdkfunc/sdkfunc/"+fn.SdkFuncSpec().Id+".json", jsonz.MustFormatPretty(expected).Bytes(), os.ModePerm))
	}
}

func TestUnitParseSdkFuncs(t *testing.T) {
	ctx := t.Context()
	type testScenario struct {
		in   []sdkfuncz.Func
		want *jssdkfunc.SdkFuncs
	}

	check := func(t *testing.T, s testScenario) {
		t.Helper()
		actual := jssdkfunc.NewSdk().ParseSdkFuncs(ctx, s.in)
		jsonz.MustParse(jsonz.MustFormat(actual).String(), &actual)
		require.Equal(t, s.want, actual)
	}

	t.Run("empty sdk funcs slice", func(t *testing.T) {
		check(t, testScenario{
			in: []sdkfuncz.Func{},
			want: &jssdkfunc.SdkFuncs{
				SdkFuncs: []*jssdkfunc.SdkFunc{},
				Typedefs: typeastz.NewTypedefs(nil),
			},
		})
	})

	t.Run("nil sdk funcs slice", func(t *testing.T) {
		check(t, testScenario{
			in: nil,
			want: &jssdkfunc.SdkFuncs{
				SdkFuncs: []*jssdkfunc.SdkFunc{},
				Typedefs: typeastz.NewTypedefs(nil),
			},
		})
	})

	t.Run("single sdk func", func(t *testing.T) {
		fbt := expectedFashionFbtSdkFunc()
		check(t, testScenario{
			in: []sdkfuncz.Func{
				&inferzrelevancefashionfbt.Func{},
			},
			want: &jssdkfunc.SdkFuncs{
				SdkFuncs: []*jssdkfunc.SdkFunc{fbt},
				Typedefs: fbt.Typedefs,
			},
		})
	})

	t.Run("duplicate sdk funcs merge typedefs idempotently", func(t *testing.T) {
		fbt1 := expectedFashionFbtSdkFunc()
		fbt2 := expectedFashionFbtSdkFunc()

		check(t, testScenario{
			in: []sdkfuncz.Func{
				&inferzrelevancefashionfbt.Func{},
				&inferzrelevancefashionfbt.Func{},
			},
			want: &jssdkfunc.SdkFuncs{
				SdkFuncs: []*jssdkfunc.SdkFunc{fbt1, fbt2},
				Typedefs: fbt1.Typedefs,
			},
		})
	})
}
