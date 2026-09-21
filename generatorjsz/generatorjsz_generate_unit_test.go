package generatorjsz_test

import (
	"testing"

	"github.com/infinity6-ai/gox/routez/apiz"
	"github.com/infinity6-ai/gox/routez/sampleroutez/routezsamplefraction"
	"github.com/stretchr/testify/require"
	"go.code.infinity6.ai/sdkgo/generatorjsz"
)

// entry é um grupo de Services cujo SDK JS deve ser gerado no mesmo diretório de saída.
type entry struct {
	Services []apiz.Service
	// OutDir é relativo ao diretório deste pacote (cwd de "go test").
	OutDir string
	// ExampleParams, se não-nil, gera um example.js chamando o SDK com esses valores.
	ExampleParams map[string]any
}

// entries lista todas as Apis conhecidas para geração de SDK JS. Cada nova apiz.Api real
// entra aqui como uma nova entry — não existe descoberta automática (Go não tem plugin
// discovery em runtime), então esta lista é o único ponto de registro.
func entries() []entry {
	return []entry{
		{
			Services: routezsamplefraction.Services(),
			OutDir:   "gen/routezsamplefraction",
			ExampleParams: map[string]any{
				"numerator":     10,
				"denominator":   3,
				"precision":     3,
				"x_i6_trace_id": "xx",
				"reason":        "myreason",
			},
		},
	}
}

// TestUnitGenerateJsSdk (re)gera o SDK JS de cada entry registrada em entries() — não existe
// comando de CLI/go generate para isso, rodar os testes já regenera o gen/ de todo mundo.
func TestUnitGenerateJsSdk(t *testing.T) {
	for _, entry := range entries() {
		for _, service := range entry.Services {
			// New() é obrigatório: Services() devolve serviços com Api() nil, só New() aloca de verdade.
			svc := service.New()
			api := svc.Api()

			fnName, code, err := generatorjsz.BuildJsFetch(api.ApiSpec(), api.GetDataRefs())
			require.NoError(t, err)
			require.NoError(t, generatorjsz.WriteSdk(entry.OutDir, fnName, code))

			if entry.ExampleParams != nil {
				require.NoError(t, generatorjsz.WriteExample(entry.OutDir, fnName, entry.ExampleParams))
			}
		}
	}
}
