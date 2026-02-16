package recovery

import (
	"encoding/json"
	"fmt"
	"os"
	"path/filepath"
	"time"
)

// CreateBackup creates a backup of a lock file
func (rm *LockRecoveryManager) CreateBackup(lockPath string, reason string) (*LockBackup, error) {
	lockData, err := os.ReadFile(lockPath)
	if err != nil {
		return nil, fmt.Errorf("failed to read lock file: %w", err)
	}

	var lock ContractLock
	if err := json.Unmarshal(lockData, &lock); err != nil {
		return nil, fmt.Errorf("failed to parse lock file: %w", err)
	}

	if err := os.MkdirAll(rm.backupDir, 0755); err != nil {
		return nil, fmt.Errorf("failed to create backup directory: %w", err)
	}

	backupName := fmt.Sprintf("%s-%d.backup", filepath.Base(lockPath), time.Now().Unix())
	backupPath := filepath.Join(rm.backupDir, backupName)

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

	if err := rm.cleanOldBackups(lock.FeatureName); err != nil {
		fmt.Printf("Warning: failed to clean old backups: %v\n", err)
	}

	return backup, nil
}

// RestoreFromBackup restores a lock file from backup
func (rm *LockRecoveryManager) RestoreFromBackup(lockPath string, backupPath string) (*RecoveryResult, error) {
	backupData, err := os.ReadFile(backupPath)
	if err != nil {
		return &RecoveryResult{
			Success:      false,
			LockFile:     lockPath,
			ErrorMessage: fmt.Sprintf("failed to read backup: %v", err),
		}, fmt.Errorf("failed to read backup: %w", err)
	}

	var backup LockBackup
	if err := json.Unmarshal(backupData, &backup); err != nil {
		return &RecoveryResult{
			Success:      false,
			LockFile:     lockPath,
			BackupUsed:   backupPath,
			ErrorMessage: fmt.Sprintf("failed to parse backup: %v", err),
		}, fmt.Errorf("failed to parse backup: %w", err)
	}

	lockData, err := json.MarshalIndent(backup.Lock, "", "  ")
	if err != nil {
		return &RecoveryResult{
			Success:      false,
			LockFile:     lockPath,
			BackupUsed:   backupPath,
			ErrorMessage: fmt.Sprintf("failed to marshal lock: %v", err),
		}, fmt.Errorf("failed to marshal lock: %w", err)
	}

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

	entries, err := os.ReadDir(rm.backupDir)
	if err != nil {
		if os.IsNotExist(err) {
			return backups, nil
		}
		return nil, fmt.Errorf("failed to read backup directory: %w", err)
	}

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
			continue
		}

		var backup LockBackup
		if err := json.Unmarshal(backupData, &backup); err != nil {
			continue
		}

		backups = append(backups, &backup)
	}

	return backups, nil
}

// cleanOldBackups removes old backups beyond maxBackups limit
func (rm *LockRecoveryManager) cleanOldBackups(featureName string) error {
	backups, err := rm.ListBackups(featureName)
	if err != nil {
		return err
	}

	if len(backups) <= rm.maxBackups {
		return nil
	}

	sortedBackups := make([]*LockBackup, len(backups))
	copy(sortedBackups, backups)

	for i := 0; i < len(sortedBackups); i++ {
		for j := i + 1; j < len(sortedBackups); j++ {
			if sortedBackups[i].CreatedAt.Before(sortedBackups[j].CreatedAt) {
				sortedBackups[i], sortedBackups[j] = sortedBackups[j], sortedBackups[i]
			}
		}
	}

	for i := 0; i < len(sortedBackups)-rm.maxBackups; i++ {
		if err := os.Remove(sortedBackups[i].BackupPath); err != nil {
			return fmt.Errorf("failed to remove old backup %s: %w", sortedBackups[i].BackupPath, err)
		}
	}

	return nil
}
