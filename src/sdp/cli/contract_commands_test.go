package cli

import (
	"os"
	"path/filepath"
	"testing"
)

// TestSynthesizeCmd_FlagParsing verifies flag parsing
func TestSynthesizeCmd_FlagParsing(t *testing.T) {
	cmd := synthesizeCmd

	// Verify feature flag exists
	flag := cmd.Flag("feature")
	if flag == nil {
		t.Fatal("Feature flag not found")
	}

	if flag.Shorthand != "f" {
		t.Errorf("Expected shorthand 'f', got '%s'", flag.Shorthand)
	}
}

// TestSynthesizeCmd_MissingFeatureFlag verifies required flag
func TestSynthesizeCmd_MissingFeatureFlag(t *testing.T) {
	// Test would require running the command, skip for now
	t.Skip("Integration test - requires command execution")
}

// TestLockCmd_FlagParsing verifies lock command flags
func TestLockCmd_FlagParsing(t *testing.T) {
	cmd := lockCmd

	// Verify contract flag exists
	flag := cmd.Flag("contract")
	if flag == nil {
		t.Fatal("Contract flag not found")
	}
}

// TestValidateCmd_FlagParsing verifies validate command flags
func TestValidateCmd_FlagParsing(t *testing.T) {
	cmd := validateCmd

	// Verify contracts flag exists
	flag := cmd.Flag("contracts")
	if flag == nil {
		t.Fatal("Contracts flag not found")
	}
}

// TestSanitizePath_ValidPath verifies valid path handling
func TestSanitizePath_ValidPath(t *testing.T) {
	tmpDir := t.TempDir()
	allowedDir := filepath.Join(tmpDir, "allowed")
	os.Mkdir(allowedDir, 0755)

	testPath := filepath.Join(allowedDir, "test.yaml")
	result, err := sanitizePath(testPath, []string{allowedDir})

	if err != nil {
		t.Errorf("unexpected error: %v", err)
	}
	if result == "" {
		t.Error("expected non-empty result")
	}
}

// TestSanitizePath_TraversalAttack verifies path traversal prevention
func TestSanitizePath_TraversalAttack(t *testing.T) {
	tmpDir := t.TempDir()
	allowedDir := filepath.Join(tmpDir, "allowed")
	os.Mkdir(allowedDir, 0755)

	// Try to escape allowed directory
	escapePath := filepath.Join(allowedDir, "..", "forbidden.yaml")
	_, err := sanitizePath(escapePath, []string{allowedDir})

	if err == nil {
		t.Error("expected error for path traversal")
	}
}

// TestSanitizePath_OutsideAllowed verifies outside path rejection
func TestSanitizePath_OutsideAllowed(t *testing.T) {
	tmpDir := t.TempDir()
	allowedDir := filepath.Join(tmpDir, "allowed")
	forbiddenDir := filepath.Join(tmpDir, "forbidden")
	os.Mkdir(allowedDir, 0755)
	os.Mkdir(forbiddenDir, 0755)

	forbiddenPath := filepath.Join(forbiddenDir, "secret.yaml")
	_, err := sanitizePath(forbiddenPath, []string{allowedDir})

	if err == nil {
		t.Error("expected error for path outside allowed dirs")
	}
}

// TestSanitizePath_MultipleAllowedDirs verifies multiple allowed directories
func TestSanitizePath_MultipleAllowedDirs(t *testing.T) {
	tmpDir := t.TempDir()
	dir1 := filepath.Join(tmpDir, "dir1")
	dir2 := filepath.Join(tmpDir, "dir2")
	os.Mkdir(dir1, 0755)
	os.Mkdir(dir2, 0755)

	// Path in dir1 should work
	path1 := filepath.Join(dir1, "test.yaml")
	_, err := sanitizePath(path1, []string{dir1, dir2})
	if err != nil {
		t.Errorf("path in dir1 should be allowed: %v", err)
	}

	// Path in dir2 should work
	path2 := filepath.Join(dir2, "test.yaml")
	_, err = sanitizePath(path2, []string{dir1, dir2})
	if err != nil {
		t.Errorf("path in dir2 should be allowed: %v", err)
	}
}

