package tabledataz

import (
	"time"

	"github.com/infinity6-ai/gox/commonz/dataz"
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

func NewBqNullString(val string) dataz.NullString {
	return dataz.NullString{
		StringVal: val,
		Valid:     true,
	}
}

func ParseBqNullString(out *dataz.NullString) func() *schemaz.Parser[string] {
	return func() *schemaz.Parser[string] {
		return &schemaz.Parser[string]{
			Parse: func(v string) {
				*out = dataz.NullString{
					StringVal: v,
					Valid:     v != "",
				}
			},
			Format: func() (string, any) {
				if !out.Valid {
					return "", *out
				}
				return out.StringVal, *out
			},
		}
	}
}

func NewBqNullFloat64(val float64) dataz.NullFloat64 {
	return dataz.NullFloat64{
		Float64: val,
		Valid:   true,
	}
}

func NewBqNullInt64(val int64) dataz.NullInt64 {
	return dataz.NullInt64{
		Int64: val,
		Valid: true,
	}
}

func ParseBqNullFloat64(out *dataz.NullFloat64) func() *schemaz.Parser[string] {
	return func() *schemaz.Parser[string] {
		return &schemaz.Parser[string]{
			Parse: func(v string) {
				if v == "" {
					*out = dataz.NullFloat64{}
					return
				}
				var f float64
				strconvz.MustParseNumberInto(v, &f)
				*out = NewBqNullFloat64(f)
			},
			Format: func() (string, any) {
				if !out.Valid {
					return "", *out
				}
				return strconvz.FormatNumber(out.Float64), *out
			},
		}
	}
}

func ParseBqNullInt64(out *dataz.NullInt64) func() *schemaz.Parser[string] {
	return func() *schemaz.Parser[string] {
		return &schemaz.Parser[string]{
			Parse: func(v string) {
				if v == "" {
					*out = dataz.NullInt64{}
					return
				}
				var i int64
				strconvz.MustParseNumberInto(v, &i)
				*out = NewBqNullInt64(i)
			},
			Format: func() (string, any) {
				if !out.Valid {
					return "", *out
				}
				return strconvz.FormatNumber(out.Int64), *out
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
