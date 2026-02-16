package evidence

import (
	"testing"

	"github.com/fall-out-bug/sdp/internal/decision"
)

func TestFindSimilarDecisions_Empty(t *testing.T) {
	got := FindSimilarDecisions("", nil, nil)
	if len(got) != 0 {
		t.Errorf("expected 0, got %d", len(got))
	}
}

func TestFindSimilarDecisions_FiltersFailed(t *testing.T) {
	decisions := []decision.Decision{
		{Question: "Use MongoDB", Outcome: "FAILED", WorkstreamID: "00-042-02", Tags: []string{"database"}},
		{Question: "Use PostgreSQL", Outcome: "passed", WorkstreamID: "00-042-03", Tags: []string{"database"}},
	}
	got := FindSimilarDecisions("", nil, decisions)
	if len(got) != 1 {
		t.Fatalf("expected 1 failed, got %d", len(got))
	}
	if got[0].Question != "Use MongoDB" {
		t.Errorf("Question: want Use MongoDB, got %s", got[0].Question)
	}
}

func TestFindSimilarDecisions_KeywordMatch(t *testing.T) {
	decisions := []decision.Decision{
		{Question: "Use MongoDB for user data", Outcome: "failed", WorkstreamID: "00-042-02"},
	}
	got := FindSimilarDecisions("MongoDB", nil, decisions)
	if len(got) != 1 {
		t.Fatalf("expected 1 match, got %d", len(got))
	}
	got = FindSimilarDecisions("Postgres", nil, decisions)
	if len(got) != 0 {
		t.Errorf("expected 0 match, got %d", len(got))
	}
}

func TestFindSimilarDecisions_TagMatch(t *testing.T) {
	decisions := []decision.Decision{
		{Question: "DB choice", Outcome: "failed", WorkstreamID: "00-042-02", Tags: []string{"database", "architecture"}},
	}
	got := FindSimilarDecisions("", []string{"database"}, decisions)
	if len(got) != 1 {
		t.Fatalf("expected 1 tag match, got %d", len(got))
	}
}
