package main

import (
	"crypto/sha256"
	"encoding/hex"
	"fmt"
	"os"
	"path/filepath"
	"strings"
	"testing"
	"time"
)

// Helper function to calculate SHA256 hash
func calculateSHA256(content string) string {
	hash := sha256.Sum256([]byte(content))
	return hex.EncodeToString(hash[:])
}

// TestLockCmd_CreateLock tests successful lock creation
func TestLockCmd_CreateLock(t *testing.T) {
	// Create temp directory
	tmpDir := t.TempDir()

	// Create test contract
	contractPath := filepath.Join(tmpDir, "test.yaml")
	contractContent := "openapi: 3.0.0\ninfo:\n  title: Test API\n"
	if err := os.WriteFile(contractPath, []byte(contractContent), 0644); err != nil {
		t.Fatalf("Failed to create test contract: %v", err)
	}

	// Create lock file path
	lockPath := filepath.Join(tmpDir, "test.lock")

	// Run lock command
	featureName := "test"
	gitSHA := "abc123def456"

	err := runContractLockInternal(featureName, gitSHA, contractPath, lockPath, false)
	if err != nil {
		t.Fatalf("runContractLockInternal failed: %v", err)
	}

	// Verify lock file exists
	if _, err := os.Stat(lockPath); os.IsNotExist(err) {
		t.Fatalf("Lock file was not created: %s", lockPath)
	}

	// Read and verify lock file content
	lockContent, err := os.ReadFile(lockPath)
	if err != nil {
		t.Fatalf("Failed to read lock file: %v", err)
	}

	lockStr := string(lockContent)

	// Verify contract hash is present
	expectedHash := calculateSHA256(contractContent)
	if !strings.Contains(lockStr, expectedHash) {
		t.Errorf("Lock file does not contain expected hash. Got: %s, Want: %s", lockStr, expectedHash)
	}

	// Verify git SHA is present
	if !strings.Contains(lockStr, gitSHA) {
		t.Errorf("Lock file does not contain git SHA. Got: %s, Want: %s", lockStr, gitSHA)
	}

	// Verify timestamp is present (format: 2026-02-07T12:34:56Z)
	if !strings.Contains(lockStr, "locked_at:") {
		t.Errorf("Lock file does not contain timestamp. Got: %s", lockStr)
	}
}

// TestLockCmd_ContractNotFound tests error when contract doesn't exist
func TestLockCmd_ContractNotFound(t *testing.T) {
	tmpDir := t.TempDir()
	contractPath := filepath.Join(tmpDir, "nonexistent.yaml")
	lockPath := filepath.Join(tmpDir, "test.lock")

	err := runContractLockInternal("test", "abc123", contractPath, lockPath, false)
	if err == nil {
		t.Fatal("Expected error for non-existent contract, got nil")
	}

	if !strings.Contains(err.Error(), "not found") && !strings.Contains(err.Error(), "no such file") {
		t.Errorf("Expected 'not found' error, got: %v", err)
	}
}

// TestLockCmd_LockAlreadyExists tests error when lock already exists
func TestLockCmd_LockAlreadyExists(t *testing.T) {
	tmpDir := t.TempDir()

	// Create test contract
	contractPath := filepath.Join(tmpDir, "test.yaml")
	contractContent := "openapi: 3.0.0\ninfo:\n  title: Test API\n"
	if err := os.WriteFile(contractPath, []byte(contractContent), 0644); err != nil {
		t.Fatalf("Failed to create test contract: %v", err)
	}

	// Create existing lock file
	lockPath := filepath.Join(tmpDir, "test.lock")
	if err := os.WriteFile(lockPath, []byte("existing lock"), 0644); err != nil {
		t.Fatalf("Failed to create existing lock: %v", err)
	}

	// Try to lock again without --force
	err := runContractLockInternal("test", "abc123", contractPath, lockPath, false)
	if err == nil {
		t.Fatal("Expected error for existing lock, got nil")
	}

	if !strings.Contains(err.Error(), "already exists") {
		t.Errorf("Expected 'already exists' error, got: %v", err)
	}

	// Verify lock file was not modified
	content, _ := os.ReadFile(lockPath)
	if string(content) != "existing lock" {
		t.Error("Lock file was modified when it should have been preserved")
	}
}

