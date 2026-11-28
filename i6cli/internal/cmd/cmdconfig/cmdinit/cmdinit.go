package cmdinit

import (
	"context"
	"fmt"
	"os"
	"path/filepath"
	"strings"

	"github.com/infinity6-ai/gox/commonz/logz"
	"github.com/infinity6-ai/i6sdk/i6cli/internal/cliconfig"
	"github.com/infinity6-ai/i6sdk/i6cli/internal/exitcode"
	"github.com/spf13/cobra"
)

type tlogger logz.Type

var logger = logz.Create(tlogger(true))

// New constructs the 'init' command under 'config'.
func New(ctx context.Context) *cobra.Command {
	var (
		dataset string
		force   bool
	)

	cmd := &cobra.Command{
		Use:   "init",
		Short: "Initialize profile configuration directory and config.json",
		Long: `Initializes the active profile directory with a config.json configuration file.

Allows configuring an optional default dataset. If the configuration already exists,
the command fails unless -f/--force is provided to overwrite it.`,
		RunE: func(cmd *cobra.Command, args []string) error {
			c := cmd.Context()
			if c == nil {
				c = ctx
			}

			profileFlag, _ := cmd.Flags().GetString("profile")
			profile := cliconfig.ResolveProfile(profileFlag, c)

			pDir, err := cliconfig.EnsureProfileDir(profile, c)
			if err != nil {
				return fmt.Errorf("failed to prepare profile directory: %w", err)
			}

			configPath := filepath.Join(pDir, "config.json")
			if _, err := os.Stat(configPath); err == nil && !force {
				return exitcode.New(
					exitcode.Usage,
					fmt.Sprintf("configuration already exists at %q: use -f/--force to overwrite", configPath),
				)
			}

			cfg := &cliconfig.Config{
				DefaultDataset: strings.TrimSpace(dataset),
			}

			if err := cliconfig.SaveConfig(profile, cfg, c); err != nil {
				return fmt.Errorf("failed to save configuration: %w", err)
			}

			logger.Info(c, "config-init", map[string]any{
				"profile": profile,
				"path":    configPath,
			})

			fmt.Fprintf(cmd.OutOrStdout(), "Initialized profile %q at %s\n", profile, configPath)
			return nil
		},
	}

	cmd.Flags().StringVar(&dataset, "dataset", "", "Default dataset identifier for this profile")
	cmd.Flags().BoolVarP(&force, "force", "f", false, "Overwrite existing configuration file if present")

	return cmd
}
