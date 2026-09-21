package generatorjsz_test

import (
	"testing"

	"github.com/infinity6-ai/gox/httpz/httpzserver"
	"github.com/infinity6-ai/gox/routez/routez"
	"github.com/infinity6-ai/gox/routez/sampleroutez/routezsamplefraction"
	"github.com/stretchr/testify/require"
	"go.code.infinity6.ai/sdkgo/generatorjsz"
)

// invokerPath é o Invoker real do SDK web legacy, relativo a este pacote (cwd de "go test").
const invokerPath = "../sdkweblegacy/src/invoker/invoker.js"

// TestUnitSampleFractionJsClient gera o código JS a partir do schema atual e roda ele de verdade
// contra um servidor real, via node e pelo Invoker real do SDK — se o schema mudar sem o gerador
// ser ajustado, ou o Invoker deixar de montar query/headers/corpo direito, este teste quebra.
// A samplefraction é a única api de exemplo que usa path + query + header + corpo + header de
// resposta ao mesmo tempo.
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
	// nomes reais de wire, não os do mockup antigo (que usava "trace-id"/"req-id" hardcoded)
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

// TestUnitSampleFractionJsClientNotOk confirma que a função gerada devolve o status em vez de
// lançar quando o servidor responde não-2xx (o body não é lido, não segue o formato do RespBody).
func TestUnitSampleFractionJsClientNotOk(t *testing.T) {
	ctx := t.Context()
	s := httpzserver.New(ctx, httpzserver.Options{})
	defer s.Close()
	s.Listen()
	s.Start()

	// nenhuma rota registrada: o servidor responde 404 para tudo.
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
