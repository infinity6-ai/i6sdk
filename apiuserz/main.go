package main

import (
	"context"

	"github.com/infinity6-ai/i6sdk/apiuserz/cmdapiuserz"
)

func main() {
	ctx := context.Background()
	rootCmd := cmdapiuserz.Prepare(ctx)
	cmdapiuserz.Execute(ctx, rootCmd)
}
