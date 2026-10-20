package storeid_test

import (
	"testing"

	"github.com/infinity6-ai/gox/commonz/validation/checked/tuchecked"
	"github.com/stretchr/testify/require"
	"go.code.infinity6.ai/sdkgo/sdkgo/storeid"
)

func TestUnitStoreIdChecked(t *testing.T) {
	tuchecked.Check(tuchecked.Table[*storeid.StoreId, string]{
		Create: func(v string) *storeid.StoreId {
			ret := storeid.MustParse(v)
			return &ret
		},
		Valids:   []string{"store1", "store-123", "ABC_XYZ", "12345", ""},
		Invalids: []string{"store!", "store space", "a@b", "store.123"},
	})
}

func TestUnitStoreIdValidation(t *testing.T) {
	type testScenario struct {
		input  string
		want   string
		errMsg string
	}

	check := func(t *testing.T, s testScenario) {
		t.Helper()
		if s.errMsg != "" {
			require.Panics(t, func() {
				storeid.MustParse(s.input)
			})
			return
		}
		got := storeid.MustParse(s.input)
		require.Equal(t, s.want, got.String())
	}

	t.Run("Valid simple store id", func(t *testing.T) {
		check(t, testScenario{
			input: "store1",
			want:  "store1",
		})
	})

	t.Run("Valid store id with dashes", func(t *testing.T) {
		check(t, testScenario{
			input: "store-123",
			want:  "store-123",
		})
	})

	t.Run("Valid store id with underscores", func(t *testing.T) {
		check(t, testScenario{
			input: "store_abc",
			want:  "store_abc",
		})
	})

	t.Run("Empty store id is valid", func(t *testing.T) {
		check(t, testScenario{
			input: "",
			want:  "",
		})
	})

	t.Run("Store id with invalid characters returns error", func(t *testing.T) {
		check(t, testScenario{
			input:  "store.1",
			errMsg: "StoreId does not match storeIdValidator",
		})
	})
}
