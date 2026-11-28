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
  - box: Target box environment preset (prod | staging | sandbox) to update all service URLs
  - default_dataset (or dataset): Default dataset identifier for this profile
  - userz_url (or userz): Custom userz authentication service endpoint URL
  - decsuite_url (or decsuite): Custom decsuite dataset/lake service endpoint URL
  - tfinfer_url (or tfinfer): Custom tfinfer prediction service endpoint URL`,
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
				boxUrls, err := cliconfig.BoxPresetUrls(val)
				if err != nil {
					return err
				}
				cfg.UserzUrl = boxUrls.Userz
				cfg.DecsuiteUrl = boxUrls.Decsuite
				cfg.TfinferUrl = boxUrls.Tfinfer
			case "default_dataset", "dataset":
				cfg.DefaultDataset = val
			case "userz_url", "userz":
				cfg.UserzUrl = val
			case "decsuite_url", "decsuite":
				cfg.DecsuiteUrl = val
			case "tfinfer_url", "tfinfer":
				cfg.TfinferUrl = val
			default:
				return exitcode.New(
					exitcode.Usage,
					fmt.Sprintf("unknown configuration key %q: expected box, default_dataset, userz_url, decsuite_url, or tfinfer_url", key),
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
