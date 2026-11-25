package tabledataz

import (
	"time"

	"github.com/infinity6-ai/gox/commonz/strconvz"
	"github.com/infinity6-ai/gox/schemaz/schemaz"
	"github.com/infinity6-ai/i6sdk/sdkcorez/dsz/lakedsz"
	"github.com/infinity6-ai/i6sdk/sdkcorez/dsz/laketbnamez"
)

type TableData interface {
	SchemaQuery() *schemaz.Schema
	SchemaBody() *schemaz.Schema
	Table() laketbnamez.TableName
	SetDataset(dataset lakedsz.Dataset)
	SetIngestedAt(ingestedAt string)
	Validate() error
	GetCreatedAt() *time.Time
}

func NewBqNullString(val string) *string {
	return &val
}

func ParseBqNullString(out **string) func() *schemaz.Parser[string] {
	return func() *schemaz.Parser[string] {
		return &schemaz.Parser[string]{
			Parse: func(v string) {
				if v == "" {
					*out = nil
				} else {
					*out = &v
				}
			},
			Format: func() (string, any) {
				if *out == nil {
					return "", *out
				}
				return **out, *out
			},
		}
	}
}

func NewBqNullFloat64(val float64) *float64 {
	return &val
}

func NewBqNullInt64(val int64) *int64 {
	return &val
}

func ParseBqNullFloat64(out **float64) func() *schemaz.Parser[string] {
	return func() *schemaz.Parser[string] {
		return &schemaz.Parser[string]{
			Parse: func(v string) {
				if v == "" {
					*out = nil
					return
				}
				var f float64
				strconvz.MustParseNumberInto(v, &f)
				*out = &f
			},
			Format: func() (string, any) {
				if *out == nil {
					return "", *out
				}
				return strconvz.FormatNumber(**out), *out
			},
		}
	}
}

func ParseBqNullInt64(out **int64) func() *schemaz.Parser[string] {
	return func() *schemaz.Parser[string] {
		return &schemaz.Parser[string]{
			Parse: func(v string) {
				if v == "" {
					*out = nil
					return
				}
				var i int64
				strconvz.MustParseNumberInto(v, &i)
				*out = &i
			},
			Format: func() (string, any) {
				if *out == nil {
					return "", *out
				}
				return strconvz.FormatNumber(**out), *out
			},
		}
	}
}

func ParseTime(out *time.Time) func() *schemaz.Parser[string] {
	return func() *schemaz.Parser[string] {
		return &schemaz.Parser[string]{
			Parse: func(v string) {
				t, err := time.Parse(time.RFC3339, v)
				if err == nil {
					*out = t
				}
			},
			Format: func() (string, any) {
				if out.IsZero() {
					return "", *out
				}
				return out.Format(time.RFC3339), *out
			},
		}
	}
}
