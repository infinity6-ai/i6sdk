package specz

import (
	"github.com/infinity6-ai/gox/routez/apiz"
	"github.com/infinity6-ai/gox/routez/sampleroutez/routezsamplefraction"
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
