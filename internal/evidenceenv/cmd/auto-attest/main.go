package main

import (
	"flag"
	"fmt"
	"os"

	"github.com/fall-out-bug/sdp/internal/evidenceenv"
)

func main() {
	baseBranch := flag.String("base-branch", "master", "Base branch for diff")
	prNumber := flag.String("pr-number", "", "PR number")
	prURL := flag.String("pr-url", "", "PR URL")
	output := flag.String("output", ".sdp/attestations/ci-auto.json", "Output attestation path")
	report := flag.String("report", "", "Output report path (optional)")
	flag.Parse()

	wd, err := os.Getwd()
	if err != nil {
		fmt.Fprintf(os.Stderr, "error: %v\n", err)
		os.Exit(1)
	}

	stmt, err := evidenceenv.AutoAttest(evidenceenv.AutoAttestOptions{
		BaseBranch: *baseBranch,
		PRNumber:   *prNumber,
		PRURL:      *prURL,
		RepoRoot:   wd,
	})
	if err != nil {
		fmt.Fprintf(os.Stderr, "auto-attest: %v\n", err)
		os.Exit(1)
	}

	if err := os.MkdirAll(".sdp/attestations", 0o755); err != nil {
		fmt.Fprintf(os.Stderr, "error: %v\n", err)
		os.Exit(1)
	}
	if err := evidenceenv.WriteAttestation(*output, stmt); err != nil {
		fmt.Fprintf(os.Stderr, "write attestation: %v\n", err)
		os.Exit(1)
	}
	fmt.Fprintf(os.Stderr, "attestation written to %s\n", *output)

	if *report != "" {
		if err := evidenceenv.WriteAutoAttestationReport(*report, stmt); err != nil {
			fmt.Fprintf(os.Stderr, "write report: %v\n", err)
			os.Exit(1)
		}
		fmt.Fprintf(os.Stderr, "report written to %s\n", *report)
	}
}
