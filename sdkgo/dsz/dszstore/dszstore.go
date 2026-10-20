package dszstore

import (
	"go.code.infinity6.ai/sdkgo/sdkgo/storedisplayid"
	"go.code.infinity6.ai/sdkgo/sdkgo/storeid"
)

type DszStore struct {
	Id   storeid.StoreId `json:"id"`
	Name string          `json:"name"`
}

type DszStoreDisplay struct {
	Id   storedisplayid.StoreDisplayId `json:"id"`
	Name string                        `json:"name"`
}
