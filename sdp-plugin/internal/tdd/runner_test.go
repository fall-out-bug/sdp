package tdd

import (
	"context"
	"os"
	"path/filepath"
	"testing"
)

func TestRedPhaseFails(t *testing.T) {
	// Create a runner for a Go project
	runner := NewRunner(Go)

	ctx := context.Background()
	// Use the parser package which should have passing tests
	result, err := runner.RunPhase(ctx, Red, "./internal/parser")

	// Red phase expects tests to fail
	// Since tests will pass, we should get an error
	if err == nil {
		t.Error("Expected Red phase to return error when tests pass, got nil")
	}

	if result == nil {
		t.Fatal("Expected result, got nil")
	}

	// Success should be true (tests passed), but phase should return error
	// because Red phase expects failure
	if !result.Success {
		t.Error("Expected tests to pass (Success=true), got Success=false")
	}

	if result.Phase != Red {
		t.Errorf("Expected phase Red, got %v", result.Phase)
	}
}

func TestGreenPhaseSucceeds(t *testing.T) {
	runner := NewRunner(Go)

	ctx := context.Background()
	result, err := runner.RunPhase(ctx, Green, "./internal/parser")

	if err != nil {
		t.Errorf("Expected Green phase to succeed, got error: %v\nStderr: %s\nStdout: %s", err, result.Stderr, result.Stdout)
	}

	if result == nil {
		t.Fatal("Expected result, got nil")
	}

	if !result.Success {
		t.Error("Expected Green phase to succeed (Success=true), got Success=false")
	}

	if result.Phase != Green {
		t.Errorf("Expected phase Green, got %v", result.Phase)
	}
}

func TestRefactorPhasePasses(t *testing.T) {
	runner := NewRunner(Go)

	ctx := context.Background()
	result, err := runner.RunPhase(ctx, Refactor, "./internal/parser")

	if err != nil {
		t.Errorf("Expected Refactor phase to pass, got error: %v", err)
	}

	if result == nil {
		t.Fatal("Expected result, got nil")
	}

	if !result.Success {
		t.Error("Expected Refactor phase to succeed (Success=true), got Success=false")
	}

	if result.Phase != Refactor {
		t.Errorf("Expected phase Refactor, got %v", result.Phase)
	}
}

func TestDetectGoProject(t *testing.T) {
	runner := &Runner{}

	// We need to detect the project root (sdp-plugin has go.mod)
	// But we're running from internal/tdd directory during tests
	// So we use the parent directories
	lang, err := runner.DetectLanguage("../..")
	if err != nil {
		t.Fatalf("Expected to detect Go project, got error: %v", err)
	}

	if lang != Go {
		t.Errorf("Expected language Go, got %v", lang)
	}
}

func TestDetectPythonProject(t *testing.T) {
	// Skip if no Python test project available
	t.Skip("Needs test Python project setup")

	runner := &Runner{}
	lang, err := runner.DetectLanguage("../../src/sdp")
	if err != nil {
		t.Fatalf("Expected to detect Python project, got error: %v", err)
	}

	if lang != Python {
		t.Errorf("Expected language Python, got %v", lang)
	}
}

func TestDetectJavaProject(t *testing.T) {
	// Skip if no Java test project available
	t.Skip("Needs test Java project setup")

	runner := &Runner{}
	lang, err := runner.DetectLanguage("some/java/path")
	if err != nil {
		t.Fatalf("Expected to detect Java project, got error: %v", err)
	}

	if lang != Java {
		t.Errorf("Expected language Java, got %v", lang)
	}
}

func TestContextCancellation(t *testing.T) {
	runner := NewRunner(Go)

	// Create context that's already cancelled
	ctx, cancel := context.WithCancel(context.Background())
	cancel() // Cancel immediately

	result, err := runner.RunPhase(ctx, Red, "./internal/parser")

	// Should return error due to cancellation
	if err == nil {
		t.Error("Expected error when context cancelled, got nil")
	}

	if result != nil && result.Success {
		t.Error("Expected phase to fail due to cancellation")
	}
}

