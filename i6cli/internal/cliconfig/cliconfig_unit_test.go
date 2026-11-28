package cliconfig_test

import (
	"path/filepath"
	"testing"

	"github.com/infinity6-ai/i6sdk/i6cli/internal/cliconfig"
	"github.com/infinity6-ai/i6sdk/i6cli/internal/exitcode"
	"github.com/infinity6-ai/i6sdk/sdkgo/config"
	"github.com/stretchr/testify/require"
)

func TestUnitResolveProfile(t *testing.T) {
	ctx := t.Context()

	type testScenario struct {
		flagVal string
		propVal string
		want    string
	}

	check := func(t *testing.T, s testScenario) {
		t.Helper()
		rev, err := cliconfig.PropConfigProfile.Set(ctx, "%s", s.propVal)
		require.NoError(t, err)
		defer rev.Close()

		got := cliconfig.ResolveProfile(s.flagVal, ctx)
		require.Equal(t, s.want, got)
	}

	t.Run("Flag takes highest precedence", func(t *testing.T) {
		check(t, testScenario{
			flagVal: "custom-profile",
			propVal: "env-profile",
			want:    "custom-profile",
		})
	})

	t.Run("Config property used when flag is empty", func(t *testing.T) {
		check(t, testScenario{
			flagVal: "",
			propVal: "env-profile",
			want:    "env-profile",
		})
	})

	t.Run("Default main profile used when flag and prop are empty", func(t *testing.T) {
		check(t, testScenario{
			flagVal: "",
			propVal: "",
			want:    cliconfig.DefaultProfile,
		})
	})
}

func TestUnitResolveBaseConfigDir(t *testing.T) {
	ctx := t.Context()

	t.Run("Uses PropConfigDir when set", func(t *testing.T) {
		tmpDir := t.TempDir()
		rev, err := cliconfig.PropConfigDir.Set(ctx, "%s", tmpDir)
		require.NoError(t, err)
		defer rev.Close()

		baseDir, err := cliconfig.ResolveBaseConfigDir(ctx)
		require.NoError(t, err)
		require.Equal(t, tmpDir, baseDir)
	})

	t.Run("Falls back to home directory when prop is unset", func(t *testing.T) {
		rev, err := cliconfig.PropConfigDir.Set(ctx, "")
		require.NoError(t, err)
		defer rev.Close()

		baseDir, err := cliconfig.ResolveBaseConfigDir(ctx)
		require.NoError(t, err)
		require.Contains(t, baseDir, filepath.Join(".config", "i6"))
	})
}

func TestUnitEnsureProfileDirAndSqliteUrl(t *testing.T) {
	ctx := t.Context()
	tmpDir := t.TempDir()
	rev, err := cliconfig.PropConfigDir.Set(ctx, "%s", tmpDir)
	require.NoError(t, err)
	defer rev.Close()

	profile := "test-worker"
	pDir, err := cliconfig.EnsureProfileDir(profile, ctx)
	require.NoError(t, err)
	require.DirExists(t, pDir)
	require.DirExists(t, filepath.Join(pDir, "cache"))

	// Verify SQLite URL format
	u, err := cliconfig.ResolveTokensSqliteUrl(profile, ctx)
	require.NoError(t, err)
	require.NotNil(t, u)
	require.Equal(t, "file", u.Scheme)
	require.True(t, u.Path != nil && filepath.Base(u.Path.String()) == "tokens.db")

	// Verify token store can be constructed
	store, err := cliconfig.ResolveTokenStore(profile, ctx)
	require.NoError(t, err)
	require.NotNil(t, store)

	// Verify memory token store when inline JSON is provided
	revJSON, err := cliconfig.PropConfigJSON.Set(ctx, `{"box":"sandbox"}`)
	require.NoError(t, err)
	defer revJSON.Close()

	memStore, err := cliconfig.ResolveTokenStore(profile, ctx)
	require.NoError(t, err)
	require.NotNil(t, memStore)
}

func TestUnitConfigSaveAndResolve(t *testing.T) {
	ctx := t.Context()
	tmpDir := t.TempDir()
	rev, err := cliconfig.PropConfigDir.Set(ctx, "%s", tmpDir)
	require.NoError(t, err)
	defer rev.Close()

	profile := "staging-profile"

	// Initial resolution when config.json is absent defaults to prod
	cfg, err := cliconfig.ResolveConfig(profile, ctx)
	require.NoError(t, err)
	require.Equal(t, cliconfig.BoxProd, cfg.Box)
	require.Empty(t, cfg.DefaultDataset)

	// Save modified config
	cfg.Box = cliconfig.BoxStaging
	cfg.DefaultDataset = "fashion-ds"
	err = cliconfig.SaveConfig(profile, cfg, ctx)
	require.NoError(t, err)

	// Read back saved config
	readCfg, err := cliconfig.ResolveConfig(profile, ctx)
	require.NoError(t, err)
	require.Equal(t, cliconfig.BoxStaging, readCfg.Box)
	require.Equal(t, "fashion-ds", readCfg.DefaultDataset)

	// Inline JSON override
	revJSON, err := cliconfig.PropConfigJSON.Set(ctx, `{"box":"sandbox","default_dataset":"inline-ds"}`)
	require.NoError(t, err)
	defer revJSON.Close()

	inlineCfg, err := cliconfig.ResolveConfig(profile, ctx)
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
	ctx := t.Context()

	type testScenario struct {
		flagVal string
		cfg     *cliconfig.Config
		propVal string
		want    string
		wantErr bool
	}

	check := func(t *testing.T, s testScenario) {
		t.Helper()
		rev, err := cliconfig.PropDataset.Set(ctx, "%s", s.propVal)
		require.NoError(t, err)
		defer rev.Close()

		got, err := cliconfig.ResolveDataset(s.flagVal, s.cfg, ctx)
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
			propVal: "env-ds",
			want:    "flag-ds",
		})
	})

	t.Run("Config default_dataset used when flag is empty", func(t *testing.T) {
		check(t, testScenario{
			flagVal: "",
			cfg:     &cliconfig.Config{DefaultDataset: "cfg-ds"},
			propVal: "env-ds",
			want:    "cfg-ds",
		})
	})

	t.Run("Config property used when flag and config are empty", func(t *testing.T) {
		check(t, testScenario{
			flagVal: "",
			cfg:     &cliconfig.Config{},
			propVal: "env-ds",
			want:    "env-ds",
		})
	})

	t.Run("Returns error when dataset cannot be resolved", func(t *testing.T) {
		check(t, testScenario{
			flagVal: "",
			cfg:     &cliconfig.Config{},
			propVal: "",
			wantErr: true,
		})
	})
}
