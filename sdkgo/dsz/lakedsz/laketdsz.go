package lakedsz

import (
	"regexp"

	"github.com/infinity6-ai/gox/commonz/constraintz/optionalz"
	"github.com/infinity6-ai/gox/commonz/validation"
	"github.com/infinity6-ai/gox/commonz/validation/checked"
	"github.com/infinity6-ai/gox/schemaz/schemaz"
)

var validator = regexp.MustCompile(`^[a-z][a-z0-9\-]{1,}[a-z0-9]$`)

type Dataset checked.Value[string]

func (d Dataset) Optional() optionalz.Optional[string] {
	return checked.Value[string](d).Optional()
}

func (d Dataset) Get() string {
	return checked.Value[string](d).Get()
}

func (d Dataset) String() string {
	return d.Get()
}

func (d Dataset) MarshalJSON() ([]byte, error) {
	return checked.Marshal[string](d)
}

func (d *Dataset) UnmarshalJSON(data []byte) error {
	return checked.Unmarshal(d, data)
}

func (d Dataset) Validate(v string) error {
	if err := validation.StrNotEmpty(v, "TableId cannot be empty"); err != nil {
		return err
	}
	if err := validation.RegexMatch(validator, v, "TableId"); err != nil {
		return err
	}
	return nil
}

func MustParse(id string) Dataset {
	var ret Dataset
	checked.MustSet(&ret, id)
	return ret
}

func ParseStr(out *Dataset) func() *schemaz.Parser[string] {
	return func() *schemaz.Parser[string] {
		return &schemaz.Parser[string]{
			Parse: func(v string) {
				*out = MustParse(v)
			},
			Format: func() (string, any) {
				return out.Optional().Or(), *out
			},
		}
	}
}
