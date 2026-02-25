package verify

import (
	"context"
	"os"
	"path/filepath"
	"testing"
)

func TestNewParser(t *testing.T) {
	parser := NewParser("/tmp/workstreams")
	if parser.wsDir != "/tmp/workstreams" {
		t.Errorf("Expected wsDir to be /tmp/workstreams, got %s", parser.wsDir)
	}
}

func TestNewVerifier(t *testing.T) {
	verifier := NewVerifier("/tmp/workstreams")
	if verifier == nil {
		t.Fatal("Expected verifier to be created")
	}
	if verifier.parser == nil {
		t.Error("Expected parser to be created")
	}
}

func TestParserFindWSFileNotFound(t *testing.T) {
	parser := NewParser("/nonexistent/path")
	_, err := parser.FindWSFile("00-000-00")
	if err == nil {
		t.Error("Expected error when finding nonexistent workstream")
	}
}

func TestParserParseWSFile(t *testing.T) {
	// Create temp file
	tmpDir := t.TempDir()
	wsFile := filepath.Join(tmpDir, "test-ws.md")

	content := `---
ws_id: 00-067-15
title: Test Workstream
status: in_progress
coverage_threshold: 80.0
scope_files:
  - file1.go
  - file2.go
verification_commands:
  - go test ./...
---

## Goal
Test workstream for unit tests.
`

	if err := os.WriteFile(wsFile, []byte(content), 0644); err != nil {
		t.Fatalf("Failed to create test file: %v", err)
	}

	parser := NewParser(tmpDir)
	data, err := parser.ParseWSFile(wsFile)
	if err != nil {
		t.Fatalf("ParseWSFile failed: %v", err)
	}

	if data.WSID != "00-067-15" {
		t.Errorf("Expected WSID '00-067-15', got %s", data.WSID)
	}
	if data.Title != "Test Workstream" {
		t.Errorf("Expected Title 'Test Workstream', got %s", data.Title)
	}
	if data.Status != "in_progress" {
		t.Errorf("Expected Status 'in_progress', got %s", data.Status)
	}
	if data.CoverageThreshold != 80.0 {
		t.Errorf("Expected CoverageThreshold 80.0, got %f", data.CoverageThreshold)
	}
	if len(data.ScopeFiles) != 2 {
		t.Errorf("Expected 2 scope files, got %d", len(data.ScopeFiles))
	}
	if len(data.VerificationCommands) != 1 {
		t.Errorf("Expected 1 verification command, got %d", len(data.VerificationCommands))
	}
}

func TestParserParseWSFileNoFrontmatter(t *testing.T) {
	tmpDir := t.TempDir()
	wsFile := filepath.Join(tmpDir, "no-frontmatter.md")

	content := `# No frontmatter here
Just content`

	if err := os.WriteFile(wsFile, []byte(content), 0644); err != nil {
		t.Fatalf("Failed to create test file: %v", err)
	}

	parser := NewParser(tmpDir)
	_, err := parser.ParseWSFile(wsFile)
	if err == nil {
		t.Error("Expected error when parsing file without frontmatter")
	}
}

func TestVerifierVerifyOutputFiles(t *testing.T) {
	tmpDir := t.TempDir()

	// Create existing file
	existingFile := filepath.Join(tmpDir, "exists.go")
	if err := os.WriteFile(existingFile, []byte("package test"), 0644); err != nil {
		t.Fatalf("Failed to create test file: %v", err)
	}

	// Use mock PathValidator so path traversal check passes for temp dir (projectRoot may differ)
	verifier := NewVerifierWithOptions(tmpDir, WithPathValidator(mockPathValidator{}))
	wsData := &WorkstreamData{
		ScopeFiles: []string{existingFile, filepath.Join(tmpDir, "missing.go")},
	}

	checks := verifier.VerifyOutputFiles(wsData)

	if len(checks) != 2 {
		t.Fatalf("Expected 2 checks, got %d", len(checks))
	}

	// First file exists
	if !checks[0].Passed {
		t.Error("Expected first file check to pass")
	}

	// Second file doesn't exist
	if checks[1].Passed {
		t.Error("Expected second file check to fail")
	}
}

