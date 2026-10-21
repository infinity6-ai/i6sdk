package dszstore

import (
	"go.code.infinity6.ai/sdkgo/sdkgo/dsz/solutiondomain"
	"go.code.infinity6.ai/sdkgo/sdkgo/dsz/storedisplayid"
	"go.code.infinity6.ai/sdkgo/sdkgo/dsz/storeid"
)

type DszStore struct {
	Id              storeid.StoreId
	Name            string
	SolutionDomains []solutiondomain.SolutionDomain
}

type DszStoreDisplay struct {
	Id   storedisplayid.StoreDisplayId
	Name string
}
