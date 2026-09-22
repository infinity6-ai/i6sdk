package generatorjsz_test

import (
	"testing"

	"github.com/infinity6-ai/gox/httpz/httpzserver"
	"github.com/infinity6-ai/gox/routez/routez"
	"github.com/infinity6-ai/gox/routez/sampleroutez/routezsamplefraction"
	"github.com/stretchr/testify/require"
	"go.code.infinity6.ai/sdkgo/generatorjsz"
)

const invokerPath = "../sdkweblegacy/src/invoker/invoker.js"

func TestUnitSampleFractionJsClient(t *testing.T) {
	ctx := t.Context()
	s := httpzserver.New(ctx, httpzserver.Options{})
	defer s.Close()
	s.Listen()
	s.Start()

	routez.Register(s, routezsamplefraction.Services()...)

	svc := routezsamplefraction.Services()[0].New()
	fnName, code, err := generatorjsz.BuildJsFetch(svc.Api().ApiSpec(), svc.Api().GetDataRefs())
	require.NoError(t, err)

	require.Equal(t, "samplefraction", fnName)
	require.Contains(t, code, "@property {number} numerator")
	require.Contains(t, code, "@property {string} reason")
	require.Contains(t, code, "@property {{x_i6_trace_message: string}} headers")
	require.Contains(t, code, "@property {{display: string, result: string}} [body]")
	require.Contains(t, code, `headers: {"x-i6-trace-id": String(params.x_i6_trace_id)}`)
	require.Contains(t, code, `query: {"precision": params.precision}`)
	require.Contains(t, code, `resp.headers.get("x-i6-trace-message")`)

	sampleParams := map[string]any{
		"numerator":     10,
		"denominator":   3,
		"precision":     3,
		"x_i6_trace_id": "xx",
		"reason":        "myreason",
	}

	result := runJs(t, fnName, code, invokerPath, s.Base().String(), sampleParams)

	require.InDelta(t, 201, result["status"], 0)
	require.Equal(t, true, result["ok"])
	require.Equal(t, "reason: myreason, trace: xx", result["headers"].(map[string]any)["x_i6_trace_message"])
	body := result["body"].(map[string]any)
	require.Equal(t, "10.000/3.000", body["display"])
	require.Equal(t, "3.333", body["result"])
}

func TestUnitSampleFractionJsClientNotOk(t *testing.T) {
	ctx := t.Context()
	s := httpzserver.New(ctx, httpzserver.Options{})
	defer s.Close()
	s.Listen()
	s.Start()

	svc := routezsamplefraction.Services()[0].New()
	fnName, code, err := generatorjsz.BuildJsFetch(svc.Api().ApiSpec(), svc.Api().GetDataRefs())
	require.NoError(t, err)

	result := runJs(t, fnName, code, invokerPath, s.Base().String(), map[string]any{
		"numerator":     10,
		"denominator":   3,
		"precision":     3,
		"x_i6_trace_id": "xx",
		"reason":        "myreason",
	})

	require.InDelta(t, 404, result["status"], 0)
	require.Equal(t, false, result["ok"])
	require.NotContains(t, result, "body")
}
