package sdkfuncz

import (
	"github.com/infinity6-ai/gox/schemaz/schemaz"
)

type Typ string

const (
	TypClientInfer Typ = "client_infer"
	TypClientQuery Typ = "client_query"
)

type DataRefs struct {
	Input  *schemaz.Schema
	Output *schemaz.Schema
}

type Spec struct {
	Id       string
	Desc     func() *schemaz.Desc
	Typ      Typ
	FuncName string
}

type Infer interface {
	InferSpec() Spec
	InferDataRefs() *DataRefs
}

type SdkFunc = Infer
type SdkFuncSpec = Spec
