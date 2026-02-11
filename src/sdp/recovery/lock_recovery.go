package recovery

import (
	"crypto/sha256"
	"encoding/hex"
	"encoding/json"
	"fmt"
	"os"
	"path/filepath"
	"time"
)

// ContractLock represents a contract lock file
type ContractLock struct {
	FeatureName    string    `json:"feature_name"`
	ContractSHA    string    `json:"contract_sha"`
	LockedAt       time.Time `json:"locked_at"`
	LockedBy       string    `json:"locked_by"`
	ContractPath   string    `json:"contract_path"`
	ValidationHash string    `json:"validation_hash"` // Hash of contract content
}

// LockBackup represents a backup of a lock file
type LockBackup struct {
	Lock         ContractLock `json:"lock"`
	BackupPath   string       `json:"backup_path"`
	CreatedAt    time.Time    `json:"created_at"`
	BackupReason string       `json:"backup_reason"` // "pre_modify", "auto", "manual"
}

// RecoveryResult represents the result of a recovery operation
type RecoveryResult struct {
	Success      bool      `json:"success"`
	LockFile     string    `json:"lock_file"`
	BackupUsed   string    `json:"backup_used,omitempty"`
	RestoredAt   time.Time `json:"restored_at"`
	ErrorMessage string    `json:"error_message,omitempty"`
}

// LockRecoveryManager manages contract lock disaster recovery
type LockRecoveryManager struct {
	backupDir string
	maxBackups int
}

// NewLockRecoveryManager creates a new lock recovery manager
func NewLockRecoveryManager(backupDir string, maxBackups int) *LockRecoveryManager {
	return &LockRecoveryManager{
		backupDir:  backupDir,
		maxBackups: maxBackups,
	}
}

// CreateBackup creates a backup of a lock file
func (rm *LockRecoveryManager) CreateBackup(lockPath string, reason string) (*LockBackup, error) {
	// Read lock file
	lockData, err := os.ReadFile(lockPath)
	if err != nil {
		return nil, fmt.Errorf("failed to read lock file: %w", err)
	}

	// Parse lock
	var lock ContractLock
	if err := json.Unmarshal(lockData, &lock); err != nil {
		return nil, fmt.Errorf("failed to parse lock file: %w", err)
	}

	// Create backup directory if needed
	if err := os.MkdirAll(rm.backupDir, 0755); err != nil {
		return nil, fmt.Errorf("failed to create backup directory: %w", err)
	}

	// Generate backup filename
	backupName := fmt.Sprintf("%s-%d.backup", filepath.Base(lockPath), time.Now().Unix())
	backupPath := filepath.Join(rm.backupDir, backupName)

	// Create backup
	backup := &LockBackup{
		Lock:         lock,
		BackupPath:   backupPath,
		CreatedAt:    time.Now(),
		BackupReason: reason,
	}

	backupData, err := json.MarshalIndent(backup, "", "  ")
	if err != nil {
		return nil, fmt.Errorf("failed to marshal backup: %w", err)
	}

	if err := os.WriteFile(backupPath, backupData, 0644); err != nil {
		return nil, fmt.Errorf("failed to write backup: %w", err)
	}

	// Clean old backups
	if err := rm.cleanOldBackups(lock.FeatureName); err != nil {
		// Log but don't fail - cleanup is secondary
		fmt.Printf("Warning: failed to clean old backups: %v\n", err)
	}

	return backup, nil
}

// RestoreFromBackup restores a lock file from backup
func (rm *LockRecoveryManager) RestoreFromBackup(lockPath string, backupPath string) (*RecoveryResult, error) {
	// Read backup
	backupData, err := os.ReadFile(backupPath)
	if err != nil {
		return &RecoveryResult{
			Success:      false,
			LockFile:     lockPath,
			ErrorMessage: fmt.Sprintf("failed to read backup: %v", err),
		}, fmt.Errorf("failed to read backup: %w", err)
	}

	// Parse backup
	var backup LockBackup
	if err := json.Unmarshal(backupData, &backup); err != nil {
		return &RecoveryResult{
			Success:      false,
			LockFile:     lockPath,
			BackupUsed:   backupPath,
			ErrorMessage: fmt.Sprintf("failed to parse backup: %v", err),
		}, fmt.Errorf("failed to parse backup: %w", err)
	}

	// Marshal lock to JSON
	lockData, err := json.MarshalIndent(backup.Lock, "", "  ")
	if err != nil {
		return &RecoveryResult{
			Success:      false,
			LockFile:     lockPath,
			BackupUsed:   backupPath,
			ErrorMessage: fmt.Sprintf("failed to marshal lock: %v", err),
		}, fmt.Errorf("failed to marshal lock: %w", err)
	}

	// Restore lock file
	if err := os.WriteFile(lockPath, lockData, 0644); err != nil {
		return &RecoveryResult{
			Success:      false,
			LockFile:     lockPath,
			BackupUsed:   backupPath,
			ErrorMessage: fmt.Sprintf("failed to write lock file: %v", err),
		}, fmt.Errorf("failed to write lock file: %w", err)
	}

	return &RecoveryResult{
		Success:    true,
		LockFile:   lockPath,
		BackupUsed: backupPath,
		RestoredAt: time.Now(),
	}, nil
}

