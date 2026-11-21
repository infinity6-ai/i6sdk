package main

import (
	"context"

	"github.com/infinity6-ai/i6sdk/sdkcorez/cmdsdkcorez"
)

func main() {
	ctx := context.Background()
	rootCmd := cmdsdkcorez.Prepare(ctx)
	cmdsdkcorez.Execute(ctx, rootCmd)
}
