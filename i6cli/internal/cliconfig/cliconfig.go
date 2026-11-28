package cliconfig

import (
	"encoding/json"
	"fmt"
	"os"
	"path/filepath"
	"strings"

	"github.com/infinity6-ai/gox/commonz/urlz"
	"github.com/infinity6-ai/i6sdk/i6cli/internal/exitcode"
	"github.com/infinity6-ai/i6sdk/sdkgo/config"
	"github.com/infinity6-ai/i6sdk/sdkgo/tokenstore"
)

const (
	// DefaultProfile is the fallback profile name if none is explicitly specified.
	DefaultProfile = "main"

	// EnvConfigDir overrides the base configuration directory ($HOME/.config/i6).
	EnvConfigDir = "I6_CLI_CONFIG_DIR"

	// EnvConfigProfile sets the active profile via environment variable.
	EnvConfigProfile = "I6_CLI_CONFIG_PROFILE"

	// EnvConfigJSON provides an inline JSON configuration for zero-disk CI/CD execution.
	EnvConfigJSON = "I6_CLI_CONFIG_JSON"

	// EnvDataset overrides the target dataset identifier.
	EnvDataset = "I6_DATASET"

	// BoxProd selects the production service endpoints.
	BoxProd = "prod"

	// BoxStaging selects the staging service endpoints.
	BoxStaging = "staging"

	// BoxSandbox selects the sandbox service endpoints.
	BoxSandbox = "sandbox"
)

// Config represents the persisted profile configuration (config.json).
type Config struct {
	Box            string `json:"box"`
	DefaultDataset string `json:"default_dataset,omitempty"`
}

// ResolveProfile determines the active profile name using flag, env, and default precedence.
func ResolveProfile(flagVal string) string {
	if trimmed := strings.TrimSpace(flagVal); trimmed != "" {
		return trimmed
	}
	if envVal := strings.TrimSpace(os.Getenv(EnvConfigProfile)); envVal != "" {
		return envVal
	}
	return DefaultProfile
}

// ResolveBaseConfigDir returns the root CLI configuration directory ($HOME/.config/i6 or I6_CLI_CONFIG_DIR).
func ResolveBaseConfigDir() (string, error) {
	if envDir := strings.TrimSpace(os.Getenv(EnvConfigDir)); envDir != "" {
		absDir, err := filepath.Abs(envDir)
		if err != nil {
			return "", fmt.Errorf("failed to resolve absolute config dir from %s=%q: %w", EnvConfigDir, envDir, err)
		}
		return absDir, nil
	}

	homeDir, err := os.UserHomeDir()
	if err != nil {
		homeDir = os.TempDir()
	}
	return filepath.Join(homeDir, ".config", "i6"), nil
}

// ResolveProfileDir returns the directory path for the given profile ($BASE_DIR/<profile>).
func ResolveProfileDir(profile string) (string, error) {
	baseDir, err := ResolveBaseConfigDir()
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
func EnsureProfileDir(profile string) (string, error) {
	pDir, err := ResolveProfileDir(profile)
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
func ResolveTokensSqliteUrl(profile string) (*urlz.Url, error) {
	pDir, err := ResolveProfileDir(profile)
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

// ResolveTokenStore returns an in-memory Store if I6_CLI_CONFIG_JSON is present, or a persistent SQLite Store.
func ResolveTokenStore(profile string) (tokenstore.Store, error) {
	if strings.TrimSpace(os.Getenv(EnvConfigJSON)) != "" {
		return tokenstore.NewMemory(), nil
	}
	targetUrl, err := ResolveTokensSqliteUrl(profile)
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
func ResolveCacheDir(profile string) (string, error) {
	pDir, err := ResolveProfileDir(profile)
	if err != nil {
		return "", err
	}
	return filepath.Join(pDir, "cache"), nil
}

// ResolveConfig loads the profile's Config, checking inline JSON, file, and defaults.
func ResolveConfig(profile string) (*Config, error) {
	// 1. Check inline JSON override first (zero-disk execution)
	if inline := strings.TrimSpace(os.Getenv(EnvConfigJSON)); inline != "" {
		var cfg Config
		if err := json.Unmarshal([]byte(inline), &cfg); err != nil {
			return nil, exitcode.Wrap(exitcode.Usage, fmt.Errorf("failed to parse %s JSON: %w", EnvConfigJSON, err))
		}
		if cfg.Box == "" {
			cfg.Box = BoxProd
		}
		return &cfg, nil
	}

	// 2. Read config.json from profile directory
	pDir, err := ResolveProfileDir(profile)
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
func SaveConfig(profile string, cfg *Config) error {
	if cfg == nil {
		return fmt.Errorf("config cannot be nil")
	}
	pDir, err := EnsureProfileDir(profile)
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
func ResolveDataset(flagVal string, cfg *Config) (string, error) {
	if trimmed := strings.TrimSpace(flagVal); trimmed != "" {
		return trimmed, nil
	}
	if cfg != nil && strings.TrimSpace(cfg.DefaultDataset) != "" {
		return strings.TrimSpace(cfg.DefaultDataset), nil
	}
	if envDs := strings.TrimSpace(os.Getenv(EnvDataset)); envDs != "" {
		return envDs, nil
	}
	return "", exitcode.New(
		exitcode.Usage,
		"dataset required. Provide --dataset <name> or configure default via 'i6 config init --dataset <name>'",
	)
}
