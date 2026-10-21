package dszstore

import (
	"github.com/infinity6-ai/gox/schemaz/schemaz"
	"go.code.infinity6.ai/sdkgo/sdkgo/dsz/solutiondomain"
	"go.code.infinity6.ai/sdkgo/sdkgo/dsz/storedisplayid"
	"go.code.infinity6.ai/sdkgo/sdkgo/dsz/storeid"
)

type DszStore struct {
	Id              storeid.StoreId
	Name            string
	SolutionDomains []solutiondomain.SolutionDomain
}

func (dss *DszStore) Schema() *schemaz.Schema {
	return &schemaz.Schema{
		Object: func(read bool) map[string]*schemaz.Schema {
			return map[string]*schemaz.Schema{
				"id": {
					Field: func() any { return &dss.Id },
				},
				"name": {
					Field: func() any { return &dss.Name },
				},
				"solution_domains": {
					Array: schemaz.NewArray(schemaz.ArrayOptions[[]solutiondomain.SolutionDomain, solutiondomain.SolutionDomain]{
						Slice: func() *[]solutiondomain.SolutionDomain { return &dss.SolutionDomains },
						Get: func(idx int, element *solutiondomain.SolutionDomain) *schemaz.Schema {
							return &schemaz.Schema{
								Field: func() any { return element },
							}
						},
					}),
				},
			}
		},
	}
}

type DszStoreDisplay struct {
	Id   storedisplayid.StoreDisplayId
	Name string
}