func TestVerifierVerifyCoverage(t *testing.T) {
	// Use mock to avoid running real coverage (which may fail in /tmp)
	mock := &mockCoverageChecker{
		result: &CoverageResult{Coverage: 85.0, Threshold: 80.0, Report: "ok"},
	}
	verifier := NewVerifierWithOptions("/tmp", WithCoverageChecker(mock))

	// No coverage threshold
	wsData := &WorkstreamData{CoverageThreshold: 0}
	result := verifier.VerifyCoverage(context.Background(), wsData)
	if result != nil {
		t.Error("Expected nil when no coverage threshold")
	}

	// With coverage threshold
	wsData = &WorkstreamData{CoverageThreshold: 80.0}
	result = verifier.VerifyCoverage(context.Background(), wsData)
	if result == nil {
		t.Fatal("Expected result with coverage threshold")
	}
	if !result.Passed {
		t.Error("Expected coverage check to pass with mock")
	}
}

// mockCoverageChecker is a test double for CoverageChecker.
type mockCoverageChecker struct {
	result *CoverageResult
	err    error
}

func (m *mockCoverageChecker) CheckCoverage(ctx context.Context) (*CoverageResult, error) {
	if m.err != nil {
		return nil, m.err
	}
	return m.result, nil
}

// mockPathValidator is a test double that always passes (no path traversal check).
type mockPathValidator struct{}

func (mockPathValidator) ValidatePathInDirectory(baseDir, targetPath string) error {
	return nil
}

func TestVerifierVerifyCoverageWithMock(t *testing.T) {
	// Inject mock that returns 90% coverage (above threshold)
	mock := &mockCoverageChecker{
		result: &CoverageResult{Coverage: 90.0, Threshold: 80.0, Report: "mock report"},
	}
	verifier := NewVerifierWithOptions("/tmp", WithCoverageChecker(mock))

	wsData := &WorkstreamData{CoverageThreshold: 80.0}
	result := verifier.VerifyCoverage(context.Background(), wsData)
	if result == nil {
		t.Fatal("Expected result")
	}
	if !result.Passed {
		t.Errorf("Expected pass with 90%% coverage, got %s", result.Message)
	}
	if result.Evidence != "mock report" {
		t.Errorf("Expected 'mock report' evidence, got %s", result.Evidence)
	}

	// Inject mock that returns 50% (below threshold)
	mock.result = &CoverageResult{Coverage: 50.0, Threshold: 80.0, Report: "low"}
	result = verifier.VerifyCoverage(context.Background(), wsData)
	if result.Passed {
		t.Error("Expected fail with 50% coverage")
	}
}

func TestTruncate(t *testing.T) {
	tests := []struct {
		input    string
		maxLen   int
		expected string
	}{
		{"short", 10, "short"},
		{"this is a long string", 10, "this is a ..."},
		{"exact", 5, "exact"},
	}

	for _, tt := range tests {
		result := truncate(tt.input, tt.maxLen)
		if result != tt.expected {
			t.Errorf("truncate(%q, %d) = %q, want %q", tt.input, tt.maxLen, result, tt.expected)
		}
	}
}

func TestVerifierVerifyWSNotFound(t *testing.T) {
	verifier := NewVerifier("/nonexistent/path")
	result := verifier.Verify(context.Background(), "00-000-00")

	if result.Passed {
		t.Error("Expected verification to fail for nonexistent WS")
	}
	if len(result.Checks) == 0 {
		t.Error("Expected at least one check")
	}
}

func TestVerifierCommandsEmptyCommand(t *testing.T) {
	verifier := NewVerifier("/tmp")

	wsData := &WorkstreamData{
		VerificationCommands: []string{""}, // Empty command
	}

	checks := verifier.VerifyCommands(context.Background(), wsData)

	if len(checks) != 1 {
		t.Fatalf("Expected 1 check, got %d", len(checks))
	}

	if checks[0].Passed {
		t.Error("Empty command should fail")
	}
	if checks[0].Message != "Empty command" {
		t.Errorf("Expected 'Empty command' message, got %s", checks[0].Message)
	}
}

func TestVerifierCommandsMultipleCommands(t *testing.T) {
	verifier := NewVerifier("/tmp")

	wsData := &WorkstreamData{
		VerificationCommands: []string{
			"",                        // Empty - should fail
			"nonexistent_command_xyz", // Should fail security validation
		},
	}

	checks := verifier.VerifyCommands(context.Background(), wsData)

	if len(checks) != 2 {
		t.Fatalf("Expected 2 checks, got %d", len(checks))
	}

	// First check should fail (empty)
	if checks[0].Passed {
		t.Error("Empty command should fail")
	}

	// Second check should fail (command doesn't exist)
	if checks[1].Passed {
		t.Error("Nonexistent command should fail")
	}
}