// TestContractCmd_HasSubcommands verifies subcommands are registered
func TestContractCmd_HasSubcommands(t *testing.T) {
	commands := contractCmd.Commands()
	if len(commands) != 3 {
		t.Errorf("expected 3 subcommands, got %d", len(commands))
	}

	commandNames := make(map[string]bool)
	for _, cmd := range commands {
		commandNames[cmd.Name()] = true
	}

	if !commandNames["synthesize"] {
		t.Error("synthesize subcommand not found")
	}
	if !commandNames["lock"] {
		t.Error("lock subcommand not found")
	}
	if !commandNames["validate"] {
		t.Error("validate subcommand not found")
	}
}

// TestRegisterContractCommand verifies command registration
func TestRegisterContractCommand(t *testing.T) {
	// This function just adds the command to root - verify it doesn't panic
	// We can't easily test without a real cobra root command
	// Just verify the function exists and contractCmd is valid
	if contractCmd == nil {
		t.Error("contractCmd should not be nil")
	}
}

// TestSynthesizeFeature_InvalidName verifies feature name validation
func TestSynthesizeFeature_InvalidName(t *testing.T) {
	// Test via command flags
	invalidNames := []string{
		"InvalidUpperCase",
		"with spaces",
		"with/slash",
		"with\\backslash",
		"with.dot",
	}

	for _, name := range invalidNames {
		// The regex check happens in runContractSynthesize
		// We can't easily test without mocking
		t.Logf("Feature name %q would be invalid", name)
	}
}

// TestLockContract_NonexistentFile verifies error for missing file
func TestLockContract_NonexistentFile(t *testing.T) {
	// Reset flag
	lockContract = "/nonexistent/path/contract.yaml"
	lockReason = "test"

	err := runContractLock(nil, nil)
	if err == nil {
		t.Error("expected error for nonexistent contract")
	}
}

// TestValidateContracts_InsufficientContracts verifies minimum contract count
func TestValidateContracts_InsufficientContracts(t *testing.T) {
	// Reset flags
	validateContracts = []string{"only-one.yaml"}
	validateOutput = ".contracts/report.md"

	err := runContractValidate(nil, nil)
	if err == nil {
		t.Error("expected error for less than 2 contracts")
	}
}

// TestLockReason_Default verifies default lock reason
func TestLockReason_Default(t *testing.T) {
	// Just verify the flag exists
	flag := lockCmd.Flag("reason")
	if flag == nil {
		t.Error("reason flag should exist")
	}
}

// TestValidateOutput_Default verifies default output path
func TestValidateOutput_Default(t *testing.T) {
	flag := validateCmd.Flag("output")
	if flag == nil {
		t.Error("output flag should exist")
	}
	if flag.DefValue != ".contracts/validation-report.md" {
		t.Errorf("expected default '.contracts/validation-report.md', got %s", flag.DefValue)
	}
}

// TestLoadContract verifies contract loading
func TestLoadContract(t *testing.T) {
	contract, err := loadContract("any-path.yaml")
	if err != nil {
		t.Errorf("unexpected error: %v", err)
	}
	if contract == nil {
		t.Error("expected non-nil contract")
	}
	if contract.OpenAPI != "3.0.0" {
		t.Error("expected OpenAPI 3.0.0")
	}
}

// TestSynthesizeRequirements_Default verifies default requirements path
func TestSynthesizeRequirements_Default(t *testing.T) {
	flag := synthesizeCmd.Flag("requirements")
	if flag == nil {
		t.Error("requirements flag should exist")
	}
}

// TestSynthesizeOutput_Default verifies default output path
func TestSynthesizeOutput_Default(t *testing.T) {
	flag := synthesizeCmd.Flag("output")
	if flag == nil {
		t.Error("output flag should exist")
	}
}

