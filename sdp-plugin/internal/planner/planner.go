// Package planner provides feature decomposition into workstreams.
//
// The planner decomposes feature descriptions into structured workstreams
// with dependencies, supporting interactive planning, auto-apply execution,
// and multiple output formats.
package planner

// All core types and functionality are split into separate files:
// - types.go: Core data structures (Planner, Workstream, Dependency, etc.)
// - decomposition.go: Feature decomposition logic
// - workstream.go: Workstream file creation and content generation
// - evidence.go: Evidence event emission
// - output.go: Output formatting (JSON, human-readable)
// - mode.go: Interactive and auto-apply mode handling
