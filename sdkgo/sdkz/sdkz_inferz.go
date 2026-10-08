package sdkz

import (
	"fmt"
	"maps"
	"slices"

	"github.com/infinity6-ai/gox/commonz/validation/checker"
	"go.code.infinity6.ai/sdkgo/sdkgo/inferz"
	"go.code.infinity6.ai/sdkgo/sdkgo/productz/relevance/inferz/inferzrelevancefashionfbt"
)

var infers map[string]inferzspec.Infer

func init() {
	infers = make(map[string]inferzspec.Infer)

	add := func(i inferzspec.Infer) {
		id := i.InferSpec().Id
		old, ok := infers[id]
		if ok {
			panic(fmt.Sprintf("infer id collision %s: old %T, new %T", id, old, i))
		}
		infers[id] = i
	}

	add(&inferzrelevancefashionfbt.Api{})
}

func Infers() []inferzspec.Infer {
	return slices.Collect(maps.Values(infers))
}

func Infer(id string) inferzspec.Infer {
	ret, ok := infers[id]
	checker.True(ok, "infer not found: %s", id)
	return ret
}
