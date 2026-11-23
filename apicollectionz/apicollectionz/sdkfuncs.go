package apicollectionz

import (
	"fmt"
	"maps"
	"slices"

	"github.com/infinity6-ai/gox/commonz/validation/checker"
	"github.com/infinity6-ai/i6sdk/sdkcorez/productz/relevance/dash/relevancefashiondashattachedrevenue"
	"github.com/infinity6-ai/i6sdk/sdkcorez/productz/relevance/dash/relevancefashiondashcoverage"
	"github.com/infinity6-ai/i6sdk/sdkcorez/productz/relevance/dash/relevancefashiondashdisplaymetricsstatus"
	"github.com/infinity6-ai/i6sdk/sdkcorez/productz/relevance/dash/relevancefashiondashfunnel"
	"github.com/infinity6-ai/i6sdk/sdkcorez/productz/relevance/dash/relevancefashiondashimpactedmetricsbyshowcasecategory"
	"github.com/infinity6-ai/i6sdk/sdkcorez/productz/relevance/dash/relevancefashiondashimpactedmetricsbysku"
	"github.com/infinity6-ai/i6sdk/sdkcorez/productz/relevance/dash/relevancefashiondashlistshowcasecategories"
	"github.com/infinity6-ai/i6sdk/sdkcorez/productz/relevance/inferz/inferzrelevancefashionfbt"
	"github.com/infinity6-ai/i6sdk/sdkcorez/sdkfuncz"
)

var sdkFuncs map[string]sdkfuncz.Func

func init() {
	sdkFuncs = make(map[string]sdkfuncz.Func)

	add := func(i sdkfuncz.Func) {
		id := i.SdkFuncSpec().Id
		old, ok := sdkFuncs[id]
		if ok {
			panic(fmt.Sprintf("sdkfunc id collision %s: old %T, new %T", id, old, i))
		}
		sdkFuncs[id] = i
	}

	add(&inferzrelevancefashionfbt.Func{})
	add(&relevancefashiondashfunnel.Func{})
	add(&relevancefashiondashattachedrevenue.Func{})
	add(&relevancefashiondashlistshowcasecategories.Func{})
	add(&relevancefashiondashimpactedmetricsbyshowcasecategory.Func{})
	add(&relevancefashiondashimpactedmetricsbysku.Func{})
	add(&relevancefashiondashcoverage.Func{})
	add(&relevancefashiondashdisplaymetricsstatus.Func{})
}

func SdkFuncs() []sdkfuncz.Func {
	return slices.Collect(maps.Values(sdkFuncs))
}

func SdkFunc(id string) sdkfuncz.Func {
	ret, ok := sdkFuncs[id]
	checker.True(ok, "sdkfunc not found: %s", id)
	return ret
}