func TestVerifierVerifyWithFilesAndCommands(t *testing.T) {
	// Create temp directory with proper WS directory structure
	tmpDir := t.TempDir()
	backlogDir := filepath.Join(tmpDir, "backlog")
	if err := os.MkdirAll(backlogDir, 0755); err != nil {
		t.Fatalf("Failed to create backlog dir: %v", err)
	}

	wsFile := filepath.Join(backlogDir, "00-999-99-ws.md")

	content := `---
ws_id: 00-999-99
title: Test WS
status: in_progress
scope_files:
  - /nonexistent/file.go
verification_commands:
  - ""
---
## Goal
Test
`
	if err := os.WriteFile(wsFile, []byte(content), 0644); err != nil {
		t.Fatalf("Failed to create test file: %v", err)
	}

	// Use mocks to avoid real coverage run and path validation (temp dir may be outside project root)
	mockCC := &mockCoverageChecker{result: &CoverageResult{Coverage: 100, Threshold: 80, Report: "ok"}}
	verifier := NewVerifierWithOptions(tmpDir, WithCoverageChecker(mockCC), WithPathValidator(mockPathValidator{}))
	result := verifier.Verify(context.Background(), "00-999-99")

	// Should fail because file doesn't exist and command is empty
	if result.Passed {
		t.Error("Expected verification to fail")
	}

	// Should have checks for: file check, command check
	if len(result.Checks) < 1 {
		t.Errorf("Expected at least 1 check, got %d", len(result.Checks))
	}

	// Should have missing files (the nonexistent file)
	if len(result.MissingFiles) == 0 {
		t.Errorf("Expected missing files to be recorded, got: %v", result.MissingFiles)
	}
}

func TestVerifierCheckResultStructure(t *testing.T) {
	check := CheckResult{
		Name:     "Test Check",
		Passed:   true,
		Message:  "Test message",
		Evidence: "Test evidence",
	}

	if check.Name != "Test Check" {
		t.Errorf("Name = %v, want 'Test Check'", check.Name)
	}
	if !check.Passed {
		t.Error("Passed should be true")
	}
	if check.Message != "Test message" {
		t.Errorf("Message = %v, want 'Test message'", check.Message)
	}
	if check.Evidence != "Test evidence" {
		t.Errorf("Evidence = %v, want 'Test evidence'", check.Evidence)
	}
}

func TestVerifierVerificationResultStructure(t *testing.T) {
	result := &VerificationResult{
		WSID:           "00-999-99",
		Passed:         true,
		Checks:         []CheckResult{{Name: "test", Passed: true}},
		MissingFiles:   []string{},
		FailedCommands: []string{},
		Duration:       0,
	}

	if result.WSID != "00-999-99" {
		t.Errorf("WSID = %v, want '00-999-99'", result.WSID)
	}
	if !result.Passed {
		t.Error("Passed should be true")
	}
	if len(result.Checks) != 1 {
		t.Errorf("Expected 1 check, got %d", len(result.Checks))
	}
}

func TestVerifierOutputFilesAbsolutePath(t *testing.T) {
	tmpDir := t.TempDir()

	// Create file with known path
	testFile := filepath.Join(tmpDir, "test.go")
	if err := os.WriteFile(testFile, []byte("package test"), 0644); err != nil {
		t.Fatalf("Failed to create test file: %v", err)
	}

	verifier := NewVerifierWithOptions(tmpDir, WithPathValidator(mockPathValidator{}))
	wsData := &WorkstreamData{
		ScopeFiles: []string{testFile},
	}

	checks := verifier.VerifyOutputFiles(wsData)

	if len(checks) != 1 {
		t.Fatalf("Expected 1 check, got %d", len(checks))
	}

	if !checks[0].Passed {
		t.Error("File check should pass")
	}

	// Evidence should be absolute path
	if !filepath.IsAbs(checks[0].Evidence) {
		t.Errorf("Evidence should be absolute path, got %s", checks[0].Evidence)
	}
}

func TestTruncateEdgeCases(t *testing.T) {
	tests := []struct {
		input    string
		maxLen   int
		expected string
	}{
		{"", 10, ""},
		{"a", 10, "a"},
		{"ab", 1, "a..."},
		{"exact", 5, "exact"},
		{"日本語", 10, "日本語"}, // Unicode within limit
	}

	for _, tt := range tests {
		t.Run(tt.input, func(t *testing.T) {
			result := truncate(tt.input, tt.maxLen)
			if result != tt.expected {
				t.Errorf("truncate(%q, %d) = %q, want %q", tt.input, tt.maxLen, result, tt.expected)
			}
		})
	}
}