func TestRunAllPhases(t *testing.T) {
	runner := NewRunner(Go)

	ctx := context.Background()
	results, err := runner.RunAllPhases(ctx, "./internal/parser")

	if err != nil {
		t.Errorf("Expected all phases to complete, got error: %v", err)
	}

	if len(results) != 3 {
		t.Errorf("Expected 3 phase results, got %d", len(results))
		return
	}

	// Check Red phase
	if results[0].Phase != Red {
		t.Errorf("Expected first phase to be Red, got %v", results[0].Phase)
	}

	// Check Green phase
	if results[1].Phase != Green {
		t.Errorf("Expected second phase to be Green, got %v", results[1].Phase)
	}

	// Check Refactor phase
	if results[2].Phase != Refactor {
		t.Errorf("Expected third phase to be Refactor, got %v", results[2].Phase)
	}
}

func TestPhaseResultFields(t *testing.T) {
	runner := NewRunner(Go)

	ctx := context.Background()
	result, err := runner.RunPhase(ctx, Green, "./internal/parser")

	if err != nil {
		t.Fatalf("Expected Green phase to succeed, got error: %v", err)
	}

	// Check all fields are populated
	if result.Phase != Green {
		t.Error("Phase field not set correctly")
	}

	if result.Duration == 0 {
		t.Error("Duration field not set")
	}

	if result.Stdout == "" && result.Stderr == "" {
		t.Error("Stdout and Stderr both empty - expected some output")
	}
}

func TestDetectLanguageByAbsolutePath(t *testing.T) {
	// Create temporary directory with go.mod
	tmpDir := t.TempDir()
	goModPath := filepath.Join(tmpDir, "go.mod")
	err := os.WriteFile(goModPath, []byte("module test\n"), 0644)
	if err != nil {
		t.Fatalf("Failed to create go.mod: %v", err)
	}

	runner := &Runner{}
	lang, err := runner.DetectLanguage(tmpDir)

	if err != nil {
		t.Errorf("Expected to detect Go project, got error: %v", err)
	}

	if lang != Go {
		t.Errorf("Expected language Go, got %v", lang)
	}
}

func TestPhaseString(t *testing.T) {
	tests := []struct {
		phase    Phase
		expected string
	}{
		{Red, "Red"},
		{Green, "Green"},
		{Refactor, "Refactor"},
		{Phase(99), "Unknown"},
	}

	for _, tt := range tests {
		t.Run(tt.expected, func(t *testing.T) {
			if got := tt.phase.String(); got != tt.expected {
				t.Errorf("Phase.String() = %v, want %v", got, tt.expected)
			}
		})
	}
}

func TestLanguageString(t *testing.T) {
	tests := []struct {
		language Language
		expected string
	}{
		{Python, "Python"},
		{Go, "Go"},
		{Java, "Java"},
		{Language(99), "Unknown"},
	}

	for _, tt := range tests {
		t.Run(tt.expected, func(t *testing.T) {
			if got := tt.language.String(); got != tt.expected {
				t.Errorf("Language.String() = %v, want %v", got, tt.expected)
			}
		})
	}
}

func TestNewRunnerDefaults(t *testing.T) {
	runner := NewRunner(Go)

	if runner.language != Go {
		t.Errorf("Expected language Go, got %v", runner.language)
	}

	if runner.testCmd != "go test" {
		t.Errorf("Expected testCmd 'go test', got %v", runner.testCmd)
	}
}

func TestBuildTestCommandPython(t *testing.T) {
	runner := NewRunner(Python)
	cmd := runner.buildTestCommand("some/path")

	// Args[0] is the command (pytest), Args[1] is the path, Args[2] is -v
	if cmd.Args[0] != "pytest" {
		t.Errorf("Expected pytest, got %v", cmd.Args[0])
	}

	if cmd.Args[1] != "some/path" {
		t.Errorf("Expected path 'some/path', got %v", cmd.Args[1])
	}
}

func TestBuildTestCommandJava(t *testing.T) {
	runner := NewRunner(Java)
	cmd := runner.buildTestCommand("some/path")

	// Args[0] is mvn, Args[1] is test, Args[2] is -f, Args[3] is path
	if cmd.Args[0] != "mvn" {
		t.Errorf("Expected mvn, got %v", cmd.Args[0])
	}

	if cmd.Args[1] != "test" {
		t.Errorf("Expected test, got %v", cmd.Args[1])
	}

	if cmd.Args[3] != "some/path" {
		t.Errorf("Expected path 'some/path', got %v", cmd.Args[3])
	}
}

func TestFindProjectRoot(t *testing.T) {
	runner := NewRunner(Go)
	if runner.projectRoot == "" {
		t.Error("Project root is empty")
	}
}
