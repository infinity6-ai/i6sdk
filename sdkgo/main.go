package main

import (
	"context"

	"github.com/infinity6-ai/i6sdk/sdkgo/cmdsdkgo"
)

func main() {
	ctx := context.Background()
	rootCmd := cmdsdkgo.Prepare(ctx)
	cmdsdkgo.Execute(ctx, rootCmd)
}
