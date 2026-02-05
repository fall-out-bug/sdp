package main

import (
	"fmt"

	"github.com/ai-masters/sdp/internal/ui"
	"github.com/spf13/cobra"
)

var completionCmd = &cobra.Command{
	Use:   "completion [bash|zsh|fish]",
	Short: "Generate shell completion script",
	Long: `Generate shell completion script for SDP.

Supported shells:
  bash   Generate bash completion script
  zsh    Generate zsh completion script
  fish   Generate fish completion script

Examples:
  # Generate bash completion
  sdp completion bash > ~/.bash_completion.d/sdp
  source ~/.bash_completion.d/sdp

  # Generate zsh completion
  mkdir -p ~/.zsh/completion
  sdp completion zsh > ~/.zsh/completion/_sdp
  # Add to ~/.zshrc: fpath=(~/.zsh/completion $fpath)
  autoload -U compinit && compinit

  # Generate fish completion
  sdp completion fish > ~/.config/fish/completions/sdp.fish

  # Generate to stdout (for testing)
  sdp completion bash`,
	Args: cobra.ExactArgs(1),
	RunE: func(cmd *cobra.Command, args []string) error {
		shell := ui.CompletionType(args[0])

		err := ui.GenerateCompletion(shell)
		if err != nil {
			return fmt.Errorf("failed to generate completion: %w", err)
		}

		return nil
	},
	ValidArgsFunction: func(cmd *cobra.Command, args []string, toComplete string) ([]string, cobra.ShellCompDirective) {
		if len(args) != 0 {
			return nil, cobra.ShellCompDirectiveNoFileComp
		}
		return []string{"bash", "zsh", "fish"}, cobra.ShellCompDirectiveNoFileComp
	},
}
