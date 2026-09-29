package sdkz

import (
	"fmt"
	"maps"
	"slices"

	"github.com/infinity6-ai/gox/commonz/validation/checker"
	"github.com/infinity6-ai/gox/routez/apiz"
	"github.com/infinity6-ai/gox/routez/sampleroutez/routezsamplefraction"
	"go.code.infinity6.ai/sdkgo/sdkgo/apiz/ingestzapi/ingestzgeturlapi"
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
	// add(&pipezapi.PipezStartApi{})
	// add(&dszapi.DszSolutionDomainSelectApi{})
}

func All() []apiz.Api {
	return slices.Collect(maps.Values(all))
}

func Get(id string) apiz.Api {
	ret, ok := all[id]
	checker.True(ok, "api not found: %s", id)
	return ret
}
