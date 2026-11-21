package version_test

import (
	"testing"

	"github.com/infinity6-ai/i6sdk/versionz/version"
	"github.com/stretchr/testify/require"
)

func TestUnitVersion(t *testing.T) {
	v := version.Version()
	require.NotEmpty(t, v)
	require.Regexp(t, `^v\d+\.\d+\.\d+$`, v)
}
