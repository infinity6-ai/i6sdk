package jsschema

import (
	"cloud.google.com/go/bigquery"
	"github.com/infinity6-ai/gox/fsz/fsz"
	"github.com/infinity6-ai/gox/schemaz/jsschemaz"
	"github.com/infinity6-ai/gox/schemaz/jstypez"
	"github.com/infinity6-ai/gox/schemaz/schemaz"
)

type Parser = jsschemaz.Parser

func New() *Parser {
	return jsschemaz.New()
}

func NewSdk() *Parser {
	return jsschemaz.New().
		AddKnownType(fsz.FileStat{}, jstypez.TypFileStat).
		AddKnownType(bigquery.NullString{}, jstypez.TypString).
		AddKnownType(bigquery.NullInt64{}, jstypez.TypNumber).
		AddKnownType(bigquery.NullFloat64{}, jstypez.TypNumber)
}

func DescName(d *schemaz.Desc) string {
	return jsschemaz.DescName(d)
}

func DescSummary(d *schemaz.Desc) string {
	return jsschemaz.DescSummary(d)
}

func FormatTypedefs(pipedefs *jstypez.Typedefs) string {
	return jsschemaz.FormatTypedefs(pipedefs)
}
