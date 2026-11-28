package cliconfig

import (
	"context"
	"encoding/json"
	"fmt"
	"os"
	"path/filepath"
	"strings"

	"github.com/infinity6-ai/gox/commonz/configz"
	"github.com/infinity6-ai/gox/commonz/urlz"
	"github.com/infinity6-ai/i6sdk/i6cli/internal/exitcode"
	"github.com/infinity6-ai/i6sdk/sdkgo/config"
	"github.com/infinity6-ai/i6sdk/sdkgo/tokenstore"
)

const (
	// DefaultProfile is the fallback profile name if none is explicitly specified.
	DefaultProfile = "main"

	// BoxProd selects the production service endpoints.
	BoxProd = "prod"

	// BoxStaging selects the staging service endpoints.
	BoxStaging = "staging"

	// BoxSandbox selects the sandbox service endpoints.
	BoxSandbox = "sandbox"
)

var (
	// PropConfigDir overrides the base configuration directory ($HOME/.config/i6).
	PropConfigDir = configz.Create("I6_CLI_CONFIG_DIR", "")

	// PropConfigProfile sets the active profile via environment property.
	PropConfigProfile = configz.Create("I6_CLI_CONFIG_PROFILE", DefaultProfile)

	// PropConfigJSON provides an inline JSON configuration for zero-disk CI/CD execution.
	PropConfigJSON = configz.Create("I6_CLI_CONFIG_JSON", "")

	// PropDataset overrides the target dataset identifier.
	PropDataset = configz.Create("I6_DATASET", "")
)

func getCtx(ctx ...context.Context) context.Context {
	if len(ctx) > 0 && ctx[0] != nil {
		return ctx[0]
	}
	return context.Background()
}

// Config represents the persisted profile configuration (config.json).
type Config struct {
	Box            string `json:"box"`
	DefaultDataset string `json:"default_dataset,omitempty"`
}

// ResolveProfile determines the active profile name using flag, configz property, and default precedence.
func ResolveProfile(flagVal string, ctx ...context.Context) string {
	if trimmed := strings.TrimSpace(flagVal); trimmed != "" {
		return trimmed
	}
	c := getCtx(ctx...)
	if propVal, err := PropConfigProfile.Get(c); err == nil {
		if trimmed := strings.TrimSpace(propVal); trimmed != "" {
			return trimmed
		}
	}
	return DefaultProfile
}

// ResolveBaseConfigDir returns the root CLI configuration directory ($HOME/.config/i6 or PropConfigDir).
func ResolveBaseConfigDir(ctx ...context.Context) (string, error) {
	c := getCtx(ctx...)
	if propDir, err := PropConfigDir.Get(c); err == nil {
		if trimmed := strings.TrimSpace(propDir); trimmed != "" {
			absDir, err := filepath.Abs(trimmed)
			if err != nil {
				return "", fmt.Errorf("failed to resolve absolute config dir from %s=%q: %w", PropConfigDir, trimmed, err)
			}
			return absDir, nil
		}
	}

	homeDir, err := os.UserHomeDir()
	if err != nil {
		homeDir = os.TempDir()
	}
	return filepath.Join(homeDir, ".config", "i6"), nil
}

// ResolveProfileDir returns the directory path for the given profile ($BASE_DIR/<profile>).
func ResolveProfileDir(profile string, ctx ...context.Context) (string, error) {
	baseDir, err := ResolveBaseConfigDir(ctx...)
	if err != nil {
		return "", err
	}
	p := strings.TrimSpace(profile)
	if p == "" {
		p = DefaultProfile
	}
	return filepath.Join(baseDir, p), nil
}

// EnsureProfileDir guarantees the profile directory and cache directory exist with 0700 permissions.
func EnsureProfileDir(profile string, ctx ...context.Context) (string, error) {
	pDir, err := ResolveProfileDir(profile, ctx...)
	if err != nil {
		return "", err
	}
	if err := os.MkdirAll(pDir, 0700); err != nil {
		return "", fmt.Errorf("failed to create profile directory %q: %w", pDir, err)
	}
	cacheDir := filepath.Join(pDir, "cache")
	if err := os.MkdirAll(cacheDir, 0700); err != nil {
		return "", fmt.Errorf("failed to create profile cache directory %q: %w", cacheDir, err)
	}
	return pDir, nil
}

// ResolveTokensSqliteUrl constructs the profile-isolated SQLite URL (file://$BASE_DIR/<profile>/tokens.db).
func ResolveTokensSqliteUrl(profile string, ctx ...context.Context) (*urlz.Url, error) {
	pDir, err := ResolveProfileDir(profile, ctx...)
	if err != nil {
		return nil, err
	}
	dbPath := filepath.Join(pDir, "tokens.db")
	u, err := urlz.Parse("file://" + filepath.ToSlash(dbPath))
	if err != nil {
		return nil, fmt.Errorf("failed to parse profile tokens sqlite url: %w", err)
	}
	return u, nil
}

