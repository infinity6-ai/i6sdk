//go:build !testcover

package main

import (
	"context"

	"github.com/infinity6-ai/i6sdk/cmdz/cmdcmdz"
)

func main() {
	ctx := context.Background()
	rootCmd := cmdcmdz.Prepare(ctx)
	cmdcmdz.Execute(ctx, rootCmd)
}
