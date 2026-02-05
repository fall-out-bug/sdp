package ui

import (
	"fmt"
	"os"
	"strings"
)

// CompletionType represents the type of shell completion
type CompletionType string

const (
	Bash CompletionType = "bash"
	Zsh  CompletionType = "zsh"
	Fish CompletionType = "fish"
)

// GenerateCompletion generates shell completion script
func GenerateCompletion(shell CompletionType) error {
	var script string
	var err error

	switch shell {
	case Bash:
		script, err = generateBashCompletion()
	case Zsh:
		script, err = generateZshCompletion()
	case Fish:
		script, err = generateFishCompletion()
	default:
		return fmt.Errorf("unsupported shell: %s (supported: bash, zsh, fish)", shell)
	}

	if err != nil {
		return err
	}

	fmt.Println(script)
	return nil
}

// generateBashCompletion generates bash completion script
func generateBashCompletion() (string, error) {
	return `# Bash completion for SDP
# Source this file in your .bashrc or .bash_profile:
#   source <(sdp completion bash)

_sdp_completion() {
	local cur prev words cword
	_init_completion || return

	local commands="init doctor hooks parse beads tdd drift quality watch telemetry checkpoint orchestrate"
	local checkpoint_commands="create resume list clean"
	local orchestrate_commands="start status stop"

	case ${prev} in
		checkpoint)
			COMPREPLY=($(compgen -W "${checkpoint_commands}" -- "${cur}"))
			return
			;;
		orchestrate)
			COMPREPLY=($(compgen -W "${orchestrate_commands}" -- "${cur}"))
			return
			;;
		beads)
			COMPREPLY=($(compgen -W "ready show update sync" -- "${cur}"))
			return
			;;
		quality)
			COMPREPLY=($(compgen -W "check gate scan report" -- "${cur}"))
			return
			;;
		checkpoint)
			case ${words[2]} in
				create|resume)
					COMPREPLY=($(compgen -f -- "${cur}"))
					return
					;;
			esac
			;;
		*)
			;;
	esac

	if [[ ${cword} -eq 1 ]]; then
		COMPREPLY=($(compgen -W "${commands}" -- "${cur}"))
	fi

	# File completion for certain commands
	if [[ ${cword} -ge 2 ]]; then
		case ${words[1]} in
			init|parse|drift|watch)
				COMPREPLY=($(compgen -f -- "${cur}"))
				;;
		esac
	fi
}

complete -F _sdp_completion sdp
`, nil
}

// generateZshCompletion generates zsh completion script
func generateZshCompletion() (string, error) {
	return `#compdef sdp

# Zsh completion for SDP
# Place this file in ~/.zsh/completion/ or source it in your .zshrc:
#   source <(sdp completion zsh)

_sdp() {
	local -a commands
	local -a checkpoint_commands
	local -a orchestrate_commands
	local -a beads_commands
	local -a quality_commands

	commands=(
		'init:Initialize project with SDP prompts'
		'doctor:Check environment (Git, Claude Code, .claude/)'
		'hooks:Manage Git hooks for SDP'
		'parse:Parse SDP workstream files'
		'beads:Interact with Beads task tracker'
		'tdd:Run TDD cycle (Red-Green-Refactor)'
		'drift:Detect code drift from specification'
		'quality:Check code quality gates'
		'watch:Watch files for quality violations'
		'telemetry:Manage telemetry data'
		'checkpoint:Manage checkpoints for long-running features'
		'orchestrate:Orchestrate multi-agent execution'
	)

	checkpoint_commands=(
		'create:Create a new checkpoint'
		'resume:Resume from an existing checkpoint'
		'list:List all checkpoints'
		'clean:Clean old checkpoints'
	)

	orchestrate_commands=(
		'start:Start orchestration'
		'status:Check orchestration status'
		'stop:Stop orchestration'
	)

	beads_commands=(
		'ready:List available tasks'
		'show:Show task details'
		'update:Update task status'
		'sync:Synchronize Beads state'
	)

	quality_commands=(
		'check:Run quality checks'
		'gate:Verify quality gates pass'
		'scan:Scan for quality issues'
		'report:Generate quality report'
	)

	case $state in
		command)
			_describe 'command' commands
			;;
		checkpoint)
			_describe 'checkpoint command' checkpoint_commands
			;;
		orchestrate)
			_describe 'orchestrate command' orchestrate_commands
			;;
		beads)
			_describe 'beads command' beads_commands
			;;
		quality)
			_describe 'quality command' quality_commands
			;;
	esac
}

_arguments \\
	'1: :->command' \\
	'*::arg:->args'

case $line[1] in
	checkpoint)
		_sdp checkpoint "$words[2,-1]"
		;;
	orchestrate)
		_sdp orchestrate "$words[2,-1]"
		;;
	beads)
		_sdp beads "$words[2,-1]"
		;;
	quality)
		_sdp quality "$words[2,-1]"
		;;
	*)
		_sdp
		;;
esac
`, nil
}

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
		script, _ = generateBashCompletion()
	case Zsh:
		// Try common zsh completion directories
		completionDir := homeDir + "/.zsh/completion"
		if _, err := os.Stat(completionDir); os.IsNotExist(err) {
			completionDir = homeDir + "/.zsh/completions"
		}
		scriptPath = completionDir + "/_sdp"
		script, _ = generateZshCompletion()
	case Fish:
		completionDir := homeDir + "/.config/fish/completions"
		scriptPath = completionDir + "/sdp.fish"
		script, _ = generateFishCompletion()
	default:
		return fmt.Errorf("unsupported shell: %s", shell)
	}

	// Create directory if it doesn't exist
	dir := strings.TrimSuffix(scriptPath, "/"+scriptPath[strings.LastIndex(scriptPath, "/")+1:])
	if err := os.MkdirAll(dir, 0755); err != nil {
		return fmt.Errorf("failed to create completion directory: %w", err)
	}

	// Write completion script
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
