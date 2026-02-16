package recovery

import "time"

// ContractLock represents a contract lock file
type ContractLock struct {
	FeatureName    string    `json:"feature_name"`
	ContractSHA    string    `json:"contract_sha"`
	LockedAt       time.Time `json:"locked_at"`
	LockedBy       string    `json:"locked_by"`
	ContractPath   string    `json:"contract_path"`
	ValidationHash string    `json:"validation_hash"`
}

// LockBackup represents a backup of a lock file
type LockBackup struct {
	Lock         ContractLock `json:"lock"`
	BackupPath   string       `json:"backup_path"`
	CreatedAt    time.Time    `json:"created_at"`
	BackupReason string       `json:"backup_reason"`
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
	backupDir  string
	maxBackups int
}

// NewLockRecoveryManager creates a new lock recovery manager
func NewLockRecoveryManager(backupDir string, maxBackups int) *LockRecoveryManager {
	return &LockRecoveryManager{
		backupDir:  backupDir,
		maxBackups: maxBackups,
	}
}
