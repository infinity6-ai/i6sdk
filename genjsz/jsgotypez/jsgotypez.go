package jsgotypez

type GoTypKind string

const (
	GoTypKindRaw    = "raw"
	GoTypKindSlice  = "slice"
	GoTypKindMap    = "map"
	GoTypKindObject = "object"
)

type GoTyp struct {
	Kind GoTypKind `json:"kind"`
	Raw  string    `json:"raw"`
}

func NewSlice(raw string) *GoTyp {
	return &GoTyp{
		Kind: GoTypKindSlice,
		Raw:  raw,
	}
}

func NewObect(raw string) *GoTyp {
	return &GoTyp{
		Kind: GoTypKindObject,
		Raw:  raw,
	}
}

func NewMap(k string, v string) *GoTyp {
	return &GoTyp{
		Kind: GoTypKindMap,
		Raw:  k + "," + v,
	}
}
