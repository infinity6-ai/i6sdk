package cmdsdkweb_test

import (
	"context"
	"flag"
	"os"
	"path/filepath"
	"testing"

	"github.com/infinity6-ai/gox/routez/apispecz"
	"github.com/infinity6-ai/gox/schemaz/schemaz"
	"github.com/infinity6-ai/gox/schemaz/typeastz"
	"github.com/infinity6-ai/i6sdk/mold/moldsdkfunc"
	"github.com/infinity6-ai/i6sdk/sdkweb/cmdsdkweb"
	"github.com/stretchr/testify/require"
)

var updateGolden = flag.Bool("update", false, "update golden test files")

func TestUnitDocsGolden(t *testing.T) {
	type testScenario struct {
		name       string
		generated  string
		goldenPath string
	}

	ctx := context.Background()
	outDir := t.TempDir()

	err := cmdsdkweb.RunDocs(ctx, cmdsdkweb.DocsOptions{
		SdkJsonPath:  filepath.Join("..", "target", "sdk.json"),
		OutDir:       outDir,
		BaseUrl:      "https://cdn.infinity6.ai/docs",
		AllowMissing: true,
	})
	require.NoError(t, err)

	check := func(t *testing.T, s testScenario) {
		t.Helper()

		actualBytes, err := os.ReadFile(filepath.Join(outDir, s.generated))
		require.NoError(t, err)

		if *updateGolden {
			require.NoError(t, os.WriteFile(s.goldenPath, actualBytes, 0644))
		}

		expectedBytes, err := os.ReadFile(s.goldenPath)
		require.NoError(t, err)

		require.Equal(t, string(expectedBytes), string(actualBytes))
	}

	t.Run("API method page matches golden file", func(t *testing.T) {
		check(t, testScenario{
			name:       "dsz-create-customer-portal",
			generated:  filepath.Join("apis", "dsz-create-customer-portal.html"),
			goldenPath: filepath.Join("testdata", "dsz-create-customer-portal.golden.html"),
		})
	})

	t.Run("Handwritten method page matches golden file", func(t *testing.T) {
		check(t, testScenario{
			name:       "infer-run",
			generated:  "infer-run.html",
			goldenPath: filepath.Join("testdata", "infer-run.golden.html"),
		})
	})
}

func TestUnitDocsCheck(t *testing.T) {
	type testScenario struct {
		name        string
		api         *apispecz.Api
		sdkFunc     *moldsdkfunc.SdkFunc
		errContains string
	}

	check := func(t *testing.T, s testScenario) {
		t.Helper()

		var missingDocs []string
		checkFieldDoc := func(id, field, summary string) {
			if summary == "" {
				missingDocs = append(missingDocs, id+": "+field+": missing doc")
			}
		}

		if s.api != nil {
			if s.api.Spec.Desc == nil || s.api.Spec.Desc.Summary == "" {
				missingDocs = append(missingDocs, s.api.Spec.Id+": missing doc")
			}
			cmdsdkweb.CheckApiTypedefFields(s.api, checkFieldDoc)
		}

		if s.sdkFunc != nil {
			if s.sdkFunc.Spec.Desc == nil || s.sdkFunc.Spec.Desc.Summary == "" {
				missingDocs = append(missingDocs, s.sdkFunc.Spec.Id+": missing doc")
			}
			cmdsdkweb.CheckSdkFuncTypedefFields(s.sdkFunc, checkFieldDoc)
		}

		require.NotEmpty(t, missingDocs)
		require.Contains(t, missingDocs[0], s.errContains)
	}

	t.Run("Rule 1: Missing API summary is reported", func(t *testing.T) {
		check(t, testScenario{
			name: "missing api summary",
			api: &apispecz.Api{
				Spec: apispecz.ApiSpec{
					Id:   "test-api",
					Desc: &schemaz.Desc{Summary: ""},
				},
			},
			errContains: "test-api: missing doc",
		})
	})

	t.Run("Rule 2: Missing SdkFunc summary is reported", func(t *testing.T) {
		check(t, testScenario{
			name: "missing func summary",
			sdkFunc: &moldsdkfunc.SdkFunc{
				Spec: moldsdkfunc.SdkFuncSpec{
					Id:   "test-func",
					Desc: &schemaz.Desc{Summary: ""},
				},
			},
			errContains: "test-func: missing doc",
		})
	})

	t.Run("Rule 3: Missing leaf field summary is reported", func(t *testing.T) {
		paramsTdf := &typeastz.Typedef{
			Head: &typeastz.Tuple{Typ: "object", Name: "apiTestParams"},
			Properties: typeastz.NewTuples([]typeastz.Tuple{
				{Name: "missing_field", Typ: "string", Summary: ""},
			}),
		}
		tdfs := typeastz.NewTypedefs([]*typeastz.Typedef{paramsTdf})

		check(t, testScenario{
			name: "missing leaf field summary",
			api: &apispecz.Api{
				Spec: apispecz.ApiSpec{
					Id:   "test-api-field",
					Desc: &schemaz.Desc{Summary: "Valid API summary"},
				},
				Typedefs: tdfs,
				Req: apispecz.Req{
					Params: typeastz.Typedef{
						Head: &typeastz.Tuple{Typ: "apiTestParams"},
					},
				},
			},
			errContains: "test-api-field: params.missing_field: missing doc",
		})
	})
}
