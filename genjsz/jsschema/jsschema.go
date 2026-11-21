package jsschema

import (
	"cloud.google.com/go/bigquery"
	"github.com/infinity6-ai/gox/fsz/fsz"
	"github.com/infinity6-ai/gox/schemaz/jsschemaz"
	"github.com/infinity6-ai/gox/schemaz/jstypez"
)

func NewSdk() *jsschemaz.Parser {
	return jsschemaz.New().
		AddKnownType(fsz.FileStat{}, jstypez.TypFileStat).
		AddKnownType(bigquery.NullString{}, jstypez.TypString).
		AddKnownType(bigquery.NullInt64{}, jstypez.TypNumber).
		AddKnownType(bigquery.NullFloat64{}, jstypez.TypNumber)
}
