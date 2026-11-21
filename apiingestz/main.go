package main

import (
	"context"

	"github.com/infinity6-ai/i6sdk/apiingestz/cmdapiingestz"
)

func main() {
	ctx := context.Background()
	rootCmd := cmdapiingestz.Prepare(ctx)
	cmdapiingestz.Execute(ctx, rootCmd)
}
