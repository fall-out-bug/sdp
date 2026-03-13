package beadscli

import (
	"os"
	"path/filepath"
	"testing"
)

func TestCreateAndSync_UsesExportFallback(t *testing.T) {
	tmpDir := t.TempDir()
	oldWd, _ := os.Getwd()
	oldPath := os.Getenv("PATH")
	t.Cleanup(func() {
		_ = os.Chdir(oldWd)
		_ = os.Setenv("PATH", oldPath)
	})
	if err := os.Chdir(tmpDir); err != nil {
		t.Fatalf("chdir: %v", err)
	}

	bdScript := `#!/bin/sh
if [ "$1" = "create" ]; then
	echo "Created issue: sdp-test1"
	exit 0
fi
if [ "$1" = "sync" ]; then
	echo 'unknown command "sync"' >&2
	exit 1
fi
if [ "$1" = "export" ] && [ "$2" = "-o" ]; then
	mkdir -p "$(dirname "$3")"
	printf '{"id":"sdp-test1"}\n' > "$3"
	exit 0
fi
exit 1
`
	bdPath := filepath.Join(tmpDir, "bd")
	if err := os.WriteFile(bdPath, []byte(bdScript), 0o755); err != nil {
		t.Fatalf("write fake bd: %v", err)
	}
	if err := os.Setenv("PATH", tmpDir+string(os.PathListSeparator)+oldPath); err != nil {
		t.Fatalf("set PATH: %v", err)
	}

	if err := CreateAndSync(CreateOptions{Title: "CI BLOCKED", Priority: "0", Labels: []string{"ci-finding"}}); err != nil {
		t.Fatalf("CreateAndSync() failed: %v", err)
	}

	if _, err := os.Stat(filepath.Join(tmpDir, ".beads", "issues.jsonl")); err != nil {
		t.Fatalf("expected .beads/issues.jsonl after sync fallback: %v", err)
	}
}

func TestSync_ReturnsExportFailure(t *testing.T) {
	tmpDir := t.TempDir()
	oldWd, _ := os.Getwd()
	oldPath := os.Getenv("PATH")
	t.Cleanup(func() {
		_ = os.Chdir(oldWd)
		_ = os.Setenv("PATH", oldPath)
	})
	if err := os.Chdir(tmpDir); err != nil {
		t.Fatalf("chdir: %v", err)
	}

	bdScript := `#!/bin/sh
if [ "$1" = "sync" ]; then
	echo 'unknown command "sync"' >&2
	exit 1
fi
if [ "$1" = "export" ]; then
	echo "export failed" >&2
	exit 1
fi
exit 1
`
	bdPath := filepath.Join(tmpDir, "bd")
	if err := os.WriteFile(bdPath, []byte(bdScript), 0o755); err != nil {
		t.Fatalf("write fake bd: %v", err)
	}
	if err := os.Setenv("PATH", tmpDir+string(os.PathListSeparator)+oldPath); err != nil {
		t.Fatalf("set PATH: %v", err)
	}

	if err := Sync(); err == nil {
		t.Fatal("expected Sync() to fail when export fallback fails")
	}
}