// ResolveTokenStore returns an in-memory Store if PropConfigJSON is present, or a persistent SQLite Store.
func ResolveTokenStore(profile string, ctx ...context.Context) (tokenstore.Store, error) {
	c := getCtx(ctx...)
	if inline, err := PropConfigJSON.Get(c); err == nil && strings.TrimSpace(inline) != "" {
		return tokenstore.NewMemory(), nil
	}
	targetUrl, err := ResolveTokensSqliteUrl(profile, ctx...)
	if err != nil {
		return nil, err
	}
	store, err := tokenstore.NewSqlite(targetUrl)
	if err != nil {
		return nil, fmt.Errorf("failed to initialize sqlite token store for profile %q: %w", profile, err)
	}
	return store, nil
}

// ResolveCacheDir returns the profile-scoped local cache directory.
func ResolveCacheDir(profile string, ctx ...context.Context) (string, error) {
	pDir, err := ResolveProfileDir(profile, ctx...)
	if err != nil {
		return "", err
	}
	return filepath.Join(pDir, "cache"), nil
}

// ResolveConfig loads the profile's Config, checking inline JSON, file, and defaults.
func ResolveConfig(profile string, ctx ...context.Context) (*Config, error) {
	c := getCtx(ctx...)
	// 1. Check inline JSON override first (zero-disk execution)
	if inline, err := PropConfigJSON.Get(c); err == nil && strings.TrimSpace(inline) != "" {
		var cfg Config
		if err := json.Unmarshal([]byte(inline), &cfg); err != nil {
			return nil, exitcode.Wrap(exitcode.Usage, fmt.Errorf("failed to parse %s JSON: %w", PropConfigJSON, err))
		}
		if cfg.Box == "" {
			cfg.Box = BoxProd
		}
		return &cfg, nil
	}

	// 2. Read config.json from profile directory
	pDir, err := ResolveProfileDir(profile, ctx...)
	if err != nil {
		return nil, err
	}
	configPath := filepath.Join(pDir, "config.json")
	data, err := os.ReadFile(configPath)
	if err != nil {
		if os.IsNotExist(err) {
			return &Config{Box: BoxProd}, nil
		}
		return nil, fmt.Errorf("failed to read profile config %q: %w", configPath, err)
	}

	var cfg Config
	if err := json.Unmarshal(data, &cfg); err != nil {
		return nil, fmt.Errorf("failed to parse profile config %q: %w", configPath, err)
	}
	if cfg.Box == "" {
		cfg.Box = BoxProd
	}
	return &cfg, nil
}

// SaveConfig writes the Config to the profile's config.json with 0600 permissions.
func SaveConfig(profile string, cfg *Config, ctx ...context.Context) error {
	if cfg == nil {
		return fmt.Errorf("config cannot be nil")
	}
	pDir, err := EnsureProfileDir(profile, ctx...)
	if err != nil {
		return err
	}
	if cfg.Box == "" {
		cfg.Box = BoxProd
	}
	data, err := json.MarshalIndent(cfg, "", "  ")
	if err != nil {
		return fmt.Errorf("failed to marshal config: %w", err)
	}
	configPath := filepath.Join(pDir, "config.json")
	if err := os.WriteFile(configPath, data, 0600); err != nil {
		return fmt.Errorf("failed to write config to %q: %w", configPath, err)
	}
	return nil
}

// ResolveBoxUrls maps a box environment name to its corresponding service UrlOptions.
func ResolveBoxUrls(box string) (config.UrlOptions, error) {
	switch strings.ToLower(strings.TrimSpace(box)) {
	case "", BoxProd:
		return config.DefaultUrlOptions, nil
	case BoxStaging:
		return config.StagingUrlOptions, nil
	case BoxSandbox:
		return config.SandboxUrlOptions, nil
	default:
		return config.UrlOptions{}, exitcode.Wrap(exitcode.Usage, fmt.Errorf("invalid box environment %q: must be one of prod, staging, sandbox", box))
	}
}

// ResolveDataset resolves the active dataset using flag, config.json, and environment precedence.
// If unresolved, returns an ExitError with exit code 2 and actionable instructions.
func ResolveDataset(flagVal string, cfg *Config, ctx ...context.Context) (string, error) {
	if trimmed := strings.TrimSpace(flagVal); trimmed != "" {
		return trimmed, nil
	}
	if cfg != nil && strings.TrimSpace(cfg.DefaultDataset) != "" {
		return strings.TrimSpace(cfg.DefaultDataset), nil
	}
	c := getCtx(ctx...)
	if propDs, err := PropDataset.Get(c); err == nil && strings.TrimSpace(propDs) != "" {
		return strings.TrimSpace(propDs), nil
	}
	return "", exitcode.New(
		exitcode.Usage,
		"dataset required. Provide --dataset <name> or configure default via 'i6 config init --dataset <name>'",
	)
}