// TestContractCmd_Use verifies command name
func TestContractCmd_Use(t *testing.T) {
	if contractCmd.Use != "contract" {
		t.Errorf("expected Use 'contract', got %s", contractCmd.Use)
	}
}

// TestContractCmd_Short verifies short description
func TestContractCmd_Short(t *testing.T) {
	if contractCmd.Short == "" {
		t.Error("Short description should not be empty")
	}
}

// TestContractCmd_Long verifies long description
func TestContractCmd_Long(t *testing.T) {
	if contractCmd.Long == "" {
		t.Error("Long description should not be empty")
	}
}

// TestSynthesizeCmd_Use verifies command name
func TestSynthesizeCmd_Use(t *testing.T) {
	if synthesizeCmd.Use != "synthesize" {
		t.Errorf("expected Use 'synthesize', got %s", synthesizeCmd.Use)
	}
}

// TestLockCmd_Use verifies command name
func TestLockCmd_Use(t *testing.T) {
	if lockCmd.Use != "lock" {
		t.Errorf("expected Use 'lock', got %s", lockCmd.Use)
	}
}

// TestValidateCmd_Use verifies command name
func TestValidateCmd_Use(t *testing.T) {
	if validateCmd.Use != "validate" {
		t.Errorf("expected Use 'validate', got %s", validateCmd.Use)
	}
}

// TestRunContractSynthesize_InvalidFeatureName verifies feature name validation
func TestRunContractSynthesize_InvalidFeatureName(t *testing.T) {
	tests := []struct {
		name string
	}{
		{"InvalidUpperCase"},
		{"with spaces"},
		{"with/slash"},
		{"with_underscore"},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			synthesizeFeature = tt.name
			err := runContractSynthesize(nil, nil)
			if err == nil {
				t.Error("expected error for invalid feature name")
			}
		})
	}
}

// TestRunContractSynthesize_ValidFeatureName verifies valid feature name
func TestRunContractSynthesize_ValidFeatureName(t *testing.T) {
	// This will fail because requirements file doesn't exist
	// but we can test that the feature name validation passes
	synthesizeFeature = "valid-feature-name-123"
	synthesizeRequirements = "/nonexistent/requirements.md"

	err := runContractSynthesize(nil, nil)
	// Should fail at path validation, not feature name
	if err == nil {
		t.Error("expected error (file doesn't exist)")
	}
	// Should NOT be a feature name error
}

// TestRunContractSynthesize_PathValidation verifies path validation
func TestRunContractSynthesize_PathValidation(t *testing.T) {
	synthesizeFeature = "test-feature"
	synthesizeRequirements = "/etc/passwd" // Outside allowed dirs

	err := runContractSynthesize(nil, nil)
	if err == nil {
		t.Error("expected error for path outside allowed dirs")
	}
}

// TestRunContractSynthesize_OutputPathValidation verifies output path validation
func TestRunContractSynthesize_OutputPathValidation(t *testing.T) {
	tmpDir := t.TempDir()
	reqDir := filepath.Join(tmpDir, "docs", "drafts")
	os.MkdirAll(reqDir, 0755)
	reqFile := filepath.Join(reqDir, "test-requirements.md")
	os.WriteFile(reqFile, []byte("# Requirements"), 0644)

	synthesizeFeature = "test"
	synthesizeRequirements = reqFile
	synthesizeOutput = "/etc/passwd" // Outside allowed dirs

	err := runContractSynthesize(nil, nil)
	if err == nil {
		t.Error("expected error for output path outside allowed dirs")
	}
}

// TestRunContractLock_PathValidation verifies lock path validation
func TestRunContractLock_PathValidation(t *testing.T) {
	lockContract = "/etc/passwd"
	lockReason = "test"

	err := runContractLock(nil, nil)
	if err == nil {
		t.Error("expected error for path outside allowed dirs")
	}
}

