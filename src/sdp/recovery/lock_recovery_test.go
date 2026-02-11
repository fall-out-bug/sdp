package recovery

import (
	"encoding/json"
	"os"
	"path/filepath"
	"testing"
	"time"
)

// TestNewLockRecoveryManager verifies recovery manager creation
func TestNewLockRecoveryManager(t *testing.T) {
	tmpDir := t.TempDir()
	manager := NewLockRecoveryManager(tmpDir, 5)

	if manager.backupDir != tmpDir {
		t.Errorf("Expected backupDir %s, got %s", tmpDir, manager.backupDir)
	}

	if manager.maxBackups != 5 {
		t.Errorf("Expected maxBackups 5, got %d", manager.maxBackups)
	}
}

// TestCreateBackup verifies backup creation
func TestCreateBackup(t *testing.T) {
	tmpDir := t.TempDir()
	manager := NewLockRecoveryManager(tmpDir, 5)

	// Create test lock file
	lock := ContractLock{
		FeatureName:    "telemetry",
		ContractSHA:    "abc123",
		LockedAt:       time.Now(),
		LockedBy:       "test",
		ContractPath:   "/contracts/telemetry.yaml",
		ValidationHash: "hash123",
	}

	lockPath := filepath.Join(tmpDir, "telemetry.lock")
	lockData, _ := json.MarshalIndent(lock, "", "  ")
	os.WriteFile(lockPath, lockData, 0644)

	// Create backup
	backup, err := manager.CreateBackup(lockPath, "test")
	if err != nil {
		t.Fatalf("CreateBackup failed: %v", err)
	}

	if backup.Lock.FeatureName != "telemetry" {
		t.Errorf("Expected feature name 'telemetry', got '%s'", backup.Lock.FeatureName)
	}

	if backup.BackupReason != "test" {
		t.Errorf("Expected backup reason 'test', got '%s'", backup.BackupReason)
	}

	// Verify backup file exists
	if _, err := os.Stat(backup.BackupPath); os.IsNotExist(err) {
		t.Error("Backup file was not created")
	}
}

// TestRestoreFromBackup verifies backup restoration
func TestRestoreFromBackup(t *testing.T) {
	tmpDir := t.TempDir()
	manager := NewLockRecoveryManager(tmpDir, 5)

	// Create test lock file and backup
	lock := ContractLock{
		FeatureName:    "telemetry",
		ContractSHA:    "abc123",
		LockedAt:       time.Now(),
		LockedBy:       "test",
		ContractPath:   "/contracts/telemetry.yaml",
		ValidationHash: "hash123",
	}

	lockPath := filepath.Join(tmpDir, "telemetry.lock")
	lockData, _ := json.MarshalIndent(lock, "", "  ")
	os.WriteFile(lockPath, lockData, 0644)

	backup, _ := manager.CreateBackup(lockPath, "test")

	// Delete original lock
	os.Remove(lockPath)

	// Restore from backup
	result, err := manager.RestoreFromBackup(lockPath, backup.BackupPath)
	if err != nil {
		t.Fatalf("RestoreFromBackup failed: %v", err)
	}

	if !result.Success {
		t.Errorf("Expected successful restore, got: %s", result.ErrorMessage)
	}

	if result.BackupUsed != backup.BackupPath {
		t.Errorf("Expected backup used %s, got %s", backup.BackupPath, result.BackupUsed)
	}

	// Verify restored lock
	restoreData, _ := os.ReadFile(lockPath)
	var restoredLock ContractLock
	json.Unmarshal(restoreData, &restoredLock)

	if restoredLock.FeatureName != "telemetry" {
		t.Errorf("Expected feature name 'telemetry', got '%s'", restoredLock.FeatureName)
	}
}

// TestListBackups verifies backup listing
func TestListBackups(t *testing.T) {
	tmpDir := t.TempDir()
	manager := NewLockRecoveryManager(tmpDir, 5)

	// Create test lock file
	lock := ContractLock{
		FeatureName:    "telemetry",
		ContractSHA:    "abc123",
		LockedAt:       time.Now(),
		LockedBy:       "test",
		ContractPath:   "/contracts/telemetry.yaml",
		ValidationHash: "hash123",
	}

	lockPath := filepath.Join(tmpDir, "telemetry.lock")
	lockData, _ := json.MarshalIndent(lock, "", "  ")
	os.WriteFile(lockPath, lockData, 0644)

	// Create multiple backups with longer delays
	_, _ = manager.CreateBackup(lockPath, "backup1")
	time.Sleep(100 * time.Millisecond)
	_, _ = manager.CreateBackup(lockPath, "backup2")
	time.Sleep(100 * time.Millisecond)
	_, _ = manager.CreateBackup(lockPath, "backup3")

	// List backups
	backups, err := manager.ListBackups("telemetry")
	if err != nil {
		t.Fatalf("ListBackups failed: %v", err)
	}

	if len(backups) < 2 {
		t.Errorf("Expected at least 2 backups, got %d", len(backups))
	}
}

