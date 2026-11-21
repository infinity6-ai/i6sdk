package version

import (
	"github.com/infinity6-ai/i6sdk/versionz/version"
)

// Version returns the current component version string.
func Version() string {
	return version.Version()
}
