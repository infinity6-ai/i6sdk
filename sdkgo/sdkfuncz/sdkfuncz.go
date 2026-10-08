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

type InferSpec struct {
	Id       string
	Desc     func() *schemaz.Desc
	Typ      Typ
	FuncName string
	StoreId  string
	Version  string
}

type Infer interface {
	InferSpec() InferSpec
	InferDataRefs() *DataRefs
}

type SdkFunc = Infer
type SdkFuncSpec = InferSpec
