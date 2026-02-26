//go:build windows

package evidence

import "os"

// lockFile is a no-op on Windows. Evidence layer uses flock on UNIX only.
// See README: "Evidence file lock requires UNIX. Windows is not supported."
func lockFile(f *os.File) error {
	_ = f
	return nil
}

func unlockFile(f *os.File) error {
	_ = f
	return nil
}
