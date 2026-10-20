package dszstore

import "go.code.infinity6.ai/sdkgo/sdkgo/storedisplayid"

type DszStore struct {
	Id   string `json:"id"`
	Name string `json:"name"`
}

type DszStoreDisplay struct {
	Id   storedisplayid.StoreDisplayId `json:"id"`
	Name string                        `json:"name"`
}
