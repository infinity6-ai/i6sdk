package jsgotypez

import "reflect"

type GoTypKind string

const (
	GoTypKindNone  = ""
	GoTypKindSlice = "slice"
	GoTypKindMap   = "map"
)

type GoTyp struct {
	Kind GoTypKind `json:"kind"`
	Raw  string    `json:"raw"`
}

func NewSlice(v any) *GoTyp {
	val := reflect.ValueOf(v)
	return &GoTyp{
		Kind: GoTypKindSlice,
		Raw:  val.String(),
	}
}
