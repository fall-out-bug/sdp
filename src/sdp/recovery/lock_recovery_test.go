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

// TestCreateBackup_ReadError verifies error handling for missing file
func TestCreateBackup_ReadError(t *testing.T) {
	tmpDir := t.TempDir()
	manager := NewLockRecoveryManager(tmpDir, 5)

	_, err := manager.CreateBackup("/nonexistent/path/lock.json", "test")
	if err == nil {
		t.Error("expected error for nonexistent file")
	}
	if !contains(err.Error(), "failed to read lock file") {
		t.Errorf("expected read error, got: %v", err)
	}
}

// TestCreateBackup_ParseError verifies error handling for invalid JSON
func TestCreateBackup_ParseError(t *testing.T) {
	tmpDir := t.TempDir()
	manager := NewLockRecoveryManager(tmpDir, 5)

	// Create invalid JSON file
	lockPath := filepath.Join(tmpDir, "invalid.lock")
	os.WriteFile(lockPath, []byte("not valid json"), 0644)

	_, err := manager.CreateBackup(lockPath, "test")
	if err == nil {
		t.Error("expected error for invalid JSON")
	}
	if !contains(err.Error(), "failed to parse lock file") {
		t.Errorf("expected parse error, got: %v", err)
	}
}

// TestRestoreFromBackup_ReadError verifies error handling for missing backup
func TestRestoreFromBackup_ReadError(t *testing.T) {
	tmpDir := t.TempDir()
	manager := NewLockRecoveryManager(tmpDir, 5)

	result, err := manager.RestoreFromBackup("/tmp/lock.json", "/nonexistent/backup.json")
	if err == nil {
		t.Error("expected error for nonexistent backup")
	}
	if result.Success {
		t.Error("expected failure result")
	}
	if !contains(result.ErrorMessage, "failed to read backup") {
		t.Errorf("expected read error, got: %s", result.ErrorMessage)
	}
}

// TestRestoreFromBackup_ParseError verifies error handling for corrupt backup
func TestRestoreFromBackup_ParseError(t *testing.T) {
	tmpDir := t.TempDir()
	manager := NewLockRecoveryManager(tmpDir, 5)

	// Create corrupt backup
	backupPath := filepath.Join(tmpDir, "corrupt.backup")
	os.WriteFile(backupPath, []byte("not valid json"), 0644)

	result, err := manager.RestoreFromBackup(filepath.Join(tmpDir, "lock.json"), backupPath)
	if err == nil {
		t.Error("expected error for corrupt backup")
	}
	if result.Success {
		t.Error("expected failure result")
	}
	if !contains(result.ErrorMessage, "failed to parse backup") {
		t.Errorf("expected parse error, got: %s", result.ErrorMessage)
	}
}

// TestListBackups_EmptyDirectory verifies empty backup directory
func TestListBackups_EmptyDirectory(t *testing.T) {
	tmpDir := t.TempDir()
	manager := NewLockRecoveryManager(tmpDir, 5)

	backups, err := manager.ListBackups("telemetry")
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if len(backups) != 0 {
		t.Errorf("expected 0 backups, got %d", len(backups))
	}
}

// TestListBackups_NonexistentDirectory verifies nonexistent backup directory
func TestListBackups_NonexistentDirectory(t *testing.T) {
	manager := NewLockRecoveryManager("/nonexistent/backup/dir", 5)

	backups, err := manager.ListBackups("telemetry")
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if len(backups) != 0 {
		t.Errorf("expected 0 backups for nonexistent dir, got %d", len(backups))
	}
}

// TestListBackups_SkipsDirectories verifies directories are skipped
func TestListBackups_SkipsDirectories(t *testing.T) {
	tmpDir := t.TempDir()
	manager := NewLockRecoveryManager(tmpDir, 5)

	// Create a directory that matches the prefix
	os.Mkdir(filepath.Join(tmpDir, "telemetry-subdir"), 0755)

	backups, err := manager.ListBackups("telemetry")
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if len(backups) != 0 {
		t.Errorf("expected 0 backups (directories skipped), got %d", len(backups))
	}
}

// TestListBackups_SkipsUnreadableFiles verifies unreadable files are skipped
func TestListBackups_SkipsUnreadableFiles(t *testing.T) {
	tmpDir := t.TempDir()
	manager := NewLockRecoveryManager(tmpDir, 5)

	// Create unreadable backup (corrupt JSON)
	backupPath := filepath.Join(tmpDir, "telemetry-123.backup")
	os.WriteFile(backupPath, []byte("corrupt"), 0644)

	backups, err := manager.ListBackups("telemetry")
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	// Corrupt files should be skipped, not cause error
	if len(backups) != 0 {
		t.Errorf("expected 0 backups (corrupt skipped), got %d", len(backups))
	}
}

