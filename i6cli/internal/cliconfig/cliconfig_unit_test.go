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
	revJSON, err := cliconfig.PropConfigJSON.Set(ctx, `{"default_dataset":"mem-ds"}`)
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

	// Initial resolution when config.json is absent returns empty config
	cfg, err := cliconfig.ResolveConfig(profile, ctx)
	require.NoError(t, err)
	require.Empty(t, cfg.DefaultDataset)

	// Save modified config
	cfg.DefaultDataset = "fashion-ds"
	cfg.UserzUrl = "https://custom-userz.internal"
	err = cliconfig.SaveConfig(profile, cfg, ctx)
	require.NoError(t, err)

	// Read back saved config
	readCfg, err := cliconfig.ResolveConfig(profile, ctx)
	require.NoError(t, err)
	require.Equal(t, "fashion-ds", readCfg.DefaultDataset)
	require.Equal(t, "https://custom-userz.internal", readCfg.UserzUrl)

	// Inline JSON override
	revJSON, err := cliconfig.PropConfigJSON.Set(ctx, `{"default_dataset":"inline-ds"}`)
	require.NoError(t, err)
	defer revJSON.Close()

	inlineCfg, err := cliconfig.ResolveConfig(profile, ctx)
	require.NoError(t, err)
	require.Equal(t, "inline-ds", inlineCfg.DefaultDataset)
}

func TestUnitResolveUrls(t *testing.T) {
	type testScenario struct {
		cfg      *cliconfig.Config
		wantUser string
		wantDec  string
		wantTf   string
	}

	check := func(t *testing.T, s testScenario) {
		t.Helper()
		urls := cliconfig.ResolveUrls(s.cfg)
		require.Equal(t, s.wantUser, urls.Userz)
		require.Equal(t, s.wantDec, urls.Decsuite)
		require.Equal(t, s.wantTf, urls.Tfinfer)
	}

	t.Run("Nil config resolves to production defaults", func(t *testing.T) {
		check(t, testScenario{
			cfg:      nil,
			wantUser: config.DefaultUrlOptions.Userz,
			wantDec:  config.DefaultUrlOptions.Decsuite,
			wantTf:   config.DefaultUrlOptions.Tfinfer,
		})
	})

	t.Run("Empty config resolves to production defaults", func(t *testing.T) {
		check(t, testScenario{
			cfg:      &cliconfig.Config{},
			wantUser: config.DefaultUrlOptions.Userz,
			wantDec:  config.DefaultUrlOptions.Decsuite,
			wantTf:   config.DefaultUrlOptions.Tfinfer,
		})
	})

	t.Run("Custom URLs override defaults", func(t *testing.T) {
		check(t, testScenario{
			cfg: &cliconfig.Config{
				UserzUrl:    "https://custom-user.example.com",
				DecsuiteUrl: "https://custom-dec.example.com",
			},
			wantUser: "https://custom-user.example.com",
			wantDec:  "https://custom-dec.example.com",
			wantTf:   config.DefaultUrlOptions.Tfinfer,
		})
	})
}

func TestUnitListProfiles(t *testing.T) {
	ctx := cliconfig.WithConfigDir(t.Context(), t.TempDir())

	// Initially, returns empty list when no profiles exist on disk
	profiles, err := cliconfig.ListProfiles(ctx)
	require.NoError(t, err)
	require.Empty(t, profiles)

	// Create a couple profiles on disk
	require.NoError(t, cliconfig.SaveConfig("staging", &cliconfig.Config{DefaultDataset: "ds-stg"}, ctx))
	require.NoError(t, cliconfig.SaveConfig("prod", &cliconfig.Config{DefaultDataset: "ds-prd"}, ctx))

	profiles, err = cliconfig.ListProfiles(ctx)
	require.NoError(t, err)
	require.Len(t, profiles, 2) // prod, staging
	require.Equal(t, "prod", profiles[0].Name)
	require.Equal(t, "ds-prd", profiles[0].DefaultDataset)
	require.True(t, profiles[0].HasConfig)
	require.Equal(t, "staging", profiles[1].Name)
	require.Equal(t, "ds-stg", profiles[1].DefaultDataset)
	require.True(t, profiles[1].HasConfig)
}

func TestUnitNewClient(t *testing.T) {
	ctx := t.Context()
	tmpDir := t.TempDir()
	rev, err := cliconfig.PropConfigDir.Set(ctx, "%s", tmpDir)
	require.NoError(t, err)
	defer rev.Close()

	cli, err := cliconfig.NewClient("test-cli", ctx)
	require.NoError(t, err)
	require.NotNil(t, cli)
	require.NotNil(t, cli.Tokens())
	require.NotNil(t, cli.Apis())
}

func TestUnitWithConfigDirParallel(t *testing.T) {
	t.Run("Parallel subtests run in isolation with WithConfigDir", func(t *testing.T) {
		for i := 0; i < 5; i++ {
			t.Run(string(rune('A'+i)), func(t *testing.T) {
				t.Parallel()
				tmpDir := t.TempDir()
				ctx := cliconfig.WithConfigDir(t.Context(), tmpDir)

				baseDir, err := cliconfig.ResolveBaseConfigDir(ctx)
				require.NoError(t, err)
				require.Equal(t, tmpDir, baseDir)

				// Save unique profile in isolated directory
				profileName := "worker"
				cfg := &cliconfig.Config{DefaultDataset: "ds-" + t.Name()}
				require.NoError(t, cliconfig.SaveConfig(profileName, cfg, ctx))

				readCfg, err := cliconfig.ResolveConfig(profileName, ctx)
				require.NoError(t, err)
				require.Equal(t, "ds-"+t.Name(), readCfg.DefaultDataset)
			})
		}
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
