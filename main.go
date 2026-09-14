//go:build !testcover

package main

import (
	"context"

	"go.code.infinity6.ai/sdkgo/cmdsdkgo"
)

func main() {
	ctx := context.Background()
	rootCmd := cmdsdkgo.Prepare(ctx)
	cmdsdkgo.Execute(rootCmd)
}
