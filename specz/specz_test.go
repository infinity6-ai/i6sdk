package specz_test

import (
	"testing"

	"github.com/stretchr/testify/require"
	"go.code.infinity6.ai/sdkgo/specz"
)

func TestUnitAll(t *testing.T) {
	all := specz.All()
	require.NotEmpty(t, all)

	found := false
	for _, s := range all {
		if s.Id == "samplefraction" {
			found = true
		}
	}
	require.True(t, found, "expected the fraction spec to be present")
}
