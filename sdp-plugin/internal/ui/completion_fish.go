package ui

import (
	"fmt"
	"os"
	"strings"
)

// generateFishCompletion generates fish completion script
func generateFishCompletion() (string, error) {
	return `# Fish completion for SDP
# Source this file in your config.fish:
#   source (sdp completion fish | psub)

function __sdp_checkpoint_commands
	echo -e "create\tCreate a new checkpoint"
	echo -e "resume\tResume from an existing checkpoint"
	echo -e "list\tList all checkpoints"
	echo -e "clean\tClean old checkpoints"
end

function __sdp_orchestrate_commands
	echo -e "start\tStart orchestration"
	echo -e "status\tCheck orchestration status"
	echo -e "stop\tStop orchestration"
end

function __sdp_beads_commands
	echo -e "ready\tList available tasks"
	echo -e "show\tShow task details"
	echo -e "update\tUpdate task status"
	echo -e "sync\tSynchronize Beads state"
end

function __sdp_quality_commands
	echo -e "check\tRun quality checks"
	echo -e "gate\tVerify quality gates pass"
	echo -e "scan\tScan for quality issues"
	echo -e "report\tGenerate quality report"
end

complete -c sdp -f

complete -c sdp -n "__fish_use_subcommand" -a init -d "Initialize project with SDP prompts"
complete -c sdp -n "__fish_use_subcommand" -a doctor -d "Check environment"
complete -c sdp -n "__fish_use_subcommand" -a hooks -d "Manage Git hooks"
complete -c sdp -n "__fish_use_subcommand" -a parse -d "Parse SDP workstream files"
complete -c sdp -n "__fish_use_subcommand" -a beads -d "Interact with Beads task tracker"
complete -c sdp -n "__fish_use_subcommand" -a tdd -d "Run TDD cycle"
complete -c sdp -n "__fish_use_subcommand" -a drift -d "Detect code drift"
complete -c sdp -n "__fish_use_subcommand" -a quality -d "Check code quality gates"
complete -c sdp -n "__fish_use_subcommand" -a watch -d "Watch files for quality violations"
complete -c sdp -n "__fish_use_subcommand" -a telemetry -d "Manage telemetry data"
complete -c sdp -n "__fish_use_subcommand" -a checkpoint -d "Manage checkpoints"
complete -c sdp -n "__fish_use_subcommand" -a orchestrate -d "Orchestrate multi-agent execution"

# Checkpoint subcommands
complete -c sdp -n "__fish_seen_subcommand_from checkpoint" -a "(__sdp_checkpoint_commands)"

# Orchestrate subcommands
complete -c sdp -n "__fish_seen_subcommand_from orchestrate" -a "(__sdp_orchestrate_commands)"

# Beads subcommands
complete -c sdp -n "__fish_seen_subcommand_from beads" -a "(__sdp_beads_commands)"

# Quality subcommands
complete -c sdp -n "__fish_seen_subcommand_from quality" -a "(__sdp_quality_commands)"
`, nil
}

// InstallCompletion installs completion for the specified shell
func InstallCompletion(shell CompletionType) error {
	var homeDir string
	var scriptPath string
	var script string

	homeDir = os.Getenv("HOME")
	if homeDir == "" {
		return fmt.Errorf("HOME environment variable not set")
	}

	switch shell {
	case Bash:
		scriptPath = homeDir + "/.bash_completion.d/sdp"
		var err error
		script, err = generateBashCompletion()
		if err != nil {
			return fmt.Errorf("failed to generate bash completion: %w", err)
		}
	case Zsh:
		// Try common zsh completion directories
		completionDir := homeDir + "/.zsh/completion"
		if _, err := os.Stat(completionDir); os.IsNotExist(err) {
			completionDir = homeDir + "/.zsh/completions"
		}
		scriptPath = completionDir + "/_sdp"
		var err error
		script, err = generateZshCompletion()
		if err != nil {
			return fmt.Errorf("failed to generate zsh completion: %w", err)
		}
	case Fish:
		completionDir := homeDir + "/.config/fish/completions"
		scriptPath = completionDir + "/sdp.fish"
		var err error
		script, err = generateFishCompletion()
		if err != nil {
			return fmt.Errorf("failed to generate fish completion: %w", err)
		}
	default:
		return fmt.Errorf("unsupported shell: %s", shell)
	}

	// Create directory if it doesn't exist
	dir := strings.TrimSuffix(scriptPath, "/"+scriptPath[strings.LastIndex(scriptPath, "/")+1:])
	if err := os.MkdirAll(dir, 0755); err != nil {
		return fmt.Errorf("failed to create completion directory: %w", err)
	}

	// Write completion script (0644: non-sensitive, sourced by shell)
	if err := os.WriteFile(scriptPath, []byte(script), 0644); err != nil {
		return fmt.Errorf("failed to write completion script: %w", err)
	}

	fmt.Printf("✅ Completion script installed to: %s\n", scriptPath)
	fmt.Printf("\nTo enable completion, add to your config:\n")

	switch shell {
	case Bash:
		fmt.Printf("  source ~/.bash_completion.d/sdp\n")
	case Zsh:
		fmt.Printf("  autoload -U compinit && compinit\n")
	case Fish:
		fmt.Printf("  Completion will be loaded automatically\n")
	}

	return nil
}
