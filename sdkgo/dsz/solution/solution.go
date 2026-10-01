package solution

import (
	"github.com/infinity6-ai/gox/commonz/constraintz/optionalz"
	"github.com/infinity6-ai/gox/commonz/validation"
	"github.com/infinity6-ai/gox/commonz/validation/checked"
	"github.com/infinity6-ai/gox/schemaz/schemaz"
	"go.code.infinity6.ai/sdkgo/genjsz/jstypez"
)

type Solution checked.Value[string]

func (*Solution) JSType() string {
	return jstypez.TypString
}

func (d Solution) Optional() optionalz.Optional[string] {
	return checked.Value[string](d).Optional()
}

func (d Solution) Get() string {
	return checked.Value[string](d).Get()
}

func (d Solution) String() string {
	return d.Get()
}

func (d Solution) MarshalJSON() ([]byte, error) {
	return checked.Marshal[string](d)
}

func (d *Solution) UnmarshalJSON(data []byte) error {
	return checked.Unmarshal(d, data)
}

func (d Solution) Validate(v string) error {
	if err := validation.StrNotEmpty(v, "Solution cannot be empty"); err != nil {
		return err
	}
	if err := validation.OneOf([]string{"relevance", "forecast"}, v, "Solution"); err != nil {
		return err
	}
	return nil
}

func MustParse(id string) Solution {
	var ret Solution
	checked.MustSet(&ret, id)
	return ret
}

func ParseStr(out *Solution) func() *schemaz.Parser[string] {
	return func() *schemaz.Parser[string] {
		return &schemaz.Parser[string]{
			Parse: func(v string) {
				*out = MustParse(v)
			},
			Format: func() (string, any) {
				return out.Get(), *out
			},
		}
	}
}
