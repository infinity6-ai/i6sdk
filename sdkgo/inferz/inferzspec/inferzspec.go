package inferzspec

import (
	"context"

	"github.com/infinity6-ai/gox/schemaz/schemaz"
)

type DataRefs struct {
	Input  *schemaz.Schema
	Output *schemaz.Schema
}

type InferSpec struct {
	Id   string
	Desc func() *schemaz.Desc
	Name string
}

type Infer interface {
	InferSpec() InferSpec
	InferDataRefs() *DataRefs
}

type Handler[T Infer] func(ctx context.Context, reqResp T) (int, error)

type Service interface {
	New() Service
	Api() Infer
	Handler(ctx context.Context) (int, error)
}

type ServiceT[A Infer] interface {
	Service
	SetInfer(api A)
}

func Do[S ServiceT[A], A Infer](ctx context.Context, service S, reqResp A) (int, error) {
	service = service.New().(S)
	service.SetInfer(reqResp)
	reqResp.InferDataRefs()

	ret, err := service.Handler(ctx)
	return ret, err
}
