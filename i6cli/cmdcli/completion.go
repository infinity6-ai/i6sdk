package cmdcli

import (
	"fmt"

	"github.com/infinity6-ai/i6sdk/i6cli/internal/exitcode"
	"github.com/spf13/cobra"
)

// newCompletionCmd creates the completion subcommand for generating shell completion scripts.
func newCompletionCmd() *cobra.Command {
	var noDesc bool

	cmd := &cobra.Command{
		Use:   "completion [bash|zsh|fish|powershell]",
		Short: "Generate shell completion scripts",
		Long: `Generate shell completion scripts for i6 CLI.

To load completions:

Bash:
  $ source <(i6 completion bash)
  # To load completions for each session, execute once:
  # Linux:
  $ i6 completion bash > /etc/bash_completion.d/i6
  # macOS:
  $ i6 completion bash > $(brew --prefix)/etc/bash_completion.d/i6

Zsh:
  # If shell completion is not already enabled in your environment,
  # you will need to enable it. You can execute the following once:
  $ echo "autoload -U compinit; compinit" >> ~/.zshrc

  # To load completions for each session, execute once:
  $ i6 completion zsh > "${fpath[1]}/_i6"

Fish:
  $ i6 completion fish | source
  # To load completions for each session, execute once:
  $ i6 completion fish > ~/.config/fish/completions/i6.fish

PowerShell:
  PS> i6 completion powershell | Out-String | Invoke-Expression
`,
		DisableFlagsInUseLine: true,
		ValidArgs:             []string{"bash", "zsh", "fish", "powershell"},
		Args:                  cobra.ExactArgs(1),
		RunE: func(cmd *cobra.Command, args []string) error {
			shell := args[0]
			out := cmd.OutOrStdout()

			switch shell {
			case "bash":
				if err := cmd.Root().GenBashCompletionV2(out, !noDesc); err != nil {
					return fmt.Errorf("failed to generate bash completion: %w", err)
				}
			case "zsh":
				if noDesc {
					if err := cmd.Root().GenZshCompletionNoDesc(out); err != nil {
						return fmt.Errorf("failed to generate zsh completion: %w", err)
					}
				} else {
					if err := cmd.Root().GenZshCompletion(out); err != nil {
						return fmt.Errorf("failed to generate zsh completion: %w", err)
					}
				}
			case "fish":
				if err := cmd.Root().GenFishCompletion(out, !noDesc); err != nil {
					return fmt.Errorf("failed to generate fish completion: %w", err)
				}
			case "powershell":
				if noDesc {
					if err := cmd.Root().GenPowerShellCompletion(out); err != nil {
						return fmt.Errorf("failed to generate powershell completion: %w", err)
					}
				} else {
					if err := cmd.Root().GenPowerShellCompletionWithDesc(out); err != nil {
						return fmt.Errorf("failed to generate powershell completion: %w", err)
					}
				}
			default:
				return exitcode.New(
					exitcode.Usage,
					fmt.Sprintf("unsupported shell %q: expected one of bash, zsh, fish, powershell", shell),
				)
			}
			return nil
		},
	}

	cmd.Flags().BoolVar(&noDesc, "no-descriptions", false, "Disable completion descriptions")
	return cmd
}
