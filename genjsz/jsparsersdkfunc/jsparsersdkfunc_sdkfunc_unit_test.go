package jsparsersdkfunc_test

import (
	"os"
	"testing"

	"github.com/infinity6-ai/gox/commonz/errorz"
	"github.com/infinity6-ai/gox/commonz/jsonz"
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
			Id:       "product-relevance-fashion-fbt",
			CamelId:  "sdkFuncProductRelevanceFashionFbt",
			Typ:      sdkfuncz.TypClientInfer,
			FuncName: "fashionFbt",
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
					{Typ: "Array<sdkFuncProductRelevanceFashionFbtInputSkusElement>", Name: "skus"},
				}),
			},
			{
				Head: &jstypez.Tuple{Typ: "object", Name: "sdkFuncProductRelevanceFashionFbtInputSkusElement"},
				Properties: jstypez.NewTuples([]jstypez.Tuple{
					{Typ: "string", Name: "sku_id"},
				}),
			},
			{
				Head: &jstypez.Tuple{Typ: "object", Name: "sdkFuncProductRelevanceFashionFbtPartitions", Summary: "sdkFuncProductRelevanceFashionFbt partitions"},
				Properties: jstypez.NewTuples([]jstypez.Tuple{
					{Typ: "string", Name: "store_id"},
					{Typ: "string", Name: "version"},
				}),
			},
			{
				Head: &jstypez.Tuple{Typ: "object", Name: "sdkFuncProductRelevanceFashionFbtOutput", Summary: "sdkFuncProductRelevanceFashionFbt output"},
				Properties: jstypez.NewTuples([]jstypez.Tuple{
					{Typ: "Array<sdkFuncProductRelevanceFashionFbtOutputResultsElement>", Name: "results"},
				}),
			},
			{
				Head: &jstypez.Tuple{Typ: "object", Name: "sdkFuncProductRelevanceFashionFbtOutputResultsElement"},
				Properties: jstypez.NewTuples([]jstypez.Tuple{
					{Typ: "string", Name: "rating"},
					{Typ: "string", Name: "sku_id"},
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
	fn := &inferzrelevancefashionfbt.Func{}
	actual := jsparsersdkfunc.ParseSdkFunc(fn)
	jsonz.MustParse(jsonz.MustFormat(actual).String(), &actual)
	expected := expectedFashionFbtSdkFunc()
	require.Equal(t, expected, actual)
}

func TestUnitAllSdkFuncs(t *testing.T) {
	funcs := sdkz.SdkFuncs()
	require.NotEmpty(t, funcs)
	for _, fn := range funcs {
		expected := jsparsersdkfunc.ParseSdkFunc(fn)
		var actual *jsparsersdkfunc.SdkFunc
		content := jsonz.MustFormat(expected)
		jsonz.MustParse(content.String(), &actual)
		errorz.Check(os.MkdirAll("target/jsparser/sdkfunc", os.ModePerm))
		errorz.Check(os.WriteFile("target/jsparser/sdkfunc/"+fn.SdkFuncSpec().Id+".json", content.Bytes(), os.ModePerm))
	}
}
