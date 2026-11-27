package cmdversionz_test

import (
	"context"
	"os"
	"testing"

	"github.com/infinity6-ai/i6sdk/versionz/cmdversionz"
	"github.com/infinity6-ai/i6sdk/versionz/version"
	"github.com/stretchr/testify/require"
)

func TestUnitCmdVersionzPrepare(t *testing.T) {
	ctx := context.Background()

	type testScenario struct {
		name         string
		expectedName string
	}

	check := func(t *testing.T, s testScenario) {
		t.Helper()
		cmd := cmdversionz.Prepare(ctx)
		require.NotNil(t, cmd)
		require.Equal(t, s.expectedName, cmd.Use)
		require.Equal(t, version.Version(), cmd.Version)
		require.NotEmpty(t, cmd.Commands())

		hasVersion := false
		hasInc := false
		for _, c := range cmd.Commands() {
			if c.Name() == "version" {
				hasVersion = true
			}
			if c.Name() == "inc" {
				hasInc = true
			}
		}
		require.True(t, hasVersion)
		require.True(t, hasInc)
	}

	t.Run("Prepare returns root command with version and inc subcommands", func(t *testing.T) {
		check(t, testScenario{
			name:         "Prepare returns root command with version and inc subcommands",
			expectedName: "versionz",
		})
	})
}

func TestUnitIncrementVersion(t *testing.T) {
	type testScenario struct {
		initialContent string
		wantContent    string
		errMsg         string
	}

	check := func(t *testing.T, s testScenario) {
		t.Helper()
		tmpDir := t.TempDir()
		t.Chdir(tmpDir)

		if s.initialContent != "" {
			require.NoError(t, os.MkdirAll("version", 0755))
			require.NoError(t, os.WriteFile("version/version.go", []byte(s.initialContent), 0644))
		}

		err := cmdversionz.IncrementVersion()
		if s.errMsg != "" {
			require.Error(t, err)
			require.Contains(t, err.Error(), s.errMsg)
			return
		}
		require.NoError(t, err)

		content, err := os.ReadFile("version/version.go")
		require.NoError(t, err)
		require.Equal(t, s.wantContent, string(content))
	}

	t.Run("Valid version is incremented", func(t *testing.T) {
		check(t, testScenario{
			initialContent: "package version\n\nfunc Version() string {\n\treturn \"v1.2.3\"\n}\n",
			wantContent:    "package version\n\nfunc Version() string {\n\treturn \"v1.2.4\"\n}\n",
		})
	})

	t.Run("Missing file returns error", func(t *testing.T) {
		check(t, testScenario{
			initialContent: "",
			errMsg:         "could not read file",
		})
	})

	t.Run("Invalid content without version returns error", func(t *testing.T) {
		check(t, testScenario{
			initialContent: "package version\n\nfunc Version() string {\n\treturn \"invalid\"\n}\n",
			errMsg:         "version string not found",
		})
	})
}
