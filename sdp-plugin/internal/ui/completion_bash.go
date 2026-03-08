package ui

// generateBashCompletion generates bash completion script
func generateBashCompletion() (string, error) {
	return `# Bash completion for SDP
# Source this file in your .bashrc or .bash_profile:
#   source <(sdp completion bash)

_sdp_completion() {
	local cur prev words cword
	_init_completion || return

	local commands="init doctor status next demo hooks plan apply build verify log guard parse beads tdd drift quality watch telemetry checkpoint orchestrate contract health diagnose memory decisions"
	local checkpoint_commands="create resume list clean"
	local orchestrate_commands="start status stop"
	local log_commands="show export stats trace"
	local guard_commands="activate check status deactivate context branch finding"
	local contract_commands="synthesize generate lock validate verify"

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
		log)
			COMPREPLY=($(compgen -W "${log_commands}" -- "${cur}"))
			return
			;;
		guard)
			COMPREPLY=($(compgen -W "${guard_commands}" -- "${cur}"))
			return
			;;
		contract)
			COMPREPLY=($(compgen -W "${contract_commands}" -- "${cur}"))
			return
			;;
		quality)
			COMPREPLY=($(compgen -W "coverage complexity size types all" -- "${cur}"))
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
