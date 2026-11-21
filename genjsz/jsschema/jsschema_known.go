package jsschema

import (
	"context"
	"reflect"
	"time"

	"cloud.google.com/go/bigquery"
	"github.com/infinity6-ai/gox/fsz/fsz"
	"github.com/infinity6-ai/gox/schemaz/schemaz"
	"go.code.infinity6.ai/sdkgo/genjsz/jstypez"
)

type Parser struct {
	known map[reflect.Type]*jstypez.Tuple
}

// New creates a new base Parser initialized with standard library types.
func New() *Parser {
	p := &Parser{
		known: make(map[reflect.Type]*jstypez.Tuple),
	}
	p.AddKnownType(time.Time{}, jstypez.TypTime)
	return p
}

// NewSdk creates a Parser pre-configured with SDK and cloud types (FileStat, BigQuery, etc.).
func NewSdk() *Parser {
	p := New()
	p.AddKnownType(fsz.FileStat{}, jstypez.TypFileStat)
	p.AddKnownType(bigquery.NullString{}, jstypez.TypString)
	p.AddKnownType(bigquery.NullInt64{}, jstypez.TypNumber)
	p.AddKnownType(bigquery.NullFloat64{}, jstypez.TypNumber)
	return p
}

// AddKnownType registers how an external Go type is represented in JavaScript.
// Supports fluent chaining.
func (p *Parser) AddKnownType(sample any, jsType string) *Parser {
	t := reflect.TypeOf(sample)
	p.known[t] = &jstypez.Tuple{
		Typ:   jsType,
		GoTyp: t.String(),
	}
	return p
}

func (p *Parser) knownTypes(ctx context.Context, baseId string, val reflect.Value) (*jstypez.Typedef, *jstypez.Typedefs) {
	if !val.IsValid() {
		return nil, nil
	}
	v := val.Interface()
	switch s := v.(type) {
	case schemaz.Schema:
		a, b := p.ParseSchema(ctx, baseId, &s)
		a.Head.Name = ""
		return a, b
	case *schemaz.Schema:
		a, b := p.ParseSchema(ctx, baseId, s)
		a.Head.Name = ""
		return a, b
	}

	t := reflect.TypeOf(v)
	if tuple, ok := p.known[t]; ok {
		return &jstypez.Typedef{
			Head: &jstypez.Tuple{
				Typ:   tuple.Typ,
				GoTyp: tuple.GoTyp,
			},
		}, nil
	}

	return nil, nil
}
