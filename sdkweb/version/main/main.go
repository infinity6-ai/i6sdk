//go:build !testcover

package main

import (
	"fmt"

	"github.com/infinity6-ai/i6sdk/sdkweb/version"
)

func main() {
	fmt.Println(version.Version())
}
