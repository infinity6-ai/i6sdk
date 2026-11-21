//go:generate go run . generate
//go:generate go run github.com/infinity6-ai/gox/commonz staticz generate --dir src/site
package main

import (
	"context"

	"github.com/infinity6-ai/i6sdk/sdkweb/cmdsdkweb"
)

func main() {
	ctx := context.Background()
	rootCmd := cmdsdkweb.Prepare(ctx)
	cmdsdkweb.Execute(rootCmd)
}
