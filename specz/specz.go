package specz

import (
	"github.com/infinity6-ai/gox/routez/sampleroutez/routezsamplefraction"
	"github.com/infinity6-ai/gox/schemaz/schemaz"
)

var all []*schemaz.Api

func init() {
	all = append(all, routezsamplefraction.Schema())
}

func All() []*schemaz.Api {
	return all
}
