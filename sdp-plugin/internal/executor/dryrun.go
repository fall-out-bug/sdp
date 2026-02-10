package executor

import (
	"context"
	"fmt"
	"io"
	"path/filepath"
	"strings"

	"github.com/fall-out-bug/sdp/internal/parser"
)

// executeDryRun shows what would be executed without running
func (e *Executor) executeDryRun(ctx context.Context, output io.Writer, workstreams []string, result *ExecutionResult) (*ExecutionResult, error) {
	fmt.Fprintln(output, "DRY RUN MODE - Showing execution plan:")
	fmt.Fprintln(output, "")

	for _, wsID := range workstreams {
		// Parse workstream file
		wsPath := filepath.Join(e.config.BacklogDir, fmt.Sprintf("%s-*.md", wsID))
		matches, err := filepath.Glob(wsPath)
		if err != nil || len(matches) == 0 {
			fmt.Fprintf(output, "  [✗] %s - file not found\n", wsID)
			continue
		}

		ws, err := parser.ParseWorkstream(matches[0])
		if err != nil {
			fmt.Fprintf(output, "  [✗] %s - parse error: %v\n", wsID, err)
			continue
		}

		fmt.Fprintf(output, "  [→] %s - %s\n", wsID, ws.Goal)

		// Show dependencies
		deps, _ := e.ParseDependencies(wsID)
		if len(deps) > 0 {
			fmt.Fprintf(output, "      Depends on: %s\n", strings.Join(deps, ", "))
		}
	}

	fmt.Fprintln(output, "")
	fmt.Fprintf(output, "Would execute %d workstreams\n", len(workstreams))

	return result, nil
}
