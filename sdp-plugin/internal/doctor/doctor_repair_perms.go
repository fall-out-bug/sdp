package doctor

import (
	"fmt"
	"os"
	"path/filepath"
)

// repairFilePermissions fixes insecure permissions on sensitive files
func repairFilePermissions() RepairAction {
	sensitiveFiles := []string{
		filepath.Join(os.Getenv("HOME"), ".sdp", "telemetry.jsonl"),
		".beads/beads.db",
	}

	fixed := []string{}
	failed := []string{}
	checked := 0

	for _, path := range sensitiveFiles {
		info, err := os.Stat(path)
		if err != nil {
			continue // File doesn't exist, skip
		}
		checked++

		if info.IsDir() {
			// Fix files in directory
			entries, err := os.ReadDir(path)
			if err != nil {
				continue
			}
			for _, entry := range entries {
				if entry.IsDir() {
					continue
				}
				fullPath := filepath.Join(path, entry.Name())
				wasFixed, err := fixFilePermissionsTracked(fullPath)
				if err != nil {
					failed = append(failed, fullPath)
				} else if wasFixed {
					fixed = append(fixed, fullPath)
				}
			}
		} else {
			wasFixed, err := fixFilePermissionsTracked(path)
			if err != nil {
				failed = append(failed, path)
			} else if wasFixed {
				fixed = append(fixed, path)
			}
		}
	}

	if checked == 0 {
		return RepairAction{
			Check:   "File Permissions",
			Status:  "skipped",
			Message: "No sensitive files found",
		}
	}

	if len(fixed) == 0 && len(failed) == 0 {
		return RepairAction{
			Check:   "File Permissions",
			Status:  "skipped",
			Message: "All sensitive files already secure",
		}
	}

	if len(failed) > 0 {
		return RepairAction{
			Check:   "File Permissions",
			Status:  "partial",
			Message: fmt.Sprintf("Fixed %d files, failed: %v", len(fixed), failed),
		}
	}

	return RepairAction{
		Check:   "File Permissions",
		Status:  "fixed",
		Message: fmt.Sprintf("Fixed permissions on %d files", len(fixed)),
	}
}

// fixFilePermissions sets file to 0600
func fixFilePermissions(path string) error {
	info, err := os.Stat(path)
	if err != nil {
		return err
	}

	// Only fix if world or group readable
	if info.Mode().Perm()&0077 != 0 {
		return os.Chmod(path, 0600)
	}
	return nil
}

// fixFilePermissionsTracked sets file to 0600 and reports if change was made
func fixFilePermissionsTracked(path string) (bool, error) {
	info, err := os.Stat(path)
	if err != nil {
		return false, err
	}

	// Only fix if world or group readable
	if info.Mode().Perm()&0077 != 0 {
		if err := os.Chmod(path, 0600); err != nil {
			return false, err
		}
		return true, nil
	}
	return false, nil // Already secure, no change
}
