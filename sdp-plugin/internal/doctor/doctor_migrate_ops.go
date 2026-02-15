package doctor

import (
	"fmt"
	"os"
	"path/filepath"
	"time"
)

// createConfigBackup creates a timestamped backup of config
func createConfigBackup(configPath string) (string, error) {
	content, err := os.ReadFile(configPath)
	if err != nil {
		return "", err
	}

	backupDir := ".sdp/backups"
	if err := os.MkdirAll(backupDir, 0755); err != nil {
		return "", err
	}

	timestamp := time.Now().Format("20060102-150405")
	backupPath := filepath.Join(backupDir, fmt.Sprintf("config-%s.yml", timestamp))

	if err := os.WriteFile(backupPath, content, 0644); err != nil {
		return "", err
	}

	return backupPath, nil
}

// logMigration records migration to history file
func logMigration(m *Migration) {
	logPath := ".sdp/migrations.log"
	f, err := os.OpenFile(logPath, os.O_APPEND|os.O_CREATE|os.O_WRONLY, 0644)
	if err != nil {
		return
	}
	defer f.Close()

	fmt.Fprintf(f, "%s: v%d -> v%d, success=%v, backup=%s\n",
		m.Timestamp.Format(time.RFC3339),
		m.FromVersion,
		m.ToVersion,
		m.Success,
		m.BackupPath,
	)
}

// migrateV0ToV1 adds default fields to config
func migrateV0ToV1(configPath string) error {
	content, err := os.ReadFile(configPath)
	if err != nil {
		return err
	}

	// Add version header if not present
	contentStr := string(content)
	if len(contentStr) > 0 && contentStr[0] != 'v' {
		// Prepend version
		contentStr = fmt.Sprintf("version: 1\n%s", contentStr)
	}

	return os.WriteFile(configPath, []byte(contentStr), 0644)
}

// ListBackups returns list of available config backups
func ListBackups() ([]string, error) {
	backupDir := ".sdp/backups"
	entries, err := os.ReadDir(backupDir)
	if os.IsNotExist(err) {
		return nil, nil
	}
	if err != nil {
		return nil, err
	}

	var backups []string
	for _, entry := range entries {
		if !entry.IsDir() && filepath.Ext(entry.Name()) == ".yml" {
			backups = append(backups, filepath.Join(backupDir, entry.Name()))
		}
	}

	return backups, nil
}

// detectConfigVersion extracts version from config content
func detectConfigVersion(content string) int {
	// Simple detection - look for "version: N"
	// If no version found, assume 0 (needs migration)
	var version int
	fmt.Sscanf(content, "version: %d", &version)
	if version == 0 {
		// Try finding version anywhere in file
		fmt.Sscanf(content, "%*s\nversion: %d", &version)
	}
	return version
}

// getLatestVersion returns the highest version in registry
func getLatestVersion() int {
	max := 0
	for v := range MigrationRegistry {
		if v > max {
			max = v
		}
	}
	return max
}
