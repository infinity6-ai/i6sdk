package jsgotypez

type Kind string

const (
	KindNone  Kind = ""
	KindSlice Kind = "slice"
	KindMap   Kind = "map"
)

type Raw string

const (
	RawNone    Raw = ""
	RawString  Raw = "string"
	RawFloat32 Raw = "float32"
)

type GoTyp struct {
	Raw   Raw  `json:"raw,omitempty"`
	Key   Kind `json:"key,omitempty"`
	Value Kind `json:"value"`
}
