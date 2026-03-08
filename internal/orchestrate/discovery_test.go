package orchestrate_test

import (
	"os"
	"path/filepath"
	"testing"

	"github.com/fall-out-bug/sdp/internal/orchestrate"
)

func TestDiscoverWorkstreams(t *testing.T) {
	root := writeDiscoveryProjectRoot(t)
	ws, err := orchestrate.DiscoverWorkstreams(root, "F016")
	if err != nil {
		t.Fatalf("DiscoverWorkstreams: %v", err)
	}
	if len(ws) != 4 {
		t.Errorf("expected 4 workstreams, got %d: %v", len(ws), ws)
	}
	// 00-016-01 must come before 00-016-02, 00-016-03, 00-016-04 (depends_on)
	if ws[0] != "00-016-01" {
		t.Errorf("expected first WS 00-016-01, got %s", ws[0])
	}
}

func writeDiscoveryProjectRoot(t *testing.T) string {
	t.Helper()
	root := t.TempDir()
	wsDir := filepath.Join(root, "docs", "workstreams", "backlog")
	if err := os.MkdirAll(wsDir, 0o755); err != nil {
		t.Fatal(err)
	}
	files := map[string]string{
		"00-016-01.md": "---\nws_id: 00-016-01\nfeature_id: F016\ndepends_on: []\n---\n",
		"00-016-02.md": "---\nws_id: 00-016-02\nfeature_id: F016\ndepends_on: [\"00-016-01\"]\n---\n",
		"00-016-03.md": "---\nws_id: 00-016-03\nfeature_id: F016\ndepends_on: [\"00-016-01\"]\n---\n",
		"00-016-04.md": "---\nws_id: 00-016-04\nfeature_id: F016\ndepends_on: [\"00-016-02\", \"00-016-03\"]\n---\n",
	}
	for name, content := range files {
		if err := os.WriteFile(filepath.Join(wsDir, name), []byte(content), 0o644); err != nil {
			t.Fatal(err)
		}
	}
	return root
}
