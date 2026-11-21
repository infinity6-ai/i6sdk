package main

import (
	"context"

	"github.com/infinity6-ai/i6sdk/apicollectionz/cmdapicollectionz"
)

func main() {
	ctx := context.Background()
	rootCmd := cmdapicollectionz.Prepare(ctx)
	cmdapicollectionz.Execute(ctx, rootCmd)
}
