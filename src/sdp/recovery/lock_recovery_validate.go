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

// ValidateLock validates a lock file for corruption
func (rm *LockRecoveryManager) ValidateLock(lockPath string) error {
	lockData, err := os.ReadFile(lockPath)
	if err != nil {
		return fmt.Errorf("failed to read lock file: %w", err)
	}

	var lock ContractLock
	if err := json.Unmarshal(lockData, &lock); err != nil {
		return fmt.Errorf("corrupted lock file (invalid JSON): %w", err)
	}

	if lock.FeatureName == "" {
		return fmt.Errorf("corrupted lock file (missing feature_name)")
	}

	if lock.ContractSHA == "" {
		return fmt.Errorf("corrupted lock file (missing contract_sha)")
	}

	if lock.LockedAt.IsZero() {
		return fmt.Errorf("corrupted lock file (missing locked_at)")
	}

	if lock.ContractPath != "" {
		contractData, err := os.ReadFile(lock.ContractPath)
		if err != nil {
			return fmt.Errorf("contract file not found: %s", lock.ContractPath)
		}

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
	validationErr := rm.ValidateLock(lockPath)

	if validationErr == nil {
		return &RecoveryResult{
			Success:    true,
			LockFile:   lockPath,
			RestoredAt: time.Now(),
		}, nil
	}

	featureName := filepath.Base(lockPath)
	featureName = featureName[:len(featureName)-len(".lock")]

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

	latestBackup := backups[0]
	latestTime := latestBackup.CreatedAt

	for _, backup := range backups {
		if backup.CreatedAt.After(latestTime) {
			latestBackup = backup
			latestTime = backup.CreatedAt
		}
	}

	return rm.RestoreFromBackup(lockPath, latestBackup.BackupPath)
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
