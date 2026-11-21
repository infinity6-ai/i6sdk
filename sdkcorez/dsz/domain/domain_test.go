package domain_test

import (
	"testing"

	"github.com/infinity6-ai/gox/commonz/validation/checked/tuchecked"
	"github.com/infinity6-ai/i6sdk/sdkcorez/dsz/domain"
)

func TestUnitDomainChecked(t *testing.T) {
	tuchecked.Check(tuchecked.Table[*domain.Domain, string]{
		Create: func(v string) *domain.Domain {
			ret := domain.MustParse(v)
			return &ret
		},
		Valids:   []string{"fashion"},
		Invalids: []string{"bla"},
	})
}
