package lakepathz_test

import (
	"testing"

	"github.com/infinity6-ai/gox/commonz/validation/checked/tuchecked"
	"go.code.infinity6.ai/sdkgo/sdkgo/dsz/lakepathz"
)

func TestUnitTableChecked(t *testing.T) {
	tuchecked.Check(tuchecked.Table[*lakepathz.Path, string]{
		Create: func(v string) *lakepathz.Path {
			ret := lakepathz.MustParse(v)
			return &ret
		},
		Valids:   []string{"i6tb://bla.ble/a=2/b=x/my/file.txt"},
		Invalids: []string{"", "a_a", "_1a", "a1_"},
	})
}
