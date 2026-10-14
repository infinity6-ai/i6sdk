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
