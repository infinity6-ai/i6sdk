package main

import (
	"context"

	"go.code.infinity6.ai/sdkweb/cmdsdkweb"
)

func main() {
	ctx := context.Background()
	rootCmd := cmdsdkweb.Prepare(ctx)
	cmdsdkweb.Execute(rootCmd)
}
