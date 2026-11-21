package jsapi

import (
	"github.com/infinity6-ai/gox/routez/jsapiz"
	"github.com/infinity6-ai/gox/schemaz/jsschemaz"
	"github.com/infinity6-ai/gox/schemaz/jstypez"
	"go.code.infinity6.ai/sdkgo/genjsz/jsschema"
)

type Api = jsapiz.Api
type Apis = jsapiz.Apis
type Req = jsapiz.Req
type Resp = jsapiz.Resp
type ApiSpec = jsapiz.ApiSpec
type Parser = jsapiz.Parser

func New(schemaParser *jsschemaz.Parser) *Parser {
	return jsapiz.New(schemaParser)
}

func NewSdk() *Parser {
	return jsapiz.New(jsschema.NewSdk())
}

func ExtractTypedefs(apis []*Api) jstypez.Typedefs {
	return jsapiz.ExtractTypedefs(apis)
}