// ListBackups lists all backups for a feature
func (rm *LockRecoveryManager) ListBackups(featureName string) ([]*LockBackup, error) {
	var backups []*LockBackup

	// Read backup directory
	entries, err := os.ReadDir(rm.backupDir)
	if err != nil {
		if os.IsNotExist(err) {
			return backups, nil // No backups yet
		}
		return nil, fmt.Errorf("failed to read backup directory: %w", err)
	}

	// Find backups for this feature
	for _, entry := range entries {
		if entry.IsDir() {
			continue
		}

		if !filepath.HasPrefix(entry.Name(), featureName) {
			continue
		}

		backupPath := filepath.Join(rm.backupDir, entry.Name())
		backupData, err := os.ReadFile(backupPath)
		if err != nil {
			continue // Skip corrupted backups
		}

		var backup LockBackup
		if err := json.Unmarshal(backupData, &backup); err != nil {
			continue // Skip corrupted backups
		}

		backups = append(backups, &backup)
	}

	return backups, nil
}

// ValidateLock validates a lock file for corruption
func (rm *LockRecoveryManager) ValidateLock(lockPath string) error {
	// Read lock file
	lockData, err := os.ReadFile(lockPath)
	if err != nil {
		return fmt.Errorf("failed to read lock file: %w", err)
	}

	// Parse lock
	var lock ContractLock
	if err := json.Unmarshal(lockData, &lock); err != nil {
		return fmt.Errorf("corrupted lock file (invalid JSON): %w", err)
	}

	// Validate required fields
	if lock.FeatureName == "" {
		return fmt.Errorf("corrupted lock file (missing feature_name)")
	}

	if lock.ContractSHA == "" {
		return fmt.Errorf("corrupted lock file (missing contract_sha)")
	}

	if lock.LockedAt.IsZero() {
		return fmt.Errorf("corrupted lock file (missing locked_at)")
	}

	// Validate contract file exists and matches hash
	if lock.ContractPath != "" {
		contractData, err := os.ReadFile(lock.ContractPath)
		if err != nil {
			return fmt.Errorf("contract file not found: %s", lock.ContractPath)
		}

		// Calculate hash
		hash := sha256.Sum256(contractData)
		calculatedHash := hex.EncodeToString(hash[:])

		if lock.ValidationHash != "" && calculatedHash != lock.ValidationHash {
			return fmt.Errorf("contract hash mismatch (expected %s, got %s)", lock.ValidationHash, calculatedHash)
		}
	}

	return nil
}

// RecoverLock attempts to recover a corrupted or missing lock file
func (rm *LockRecoveryManager) RecoverLock(lockPath string) (*RecoveryResult, error) {
	// Try to validate existing lock
	validationErr := rm.ValidateLock(lockPath)

	// If lock is valid, no recovery needed
	if validationErr == nil {
		return &RecoveryResult{
			Success:  true,
			LockFile: lockPath,
			RestoredAt: time.Now(),
		}, nil
	}

	// Lock is corrupted or missing - find backups
	// Extract feature name from lock path
	featureName := filepath.Base(lockPath)
	featureName = featureName[:len(featureName)-len(".lock")] // Remove .lock suffix

	backups, err := rm.ListBackups(featureName)
	if err != nil {
		return &RecoveryResult{
			Success:      false,
			LockFile:     lockPath,
			ErrorMessage: fmt.Sprintf("failed to list backups: %v", err),
		}, fmt.Errorf("failed to list backups: %w", err)
	}

	if len(backups) == 0 {
		return &RecoveryResult{
			Success:      false,
			LockFile:     lockPath,
			ErrorMessage: "no backups available for recovery",
		}, fmt.Errorf("no backups available for recovery")
	}

	// Use most recent backup
	latestBackup := backups[0]
	latestTime := latestBackup.CreatedAt

	for _, backup := range backups {
		if backup.CreatedAt.After(latestTime) {
			latestBackup = backup
			latestTime = backup.CreatedAt
		}
	}

	// Restore from latest backup
	return rm.RestoreFromBackup(lockPath, latestBackup.BackupPath)
}

// cleanOldBackups removes old backups beyond maxBackups limit
func (rm *LockRecoveryManager) cleanOldBackups(featureName string) error {
	backups, err := rm.ListBackups(featureName)
	if err != nil {
		return err
	}

	if len(backups) <= rm.maxBackups {
		return nil // Within limit
	}

	// Sort by creation time (oldest first)
	sortedBackups := make([]*LockBackup, len(backups))
	copy(sortedBackups, backups)

	for i := 0; i < len(sortedBackups); i++ {
		for j := i + 1; j < len(sortedBackups); j++ {
			if sortedBackups[i].CreatedAt.Before(sortedBackups[j].CreatedAt) {
				sortedBackups[i], sortedBackups[j] = sortedBackups[j], sortedBackups[i]
			}
		}
	}

	// Remove oldest backups beyond limit
	for i := 0; i < len(sortedBackups)-rm.maxBackups; i++ {
		if err := os.Remove(sortedBackups[i].BackupPath); err != nil {
			return fmt.Errorf("failed to remove old backup %s: %w", sortedBackups[i].BackupPath, err)
		}
	}

	return nil
}

// CalculateContractHash calculates SHA256 hash of contract content
func CalculateContractHash(contractPath string) (string, error) {
	contractData, err := os.ReadFile(contractPath)
	if err != nil {
		return "", fmt.Errorf("failed to read contract: %w", err)
	}

	hash := sha256.Sum256(contractData)
	return hex.EncodeToString(hash[:]), nil
}
