package specz

import (
	"github.com/infinity6-ai/gox/routez/apiz"
	"github.com/infinity6-ai/gox/routez/sampleroutez/routezsamplefraction"
	"go.code.infinity6.ai/sdkgo/sdkgo/apiz/dszapi"
	"go.code.infinity6.ai/sdkgo/sdkgo/apiz/ingestzgeturlapi"
	"go.code.infinity6.ai/sdkgo/sdkgo/apiz/pipezapi"
)

var all []apiz.Api

func init() {
	for _, svc := range routezsamplefraction.Services() {
		all = append(all, svc.Api())
	}
}

func All() []apiz.Api {
	return all
}

// Sdk lists the real APIs exposed to the SDK clients, each one a fresh zero-value instance
// (generators call GetDataRefs on them, which mutates the instance). It is the single place
// that decides what `generate js` emits: a new API is part of the generated SDK once it is
// added here. samplefraction, which lives in All(), is deliberately not listed.
func Sdk() []apiz.Api {
	return []apiz.Api{
		&ingestzgeturlapi.GetUrlApi{},
		&pipezapi.PipezStartApi{},
		&dszapi.DszSolutionDomainSelectApi{},
	}
}
