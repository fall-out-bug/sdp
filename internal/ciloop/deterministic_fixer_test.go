package ciloop_test

import (
	"context"
	"testing"

	"github.com/fall-out-bug/sdp/internal/ciloop"
)

func TestRunDeterministicFixersNoMatchReturnsFalse(t *testing.T) {
	dir := t.TempDir()
	reg := ciloop.NewAutofixerRegistry(dir)
	committer := &fakeCommitter{}
	changed, err := ciloop.RunDeterministicFixers(
		context.Background(), dir, "secrets detected",
		reg, committer, nil, nil,
	)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if changed {
		t.Error("expected no change when no fixers match")
	}
	if len(committer.commits) != 0 {
		t.Error("expected no commit when no fixers match")
	}
}


