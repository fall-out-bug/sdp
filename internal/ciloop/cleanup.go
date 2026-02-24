package ciloop

import (
	"os"
	"path/filepath"
)

// RemoveOrphanTmpFiles removes stale .tmp files in the given directories.
// These can remain if a process crashed between WriteFile and Rename.
func RemoveOrphanTmpFiles(dirs ...string) {
	for _, dir := range dirs {
		entries, err := os.ReadDir(dir)
		if err != nil {
			continue
		}
		for _, e := range entries {
			if e.IsDir() {
				continue
			}
			if len(e.Name()) > 4 && e.Name()[len(e.Name())-4:] == ".tmp" {
				_ = os.Remove(filepath.Join(dir, e.Name()))
			}
		}
	}
}
