package cmdshow

import (
	"context"
	"fmt"
	"os"
	"path/filepath"

	"github.com/infinity6-ai/gox/commonz/logz"
	"github.com/infinity6-ai/i6sdk/i6cli/internal/cliconfig"
	"github.com/infinity6-ai/i6sdk/i6cli/internal/outfmt"
	"github.com/spf13/cobra"
)

type tlogger logz.Type

var logger = logz.Create(tlogger(true))

// New constructs the 'show' command under 'config'.
func New(ctx context.Context) *cobra.Command {
	cmd := &cobra.Command{
		Use:   "show",
		Short: "Display active profile configuration, directories, and endpoints",
		Long: `Displays the active profile configuration details including directory paths,
box environment, service endpoint URLs, default dataset, and token store location.`,
		RunE: func(cmd *cobra.Command, args []string) error {
			c := cmd.Context()
			if c == nil {
				c = ctx
			}

			profileFlag, _ := cmd.Flags().GetString("profile")
			profile := cliconfig.ResolveProfile(profileFlag, c)

			pDir, err := cliconfig.ResolveProfileDir(profile, c)
			if err != nil {
				return fmt.Errorf("failed to resolve profile directory: %w", err)
			}

			cfg, err := cliconfig.ResolveConfig(profile, c)
			if err != nil {
				return fmt.Errorf("failed to load profile configuration: %w", err)
			}

			urls, err := cliconfig.ResolveBoxUrls(cfg.Box)
			if err != nil {
				return err
			}

			cacheDir, _ := cliconfig.ResolveCacheDir(profile, c)
			sqliteUrl, _ := cliconfig.ResolveTokensSqliteUrl(profile, c)
			configPath := filepath.Join(pDir, "config.json")

			configStatus := "not initialized (using defaults)"
			if _, statErr := os.Stat(configPath); statErr == nil {
				configStatus = "initialized"
			}

			defaultDs := cfg.DefaultDataset
			if defaultDs == "" {
				defaultDs = "(none)"
			}

			tokenStoreStr := "none"
			if sqliteUrl != nil {
				tokenStoreStr = sqliteUrl.String()
			}

			data := outfmt.Data{
				Headers: []string{"Property", "Value"},
				Rows: [][]string{
					{"Profile", profile},
					{"Status", configStatus},
					{"Profile Directory", pDir},
					{"Config File", configPath},
					{"Box Environment", cfg.Box},
					{"Default Dataset", defaultDs},
					{"Userz Service", urls.Userz},
					{"Decsuite Service", urls.Decsuite},
					{"Tfinfer Service", urls.Tfinfer},
					{"Token Store URL", tokenStoreStr},
					{"Cache Directory", cacheDir},
				},
				QuietId: []string{profile},
				Raw:     configPath,
				JSON: map[string]any{
					"profile":           profile,
					"status":            configStatus,
					"profile_directory": pDir,
					"config_file":       configPath,
					"box":               cfg.Box,
					"default_dataset":   cfg.DefaultDataset,
					"endpoints": map[string]string{
						"userz":    urls.Userz,
						"decsuite": urls.Decsuite,
						"tfinfer":  urls.Tfinfer,
					},
					"token_store_url": tokenStoreStr,
					"cache_directory": cacheDir,
				},
			}

			formatFlag, _ := cmd.Flags().GetString("format")
			format, err := outfmt.ParseFormat(formatFlag, outfmt.FormatTable)
			if err != nil {
				return err
			}

			logger.Debug(c, "config-show", map[string]any{
				"profile": profile,
				"format":  string(format),
			})

			return outfmt.Render(cmd.OutOrStdout(), format, data)
		},
	}

	return cmd
}
