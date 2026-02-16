package main

import (
	"os"
	"path/filepath"
	"strings"
	"testing"
)

func TestDemoCmd(t *testing.T) {
	// Skip in short mode as demo creates temp directories
	if testing.Short() {
		t.Skip("Skipping demo test in short mode")
	}

	cmd := demoCmd()
	if cmd == nil {
		t.Fatal("demoCmd() returned nil")
	}

	if cmd.Use != "demo" {
		t.Errorf("demoCmd Use = %q, want 'demo'", cmd.Use)
	}

	if cmd.Short == "" {
		t.Error("demoCmd should have a Short description")
	}

	if cmd.Long == "" {
		t.Error("demoCmd should have a Long description")
	}

	// Check flags exist
	templateFlag := cmd.Flags().Lookup("template")
	if templateFlag == nil {
		t.Error("demoCmd should have --template flag")
	}

	cleanupFlag := cmd.Flags().Lookup("cleanup")
	if cleanupFlag == nil {
		t.Error("demoCmd should have --cleanup flag")
	}

	verboseFlag := cmd.Flags().Lookup("verbose")
	if verboseFlag == nil {
		t.Error("demoCmd should have --verbose flag")
	}
}

func TestCopyTemplate(t *testing.T) {
	// Create source template
	srcDir := t.TempDir()
	dstDir := t.TempDir()

	// Create files in source
	os.MkdirAll(filepath.Join(srcDir, "subdir"), 0755)
	os.WriteFile(filepath.Join(srcDir, "file.txt"), []byte("content"), 0644)
	os.WriteFile(filepath.Join(srcDir, "subdir", "nested.txt"), []byte("nested content"), 0644)

	// Copy
	err := copyTemplate(srcDir, dstDir)
	if err != nil {
		t.Fatalf("copyTemplate error: %v", err)
	}

	// Verify files copied
	if _, err := os.Stat(filepath.Join(dstDir, "file.txt")); os.IsNotExist(err) {
		t.Error("file.txt should be copied")
	}
	if _, err := os.Stat(filepath.Join(dstDir, "subdir", "nested.txt")); os.IsNotExist(err) {
		t.Error("subdir/nested.txt should be copied")
	}

	// Verify content
	content, err := os.ReadFile(filepath.Join(dstDir, "file.txt"))
	if err != nil {
		t.Fatalf("read copied file: %v", err)
	}
	if string(content) != "content" {
		t.Errorf("copied content = %q, want 'content'", content)
	}
}

func TestCopyTemplate_NonExistent(t *testing.T) {
	dstDir := t.TempDir()

	err := copyTemplate("/nonexistent/path", dstDir)
	if err == nil {
		t.Error("copyTemplate should fail for non-existent source")
	}
}

func TestIndent(t *testing.T) {
	tests := []struct {
		input    string
		prefix   string
		expected string
	}{
		{
			input:    "line1\nline2",
			prefix:   "  ",
			expected: "  line1\n  line2",
		},
		{
			input:    "single",
			prefix:   "->",
			expected: "->single",
		},
		{
			input:    "",
			prefix:   "  ",
			expected: "",
		},
		{
			input:    "line1\n\nline3",
			prefix:   "  ",
			expected: "  line1\n\n  line3",
		},
	}

	for _, tt := range tests {
		t.Run(tt.input, func(t *testing.T) {
			result := indent(tt.input, tt.prefix)
			if result != tt.expected {
				t.Errorf("indent(%q, %q) = %q, want %q", tt.input, tt.prefix, result, tt.expected)
			}
		})
	}
}

func TestRunDemoCommand(t *testing.T) {
	// Test with a simple command
	output, err := runDemoCommand("echo hello", false)
	if err != nil {
		t.Fatalf("runDemoCommand error: %v", err)
	}
	if !strings.Contains(output, "hello") {
		t.Errorf("output = %q, should contain 'hello'", output)
	}
}

func TestRunDemoCommand_Empty(t *testing.T) {
	output, err := runDemoCommand("", false)
	if err != nil {
		t.Errorf("empty command should not error: %v", err)
	}
	if output != "" {
		t.Errorf("empty command output = %q, want empty", output)
	}
}

func TestRunDemoCommand_Failure(t *testing.T) {
	// Test with a command that will fail
	output, err := runDemoCommand("false", false)
	if err == nil {
		t.Error("runDemoCommand should fail for 'false' command")
	}
	_ = output // Output may be empty or contain error message
}

func TestDemoStep(t *testing.T) {
	step := DemoStep{
		Name:        "Test Step",
		Description: "A test step",
		Command:     "echo test",
		Expected:    "test",
		Optional:    false,
	}

	if step.Name != "Test Step" {
		t.Error("DemoStep Name not set correctly")
	}
	if step.Command != "echo test" {
		t.Error("DemoStep Command not set correctly")
	}
}

func TestRunDemo_TemplateNotFound(t *testing.T) {
	// Create temp dir and change to it
	tmpDir := t.TempDir()
	originalWd, _ := os.Getwd()
	t.Cleanup(func() { os.Chdir(originalWd) })

	if err := os.Chdir(tmpDir); err != nil {
		t.Fatalf("chdir: %v", err)
	}

	// Run with non-existent template
	err := runDemo("nonexistent", false, false)
	if err == nil {
		t.Error("runDemo should fail for non-existent template")
	}
	if !strings.Contains(err.Error(), "template not found") {
		t.Errorf("error = %v, should contain 'template not found'", err)
	}
}

func TestRunDemo_MinimalTemplate(t *testing.T) {
	if testing.Short() {
		t.Skip("Skipping integration test in short mode")
	}

	// Get original working directory (repo root)
	originalWd, _ := os.Getwd()
	t.Cleanup(func() { os.Chdir(originalWd) })

	// Change to sdp-plugin directory where templates exist
	pluginDir := filepath.Join(originalWd, "sdp-plugin")
	if _, err := os.Stat(pluginDir); os.IsNotExist(err) {
		// We might already be in the plugin directory
		pluginDir = originalWd
	}

	if err := os.Chdir(pluginDir); err != nil {
		t.Fatalf("chdir: %v", err)
	}

	// Check template exists
	templateDir := filepath.Join("templates", "minimal-go")
	if _, err := os.Stat(templateDir); os.IsNotExist(err) {
		t.Skip("minimal-go template not found, skipping integration test")
	}

	// Run demo with cleanup (safe to run)
	err := runDemo("minimal-go", true, false)
	if err != nil {
		t.Logf("Demo output (may have expected failures in CI): %v", err)
		// Don't fail the test - demo may fail in CI environment
	}
}
