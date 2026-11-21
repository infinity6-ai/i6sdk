package main

import (
	"context"

	"github.com/infinity6-ai/i6sdk/apidatasetz/cmdapidatasetz"
)

func main() {
	ctx := context.Background()
	rootCmd := cmdapidatasetz.Prepare(ctx)
	cmdapidatasetz.Execute(ctx, rootCmd)
}
