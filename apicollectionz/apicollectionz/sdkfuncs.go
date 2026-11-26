package apicollectionz

import (
	"context"
	"fmt"
	"reflect"

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

type SdkFuncInvokerFunc func(ctx context.Context, fn sdkfuncz.Func) (int, error)

type SdkFuncs struct {
	inferzrelevancefashionfbt.ProductRelevanceFashionFbtV1Id
	relevancefashiondashattachedrevenue.ProductRelevanceFashionDashAttachedRevenueV1Id
	relevancefashiondashcoverage.ProductRelevanceFashionDashCoverageV1Id
	relevancefashiondashdisplaymetricsstatus.ProductRelevanceFashionDashDisplayMetricsStatusV1Id
	relevancefashiondashfunnel.ProductRelevanceFashionDashFunnelV1Id
	relevancefashiondashimpactedmetricsbyshowcasecategory.ProductRelevanceFashionDashImpactedMetricsByShowcaseCategoryV1Id
	relevancefashiondashimpactedmetricsbysku.ProductRelevanceFashionDashImpactedMetricsBySkuV1Id
	relevancefashiondashlistshowcasecategories.ProductRelevanceFashionDashListShowcaseCategoriesV1Id

	Invoker SdkFuncInvokerFunc
}

func (a *SdkFuncs) List() []sdkfuncz.Func {
	var list []sdkfuncz.Func
	typ := reflect.TypeOf(a)
	for i := 0; i < typ.NumMethod(); i++ {
		m := typ.Method(i)
		if m.Name == "List" {
			continue
		}
		if m.Type.NumIn() != 3 {
			continue
		}
		fnType := m.Type.In(2)
		if fnType.Kind() == reflect.Ptr {
			fnType = fnType.Elem()
		}
		fnVal, ok := reflect.New(fnType).Interface().(sdkfuncz.Func)
		if ok {
			list = append(list, fnVal)
		}
	}
	return list
}

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

	for _, fn := range (&SdkFuncs{}).List() {
		add(fn)
	}
}

func Funcs() []sdkfuncz.Func {
	return (&SdkFuncs{}).List()
}

func SdkFunc(id string) sdkfuncz.Func {
	ret, ok := sdkFuncs[id]
	checker.True(ok, "sdkfunc not found: %s", id)
	return ret
}
