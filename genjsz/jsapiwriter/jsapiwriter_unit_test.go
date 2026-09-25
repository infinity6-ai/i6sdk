package jsapiwriter_test

import (
	"testing"

	"github.com/infinity6-ai/gox/routez/sampleroutez/routezsamplefraction"
	"github.com/stretchr/testify/require"
	"go.code.infinity6.ai/sdkgo/genjsz/jsapiwriter"
	"go.code.infinity6.ai/sdkgo/genjsz/jsparser"
	"go.code.infinity6.ai/sdkgo/sdkgo/apiz/ingestzgeturlapi"
)

func TestUnitJsApiWriterIngestz(t *testing.T) {
	api := jsparser.ParseApi(&ingestzgeturlapi.GetUrlApi{})
	code := jsapiwriter.Write(api)

	println(code)

	require.Contains(t, code, "export async function ingestzGetUrl(services, params) {")
	require.Contains(t, code, "path: `api/ingest/get-url/dataset/${encodeURIComponent(String(params.dataset))}`,")
	require.Contains(t, code, "json: params.payload,")
	require.Contains(t, code, "result.payload = await resp.json();")
	require.Contains(t, code, "result.error = text;")
}

func TestUnitJsApiWriterSampleFraction(t *testing.T) {
	api := jsparser.ParseApi(routezsamplefraction.NewService().Api())
	code := jsapiwriter.Write(api)

	require.Contains(t, code, "export async function samplefraction(services, params) {")
	require.Contains(t, code, "path: `api/gox/routez/sample/fraction/${encodeURIComponent(String(params.numerator))}/${encodeURIComponent(String(params.denominator))}`,")
	require.Contains(t, code, "json: params.payload,")
	require.Contains(t, code, "result.payload = await resp.json();")
}
