package main

import (
	"fmt"
	"os"
	"path/filepath"
	"strings"

	"github.com/spf13/cobra"
)

func doctorHooksProvenanceCmd() *cobra.Command {
	return &cobra.Command{
		Use:   "hooks-provenance",
		Short: "Verify commit provenance hooks and metadata pipeline",
		RunE: func(cmd *cobra.Command, args []string) error {
			type hookCheck struct {
				name      string
				required  []string
				requiredX bool
			}

			hooksDir := filepath.Join(".git", "hooks")
			checks := []hookCheck{
				{name: "commit-msg", required: []string{"SDP-Agent", "SDP-Model", "SDP-Task"}, requiredX: true},
				{name: "post-commit", required: []string{"sdp skill record", "commit_sha", "agent", "model"}, requiredX: true},
			}

			hasErrors := false
			fmt.Println("SDP Hooks Provenance Check")
			fmt.Println("==========================")

			for _, check := range checks {
				path := filepath.Join(hooksDir, check.name)
				info, err := os.Stat(path)
				if err != nil {
					fmt.Printf("✗ %s\n", check.name)
					fmt.Printf("    missing at %s\n\n", path)
					hasErrors = true
					continue
				}

				contentBytes, err := os.ReadFile(path)
				if err != nil {
					fmt.Printf("✗ %s\n", check.name)
					fmt.Printf("    cannot read: %v\n\n", err)
					hasErrors = true
					continue
				}

				content := string(contentBytes)
				missing := []string{}
				for _, token := range check.required {
					if !strings.Contains(content, token) {
						missing = append(missing, token)
					}
				}

				notExecutable := check.requiredX && info.Mode().Perm()&0111 == 0
				if len(missing) > 0 || notExecutable {
					fmt.Printf("✗ %s\n", check.name)
					if len(missing) > 0 {
						fmt.Printf("    missing markers: %s\n", strings.Join(missing, ", "))
					}
					if notExecutable {
						fmt.Printf("    not executable\n")
					}
					fmt.Println()
					hasErrors = true
					continue
				}

				fmt.Printf("✓ %s\n", check.name)
				fmt.Printf("    installed and contains provenance markers\n\n")
			}

			if hasErrors {
				fmt.Println("Remediation:")
				fmt.Println("  sdp hooks install --with-provenance")
				return fmt.Errorf("hooks provenance checks failed")
			}

			fmt.Println("All provenance hook checks passed!")
			return nil
		},
	}
}
