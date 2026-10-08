package sdkz

import (
	"fmt"
	"maps"
	"slices"

	"github.com/infinity6-ai/gox/commonz/validation/checker"
	"go.code.infinity6.ai/sdkgo/sdkgo/productz/relevance/inferz/inferzrelevancefashionfbt"
	"go.code.infinity6.ai/sdkgo/sdkgo/sdkfuncz"
)

var sdkFuncs map[string]sdkfuncz.Infer

func init() {
	sdkFuncs = make(map[string]sdkfuncz.Infer)

	add := func(i sdkfuncz.Infer) {
		id := i.InferSpec().Id
		old, ok := sdkFuncs[id]
		if ok {
			panic(fmt.Sprintf("sdkfunc id collision %s: old %T, new %T", id, old, i))
		}
		sdkFuncs[id] = i
	}

	add(&inferzrelevancefashionfbt.Api{})
}

func SdkFuncs() []sdkfuncz.Infer {
	return slices.Collect(maps.Values(sdkFuncs))
}

func SdkFunc(id string) sdkfuncz.Infer {
	ret, ok := sdkFuncs[id]
	checker.True(ok, "sdkfunc not found: %s", id)
	return ret
}

func Infers() []sdkfuncz.Infer {
	return SdkFuncs()
}

func Infer(id string) sdkfuncz.Infer {
	return SdkFunc(id)
}
