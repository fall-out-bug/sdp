package executor

import (
	"context"
	"fmt"
	"io"
	"path/filepath"
	"strings"

	"github.com/fall-out-bug/sdp/internal/parser"
)

func writeLine(w io.Writer, s string) error {
	_, err := fmt.Fprint(w, s+"\n")
	return err
}
func writeFmt(w io.Writer, format string, a ...interface{}) error {
	_, err := fmt.Fprintf(w, format, a...)
	return err
}

// executeDryRun shows what would be executed without running.
//
//nolint:gocognit // many write branches by design
func (e *Executor) executeDryRun(ctx context.Context, output io.Writer, workstreams []string, result *ExecutionResult) (*ExecutionResult, error) {
	if err := writeLine(output, "DRY RUN MODE - Showing execution plan:"); err != nil {
		return result, fmt.Errorf("write: %w", err)
	}
	if err := writeLine(output, ""); err != nil {
		return result, fmt.Errorf("write: %w", err)
	}

	for _, wsID := range workstreams {
		wsPath := filepath.Join(e.config.BacklogDir, fmt.Sprintf("%s-*.md", wsID))
		matches, err := filepath.Glob(wsPath)
		if err != nil || len(matches) == 0 {
			if err := writeFmt(output, "  [✗] %s - file not found\n", wsID); err != nil {
				return result, fmt.Errorf("write: %w", err)
			}
			continue
		}

		ws, err := parser.ParseWorkstream(matches[0])
		if err != nil {
			if err := writeFmt(output, "  [✗] %s - parse error: %v\n", wsID, err); err != nil {
				return result, fmt.Errorf("write: %w", err)
			}
			continue
		}

		if err := writeFmt(output, "  [→] %s - %s\n", wsID, ws.Goal); err != nil {
			return result, fmt.Errorf("write: %w", err)
		}

		deps, err := e.ParseDependencies(wsID)
		if err == nil && len(deps) > 0 {
			if err := writeFmt(output, "      Depends on: %s\n", strings.Join(deps, ", ")); err != nil {
				return result, fmt.Errorf("write: %w", err)
			}
		}
	}

	if err := writeLine(output, ""); err != nil {
		return result, fmt.Errorf("write: %w", err)
	}
	if err := writeFmt(output, "Would execute %d workstreams\n", len(workstreams)); err != nil {
		return result, fmt.Errorf("write: %w", err)
	}

	return result, nil
}
