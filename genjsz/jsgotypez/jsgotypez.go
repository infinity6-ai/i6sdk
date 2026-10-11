package jsgotypez

type Raw string

const (
	None  Raw = ""
	Slice Raw = "slice"
	Map   Raw = "map"
)

type Kind string

const (
	String  Kind = "string"
	Float32 Kind = "float32"
)

type GoTyp struct {
	Raw  Raw
	Kind Kind
}
