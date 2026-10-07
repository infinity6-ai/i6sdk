package lakedspathz_test

import (
	"testing"

	"github.com/infinity6-ai/gox/commonz/validation/checked/tuchecked"
	"go.code.infinity6.ai/sdkgo/sdkgo/dsz/lakedspathz"
)

func TestUnitTableChecked(t *testing.T) {
	tuchecked.Check(tuchecked.Table[*lakedspathz.Path, string]{
		Create: func(v string) *lakedspathz.Path {
			ret := lakedspathz.MustParse(v)
			return &ret
		},
		Valids: []string{
			"i6ds://mydataset/bla.ble/a=2/b=x/my/file.txt",
			"i6ds://mydataset/bla.ble/my/file.txt",
			"i6ds://mydataset/bla.ble/a=2/b=x",
			"i6ds://mydataset/bla.ble",
			"i6ds://mydataset/bla.ble/",
		},
		Invalids: []string{
			"", "a_a", "_1a", "a1_",
			"i6ds://mydataset/bla/", "i6ds://mydataset/bla",
			"i6ds://mydataset", "i6ds://mydataset/",
			"i6ds://", "i6ds://mydataset/bla.", "i6ds://mydataset/.ble", "i6ds://mydataset/.",
		},
	})
}
