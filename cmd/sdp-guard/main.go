package main

import (
	"flag"
	"fmt"
	"os"

	"github.com/fall-out-bug/sdp/internal/guard"
	"github.com/fall-out-bug/sdp/internal/orchestrate"
)

func main() {
	ws := flag.String("ws", "", "Workstream ID (e.g. 00-023-01)")
	cached := flag.Bool("cached", false, "Use git diff --cached (staged) instead of HEAD~1")
	checkConstraints := flag.Bool("check-constraints", false, "Check agent constraint rules for a command or file")
	phase := flag.String("phase", "build", "Phase for constraint checking (build, review, pr)")
	command := flag.String("command", "", "Command to check against constraint rules")
	file := flag.String("file", "", "File path to check against constraint rules")
	flag.Parse()

	wd, err := os.Getwd()
	if err != nil {
		fmt.Fprintf(os.Stderr, "error: %v\n", err)
		os.Exit(1)
	}
	projectRoot, err := orchestrate.FindProjectRoot(wd)
	if err != nil {
		fmt.Fprintf(os.Stderr, "error: %v\n", err)
		os.Exit(1)
	}

	if *checkConstraints {
		runConstraintCheck(projectRoot, *phase, *command, *file)
		return
	}

	if *ws == "" {
		fmt.Fprintln(os.Stderr, "error: --ws is required (or use --check-constraints)")
		flag.Usage()
		os.Exit(1)
	}

	verdict, err := guard.CheckScope(projectRoot, *ws, *cached)
	if err != nil {
		fmt.Fprintf(os.Stderr, "error: %v\n", err)
		os.Exit(1)
	}

	if len(verdict.Warnings) > 0 {
		for _, w := range verdict.Warnings {
			fmt.Fprintf(os.Stderr, "WARN: %s (allowlisted)\n", w)
		}
	}

	if verdict.Pass {
		os.Exit(0)
	}

	for _, v := range verdict.Violations {
		fmt.Fprintf(os.Stderr, "SCOPE VIOLATION: %s\n", v)
	}
	fmt.Fprintf(os.Stderr, "out-of-scope changes detected (%d files)\n", len(verdict.Violations))
	os.Exit(1)
}

func runConstraintCheck(projectRoot, phase, command, file string) {
	cfg, err := orchestrate.LoadConstraintConfig(projectRoot)
	if err != nil {
		fmt.Fprintf(os.Stderr, "warning: could not load constraints: %v\n", err)
		os.Exit(0) // graceful degradation
	}

	var violations []orchestrate.ConstraintViolation

	if command != "" {
		violations = append(violations, orchestrate.CheckCommand(cfg, phase, command)...)
	}
	if file != "" {
		violations = append(violations, orchestrate.CheckFileAccess(cfg, phase, file)...)
	}

	if len(violations) == 0 {
		fmt.Println("OK: no constraint violations")
		os.Exit(0)
	}

	maxSeverity := "warn"
	for _, v := range violations {
		fmt.Fprintf(os.Stderr, "[%s] %s: %s\n", v.Severity, v.ConstraintID, v.Message)
		if severityRank(v.Severity) > severityRank(maxSeverity) {
			maxSeverity = v.Severity
		}
	}

	switch maxSeverity {
	case "escalate", "halt":
		fmt.Fprintf(os.Stderr, "HALT: agent session must stop (%s)\n", maxSeverity)
		os.Exit(2)
	case "block":
		fmt.Fprintf(os.Stderr, "BLOCK: action rejected\n")
		os.Exit(1)
	default:
		fmt.Fprintf(os.Stderr, "WARN: %d constraint warning(s)\n", len(violations))
		os.Exit(0)
	}
}

func severityRank(s string) int {
	switch s {
	case "escalate":
		return 4
	case "halt":
		return 3
	case "block":
		return 2
	case "warn":
		return 1
	default:
		return 0
	}
}
