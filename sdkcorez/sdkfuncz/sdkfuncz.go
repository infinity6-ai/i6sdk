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
	Input      *schemaz.Schema
	Partitions *schemaz.Schema
	Output     *schemaz.Schema
}

type Spec struct {
	Id      string
	Version string
	Desc    func() *schemaz.Desc
	Typ     Typ
}

type Func interface {
	SdkFuncSpec() Spec
	SdkFuncDataRefs() *DataRefs
}

type UnimplementedFunc struct{}

func (UnimplementedFunc) SdkFuncSpec() Spec {
	panic("SdkFuncSpec not implemented: run 'go generate ./...'")
}

func (UnimplementedFunc) SdkFuncDataRefs() *DataRefs {
	panic("SdkFuncDataRefs not implemented: run 'go generate ./...'")
}
