package specz

import (
	"fmt"
	"maps"
	"slices"

	"github.com/infinity6-ai/gox/commonz/validation/checker"
	"github.com/infinity6-ai/gox/routez/apiz"
	"github.com/infinity6-ai/gox/routez/sampleroutez/routezsamplefraction"
	"go.code.infinity6.ai/sdkgo/sdkgo/apiz/dszapi"
	"go.code.infinity6.ai/sdkgo/sdkgo/apiz/ingestzgeturlapi"
	"go.code.infinity6.ai/sdkgo/sdkgo/apiz/pipezapi"
)

var all map[string]apiz.Api

func init() {
	add := func(a apiz.Api) {
		id := a.ApiSpec().Id
		old, ok := all[id]
		if ok {
			panic(fmt.Sprintf("api id collision %s: old %T, new %T", id, old, a))
		}
		all[id] = a
	}

	samples := routezsamplefraction.Services()
	all = make(map[string]apiz.Api, len(samples))
	for _, svc := range routezsamplefraction.Services() {
		add(svc.Api())
	}
	add(&ingestzgeturlapi.GetUrlApi{})
	add(&pipezapi.PipezStartApi{})
	add(&dszapi.DszSolutionDomainSelectApi{})
}

func All() map[string]apiz.Api {
	return maps.Clone(all)
}

func Sdk() []apiz.Api {
	return slices.Collect(maps.Values(all))
}

func Get(id string) apiz.Api {
	ret, ok := all[id]
	checker.True(ok, "api not found: %s", id)
	return ret
}
