package resolver

// IDType represents the type of identifier detected
type IDType string

const (
	TypeUnknown    IDType = "unknown"
	TypeWorkstream IDType = "workstream"
	TypeBeads      IDType = "beads"
	TypeIssue      IDType = "issue"
)

// Result contains the resolution result
type Result struct {
	Type   IDType // Type of identifier
	ID     string // Original identifier
	WSID   string // Workstream ID (if resolved from beads)
	Path   string // File path to the artifact
	Title  string // Title from frontmatter (optional)
	Status string // Status from frontmatter (optional)
}
