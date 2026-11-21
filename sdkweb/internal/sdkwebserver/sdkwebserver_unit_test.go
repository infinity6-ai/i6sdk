package sdkwebserver_test

import (
	"context"
	"net/http"
	"testing"

	"github.com/infinity6-ai/i6sdk/sdkweb/internal/sdkwebserver"
	"github.com/stretchr/testify/require"
)

func TestUnitSdkWebServer(t *testing.T) {
	ctx := context.Background()
	srv := sdkwebserver.Start(sdkwebserver.Options{
		Context:      ctx,
		LocalAddress: "127.0.0.1:0",
	})
	defer srv.Stop()

	require.NotEmpty(t, srv.Address())
	require.NotEmpty(t, srv.Base())

	client := &http.Client{
		CheckRedirect: func(req *http.Request, via []*http.Request) error {
			return http.ErrUseLastResponse
		},
	}
	resp, err := client.Get(srv.Base() + "/")
	require.NoError(t, err)
	defer resp.Body.Close()
	require.Equal(t, http.StatusTemporaryRedirect, resp.StatusCode)
	require.Equal(t, "/static/index.html", resp.Header.Get("Location"))

	respWebhook, err := http.Post(srv.Base()+"/api/stripe/webhook", "application/json", nil)
	require.NoError(t, err)
	defer respWebhook.Body.Close()
	require.Equal(t, http.StatusNoContent, respWebhook.StatusCode)
}
