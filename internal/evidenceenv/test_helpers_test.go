package evidenceenv

import (
	"os"
	"path/filepath"
	"testing"
)

func writeValidEvidenceFixture(t *testing.T) string {
	t.Helper()
	path := filepath.Join(t.TempDir(), "strict-evidence-template.json")
	if err := os.WriteFile(path, validEvidenceFixture, 0o644); err != nil {
		t.Fatal(err)
	}
	return path
}