// TestValidateLock_ValidLock verifies validation of valid lock
func TestValidateLock_ValidLock(t *testing.T) {
	tmpDir := t.TempDir()
	manager := NewLockRecoveryManager(tmpDir, 5)

	// Create test lock file
	lock := ContractLock{
		FeatureName:    "telemetry",
		ContractSHA:    "abc123",
		LockedAt:       time.Now(),
		LockedBy:       "test",
		ContractPath:   "", // Skip contract validation
		ValidationHash: "",
	}

	lockPath := filepath.Join(tmpDir, "telemetry.lock")
	lockData, _ := json.MarshalIndent(lock, "", "  ")
	os.WriteFile(lockPath, lockData, 0644)

	// Validate
	err := manager.ValidateLock(lockPath)
	if err != nil {
		t.Errorf("Expected valid lock, got error: %v", err)
	}
}

// TestValidateLock_InvalidJSON verifies validation detects corrupted JSON
func TestValidateLock_InvalidJSON(t *testing.T) {
	tmpDir := t.TempDir()
	manager := NewLockRecoveryManager(tmpDir, 5)

	// Create corrupted lock file
	lockPath := filepath.Join(tmpDir, "telemetry.lock")
	os.WriteFile(lockPath, []byte("{invalid json"), 0644)

	// Validate
	err := manager.ValidateLock(lockPath)
	if err == nil {
		t.Error("Expected validation error for corrupted JSON")
	}

	if !contains(err.Error(), "corrupted lock file") {
		t.Errorf("Expected 'corrupted lock file' error, got: %v", err)
	}
}

// TestValidateLock_MissingFields verifies validation detects missing fields
func TestValidateLock_MissingFields(t *testing.T) {
	tmpDir := t.TempDir()
	manager := NewLockRecoveryManager(tmpDir, 5)

	// Create lock with missing fields
	lock := ContractLock{
		FeatureName: "", // Missing
		ContractSHA: "abc123",
		LockedAt:    time.Now(),
	}

	lockPath := filepath.Join(tmpDir, "telemetry.lock")
	lockData, _ := json.MarshalIndent(lock, "", "  ")
	os.WriteFile(lockPath, lockData, 0644)

	// Validate
	err := manager.ValidateLock(lockPath)
	if err == nil {
		t.Error("Expected validation error for missing fields")
	}

	if !contains(err.Error(), "missing feature_name") {
		t.Errorf("Expected 'missing feature_name' error, got: %v", err)
	}
}

// TestRecoverLock_Success verifies successful recovery from backup
func TestRecoverLock_Success(t *testing.T) {
	tmpDir := t.TempDir()
	manager := NewLockRecoveryManager(tmpDir, 5)

	// Create test lock file and backup
	lock := ContractLock{
		FeatureName:    "telemetry",
		ContractSHA:    "abc123",
		LockedAt:       time.Now(),
		LockedBy:       "test",
		ContractPath:   "",
		ValidationHash: "",
	}

	lockPath := filepath.Join(tmpDir, "telemetry.lock")
	lockData, _ := json.MarshalIndent(lock, "", "  ")
	os.WriteFile(lockPath, lockData, 0644)

	manager.CreateBackup(lockPath, "test")

	// Corrupt lock file
	os.WriteFile(lockPath, []byte("{corrupted"), 0644)

	// Recover
	result, err := manager.RecoverLock(lockPath)
	if err != nil {
		t.Fatalf("RecoverLock failed: %v", err)
	}

	if !result.Success {
		t.Errorf("Expected successful recovery, got: %s", result.ErrorMessage)
	}

	// Verify lock is restored
	restoredData, _ := os.ReadFile(lockPath)
	var restoredLock ContractLock
	json.Unmarshal(restoredData, &restoredLock)

	if restoredLock.FeatureName != "telemetry" {
		t.Errorf("Expected feature name 'telemetry', got '%s'", restoredLock.FeatureName)
	}
}

