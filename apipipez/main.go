package main

import (
	"context"

	"github.com/infinity6-ai/i6sdk/apipipez/cmdapipipez"
)

func main() {
	ctx := context.Background()
	rootCmd := cmdapipipez.Prepare(ctx)
	cmdapipipez.Execute(ctx, rootCmd)
}
