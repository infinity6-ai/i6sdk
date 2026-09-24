package genjsz

import (
	"github.com/infinity6-ai/gox/routez/apiz"
	"go.code.infinity6.ai/sdkgo/genjsz/jstyp"
)

func ParseApi(api apiz.Api) jstyp.Typ {
	return ParseSchema(api.GetDataRefs().PathParams)
}
