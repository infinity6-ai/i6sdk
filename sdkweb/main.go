package main

import (
	"context"

	"github.com/infinity6-ai/i6sdk/sdkweb/cmdsdkweb"
)

func main() {
	ctx := context.Background()
	rootCmd := cmdsdkweb.Prepare(ctx)
	cmdsdkweb.Execute(rootCmd)
}
