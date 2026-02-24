package main

import (
	"flag"
	"fmt"
	"os"
	"path/filepath"

	"github.com/fall-out-bug/sdp/internal/eval"
)

func main() {
	skill := flag.String("skill", "", "Skill name (e.g. oneshot). If empty, run all.")
	all := flag.Bool("all", false, "Run evals for all skills")
	projectRoot := flag.String("project-root", ".", "Project root")
	casesDir := flag.String("cases-dir", "", "Cases directory (default: internal/eval/cases)")
	flag.Parse()

	if *casesDir == "" {
		*casesDir = filepath.Join(*projectRoot, "internal", "eval", "cases")
	}

	skillFilter := *skill
	if *all {
		skillFilter = ""
	}
	if !*all && skillFilter == "" {
		fmt.Fprintln(os.Stderr, "error: --skill <name> or --all required")
		flag.Usage()
		os.Exit(1)
	}

	results, err := eval.Run(*projectRoot, *casesDir, skillFilter)
	if err != nil {
		fmt.Fprintf(os.Stderr, "error: %v\n", err)
		os.Exit(1)
	}

	passed := 0
	for _, r := range results {
		status := "FAIL"
		if r.Pass {
			status = "PASS"
			passed++
		}
		fmt.Printf("  %s: %s", r.Case, status)
		if !r.Pass && r.Reason != "" {
			fmt.Printf(" (%s)", r.Reason)
		}
		fmt.Println()
	}

	skillLabel := "all"
	if skillFilter != "" {
		skillLabel = skillFilter
	}
	fmt.Printf("\n%s: %d/%d passed\n", skillLabel, passed, len(results))
	if passed < len(results) {
		os.Exit(1)
	}
}
