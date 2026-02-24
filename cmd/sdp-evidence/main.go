package main

import (
	"flag"
	"fmt"
	"os"
	"path/filepath"

	"github.com/fall-out-bug/sdp/internal/evidenceenv"
)

func main() {
	validateCmd := flag.NewFlagSet("validate", flag.ExitOnError)
	evidencePath := validateCmd.String("evidence", "", "Path to evidence file")
	requirePRURL := validateCmd.Bool("require-pr-url", true, "Require trace.pr_url (set false for prepublish)")

	inspectCmd := flag.NewFlagSet("inspect", flag.ExitOnError)
	inspectEvidence := inspectCmd.String("evidence", "", "Path to evidence file")
	inspectRequirePRURL := inspectCmd.Bool("require-pr-url", true, "Require trace.pr_url (set false for prepublish)")

	if len(os.Args) < 2 {
		printUsage()
		os.Exit(2)
	}

	switch os.Args[1] {
	case "inspect":
		inspectCmd.Parse(os.Args[2:])
		if *inspectEvidence == "" && inspectCmd.NArg() > 0 {
			*inspectEvidence = inspectCmd.Arg(0)
		}
		if *inspectEvidence == "" {
			fmt.Fprintln(os.Stderr, "inspect: --evidence or positional path required")
			inspectCmd.Usage()
			os.Exit(2)
		}
		path, absErr := filepath.Abs(*inspectEvidence)
		if absErr != nil {
			path = *inspectEvidence
		}
		summary, res, err := evidenceenv.Inspect(path, *inspectRequirePRURL)
		if err != nil {
			fmt.Fprintf(os.Stderr, "inspect: %v\n", err)
			os.Exit(1)
		}
		if !res.OK {
			fmt.Fprintf(os.Stderr, "invalid: %s\n", res.Reason)
			os.Exit(1)
		}
		fmt.Println(summary)
		os.Exit(0)
	case "validate":
		validateCmd.Parse(os.Args[2:])
		if *evidencePath == "" {
			// Allow positional: validate <path>
			if validateCmd.NArg() > 0 {
				*evidencePath = validateCmd.Arg(0)
			}
		}
		if *evidencePath == "" {
			fmt.Fprintln(os.Stderr, "validate: --evidence or positional path required")
			validateCmd.Usage()
			os.Exit(2)
		}
		path, err := filepath.Abs(*evidencePath)
		if err != nil {
			path = *evidencePath
		}
		res, err := evidenceenv.ValidateStrictFile(path, *requirePRURL)
		if err != nil {
			fmt.Fprintf(os.Stderr, "validate: %v\n", err)
			os.Exit(1)
		}
		if !res.OK {
			fmt.Fprintf(os.Stderr, "invalid: %s\n", res.Reason)
			if len(res.Missing) > 0 {
				fmt.Fprintf(os.Stderr, "missing sections: %v\n", res.Missing)
			}
			os.Exit(1)
		}
		fmt.Println("valid")
		os.Exit(0)
	default:
		printUsage()
		os.Exit(2)
	}
}

func printUsage() {
	fmt.Fprintf(os.Stderr, `sdp-evidence - validate and inspect evidence envelopes

Usage:
  sdp-evidence validate --evidence <path>   Validate evidence file
  sdp-evidence validate <path>             Same (positional)
  sdp-evidence inspect --evidence <path>   Print human-readable summary
  sdp-evidence inspect <path>              Same (positional)

Exits 0 if valid, non-zero if invalid.
`)
}
