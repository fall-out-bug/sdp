package doctor

import (
	"fmt"
	"os"
	"time"
)

// Migration represents a config migration operation
type Migration struct {
	FromVersion int
	ToVersion   int
	BackupPath  string
	Timestamp   time.Time
	Success     bool
	Message     string
}

// MigrationRegistry holds all known migrations
var MigrationRegistry = map[int]func(string) error{
	// Version 0 to 1: Add default fields
	1: migrateV0ToV1,
}

// MigrateConfig migrates config to the latest version
func MigrateConfig(dryRun bool) (*Migration, error) {
	configPath := ".sdp/config.yml"

	// Check if config exists
	if _, err := os.Stat(configPath); os.IsNotExist(err) {
		return nil, fmt.Errorf("no config file to migrate")
	}

	// Read current version (simplified - just check for version: field)
	content, err := os.ReadFile(configPath)
	if err != nil {
		return nil, fmt.Errorf("cannot read config: %w", err)
	}

	currentVersion := detectConfigVersion(string(content))
	latestVersion := getLatestVersion()

	if currentVersion >= latestVersion {
		return &Migration{
			FromVersion: currentVersion,
			ToVersion:   currentVersion,
			Timestamp:   time.Now(),
			Success:     true,
			Message:     "Config already at latest version",
		}, nil
	}

	migration := &Migration{
		FromVersion: currentVersion,
		ToVersion:   latestVersion,
		Timestamp:   time.Now(),
	}

	// Create backup
	backupPath, err := createConfigBackup(configPath)
	if err != nil {
		migration.Message = fmt.Sprintf("Failed to create backup: %v", err)
		return migration, err
	}
	migration.BackupPath = backupPath

	if dryRun {
		migration.Message = fmt.Sprintf("Would migrate from v%d to v%d (dry run)", currentVersion, latestVersion)
		return migration, nil
	}

	// Run migrations
	for v := currentVersion + 1; v <= latestVersion; v++ {
		if migrateFn, ok := MigrationRegistry[v]; ok {
			if err := migrateFn(configPath); err != nil {
				migration.Message = fmt.Sprintf("Migration v%d failed: %v", v, err)
				return migration, err
			}
		}
	}

	migration.Success = true
	migration.Message = fmt.Sprintf("Migrated from v%d to v%d", currentVersion, latestVersion)

	// Log migration
	logMigration(migration)

	return migration, nil
}

// RollbackMigration restores config from backup
func RollbackMigration(backupPath string) error {
	configPath := ".sdp/config.yml"

	// Check backup exists
	if _, err := os.Stat(backupPath); os.IsNotExist(err) {
		return fmt.Errorf("backup file not found: %s", backupPath)
	}

	// Read backup
	backupContent, err := os.ReadFile(backupPath)
	if err != nil {
		return fmt.Errorf("cannot read backup: %w", err)
	}

	// Write to config
	if err := os.WriteFile(configPath, backupContent, 0644); err != nil {
		return fmt.Errorf("cannot restore config: %w", err)
	}

	return nil
}
