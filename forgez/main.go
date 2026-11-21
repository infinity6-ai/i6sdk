package main

import (
	"context"

	"github.com/infinity6-ai/i6sdk/forgez/cmdforgez"
)

func main() {
	ctx := context.Background()
	rootCmd := cmdforgez.Prepare(ctx)
	cmdforgez.Execute(ctx, rootCmd)
}
