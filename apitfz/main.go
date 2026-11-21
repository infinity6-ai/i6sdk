package main

import (
	"context"

	"github.com/infinity6-ai/i6sdk/apitfz/cmdapitfz"
)

func main() {
	ctx := context.Background()
	rootCmd := cmdapitfz.Prepare(ctx)
	cmdapitfz.Execute(ctx, rootCmd)
}
