package ui

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
	local -a reality_commands

	commands=(
		'init:Initialize project with SDP prompts'
		'doctor:Check environment (Git, Claude Code, .claude/)'
		'status:Show current project state'
		'next:Get next-step recommendation'
		'demo:Run a guided first-success walkthrough'
		'hooks:Manage Git hooks for SDP'
		'parse:Parse SDP workstream files'
		'beads:Interact with Beads task tracker'
		'tdd:Run TDD cycle (Red-Green-Refactor)'
		'drift:Detect code drift from specification'
		'quality:Check code quality gates'
		'reality:Emit and validate OSS reality baseline artifacts'
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
		'create:Create a task and persist repo snapshot'
		'show:Show task details'
		'update:Update task status'
		'close:Close a task and persist repo snapshot'
		'sync:Synchronize Beads state'
	)

	quality_commands=(
		'check:Run quality checks'
		'gate:Verify quality gates pass'
		'scan:Scan for quality issues'
		'report:Generate quality report'
	)

	reality_commands=(
		'emit-oss:Emit the OSS reality artifact set'
		'validate:Validate emitted OSS reality artifacts'
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
		reality)
			_describe 'reality command' reality_commands
			;;
	esac
}

_arguments \
	'1: :->command' \
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
	reality)
		_sdp reality "$words[2,-1]"
		;;
	*)
		_sdp
		;;
esac
`, nil
}
