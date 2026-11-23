//go:generate go run . generate
package main

import (
	"context"

	"github.com/infinity6-ai/i6sdk/sdkpy/cmdsdkpy"
)

func main() {
	ctx := context.Background()
	rootCmd := cmdsdkpy.Prepare(ctx)
	cmdsdkpy.Execute(rootCmd)
}
