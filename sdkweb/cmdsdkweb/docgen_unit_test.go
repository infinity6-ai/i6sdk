package cmdsdkweb_test

import (
	"os"
	"os/exec"
	"path/filepath"
	"testing"

	"github.com/stretchr/testify/require"
)

func TestUnitDocgen(t *testing.T) {
	type testScenario struct {
		dtsPath      string
		expectError  bool
		errContains  string
		expectedJSON string
	}

	docgenScript := filepath.Join("..", "docgen.mjs")

	check := func(t *testing.T, s testScenario) {
		t.Helper()

		outPath := filepath.Join(t.TempDir(), "sdk.json")
		cmd := exec.Command("node", docgenScript, "--dts", s.dtsPath, "--out", outPath)
		output, err := cmd.CombinedOutput()

		if s.expectError {
			require.Error(t, err)
			require.Contains(t, string(output), s.errContains)
			return
		}

		require.NoError(t, err, "unexpected error running docgen: %s", string(output))

		actualBytes, err := os.ReadFile(outPath)
		require.NoError(t, err)

		expectedBytes, err := os.ReadFile(s.expectedJSON)
		require.NoError(t, err)

		require.JSONEq(t, string(expectedBytes), string(actualBytes))
	}

	t.Run("Valid fixture generates expected JSON", func(t *testing.T) {
		check(t, testScenario{
			dtsPath:      filepath.Join("testdata", "docgen_fixture.d.ts"),
			expectError:  false,
			expectedJSON: filepath.Join("testdata", "docgen_expected.json"),
		})
	})

	t.Run("Missing documentation exits with error reporting missing items", func(t *testing.T) {
		check(t, testScenario{
			dtsPath:     filepath.Join("testdata", "docgen_missing_fixture.d.ts"),
			expectError: true,
			errContains: "missing doc",
		})
	})
}
