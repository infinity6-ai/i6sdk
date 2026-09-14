//go:build !testcover

package main

import (
	"fmt"

	"go.code.infinity6.ai/sdkgo/version"
)

func main() {
	fmt.Println(version.Version())
}
