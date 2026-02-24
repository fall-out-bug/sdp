package orchestrate_test

import (
	"path/filepath"
	"testing"

	"github.com/fall-out-bug/sdp/internal/orchestrate"
)

func TestDiscoverWorkstreams(t *testing.T) {
	root := filepath.Join("..", "..")
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
