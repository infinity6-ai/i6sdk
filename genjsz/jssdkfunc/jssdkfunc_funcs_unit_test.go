package jssdkfunc_test

import (
	"testing"

	"github.com/infinity6-ai/gox/commonz/jsonz"
	"github.com/stretchr/testify/require"
	"go.code.infinity6.ai/sdkgo/genjsz/jssdkfunc"
	"go.code.infinity6.ai/sdkgo/genjsz/jstypez"
	"go.code.infinity6.ai/sdkgo/sdkgo/productz/relevance/inferz/inferzrelevancefashionfbt"
	"go.code.infinity6.ai/sdkgo/sdkgo/sdkfuncz"
)

func TestUnitParseSdkFuncs(t *testing.T) {
	ctx := t.Context()
	type testScenario struct {
		in   []sdkfuncz.Func
		want *jssdkfunc.SdkFuncs
	}

	check := func(t *testing.T, s testScenario) {
		t.Helper()
		actual := jssdkfunc.ParseSdkFuncs(ctx, s.in)
		jsonz.MustParse(jsonz.MustFormat(actual).String(), &actual)
		require.Equal(t, s.want, actual)
	}

	t.Run("empty sdk funcs slice", func(t *testing.T) {
		check(t, testScenario{
			in: []sdkfuncz.Func{},
			want: &jssdkfunc.SdkFuncs{
				SdkFuncs: []*jssdkfunc.SdkFunc{},
				Typedefs: jstypez.NewTypedefs(nil),
			},
		})
	})

	t.Run("nil sdk funcs slice", func(t *testing.T) {
		check(t, testScenario{
			in: nil,
			want: &jssdkfunc.SdkFuncs{
				SdkFuncs: []*jssdkfunc.SdkFunc{},
				Typedefs: jstypez.NewTypedefs(nil),
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
