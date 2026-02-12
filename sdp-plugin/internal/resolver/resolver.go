package resolver

import (
	"bufio"
	"encoding/json"
	"fmt"
	"os"
	"path/filepath"
	"regexp"
	"strings"
)

// Pattern definitions for ID detection (fast, no I/O)
var (
	// Standard workstream: 00-064-01
	workstreamPattern = regexp.MustCompile(`^\d{2}-\d{3}-\d{2}$`)

	// Fix workstream: 99-F064-0001
	fixWorkstreamPattern = regexp.MustCompile(`^\d{2}-[A-Z]\d{3}-\d{4}$`)

	// Beads ID: sdp-ushh, abc-123
	beadsPattern = regexp.MustCompile(`^[a-z]{3}-[a-z0-9]+$`)

	// Issue ID: ISSUE-0001
	issuePattern = regexp.MustCompile(`^ISSUE-\d+$`)
)

// DetectIDType determines the type of identifier from its pattern
// This is a pure function with no I/O - fast pattern matching only
func DetectIDType(id string) IDType {
	id = strings.TrimSpace(id)

	if id == "" {
		return TypeUnknown
	}

	// Check patterns in order of specificity
	if workstreamPattern.MatchString(id) || fixWorkstreamPattern.MatchString(id) {
		return TypeWorkstream
	}

	if beadsPattern.MatchString(id) {
		return TypeBeads
	}

	if issuePattern.MatchString(id) {
		return TypeIssue
	}

	return TypeUnknown
}

// Resolver resolves identifiers to their file paths
type Resolver struct {
	workstreamDir string
	issuesDir     string
	indexFile     string
}

// Option configures the resolver
type Option func(*Resolver)

// WithWorkstreamDir sets the workstream directory
func WithWorkstreamDir(dir string) Option {
	return func(r *Resolver) {
		r.workstreamDir = dir
	}
}

// WithIssuesDir sets the issues directory
func WithIssuesDir(dir string) Option {
	return func(r *Resolver) {
		r.issuesDir = dir
	}
}

// WithIndexFile sets the issues index file path
func WithIndexFile(path string) Option {
	return func(r *Resolver) {
		r.indexFile = path
	}
}

// NewResolver creates a new resolver with options
func NewResolver(opts ...Option) *Resolver {
	r := &Resolver{
		workstreamDir: "docs/workstreams/backlog",
		issuesDir:     "docs/issues",
		indexFile:     ".sdp/issues-index.jsonl",
	}

	for _, opt := range opts {
		opt(r)
	}

	return r
}

// Resolve resolves any identifier to its file path
func (r *Resolver) Resolve(id string) (*Result, error) {
	id = strings.TrimSpace(id)

	if id == "" {
		return nil, fmt.Errorf("empty identifier")
	}

	idType := DetectIDType(id)

	switch idType {
	case TypeWorkstream:
		return r.resolveWorkstream(id)
	case TypeBeads:
		return r.resolveBeads(id)
	case TypeIssue:
		return r.resolveIssue(id)
	default:
		return nil, fmt.Errorf("unknown identifier format: %s", id)
	}
}

// resolveIssue resolves an issue ID to its file path
func (r *Resolver) resolveIssue(issueID string) (*Result, error) {
	// First try index file for O(1) lookup
	if r.indexFile != "" {
		result, err := r.resolveIssueFromIndex(issueID)
		if err == nil {
			return result, nil
		}
		// SECURITY: Propagate security errors - don't fall back to filesystem
		if strings.HasPrefix(err.Error(), "security:") {
			return nil, err
		}
		// For other errors (not found, etc), fall through to filesystem
	}

	// Fallback to filesystem search
	path := filepath.Join(r.issuesDir, issueID+".md")

	if _, err := os.Stat(path); err != nil {
		return nil, fmt.Errorf("issue not found: %s", issueID)
	}

	result := &Result{
		Type: TypeIssue,
		ID:   issueID,
		Path: path,
	}

	r.extractFrontmatter(path, result)

	return result, nil
}

// issueIndexEntry represents an entry in the issues index
type issueIndexEntry struct {
	IssueID string `json:"issue_id"`
	Title   string `json:"title"`
	Status  string `json:"status"`
	File    string `json:"file"`
}

// resolveIssueFromIndex resolves an issue using the index file
func (r *Resolver) resolveIssueFromIndex(issueID string) (*Result, error) {
	file, err := os.Open(r.indexFile)
	if err != nil {
		return nil, err
	}
	defer func() {
		_ = file.Close()
	}()

	scanner := bufio.NewScanner(file)
	for scanner.Scan() {
		line := scanner.Text()
		if line == "" {
			continue
		}

		var entry issueIndexEntry
		if err := json.Unmarshal([]byte(line), &entry); err != nil {
			continue
		}

		if entry.IssueID == issueID {
			// Resolve relative path to absolute
			path := entry.File
			if !filepath.IsAbs(path) {
				// Assume relative to project root
				path = filepath.Join(filepath.Dir(r.indexFile), "..", entry.File)
			}
			path = filepath.Clean(path)

			// SECURITY: Validate path is within expected issues directory
			// to prevent path traversal attacks
			if err := r.validatePathInIssuesDir(path); err != nil {
				return nil, fmt.Errorf("security: path validation failed: %w", err)
			}

			return &Result{
				Type:   TypeIssue,
				ID:     issueID,
				Path:   path,
				Title:  entry.Title,
				Status: entry.Status,
			}, nil
		}
	}

	return nil, fmt.Errorf("issue not found in index: %s", issueID)
}

// validatePathInIssuesDir ensures the resolved path is within the issues directory
func (r *Resolver) validatePathInIssuesDir(resolvedPath string) error {
	// Get absolute paths for comparison
	absResolved, err := filepath.Abs(resolvedPath)
	if err != nil {
		return fmt.Errorf("cannot resolve absolute path: %w", err)
	}

	absIssuesDir, err := filepath.Abs(r.issuesDir)
	if err != nil {
		return fmt.Errorf("cannot resolve issues directory: %w", err)
	}

	// Ensure path is within issues directory
	// Use HasPrefix check with separator to prevent /docs/issues-other from matching
	expectedPrefix := absIssuesDir + string(filepath.Separator)
	if absResolved != absIssuesDir && !strings.HasPrefix(absResolved, expectedPrefix) {
		return fmt.Errorf("path '%s' is outside issues directory '%s'", resolvedPath, r.issuesDir)
	}

	return nil
}
