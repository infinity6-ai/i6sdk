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

func Sdk() []apiz.Api {
	return []apiz.Api{
		&ingestzgeturlapi.GetUrlApi{},
		&pipezapi.PipezStartApi{},
		&dszapi.DszSolutionDomainSelectApi{},
	}
}
