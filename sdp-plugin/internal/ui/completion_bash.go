package ui

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
