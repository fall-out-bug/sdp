package doctor

import "os"

const (
	beadsDir      = ".beads"
	beadsSnapshot = ".beads/issues.jsonl"
	beadsLegacyDB = ".beads/beads.db"
)

// findBeadsStateArtifact returns the canonical Beads snapshot when present,
// otherwise falls back to a legacy/local beads.db artifact.
func findBeadsStateArtifact() (path string, legacy bool) {
	if _, err := os.Stat(beadsSnapshot); err == nil {
		return beadsSnapshot, false
	}
	if _, err := os.Stat(beadsLegacyDB); err == nil {
		return beadsLegacyDB, true
	}
	return "", false
}