// TestCleanOldBackups_SortingOrder verifies cleanup happens when over limit
func TestCleanOldBackups_SortingOrder(t *testing.T) {
	tmpDir := t.TempDir()
	manager := NewLockRecoveryManager(tmpDir, 2) // Keep only 2

	// Create test lock
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

	// Create backups with delays to ensure different timestamps
	manager.CreateBackup(lockPath, "backup1")
	time.Sleep(100 * time.Millisecond)
	manager.CreateBackup(lockPath, "backup2")
	time.Sleep(100 * time.Millisecond)
	manager.CreateBackup(lockPath, "backup3")

	// Check that we have exactly 2 backups after cleanup
	backups, _ := manager.ListBackups("telemetry")
	if len(backups) != 2 {
		t.Errorf("expected 2 backups after cleanup, got %d", len(backups))
	}
}

// TestCleanOldBackups_NoCleanupNeeded verifies no cleanup when under limit
func TestCleanOldBackups_NoCleanupNeeded(t *testing.T) {
	tmpDir := t.TempDir()
	manager := NewLockRecoveryManager(tmpDir, 5)

	lock := ContractLock{
		FeatureName:    "telemetry",
		ContractSHA:    "abc123",
		LockedAt:       time.Now(),
		LockedBy:       "test",
	}
	lockPath := filepath.Join(tmpDir, "telemetry.lock")
	lockData, _ := json.MarshalIndent(lock, "", "  ")
	os.WriteFile(lockPath, lockData, 0644)

	backup, _ := manager.CreateBackup(lockPath, "backup1")

	// Should still exist (under limit)
	_, err := os.Stat(backup.BackupPath)
	if os.IsNotExist(err) {
		t.Error("backup should still exist when under limit")
	}
}

// TestCleanOldBackups_RemovesMultiple verifies multiple old backups removed
func TestCleanOldBackups_RemovesMultiple(t *testing.T) {
	tmpDir := t.TempDir()
	manager := NewLockRecoveryManager(tmpDir, 2) // Keep only 2

	lock := ContractLock{
		FeatureName:    "telemetry",
		ContractSHA:    "abc123",
		LockedAt:       time.Now(),
		LockedBy:       "test",
	}
	lockPath := filepath.Join(tmpDir, "telemetry.lock")
	lockData, _ := json.MarshalIndent(lock, "", "  ")
	os.WriteFile(lockPath, lockData, 0644)

	// Create 5 backups
	for i := 0; i < 5; i++ {
		manager.CreateBackup(lockPath, "backup")
		time.Sleep(50 * time.Millisecond)
	}

	// Should have only 2 after cleanup
	backups, _ := manager.ListBackups("telemetry")
	if len(backups) != 2 {
		t.Errorf("expected 2 backups after cleanup, got %d", len(backups))
	}
}

// TestRecoverLock_WithBackups verifies recovery with available backups
func TestRecoverLock_WithBackups(t *testing.T) {
	tmpDir := t.TempDir()
	manager := NewLockRecoveryManager(tmpDir, 5)

	// Create lock and backup
	lock := ContractLock{
		FeatureName:    "telemetry",
		ContractSHA:    "original",
		LockedAt:       time.Now(),
		LockedBy:       "test",
		ContractPath:   "",
		ValidationHash: "",
	}
	lockPath := filepath.Join(tmpDir, "telemetry.lock")
	lockData, _ := json.MarshalIndent(lock, "", "  ")
	os.WriteFile(lockPath, lockData, 0644)

	manager.CreateBackup(lockPath, "pre-corruption")

	// Corrupt the lock
	os.WriteFile(lockPath, []byte("{corrupted"), 0644)

	// Recover
	result, err := manager.RecoverLock(lockPath)
	if err != nil {
		t.Fatalf("RecoverLock failed: %v", err)
	}

	if !result.Success {
		t.Errorf("expected success, got: %s", result.ErrorMessage)
	}
}

// TestRecoverLock_AllBackupsCorrupt verifies handling when all backups corrupt
func TestRecoverLock_AllBackupsCorrupt(t *testing.T) {
	tmpDir := t.TempDir()
	manager := NewLockRecoveryManager(tmpDir, 5)

	// Create valid lock
	lock := ContractLock{
		FeatureName:    "telemetry",
		ContractSHA:    "abc123",
		LockedAt:       time.Now(),
		LockedBy:       "test",
	}
	lockPath := filepath.Join(tmpDir, "telemetry.lock")
	lockData, _ := json.MarshalIndent(lock, "", "  ")
	os.WriteFile(lockPath, lockData, 0644)

	manager.CreateBackup(lockPath, "test")

	// Corrupt both lock and backup
	os.WriteFile(lockPath, []byte("{corrupted"), 0644)
	// Find and corrupt backup
	entries, _ := os.ReadDir(tmpDir)
	for _, e := range entries {
		if filepath.Ext(e.Name()) == ".backup" {
			os.WriteFile(filepath.Join(tmpDir, e.Name()), []byte("{corrupted"), 0644)
		}
	}

	// Recovery should fail
	result, err := manager.RecoverLock(lockPath)
	if err == nil {
		t.Error("expected error when all backups corrupt")
	}
	if result.Success {
		t.Error("expected failure when all backups corrupt")
	}
}