// TestLockCmd_ForceRelock tests successful re-lock with --force
func TestLockCmd_ForceRelock(t *testing.T) {
	tmpDir := t.TempDir()

	// Create test contract
	contractPath := filepath.Join(tmpDir, "test.yaml")
	contractContent := "openapi: 3.0.0\ninfo:\n  title: Test API\n"
	if err := os.WriteFile(contractPath, []byte(contractContent), 0644); err != nil {
		t.Fatalf("Failed to create test contract: %v", err)
	}

	// Create existing lock file
	lockPath := filepath.Join(tmpDir, "test.lock")
	if err := os.WriteFile(lockPath, []byte("old lock"), 0644); err != nil {
		t.Fatalf("Failed to create existing lock: %v", err)
	}

	// Re-lock with --force
	err := runContractLockInternal("test", "newsha123", contractPath, lockPath, true)
	if err != nil {
		t.Fatalf("Force re-lock failed: %v", err)
	}

	// Verify lock file was updated
	content, _ := os.ReadFile(lockPath)
	if string(content) == "old lock" {
		t.Error("Lock file was not updated with --force")
	}

	// Verify new git SHA is present
	lockStr := string(content)
	if !strings.Contains(lockStr, "newsha123") {
		t.Errorf("Lock file does not contain new git SHA. Got: %s", lockStr)
	}
}

// TestVerifyCmd_Match tests successful verification
func TestVerifyCmd_Match(t *testing.T) {
	tmpDir := t.TempDir()

	// Create test contract
	contractPath := filepath.Join(tmpDir, "test.yaml")
	contractContent := "openapi: 3.0.0\ninfo:\n  title: Test API\n"
	if err := os.WriteFile(contractPath, []byte(contractContent), 0644); err != nil {
		t.Fatalf("Failed to create test contract: %v", err)
	}

	// Create lock file
	lockPath := filepath.Join(tmpDir, "test.lock")
	gitSHA := "abc123def456"
	expectedHash := calculateSHA256(contractContent)

	lockContent := fmt.Sprintf("contract_file: %s\ncontract_hash: %s\ngit_sha: %s\nlocked_at: \"%s\"\n",
		contractPath, expectedHash, gitSHA, time.Now().UTC().Format(time.RFC3339))
	if err := os.WriteFile(lockPath, []byte(lockContent), 0644); err != nil {
		t.Fatalf("Failed to create lock file: %v", err)
	}

	// Run verify
	matched, err := runContractVerifyInternal(contractPath, lockPath)
	if err != nil {
		t.Fatalf("Verify failed: %v", err)
	}

	if !matched {
		t.Error("Expected verification to succeed, got false")
	}
}

// TestVerifyCmd_Mismatch tests detection of contract modification
func TestVerifyCmd_Mismatch(t *testing.T) {
	tmpDir := t.TempDir()

	// Create test contract (modified version)
	contractPath := filepath.Join(tmpDir, "test.yaml")
	modifiedContent := "openapi: 3.0.0\ninfo:\n  title: MODIFIED API\n"
	if err := os.WriteFile(contractPath, []byte(modifiedContent), 0644); err != nil {
		t.Fatalf("Failed to create test contract: %v", err)
	}

	// Create lock file with original hash
	lockPath := filepath.Join(tmpDir, "test.lock")
	originalHash := calculateSHA256("openapi: 3.0.0\ninfo:\n  title: Test API\n")

	lockContent := fmt.Sprintf("contract_file: %s\ncontract_hash: %s\ngit_sha: abc123\nlocked_at: \"%s\"\n",
		contractPath, originalHash, time.Now().UTC().Format(time.RFC3339))
	if err := os.WriteFile(lockPath, []byte(lockContent), 0644); err != nil {
		t.Fatalf("Failed to create lock file: %v", err)
	}

	// Run verify
	matched, err := runContractVerifyInternal(contractPath, lockPath)
	if err != nil {
		t.Fatalf("Verify failed: %v", err)
	}

	if matched {
		t.Error("Expected verification to fail for modified contract, got true")
	}
}

// TestVerifyCmd_LockNotFound tests error when lock doesn't exist
func TestVerifyCmd_LockNotFound(t *testing.T) {
	tmpDir := t.TempDir()
	contractPath := filepath.Join(tmpDir, "test.yaml")
	lockPath := filepath.Join(tmpDir, "nonexistent.lock")

	_, err := runContractVerifyInternal(contractPath, lockPath)
	if err == nil {
		t.Fatal("Expected error for non-existent lock, got nil")
	}

	if !strings.Contains(err.Error(), "not found") && !strings.Contains(err.Error(), "no such file") {
		t.Errorf("Expected 'not found' error, got: %v", err)
	}
}
