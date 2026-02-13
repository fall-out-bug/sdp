package ui

import "fmt"

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
