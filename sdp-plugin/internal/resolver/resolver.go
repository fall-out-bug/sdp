package resolver

import (
	"fmt"
	"strings"
)

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
