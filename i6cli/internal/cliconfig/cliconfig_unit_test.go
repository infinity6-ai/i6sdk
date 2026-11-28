package cliconfig_test

import (
	"os"
	"path/filepath"
	"testing"

	"github.com/infinity6-ai/i6sdk/i6cli/internal/cliconfig"
	"github.com/infinity6-ai/i6sdk/i6cli/internal/exitcode"
	"github.com/infinity6-ai/i6sdk/sdkgo/config"
	"github.com/stretchr/testify/require"
)

func TestUnitResolveProfile(t *testing.T) {
	type testScenario struct {
		flagVal string
		envVal  string
		want    string
	}

	check := func(t *testing.T, s testScenario) {
		t.Helper()
		if s.envVal != "" {
			t.Setenv(cliconfig.EnvConfigProfile, s.envVal)
		} else {
			_ = os.Unsetenv(cliconfig.EnvConfigProfile)
		}
		got := cliconfig.ResolveProfile(s.flagVal)
		require.Equal(t, s.want, got)
	}

	t.Run("Flag takes highest precedence", func(t *testing.T) {
		check(t, testScenario{
			flagVal: "custom-profile",
			envVal:  "env-profile",
			want:    "custom-profile",
		})
	})

	t.Run("Environment variable used when flag is empty", func(t *testing.T) {
		check(t, testScenario{
			flagVal: "",
			envVal:  "env-profile",
			want:    "env-profile",
		})
	})

	t.Run("Default main profile used when flag and env are empty", func(t *testing.T) {
		check(t, testScenario{
			flagVal: "",
			envVal:  "",
			want:    cliconfig.DefaultProfile,
		})
	})
}

func TestUnitResolveBaseConfigDir(t *testing.T) {
	t.Run("Uses I6_CLI_CONFIG_DIR when set", func(t *testing.T) {
		tmpDir := t.TempDir()
		t.Setenv(cliconfig.EnvConfigDir, tmpDir)

		baseDir, err := cliconfig.ResolveBaseConfigDir()
		require.NoError(t, err)
		require.Equal(t, tmpDir, baseDir)
	})

	t.Run("Falls back to home directory when env is unset", func(t *testing.T) {
		_ = os.Unsetenv(cliconfig.EnvConfigDir)
		baseDir, err := cliconfig.ResolveBaseConfigDir()
		require.NoError(t, err)
		require.Contains(t, baseDir, filepath.Join(".config", "i6"))
	})
}

func TestUnitEnsureProfileDirAndSqliteUrl(t *testing.T) {
	tmpDir := t.TempDir()
	t.Setenv(cliconfig.EnvConfigDir, tmpDir)

	profile := "test-worker"
	pDir, err := cliconfig.EnsureProfileDir(profile)
	require.NoError(t, err)
	require.DirExists(t, pDir)
	require.DirExists(t, filepath.Join(pDir, "cache"))

	// Verify SQLite URL format
	u, err := cliconfig.ResolveTokensSqliteUrl(profile)
	require.NoError(t, err)
	require.NotNil(t, u)
	require.Equal(t, "file", u.Scheme)
	require.True(t, u.Path != nil && filepath.Base(u.Path.String()) == "tokens.db")

	// Verify token store can be constructed
	store, err := cliconfig.ResolveTokenStore(profile)
	require.NoError(t, err)
	require.NotNil(t, store)

	// Verify memory token store when inline JSON is provided
	t.Setenv(cliconfig.EnvConfigJSON, `{"box":"sandbox"}`)
	memStore, err := cliconfig.ResolveTokenStore(profile)
	require.NoError(t, err)
	require.NotNil(t, memStore)
}

