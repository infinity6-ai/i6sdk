package main

import (
	"context"

	"github.com/infinity6-ai/i6sdk/mold/cmdmold"
)

func main() {
	ctx := context.Background()
	rootCmd := cmdmold.Prepare(ctx)
	cmdmold.Execute(ctx, rootCmd)
}
