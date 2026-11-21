package laketbnamez_test

import (
	"testing"

	"github.com/infinity6-ai/gox/commonz/validation/checked/tuchecked"
	"github.com/infinity6-ai/i6sdk/sdkcorez/dsz/laketbnamez"
)

func TestUnitTableChecked(t *testing.T) {
	tuchecked.Check(tuchecked.Table[*laketbnamez.TableName, string]{
		Create: func(v string) *laketbnamez.TableName {
			ret := laketbnamez.MustParse(v)
			return &ret
		},
		Valids:   []string{"a1a", "b-b", "c2c"},
		Invalids: []string{"", "a_a", "_1a", "a1_"},
	})
}
