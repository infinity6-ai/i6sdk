package jsparsersdkfunc_test

import (
	"os"
	"testing"

	"github.com/infinity6-ai/gox/commonz/errorz"
	"github.com/infinity6-ai/gox/commonz/jsonz"
	"github.com/infinity6-ai/gox/schemaz/schemaz"
	"github.com/stretchr/testify/require"
	"go.code.infinity6.ai/sdkgo/genjsz/jsparsersdkfunc"
	"go.code.infinity6.ai/sdkgo/genjsz/jstypez"
	"go.code.infinity6.ai/sdkgo/sdkgo/productz/relevance/inferz/inferzrelevancefashionfbt"
	"go.code.infinity6.ai/sdkgo/sdkgo/sdkfuncz"
	"go.code.infinity6.ai/sdkgo/sdkgo/sdkz"
)

func expectedFashionFbtSdkFunc() *jsparsersdkfunc.SdkFunc {
	return &jsparsersdkfunc.SdkFunc{
		Spec: jsparsersdkfunc.SdkFuncSpec{
			Id:      "product-relevance-fashion-fbt",
			CamelId: "sdkFuncProductRelevanceFashionFbt",
			Version: "v1",
			Desc: &schemaz.Desc{
				Name:    "Get Bought Together",
				Summary: "Infers the SKUs frequently bought together with the given SKUs for a store.",
			},
			Typ: sdkfuncz.TypClientInfer,
		},
		Typedefs: jstypez.NewTypedefs([]*jstypez.Typedef{
			{
				Head: &jstypez.Tuple{Typ: "string", Name: "dszTime", Summary: "Missing Summary"},
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
				Head: &jstypez.Tuple{Typ: "object", Name: "sdkFuncProductRelevanceFashionFbtInput", Summary: "sdkFuncProductRelevanceFashionFbt input"},
				Properties: jstypez.NewTuples([]jstypez.Tuple{
					{Typ: "string", GoTyp: "lakedsz.Dataset", Name: "dataset", Summary: "Dataset identifier."},
					{Typ: "string", GoTyp: "string", Name: "session_id", Summary: "Session identifier."},
					{Typ: "Array<sdkFuncProductRelevanceFashionFbtInputSkusElement>", Name: "skus", Summary: "SKUs to find matches for."},
					{Typ: "string", GoTyp: "string", Name: "user_id", Summary: "User identifier."},
				}),
			},
			{
				Head: &jstypez.Tuple{Typ: "object", Name: "sdkFuncProductRelevanceFashionFbtInputSkusElement"},
				Properties: jstypez.NewTuples([]jstypez.Tuple{
					{Typ: "string", GoTyp: "string", Name: "sku_id", Summary: "SKU identifier."},
				}),
			},
			{
				Head: &jstypez.Tuple{Typ: "object", Name: "sdkFuncProductRelevanceFashionFbtPartitions", Summary: "sdkFuncProductRelevanceFashionFbt partitions"},
				Properties: jstypez.NewTuples([]jstypez.Tuple{
					{Typ: "string", GoTyp: "string", Name: "store_id", Summary: "Store identifier."},
				}),
			},
			{
				Head: &jstypez.Tuple{Typ: "object", Name: "sdkFuncProductRelevanceFashionFbtOutput", Summary: "sdkFuncProductRelevanceFashionFbt output"},
				Properties: jstypez.NewTuples([]jstypez.Tuple{
					{Typ: "Array<sdkFuncProductRelevanceFashionFbtOutputResultsElement>", Name: "results", Summary: "Inference results."},
				}),
			},
			{
				Head: &jstypez.Tuple{Typ: "object", Name: "sdkFuncProductRelevanceFashionFbtOutputResultsElement"},
				Properties: jstypez.NewTuples([]jstypez.Tuple{
					{Typ: "number", GoTyp: "float64", Name: "rating", Summary: "Result rating."},
					{Typ: "Array<sdkFuncProductRelevanceFashionFbtOutputResultsElementSkusElement>", Name: "skus", Summary: "Matching SKUs."},
				}),
			},
			{
				Head: &jstypez.Tuple{Typ: "object", Name: "sdkFuncProductRelevanceFashionFbtOutputResultsElementSkusElement"},
				Properties: jstypez.NewTuples([]jstypez.Tuple{
					{Typ: "string", GoTyp: "string", Name: "sku_id", Summary: "SKU identifier."},
				}),
			},
		}),
		Req: jsparsersdkfunc.Req{
			Input: jstypez.Typedef{
				Head:       &jstypez.Tuple{Typ: "sdkFuncProductRelevanceFashionFbtInput", Name: "input"},
				Properties: jstypez.NewTuples(nil),
			},
			Partitions: jstypez.Typedef{
				Head:       &jstypez.Tuple{Typ: "sdkFuncProductRelevanceFashionFbtPartitions", Name: "partitions"},
				Properties: jstypez.NewTuples(nil),
			},
		},
		Resp: jsparsersdkfunc.Resp{
			Output: jstypez.Typedef{
				Head:       &jstypez.Tuple{Typ: "sdkFuncProductRelevanceFashionFbtOutput", Name: "output"},
				Properties: jstypez.NewTuples(nil),
			},
		},
	}
}

func TestUnitFashionFbtSdkFunc(t *testing.T) {
	ctx := t.Context()
	fn := &inferzrelevancefashionfbt.Func{}
	actual := jsparsersdkfunc.ParseSdkFunc(ctx, fn)
	jsonz.MustParse(jsonz.MustFormat(actual).String(), &actual)
	expected := expectedFashionFbtSdkFunc()
	require.Equal(t, expected, actual)
}

func TestUnitAllSdkFuncs(t *testing.T) {
	ctx := t.Context()
	funcs := sdkz.SdkFuncs()
	require.NotEmpty(t, funcs)
	for _, fn := range funcs {
		expected := jsparsersdkfunc.ParseSdkFunc(ctx, fn)
		var actual *jsparsersdkfunc.SdkFunc
		jsonz.MustParse(jsonz.MustFormat(expected).String(), &actual)
		errorz.Check(os.MkdirAll("target/jsparser/sdkfunc", os.ModePerm))
		errorz.Check(os.WriteFile("target/jsparser/sdkfunc/"+fn.SdkFuncSpec().Id+".json", jsonz.MustFormatPretty(expected).Bytes(), os.ModePerm))
	}
}
