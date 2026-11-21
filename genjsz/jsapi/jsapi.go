package jsapi

import (
	"github.com/infinity6-ai/gox/routez/jsapiz"
	"go.code.infinity6.ai/sdkgo/genjsz/jsschema"
)

func NewSdk() *jsapiz.Parser {
	return jsapiz.New(jsschema.NewSdk())
}