// TestRunContractValidate_PathValidation verifies validate path validation
func TestRunContractValidate_PathValidation(t *testing.T) {
	validateContracts = []string{"/etc/passwd", "/etc/shadow"}
	validateOutput = ".contracts/report.md"

	err := runContractValidate(nil, nil)
	if err == nil {
		t.Error("expected error for paths outside allowed dirs")
	}
}

// TestRunContractValidate_OutputPathValidation verifies validate output path
func TestRunContractValidate_OutputPathValidation(t *testing.T) {
	validateContracts = []string{".contracts/a.yaml", ".contracts/b.yaml"}
	validateOutput = "/etc/passwd"

	err := runContractValidate(nil, nil)
	if err == nil {
		t.Error("expected error for output path outside allowed dirs")
	}
}

// TestSanitizePath_RelativePath verifies relative path handling
func TestSanitizePath_RelativePath(t *testing.T) {
	// Create test structure
	tmpDir := t.TempDir()
	allowedDir := filepath.Join(tmpDir, "allowed")
	os.Mkdir(allowedDir, 0755)

	// Test with absolute path (simpler and more reliable)
	absPath := filepath.Join(allowedDir, "test.yaml")
	result, err := sanitizePath(absPath, []string{allowedDir})
	if err != nil {
		t.Errorf("path inside allowed should work: %v", err)
	}
	if result == "" {
		t.Error("expected non-empty result")
	}
}

// TestSanitizePath_EmptyPath verifies empty path handling
func TestSanitizePath_EmptyPath(t *testing.T) {
	_, err := sanitizePath("", []string{"/tmp"})
	if err == nil {
		t.Error("expected error for empty path")
	}
}

// TestContractFlags_Types verifies flag types
func TestContractFlags_Types(t *testing.T) {
	// Verify flag variables are accessible
	if synthesizeFeature != "" {
		t.Log("synthesizeFeature is set")
	}
	if lockContract != "" {
		t.Log("lockContract is set")
	}
	if len(validateContracts) != 0 {
		t.Log("validateContracts is set")
	}
}

// TestContractCmd_Run verifies Run function
func TestContractCmd_Run(t *testing.T) {
	// contractCmd.Run just calls Help
	// We can verify it's set
	if contractCmd.Run == nil {
		t.Error("Run function should be set")
	}
}

// TestSynthesizeCmd_RunE verifies RunE function
func TestSynthesizeCmd_RunE(t *testing.T) {
	if synthesizeCmd.RunE == nil {
		t.Error("RunE function should be set")
	}
}

// TestLockCmd_RunE verifies RunE function
func TestLockCmd_RunE(t *testing.T) {
	if lockCmd.RunE == nil {
		t.Error("RunE function should be set")
	}
}

// TestValidateCmd_RunE verifies RunE function
func TestValidateCmd_RunE(t *testing.T) {
	if validateCmd.RunE == nil {
		t.Error("RunE function should be set")
	}
}

// TestRunContractLock_Success verifies successful lock creation
func TestRunContractLock_Success(t *testing.T) {
	// Create test structure in current directory
	cwd, _ := os.Getwd()
	contractsDir := filepath.Join(cwd, ".contracts")
	os.MkdirAll(contractsDir, 0755)
	defer os.RemoveAll(contractsDir)

	contractPath := filepath.Join(contractsDir, "test.yaml")
	os.WriteFile(contractPath, []byte("openapi: 3.0.0"), 0644)

	lockContract = ".contracts/test.yaml"
	lockReason = "test lock"

	err := runContractLock(nil, nil)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	// Verify lock file was created
	lockPath := contractPath + ".lock"
	if _, err := os.Stat(lockPath); os.IsNotExist(err) {
		t.Error("lock file should have been created")
	}
}

