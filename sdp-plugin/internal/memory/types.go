package memory

import (
	"time"
)

// Artifact represents an indexed project artifact
type Artifact struct {
	ID           string    // hash of path
	Path         string    // relative path
	Type         string    // "doc", "code", "decision"
	Title        string    // extracted from frontmatter or first heading
	Content      string    // full text content
	Embedding    []float64 // cached embedding (nil if not computed)
	FeatureID    string    // extracted from frontmatter
	WorkstreamID string    // extracted from filename
	Tags         []string  // extracted from frontmatter
	FileHash     string    // SHA256 of file content
	IndexedAt    time.Time
}

// IndexStats holds statistics about indexing operations
type IndexStats struct {
	TotalFiles int
	Indexed    int
	Updated    int
	Skipped    int
	Errors     int
}

// StoreStats holds statistics about the artifact store
type StoreStats struct {
	TotalArtifacts int
	ByType         map[string]int
	LastIndexed    time.Time
}
