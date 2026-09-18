package appstack_test

import (
	"context"
	"testing"

	"go.code.infinity6.ai/sdkweb/internal/appstack"
	"go.code.infinity6.ai/web/api/webchecker"
)

func TestWebChecker(t *testing.T) {
	ctx := context.Background()

	server := appstack.StartForTestWeb(ctx)
	defer server.Stop()

	webchecker.RunWeb(ctx, server.Client().Base(), "test.html")
}
