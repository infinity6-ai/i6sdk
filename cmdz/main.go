//go:build !testcover

package main

import (
	"context"

	"github.com/infinity6-ai/i6sdk/cmdz/cmdsdk"
)

func main() {
	ctx := context.Background()
	rootCmd := cmdsdk.Prepare(ctx)
	cmdsdk.Execute(rootCmd)
}
