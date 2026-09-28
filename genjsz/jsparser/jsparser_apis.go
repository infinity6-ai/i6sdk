package jsparser

import (
	"github.com/infinity6-ai/gox/routez/apiz"
	"go.code.infinity6.ai/sdkgo/genjsz/jstypez"
)

type Apis struct {
	Apis     []*Api           `json:"apis"`
	Typedefs jstypez.Typedefs `json:"typedefs"`
}

func ParseApis(apis []apiz.Api) *Apis {
	ret := Apis{
		Apis:     make([]*Api, len(apis)),
		Typedefs: jstypez.NewTypedefs(nil),
	}
	for idx, api := range apis {
		parsed := ParseApi(api)
		ret.Typedefs.Merge(parsed.Typedefs)
		ret.Apis[idx] = parsed
	}
	return &ret
}