func TestUnitConfigSaveAndResolve(t *testing.T) {
	tmpDir := t.TempDir()
	t.Setenv(cliconfig.EnvConfigDir, tmpDir)
	profile := "staging-profile"

	// Initial resolution when config.json is absent defaults to prod
	cfg, err := cliconfig.ResolveConfig(profile)
	require.NoError(t, err)
	require.Equal(t, cliconfig.BoxProd, cfg.Box)
	require.Empty(t, cfg.DefaultDataset)

	// Save modified config
	cfg.Box = cliconfig.BoxStaging
	cfg.DefaultDataset = "fashion-ds"
	err = cliconfig.SaveConfig(profile, cfg)
	require.NoError(t, err)

	// Read back saved config
	readCfg, err := cliconfig.ResolveConfig(profile)
	require.NoError(t, err)
	require.Equal(t, cliconfig.BoxStaging, readCfg.Box)
	require.Equal(t, "fashion-ds", readCfg.DefaultDataset)

	// Inline JSON override
	t.Setenv(cliconfig.EnvConfigJSON, `{"box":"sandbox","default_dataset":"inline-ds"}`)
	inlineCfg, err := cliconfig.ResolveConfig(profile)
	require.NoError(t, err)
	require.Equal(t, cliconfig.BoxSandbox, inlineCfg.Box)
	require.Equal(t, "inline-ds", inlineCfg.DefaultDataset)
}

func TestUnitResolveBoxUrls(t *testing.T) {
	type testScenario struct {
		box      string
		wantUser string
		wantErr  bool
	}

	check := func(t *testing.T, s testScenario) {
		t.Helper()
		urls, err := cliconfig.ResolveBoxUrls(s.box)
		if s.wantErr {
			require.Error(t, err)
			require.Equal(t, exitcode.Usage, exitcode.Code(err))
			return
		}
		require.NoError(t, err)
		require.Equal(t, s.wantUser, urls.Userz)
	}

	t.Run("Prod / default", func(t *testing.T) {
		check(t, testScenario{
			box:      "prod",
			wantUser: config.DefaultUrlOptions.Userz,
		})
	})

	t.Run("Staging", func(t *testing.T) {
		check(t, testScenario{
			box:      "staging",
			wantUser: config.StagingUrlOptions.Userz,
		})
	})

	t.Run("Sandbox", func(t *testing.T) {
		check(t, testScenario{
			box:      "sandbox",
			wantUser: config.SandboxUrlOptions.Userz,
		})
	})

	t.Run("Invalid box returns usage error", func(t *testing.T) {
		check(t, testScenario{
			box:     "unknown-env",
			wantErr: true,
		})
	})
}

func TestUnitResolveDataset(t *testing.T) {
	type testScenario struct {
		flagVal string
		cfg     *cliconfig.Config
		envVal  string
		want    string
		wantErr bool
	}

	check := func(t *testing.T, s testScenario) {
		t.Helper()
		if s.envVal != "" {
			t.Setenv(cliconfig.EnvDataset, s.envVal)
		} else {
			_ = os.Unsetenv(cliconfig.EnvDataset)
		}
		got, err := cliconfig.ResolveDataset(s.flagVal, s.cfg)
		if s.wantErr {
			require.Error(t, err)
			require.Equal(t, exitcode.Usage, exitcode.Code(err))
			require.Contains(t, err.Error(), "dataset required")
			return
		}
		require.NoError(t, err)
		require.Equal(t, s.want, got)
	}

	t.Run("Flag takes highest precedence", func(t *testing.T) {
		check(t, testScenario{
			flagVal: "flag-ds",
			cfg:     &cliconfig.Config{DefaultDataset: "cfg-ds"},
			envVal:  "env-ds",
			want:    "flag-ds",
		})
	})

	t.Run("Config default_dataset used when flag is empty", func(t *testing.T) {
		check(t, testScenario{
			flagVal: "",
			cfg:     &cliconfig.Config{DefaultDataset: "cfg-ds"},
			envVal:  "env-ds",
			want:    "cfg-ds",
		})
	})

	t.Run("Environment variable used when flag and config are empty", func(t *testing.T) {
		check(t, testScenario{
			flagVal: "",
			cfg:     &cliconfig.Config{},
			envVal:  "env-ds",
			want:    "env-ds",
		})
	})

	t.Run("Returns error when dataset cannot be resolved", func(t *testing.T) {
		check(t, testScenario{
			flagVal: "",
			cfg:     &cliconfig.Config{},
			envVal:  "",
			wantErr: true,
		})
	})
}
