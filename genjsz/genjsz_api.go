package genjsz

import "github.com/infinity6-ai/gox/routez/apiz"

func ParseApi(api apiz.Api) {
	ParseSchema(api.GetDataRefs().PathParams)
}
