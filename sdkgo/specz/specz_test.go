package specz_test

import (
	"testing"

	"github.com/stretchr/testify/require"
	"go.code.infinity6.ai/sdkgo/sdkgo/specz"
)

func TestUnitAll(t *testing.T) {
	all := specz.All()
	require.NotEmpty(t, all)

	found := false
	for _, s := range all {
		if s.ApiSpec().Id == "samplefraction" {
			found = true
		}
	}
	require.True(t, found, "expected the fraction spec to be present")
}

func TestUnitSdk(t *testing.T) {
	ids := []string{}
	for _, api := range specz.Sdk() {
		ids = append(ids, api.ApiSpec().Id)
	}
	require.ElementsMatch(t, []string{"ingestz-get-url", "pipez-start", "dsz-domain-select"}, ids)
}

func TestUnitSdkReturnsFreshInstances(t *testing.T) {
	first := specz.Sdk()
	second := specz.Sdk()
	require.Len(t, second, len(first))
	for i := range first {
		require.NotSame(t, first[i], second[i], "Sdk() must not share instances between calls")
	}
}
