package cmdset

import (
	"context"
	"fmt"
	"strings"

	"github.com/infinity6-ai/gox/commonz/logz"
	"github.com/infinity6-ai/i6sdk/i6cli/internal/cliconfig"
	"github.com/infinity6-ai/i6sdk/i6cli/internal/exitcode"
	"github.com/spf13/cobra"
)

type tlogger logz.Type

var logger = logz.Create(tlogger(true))

// New constructs the 'set' command under 'config'.
func New(ctx context.Context) *cobra.Command {
	cmd := &cobra.Command{
		Use:   "set <key> <value>",
		Short: "Update a profile configuration setting in config.json",
		Long: `Updates a configuration property in the active profile's config.json.

Supported keys:
  - box: Target box environment (prod | staging | sandbox)
  - default_dataset (or dataset): Default dataset identifier for this profile`,
		Args: cobra.ExactArgs(2),
		RunE: func(cmd *cobra.Command, args []string) error {
			c := cmd.Context()
			if c == nil {
				c = ctx
			}

			key := strings.ToLower(strings.TrimSpace(args[0]))
			val := strings.TrimSpace(args[1])

			profileFlag, _ := cmd.Flags().GetString("profile")
			profile := cliconfig.ResolveProfile(profileFlag, c)

			cfg, err := cliconfig.ResolveConfig(profile, c)
			if err != nil {
				return fmt.Errorf("failed to load configuration for profile %q: %w", profile, err)
			}

			switch key {
			case "box":
				boxVal := strings.ToLower(val)
				if _, err := cliconfig.ResolveBoxUrls(boxVal); err != nil {
					return err
				}
				cfg.Box = boxVal
			case "default_dataset", "dataset":
				cfg.DefaultDataset = val
			default:
				return exitcode.New(
					exitcode.Usage,
					fmt.Sprintf("unknown configuration key %q: expected box or default_dataset", key),
				)
			}

			if err := cliconfig.SaveConfig(profile, cfg, c); err != nil {
				return fmt.Errorf("failed to save configuration: %w", err)
			}

			logger.Info(c, "config-set", map[string]any{
				"profile": profile,
				"key":     key,
				"value":   val,
			})

			fmt.Fprintf(cmd.OutOrStdout(), "Updated %s = %q (profile: %s)\n", key, val, profile)
			return nil
		},
	}

	return cmd
}
