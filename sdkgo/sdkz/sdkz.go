package sdkz

import (
	"fmt"
	"maps"
	"slices"

	"github.com/infinity6-ai/gox/commonz/validation/checker"
	"github.com/infinity6-ai/gox/routez/apiz"
	"github.com/infinity6-ai/gox/routez/sampleroutez/routezsamplefraction"
	"go.code.infinity6.ai/sdkgo/sdkgo/apiz/dszapi/dszapidatasetlist"
	"go.code.infinity6.ai/sdkgo/sdkgo/apiz/dszapi/dszapidomainsolutionget"
	"go.code.infinity6.ai/sdkgo/sdkgo/apiz/dszapi/dszapidomainsolutionselect"
	"go.code.infinity6.ai/sdkgo/sdkgo/apiz/dszapi/dszapiingestiontokencreate"
	"go.code.infinity6.ai/sdkgo/sdkgo/apiz/dszapi/dszapiingestiontokenslist"
	"go.code.infinity6.ai/sdkgo/sdkgo/apiz/dszapi/dszapimodelget"
	"go.code.infinity6.ai/sdkgo/sdkgo/apiz/dszapi/dszapirecsystokencreate"
	"go.code.infinity6.ai/sdkgo/sdkgo/apiz/dszapi/dszapirecsystokenslist"
	"go.code.infinity6.ai/sdkgo/sdkgo/apiz/dszapi/dszapisignedurl"
	"go.code.infinity6.ai/sdkgo/sdkgo/apiz/dszapi/dszapistorecreate"
	"go.code.infinity6.ai/sdkgo/sdkgo/apiz/dszapi/dszapistorelist"
	"go.code.infinity6.ai/sdkgo/sdkgo/apiz/dszapi/dszapitableget"
	"go.code.infinity6.ai/sdkgo/sdkgo/apiz/dszapi/dszapitablelist"
	"go.code.infinity6.ai/sdkgo/sdkgo/apiz/dszapi/dszapitokendelete"
	"go.code.infinity6.ai/sdkgo/sdkgo/apiz/ingestzapi/ingestzapigeturl"
	"go.code.infinity6.ai/sdkgo/sdkgo/apiz/pipezapi/pipezstart"
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
	add(&ingestzapigeturl.Api{})
	add(&dszapisignedurl.Api{})
	add(&dszapitableget.Api{})
	add(&dszapitablelist.Api{})
	add(&pipezstart.Api{})
	add(&dszapidomainsolutionselect.Api{})
	add(&dszapidomainsolutionget.Api{})
	add(&dszapistorecreate.Api{})
	add(&dszapistorelist.Api{})
	add(&dszapidatasetlist.Api{})
	add(&dszapiingestiontokenslist.Api{})
	add(&dszapiingestiontokencreate.Api{})
	add(&dszapirecsystokenslist.Api{})
	add(&dszapirecsystokencreate.Api{})
	add(&dszapitokendelete.Api{})
	add(&dszapimodelget.Api{})
}

func All() []apiz.Api {
	return slices.Collect(maps.Values(all))
}

func Get(id string) apiz.Api {
	ret, ok := all[id]
	checker.True(ok, "api not found: %s", id)
	return ret
}
