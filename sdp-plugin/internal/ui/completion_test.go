package ui

import (
	"strings"
	"testing"
)

func TestGenerateBashCompletion(t *testing.T) {
	script, err := generateBashCompletion()
	if err != nil {
		t.Fatalf("generateBashCompletion() error = %v", err)
	}

	if !strings.Contains(script, "_sdp_completion") {
		t.Error("Bash completion should contain _sdp_completion function")
	}
	if !strings.Contains(script, "complete -F _sdp_completion sdp") {
		t.Error("Bash completion should register completion for sdp command")
	}
	if !strings.Contains(script, "checkpoint") {
		t.Error("Bash completion should include checkpoint commands")
	}
	if !strings.Contains(script, "orchestrate") {
		t.Error("Bash completion should include orchestrate commands")
	}
}

func TestGenerateZshCompletion(t *testing.T) {
	script, err := generateZshCompletion()
	if err != nil {
		t.Fatalf("generateZshCompletion() error = %v", err)
	}

	if !strings.Contains(script, "#compdef sdp") {
		t.Error("Zsh completion should contain #compdef sdp directive")
	}
	if !strings.Contains(script, "_sdp()") {
		t.Error("Zsh completion should contain _sdp function")
	}
	if !strings.Contains(script, "init:") {
		t.Error("Zsh completion should include init command description")
	}
	if !strings.Contains(script, "checkpoint:") {
		t.Error("Zsh completion should include checkpoint command description")
	}
}

func TestGenerateFishCompletion(t *testing.T) {
	script, err := generateFishCompletion()
	if err != nil {
		t.Fatalf("generateFishCompletion() error = %v", err)
	}

	if !strings.Contains(script, "complete -c sdp") {
		t.Error("Fish completion should register completion for sdp command")
	}
	if !strings.Contains(script, "complete -c sdp -n") {
		t.Error("Fish completion should contain subcommand completions")
	}
	if !strings.Contains(script, "__sdp_checkpoint_commands") {
		t.Error("Fish completion should define checkpoint command helper")
	}
}

func TestGenerateCompletionInvalidShell(t *testing.T) {
	invalidShell := CompletionType("invalid")
	err := GenerateCompletion(invalidShell)

	if err == nil {
		t.Error("GenerateCompletion() should return error for invalid shell")
	}
	if !strings.Contains(err.Error(), "unsupported shell") {
		t.Errorf("Error message should mention unsupported shell, got: %v", err)
	}
}

func TestGenerateCompletionBash(t *testing.T) {
	err := GenerateCompletion(Bash)
	if err != nil {
		t.Fatalf("GenerateCompletion(Bash) error = %v", err)
	}
}

func TestGenerateCompletionZsh(t *testing.T) {
	err := GenerateCompletion(Zsh)
	if err != nil {
		t.Fatalf("GenerateCompletion(Zsh) error = %v", err)
	}
}

func TestGenerateCompletionFish(t *testing.T) {
	err := GenerateCompletion(Fish)
	if err != nil {
		t.Fatalf("GenerateCompletion(Fish) error = %v", err)
	}
}

func TestBashCompletionIncludesAllCommands(t *testing.T) {
	script, err := generateBashCompletion()
	if err != nil {
		t.Fatalf("generateBashCompletion() error = %v", err)
	}

	requiredCommands := []string{
		"init", "doctor", "hooks", "parse", "beads",
		"tdd", "drift", "quality", "watch", "telemetry",
		"checkpoint", "orchestrate",
	}

	for _, cmd := range requiredCommands {
		if !strings.Contains(script, cmd) {
			t.Errorf("Bash completion should include %s command", cmd)
		}
	}
}

func TestZshCompletionIncludesDescriptions(t *testing.T) {
	script, err := generateZshCompletion()
	if err != nil {
		t.Fatalf("generateZshCompletion() error = %v", err)
	}

	// Check that commands have descriptions
	if !strings.Contains(script, "init:Initialize") {
		t.Error("Zsh completion should include descriptions for commands")
	}
	if !strings.Contains(script, "doctor:Check") {
		t.Error("Zsh completion should include descriptions for doctor")
	}
}

func TestFishCompletionIncludesSubcommands(t *testing.T) {
	script, err := generateFishCompletion()
	if err != nil {
		t.Fatalf("generateFishCompletion() error = %v", err)
	}

	// Check for checkpoint subcommands
	if !strings.Contains(script, "create\\tCreate a new checkpoint") {
		t.Error("Fish completion should include create subcommand with description")
	}
	if !strings.Contains(script, "resume\\tResume from") {
		t.Error("Fish completion should include resume subcommand")
	}
	if !strings.Contains(script, "list\\tList all") {
		t.Error("Fish completion should include list subcommand")
	}
	if !strings.Contains(script, "clean\\tClean old") {
		t.Error("Fish completion should include clean subcommand")
	}
}

func TestCompletionTypeValues(t *testing.T) {
	tests := []struct {
		name  string
		shell CompletionType
	}{
		{"Bash", Bash},
		{"Zsh", Zsh},
		{"Fish", Fish},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if tt.shell == "" {
				t.Errorf("CompletionType %s should not be empty", tt.name)
			}
		})
	}
}
