package relavancefashiondashfunnel

import "go.code.infinity6.ai/sdkgo/sdkgo/dsz/lakedsz"

type Params struct {
	MinYearMonth string
	MaxYearMonth string
}

type Partitions struct {
	Version string
	StoreId string
}

type Row struct {
	EventType string
	FromRec   string
	QtEvents  string
}

type Result struct {
	Rows []*Row
}

type Req struct {
	Dataset        lakedsz.Dataset
	DashPartitions Partitions
	Params         []*Params
}

type Resp struct {
	Results Result
}

type Func struct {
	Req  Req
	Resp Resp
}
