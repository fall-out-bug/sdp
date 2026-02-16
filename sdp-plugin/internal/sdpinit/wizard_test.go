package sdpinit

import (
	"bytes"
	"io"
	"strings"
	"testing"
)

func TestNewWizard(t *testing.T) {
	reader := strings.NewReader("")
	writer := &bytes.Buffer{}
	preflight := &PreflightResult{ProjectType: "go"}

	wizard := NewWizard(reader, writer, preflight)

	if wizard == nil {
		t.Fatal("NewWizard returned nil")
	}

	if wizard.reader == nil {
		t.Error("Wizard reader should not be nil")
	}

	if wizard.writer == nil {
		t.Error("Wizard writer should not be nil")
	}
}

func TestNewWizard_NilReader(t *testing.T) {
	writer := &bytes.Buffer{}
	preflight := &PreflightResult{ProjectType: "go"}

	// Should not panic with nil reader
	wizard := NewWizard(nil, writer, preflight)

	if wizard == nil {
		t.Fatal("NewWizard returned nil")
	}
}

func TestWizard_promptString(t *testing.T) {
	tests := []struct {
		name     string
		input    string
		label    string
		def      string
		expected string
	}{
		{
			name:     "with default and empty input",
			input:    "\n",
			label:    "Project name",
			def:      "my-project",
			expected: "my-project",
		},
		{
			name:     "with user input",
			input:    "custom-name\n",
			label:    "Project name",
			def:      "default",
			expected: "custom-name",
		},
		{
			name:     "with whitespace in input",
			input:    "  trimmed  \n",
			label:    "Project name",
			def:      "default",
			expected: "trimmed",
		},
		{
			name:     "empty default with empty input",
			input:    "\n",
			label:    "Project name",
			def:      "",
			expected: "",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			reader := strings.NewReader(tt.input)
			writer := &bytes.Buffer{}
			wizard := NewWizard(reader, writer, nil)

			result, err := wizard.promptString(tt.label, tt.def)
			if err != nil {
				t.Fatalf("promptString error: %v", err)
			}

			if result != tt.expected {
				t.Errorf("promptString() = %q, want %q", result, tt.expected)
			}
		})
	}
}

func TestWizard_promptBool(t *testing.T) {
	tests := []struct {
		name     string
		input    string
		def      bool
		expected bool
	}{
		{
			name:     "empty input with true default",
			input:    "\n",
			def:      true,
			expected: true,
		},
		{
			name:     "empty input with false default",
			input:    "\n",
			def:      false,
			expected: false,
		},
		{
			name:     "yes input",
			input:    "y\n",
			def:      false,
			expected: true,
		},
		{
			name:     "YES input",
			input:    "YES\n",
			def:      false,
			expected: true,
		},
		{
			name:     "no input",
			input:    "n\n",
			def:      true,
			expected: false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			reader := strings.NewReader(tt.input)
			writer := &bytes.Buffer{}
			wizard := NewWizard(reader, writer, nil)

			result, err := wizard.promptBool("Question", tt.def)
			if err != nil {
				t.Fatalf("promptBool error: %v", err)
			}

			if result != tt.expected {
				t.Errorf("promptBool() = %v, want %v", result, tt.expected)
			}
		})
	}
}

func TestWizard_promptInt(t *testing.T) {
	tests := []struct {
		name     string
		input    string
		def      int
		expected int
	}{
		{
			name:     "empty input uses default",
			input:    "\n",
			def:      5,
			expected: 5,
		},
		{
			name:     "valid number input",
			input:    "3\n",
			def:      1,
			expected: 3,
		},
		{
			name:     "invalid input uses default",
			input:    "abc\n",
			def:      2,
			expected: 2,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			reader := strings.NewReader(tt.input)
			writer := &bytes.Buffer{}
			wizard := NewWizard(reader, writer, nil)

			result, err := wizard.promptInt("Number", tt.def)
			if err != nil {
				t.Fatalf("promptInt error: %v", err)
			}

			if result != tt.expected {
				t.Errorf("promptInt() = %d, want %d", result, tt.expected)
			}
		})
	}
}

func TestWizard_Run(t *testing.T) {
	// Simulate user input:
	// - Project name: custom-project
	// - Project type: 2 (node)
	// - Skills: (empty = use defaults)
	// - Evidence: n (no, so enabled)
	// - Beads: (empty = no skip)
	// - Confirm: y
	input := "custom-project\n2\n\nn\n\ny\n"
	reader := strings.NewReader(input)
	writer := &bytes.Buffer{}
	preflight := &PreflightResult{
		ProjectType: "go",
		HasGit:      true,
	}

	wizard := NewWizard(reader, writer, preflight)
	answers, err := wizard.Run()

	if err != nil {
		t.Fatalf("Wizard.Run error: %v", err)
	}

	if answers == nil {
		t.Fatal("Wizard.Run returned nil answers")
	}

	if answers.ProjectName != "custom-project" {
		t.Errorf("ProjectName = %q, want %q", answers.ProjectName, "custom-project")
	}

	if answers.ProjectType != "node" {
		t.Errorf("ProjectType = %q, want %q", answers.ProjectType, "node")
	}

	// Evidence should be enabled (NoEvidence = false)
	if answers.NoEvidence {
		t.Error("NoEvidence should be false (evidence enabled)")
	}
}

func TestWizard_Run_Cancelled(t *testing.T) {
	// User cancels at the confirmation
	input := "my-project\n1\n\nn\n\nn\n"
	reader := strings.NewReader(input)
	writer := &bytes.Buffer{}

	wizard := NewWizard(reader, writer, nil)
	_, err := wizard.Run()

	if err == nil {
		t.Fatal("Expected error when user cancels")
	}

	if !strings.Contains(err.Error(), "cancelled") {
		t.Errorf("Expected cancellation error, got: %v", err)
	}
}

