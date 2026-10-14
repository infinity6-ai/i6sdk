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
