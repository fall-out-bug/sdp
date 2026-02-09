package evidence

import (
	"strings"

	"github.com/fall-out-bug/sdp/internal/decision"
)

// Match represents a similar past decision (AC7, AC8).
type Match struct {
	Question   string
	Outcome    string
	WorkstreamID string
	Tags       []string
}

// FindSimilarDecisions returns past decisions similar to query by tags and keywords (AC7).
// Filters to failed/mixed outcomes. No embeddings — keyword + tag match only.
func FindSimilarDecisions(query string, tags []string, decisions []decision.Decision) []Match {
	var out []Match
	queryLower := strings.ToLower(query)
	tagSet := make(map[string]bool)
	for _, t := range tags {
		tagSet[strings.ToLower(t)] = true
	}
	for _, d := range decisions {
		outcome := strings.ToLower(strings.TrimSpace(d.Outcome))
		if outcome != "failed" && !strings.Contains(outcome, "fail") && outcome != "mixed" {
			continue
		}
		score := 0
		if query != "" {
			if strings.Contains(strings.ToLower(d.Question), queryLower) {
				score++
			}
		}
		for _, t := range d.Tags {
			if tagSet[strings.ToLower(t)] {
				score++
				break
			}
		}
		if query == "" && len(tags) == 0 {
			score = 1
		}
		if score > 0 {
			out = append(out, Match{
				Question:      d.Question,
				Outcome:       d.Outcome,
				WorkstreamID:  d.WorkstreamID,
				Tags:          d.Tags,
			})
		}
	}
	return out
}
