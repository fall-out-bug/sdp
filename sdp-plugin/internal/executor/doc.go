// Package executor provides workstream execution with progress tracking and retry.
//
// # WorkstreamRunner
//
// Executor requires a WorkstreamRunner to execute individual workstreams. The interface
// has one method: Run(ctx, wsID) error. Callers must provide an implementation.
//
// # Production implementation: CLIRunner
//
// CLIRunner is the production implementation for CLI-driven workstream execution.
// It runs a configurable command (e.g. "sdp build") for each workstream.
// Use it when workstream execution is delegated to a CLI or agent:
//
//	r := executor.NewCLIRunner("sdp", "build")
//	exec := executor.NewExecutor(cfg, r)
//
// For sdp apply, use CLIRunner with "sdp" "build" to invoke the @build skill via CLI.
//
// # Custom implementations
//
// Implement WorkstreamRunner for:
//   - Agent integration: spawn LLM with @build prompt
//   - CI: run go test, apply patches, etc.
//   - Dry-run/testing: no-op or mock
//
// See runner_test.go for test doubles (testRunner, blockingRunner).
package executor