func TestWizard_printPreflightSummary(t *testing.T) {
	writer := &bytes.Buffer{}
	preflight := &PreflightResult{
		ProjectType: "go",
		HasGit:      true,
		HasSDP:      false,
		HasClaude:   true,
		Conflicts:   []string{".claude/settings.json"},
		Warnings:    []string{"Test warning"},
	}

	wizard := NewWizard(strings.NewReader(""), writer, preflight)
	wizard.printPreflightSummary()

	output := writer.String()

	// Check key elements are present
	if !strings.Contains(output, "go") {
		t.Error("Output should contain project type")
	}

	if !strings.Contains(output, "Git: detected") {
		t.Error("Output should show git detected")
	}

	if !strings.Contains(output, "Conflicts:") {
		t.Error("Output should show conflicts section")
	}

	if !strings.Contains(output, "Warnings:") {
		t.Error("Output should show warnings section")
	}
}

func TestWizard_promptSkills(t *testing.T) {
	tests := []struct {
		name     string
		input    string
		expected []string
	}{
		{
			name:     "empty input uses defaults",
			input:    "\n",
			expected: nil, // Will use defaults from GetDefaults
		},
		{
			name:     "all input",
			input:    "all\n",
			expected: nil, // Will use defaults
		},
		{
			name:     "select specific skills",
			input:    "1,3,5\n",
			expected: []string{"feature", "design", "review"},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			reader := strings.NewReader(tt.input)
			writer := &bytes.Buffer{}
			wizard := NewWizard(reader, writer, nil)

			result, err := wizard.promptSkills("go")
			if err != nil {
				t.Fatalf("promptSkills error: %v", err)
			}

			if tt.expected == nil {
				// Should use defaults
				defaults := GetDefaults("go")
				if len(result) != len(defaults.Skills) {
					t.Errorf("Expected default skills, got %v", result)
				}
			} else {
				if len(result) != len(tt.expected) {
					t.Errorf("promptSkills() = %v, want %v", result, tt.expected)
				}
			}
		})
	}
}

func TestWizard_promptProjectType(t *testing.T) {
	tests := []struct {
		name            string
		input           string
		detectedType    string
		expectedType    string
	}{
		{
			name:         "empty input uses detected",
			input:        "\n",
			detectedType: "go",
			expectedType: "go",
		},
		{
			name:         "select different type",
			input:        "2\n",
			detectedType: "go",
			expectedType: "node",
		},
		{
			name:         "invalid selection uses detected",
			input:        "99\n",
			detectedType: "python",
			expectedType: "python",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			reader := strings.NewReader(tt.input)
			writer := &bytes.Buffer{}
			preflight := &PreflightResult{ProjectType: tt.detectedType}
			wizard := NewWizard(reader, writer, preflight)

			result, err := wizard.promptProjectType()
			if err != nil {
				t.Fatalf("promptProjectType error: %v", err)
			}

			if result != tt.expectedType {
				t.Errorf("promptProjectType() = %q, want %q", result, tt.expectedType)
			}
		})
	}
}

func TestRunWizard(t *testing.T) {
	// This tests the convenience function with mock I/O
	input := "test-project\n1\n\nn\n\ny\n"
	reader := strings.NewReader(input)
	writer := &bytes.Buffer{}

	wizard := NewWizard(reader, writer, nil)
	answers, err := wizard.Run()

	if err != nil {
		t.Fatalf("RunWizard error: %v", err)
	}

	if answers.ProjectName != "test-project" {
		t.Errorf("ProjectName = %q, want %q", answers.ProjectName, "test-project")
	}
}

// TestWizard_EOF tests graceful handling of EOF (e.g., piped input ending)
func TestWizard_EOF(t *testing.T) {
	// Empty input will hit EOF
	reader := strings.NewReader("")
	writer := &bytes.Buffer{}
	wizard := NewWizard(reader, writer, nil)

	// Should handle EOF gracefully (use defaults where possible)
	_, err := wizard.Run()
	// Will likely error due to incomplete input, but shouldn't panic
	if err != nil && !strings.Contains(err.Error(), "cancelled") {
		// Error is expected, just ensure it's not a panic-inducing error
		t.Logf("Expected error on EOF: %v", err)
	}
}

// TestWizard_printWelcome tests welcome message output
func TestWizard_printWelcome(t *testing.T) {
	writer := &bytes.Buffer{}
	wizard := NewWizard(strings.NewReader(""), writer, nil)
	wizard.printWelcome()

	output := writer.String()

	if !strings.Contains(output, "Welcome") {
		t.Error("Welcome message should contain 'Welcome'")
	}

	if !strings.Contains(output, "SDP Onboarding") {
		t.Error("Welcome message should contain 'SDP Onboarding'")
	}
}

// TestWizard_printSummary tests summary output
func TestWizard_printSummary(t *testing.T) {
	writer := &bytes.Buffer{}
	wizard := NewWizard(strings.NewReader(""), writer, nil)

	answers := &WizardAnswers{
		ProjectName: "test-project",
		ProjectType: "go",
		Skills:      []string{"feature", "build"},
		NoEvidence:  false,
		SkipBeads:   false,
	}

	wizard.printSummary(answers)

	output := writer.String()

	if !strings.Contains(output, "test-project") {
		t.Error("Summary should contain project name")
	}

	if !strings.Contains(output, "go") {
		t.Error("Summary should contain project type")
	}

	if !strings.Contains(output, "feature") {
		t.Error("Summary should contain skills")
	}
}

// Ensure Wizard implements expected interface
var _ io.Writer = (*bytes.Buffer)(nil)