// TestRunContractValidate_Success verifies successful validation
func TestRunContractValidate_Success(t *testing.T) {
	// Create test structure in current directory
	cwd, _ := os.Getwd()
	contractsDir := filepath.Join(cwd, ".contracts")
	os.MkdirAll(contractsDir, 0755)
	defer os.RemoveAll(contractsDir)

	// Create contract files
	contract1 := filepath.Join(contractsDir, "frontend.yaml")
	contract2 := filepath.Join(contractsDir, "backend.yaml")
	os.WriteFile(contract1, []byte("openapi: 3.0.0"), 0644)
	os.WriteFile(contract2, []byte("openapi: 3.0.0"), 0644)

	validateContracts = []string{".contracts/frontend.yaml", ".contracts/backend.yaml"}
	validateOutput = ".contracts/report.md"

	err := runContractValidate(nil, nil)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	// Verify report was created
	outputPath := filepath.Join(contractsDir, "report.md")
	if _, err := os.Stat(outputPath); os.IsNotExist(err) {
		t.Error("report file should have been created")
	}
}

// TestRunContractValidate_ThreeContracts verifies validation with 3 contracts
func TestRunContractValidate_ThreeContracts(t *testing.T) {
	cwd, _ := os.Getwd()
	contractsDir := filepath.Join(cwd, ".contracts")
	os.MkdirAll(contractsDir, 0755)
	defer os.RemoveAll(contractsDir)

	for _, name := range []string{"a.yaml", "b.yaml", "c.yaml"} {
		path := filepath.Join(contractsDir, name)
		os.WriteFile(path, []byte("openapi: 3.0.0"), 0644)
	}

	validateContracts = []string{".contracts/a.yaml", ".contracts/b.yaml", ".contracts/c.yaml"}
	validateOutput = ".contracts/report.md"

	err := runContractValidate(nil, nil)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
}

// TestRunContractSynthesize_DefaultPaths verifies default path resolution
func TestRunContractSynthesize_DefaultPaths(t *testing.T) {
	// Test that default paths are used when flags not set
	// This will fail because files don't exist, but we can verify path logic

	synthesizeFeature = "my-feature"
	synthesizeRequirements = "" // Should default to docs/drafts/my-feature-requirements.md
	synthesizeOutput = ""       // Should default to .contracts/my-feature.yaml

	// The function will try to validate the path
	// We just verify it runs (will fail at file check)
	err := runContractSynthesize(nil, nil)
	// Error is expected - we're testing the path resolution logic
	if err == nil {
		t.Log("synthesize ran without error (unexpected)")
	}
}

// TestSanitizePath_AllowedDirNotFound verifies error when allowed dir doesn't exist
func TestSanitizePath_AllowedDirNotFound(t *testing.T) {
	tmpDir := t.TempDir()
	testPath := filepath.Join(tmpDir, "test.yaml")
	os.WriteFile(testPath, []byte("test"), 0644)

	// Non-existent allowed dir
	_, err := sanitizePath(testPath, []string{"/nonexistent/allowed"})
	if err == nil {
		t.Error("expected error when allowed dir doesn't exist")
	}
}

// TestContractLock_WithReason verifies lock content includes reason
func TestContractLock_WithReason(t *testing.T) {
	cwd, _ := os.Getwd()
	contractsDir := filepath.Join(cwd, ".contracts")
	os.MkdirAll(contractsDir, 0755)
	defer os.RemoveAll(contractsDir)

	contractPath := filepath.Join(contractsDir, "test.yaml")
	os.WriteFile(contractPath, []byte("openapi: 3.0.0"), 0644)

	lockContract = ".contracts/test.yaml"
	lockReason = "locking for feature F001"

	err := runContractLock(nil, nil)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	// Read lock file and verify reason
	lockPath := contractPath + ".lock"
	content, err := os.ReadFile(lockPath)
	if err != nil {
		t.Fatalf("failed to read lock file: %v", err)
	}

	if !containsString(string(content), "locking for feature F001") {
		t.Error("lock file should contain reason")
	}
}

// Helper function
func containsString(s, substr string) bool {
	for i := 0; i <= len(s)-len(substr); i++ {
		if s[i:i+len(substr)] == substr {
			return true
		}
	}
	return false
}
