package solution_test

import (
	"testing"

	"github.com/infinity6-ai/gox/commonz/validation/checked/tuchecked"
	"go.code.infinity6.ai/sdkgo/sdkgo/dsz/solution"
)

func TestUnitSolutionChecked(t *testing.T) {
	tuchecked.Check(tuchecked.Table[*solution.Solution, string]{
		Create: func(v string) *solution.Solution {
			ret := solution.MustParse(v)
			return &ret
		},
		Valids:   []string{"relevance", "forecast"},
		Invalids: []string{"bla"},
	})
}
