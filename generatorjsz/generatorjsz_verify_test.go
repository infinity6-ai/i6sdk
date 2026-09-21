package generatorjsz_test

import (
	"bytes"
	"encoding/json"
	"fmt"
	"net/url"
	"os/exec"
	"path/filepath"
	"testing"

	"github.com/infinity6-ai/gox/commonz/filez"
	"github.com/stretchr/testify/require"
)

// runJs executa o código JS gerado com o Node contra um servidor real, e devolve o resultado —
// é como um teste confirma que o código gerado realmente funciona, não só compila. A função
// gerada roda com o Invoker de verdade (invokerPath aponta para invoker.js do SDK) e um Services
// mínimo que só sabe a baseUrl, então o que é testado é gerador + Invoker juntos.
func runJs(t testing.TB, fnName, code, invokerPath, baseUrl string, params map[string]any) map[string]any {
	t.Helper()
	if _, err := exec.LookPath("node"); err != nil {
		t.Skip("node is not installed")
	}
	invokerAbs, err := filepath.Abs(invokerPath)
	require.NoError(t, err)
	invokerUrl := (&url.URL{Scheme: "file", Path: invokerAbs}).String()

	paramsJson, err := json.Marshal(params)
	require.NoError(t, err)

	driver := fmt.Sprintf(`import { Invoker } from %q;
%s
const services = { config: () => ({ baseUrl: %q }), invoker: () => new Invoker(services) };
%s(services, %s).then(r => console.log(JSON.stringify(r))).catch(e => { console.error(e); process.exit(1); });
`, invokerUrl, code, baseUrl, fnName, paramsJson)
	scriptPath := filepath.Join(t.TempDir(), "client.mjs")
	require.NoError(t, filez.WriteFile(scriptPath, []byte(driver)))

	out, err := exec.Command("node", scriptPath).CombinedOutput()
	require.NoError(t, err, string(out))

	var result map[string]any
	require.NoError(t, json.Unmarshal(bytes.TrimSpace(out), &result))
	return result
}
