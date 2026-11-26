package main

import (
	"context"

	"github.com/infinity6-ai/i6sdk/i6cli/cmdcli"
)

func main() {
	ctx := context.Background()
	rootCmd := cmdcli.Prepare(ctx)
	cmdcli.Execute(ctx, rootCmd)
}
