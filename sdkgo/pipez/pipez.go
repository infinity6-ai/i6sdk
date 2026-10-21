package pipez

import (
	"github.com/infinity6-ai/gox/schemaz/schemaz"
	"go.code.infinity6.ai/sdkgo/sdkgo/dsz/lakedsz"
	"go.code.infinity6.ai/sdkgo/sdkgo/dsz/solutiondomain"
)

type Pipe interface {
	SolutionDomain() solutiondomain.SolutionDomain
	Pipename() string
	Concurrent() string
	Schema() *schemaz.Schema
	SetDataset(dataset lakedsz.Dataset)
	Params() map[string]string
}

type PipeState struct {
	Status    string         `json:"status"`
	Params    map[string]any `json:"params"`
	CreatedAt string         `json:"created_at"`
}

type Pipeline struct {
	Id             string                        `json:"id"`
	ConcurrentId   string                        `json:"concurrent_id"`
	SolutionDomain solutiondomain.SolutionDomain `json:"solution_domain"`
	Name           string                        `json:"name"`
	Dataset        string                        `json:"dataset"`
	States         []*PipeState                  `json:"states"`
}

func (p *Pipeline) Schema() *schemaz.Schema {
	return &schemaz.Schema{
		Object: func(read bool) map[string]*schemaz.Schema {
			return map[string]*schemaz.Schema{
				"id": {
					Field: func() any { return &p.Id },
				},
				"concurrent_id": {
					Field: func() any { return &p.ConcurrentId },
				},
				"solution_domain": {
					Field: func() any { return &p.SolutionDomain },
				},
				"name": {
					Field: func() any { return &p.Name },
				},
				"dataset": {
					Field: func() any { return &p.Dataset },
				},
				"states": {
					Array: schemaz.NewArray(schemaz.ArrayOptions[[]*PipeState, *PipeState]{
						Slice:      func() *[]*PipeState { return &p.States },
						NewElement: func(idx int) *PipeState { return &PipeState{} },
						Get: func(idx int, element **PipeState) *schemaz.Schema {
							return &schemaz.Schema{
								Object: func(_ bool) map[string]*schemaz.Schema {
									return map[string]*schemaz.Schema{
										"status": {
											Field: func() any { return &(*element).Status },
										},
										"params": {
											Field: func() any { return &(*element).Params },
										},
										"created_at": {
											Field: func() any { return &(*element).CreatedAt },
										},
									}
								},
							}
						},
					}),
				},
			}
		},
	}
}