// TestRecoverLock_NoBackups verifies recovery fails when no backups available
func TestRecoverLock_NoBackups(t *testing.T) {
	tmpDir := t.TempDir()
	manager := NewLockRecoveryManager(tmpDir, 5)

	// Create corrupted lock file without backup
	lockPath := filepath.Join(tmpDir, "telemetry.lock")
	os.WriteFile(lockPath, []byte("{corrupted"), 0644)

	// Attempt recovery
	result, err := manager.RecoverLock(lockPath)
	if err == nil {
		t.Error("Expected recovery error when no backups available")
	}

	if result.Success {
		t.Error("Expected failed recovery result")
	}

	if !contains(result.ErrorMessage, "no backups") {
		t.Errorf("Expected 'no backups' error, got: %s", result.ErrorMessage)
	}
}

// TestCalculateContractHash verifies contract hash calculation
func TestCalculateContractHash(t *testing.T) {
	tmpDir := t.TempDir()

	// Create test contract file
	contractPath := filepath.Join(tmpDir, "contract.yaml")
	contractContent := "openapi: 3.0.0\ninfo:\n  title: Test API"
	os.WriteFile(contractPath, []byte(contractContent), 0644)

	// Calculate hash
	hash, err := CalculateContractHash(contractPath)
	if err != nil {
		t.Fatalf("CalculateContractHash failed: %v", err)
	}

	if len(hash) == 0 {
		t.Error("Expected non-empty hash")
	}

	if len(hash) != 64 { // SHA256 hex encoding is 64 characters
		t.Errorf("Expected hash length 64, got %d", len(hash))
	}

	// Verify hash is consistent
	hash2, _ := CalculateContractHash(contractPath)
	if hash != hash2 {
		t.Error("Hash calculation should be deterministic")
	}
}

// TestCleanOldBackups verifies old backup cleanup
func TestCleanOldBackups(t *testing.T) {
	tmpDir := t.TempDir()
	manager := NewLockRecoveryManager(tmpDir, 2) // Max 2 backups

	// Create test lock file
	lock := ContractLock{
		FeatureName:    "telemetry",
		ContractSHA:    "abc123",
		LockedAt:       time.Now(),
		LockedBy:       "test",
		ContractPath:   "",
		ValidationHash: "",
	}

	lockPath := filepath.Join(tmpDir, "telemetry.lock")
	lockData, _ := json.MarshalIndent(lock, "", "  ")
	os.WriteFile(lockPath, lockData, 0644)

	// Create 3 backups
	manager.CreateBackup(lockPath, "backup1")
	time.Sleep(10 * time.Millisecond)
	manager.CreateBackup(lockPath, "backup2")
	time.Sleep(10 * time.Millisecond)
	manager.CreateBackup(lockPath, "backup3")

	// List backups - should only have 2 (oldest removed)
	backups, _ := manager.ListBackups("telemetry")
	if len(backups) != 2 {
		t.Errorf("Expected 2 backups after cleanup, got %d", len(backups))
	}
}

// TestValidateLock_ContractHashMismatch verifies contract hash validation
func TestValidateLock_ContractHashMismatch(t *testing.T) {
	tmpDir := t.TempDir()
	manager := NewLockRecoveryManager(tmpDir, 5)

	// Create contract file
	contractPath := filepath.Join(tmpDir, "contract.yaml")
	os.WriteFile(contractPath, []byte("original content"), 0644)

	// Calculate hash
	hash, _ := CalculateContractHash(contractPath)

	// Create lock with hash
	lock := ContractLock{
		FeatureName:    "telemetry",
		ContractSHA:    "abc123",
		LockedAt:       time.Now(),
		LockedBy:       "test",
		ContractPath:   contractPath,
		ValidationHash: hash,
	}

	lockPath := filepath.Join(tmpDir, "telemetry.lock")
	lockData, _ := json.MarshalIndent(lock, "", "  ")
	os.WriteFile(lockPath, lockData, 0644)

	// Modify contract file
	os.WriteFile(contractPath, []byte("modified content"), 0644)

	// Validate - should detect hash mismatch
	err := manager.ValidateLock(lockPath)
	if err == nil {
		t.Error("Expected validation error for hash mismatch")
	}

	if !contains(err.Error(), "hash mismatch") {
		t.Errorf("Expected 'hash mismatch' error, got: %v", err)
	}
}

// Helper function
func contains(s, substr string) bool {
	return len(s) >= len(substr) && indexOfStr(s, substr) >= 0
}

func indexOfStr(s, substr string) int {
	for i := 0; i <= len(s)-len(substr); i++ {
		if s[i:i+len(substr)] == substr {
			return i
		}
	}
	return -1
}
