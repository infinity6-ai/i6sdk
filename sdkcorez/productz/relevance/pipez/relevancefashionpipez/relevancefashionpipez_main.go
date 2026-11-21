package relevancefashionpipez

import (
	"fmt"

	"github.com/infinity6-ai/gox/schemaz/schemaz"
	"github.com/infinity6-ai/i6sdk/sdkcorez/dsz/lakedsz"
	"github.com/infinity6-ai/i6sdk/sdkcorez/dsz/solutiondomain"
)

type Main struct {
	Dataset string `json:"dataset"`
}

func (m *Main) SolutionDomain() solutiondomain.SolutionDomain {
	return solutiondomain.MustParseString("relevance_fashion")
}

func (m *Main) Pipename() string {
	return "relevance-fashion-main"
}

func (m *Main) Concurrent() string {
	return fmt.Sprintf("a.%s.d.%s", m.Pipename(), m.Dataset)
}

func (m *Main) SetDataset(ds lakedsz.Dataset) {
	m.Dataset = ds.String()
}

func (m *Main) Params() map[string]string {
	return map[string]string{
		"dataset": m.Dataset,
	}
}

func (m *Main) Schema() *schemaz.Schema {
	return nil
}
