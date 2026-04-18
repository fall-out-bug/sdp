package main

import (
	"fmt"
	"os/exec"
	"strings"

	"github.com/fall-out-bug/sdp/internal/evidence"
	"github.com/spf13/cobra"
)

const deployWSID = "00-000-00" // repo-level approval

var (
	deployResolveSHA = func() (string, error) {
		out, err := exec.Command("git", "rev-parse", "HEAD").Output()
		if err != nil {
			return "", err
		}
		return strings.TrimSpace(string(out)), nil
	}
	deployResolveApprover = func() (string, error) {
		out, err := exec.Command("git", "config", "user.name").Output()
		if err != nil {
			return "", err
		}
		return strings.TrimSpace(string(out)), nil
	}
)

func deployCmd() *cobra.Command {
	var targetBranch, sha, who string

	cmd := &cobra.Command{
		Use:   "deploy",
		Short: "Record deployment approval in evidence log",
		Long: `Emit an approval event after merge. Call after git merge (e.g. from @deploy skill).

  sdp deploy --target main
  sdp deploy --target main --sha abc123 --who "CI"`,
		RunE: func(cmd *cobra.Command, args []string) error {
			if targetBranch == "" {
				targetBranch = "main"
			}
			if sha == "" {
				resolvedSHA, err := deployResolveSHA()
				if err != nil {
					return fmt.Errorf("git rev-parse HEAD: %w", err)
				}
				sha = resolvedSHA
			}
			if who == "" {
				resolvedWho, err := deployResolveApprover()
				if err == nil {
					who = resolvedWho
				}
				if who == "" {
					who = "unknown"
				}
			}
			if evidence.Enabled() {
				if err := evidence.EmitSync(evidence.ApprovalEvent(deployWSID, targetBranch, sha, who)); err != nil {
					return err
				}
			}
			fmt.Printf("Approval recorded: %s -> %s (%s)\n", sha[:min(7, len(sha))], targetBranch, who)
			return nil
		},
	}
	cmd.Flags().StringVar(&targetBranch, "target", "main", "Merge target branch")
	cmd.Flags().StringVar(&sha, "sha", "", "Commit SHA (default: git rev-parse HEAD)")
	cmd.Flags().StringVar(&who, "who", "", "Approved by (default: git config user.name)")
	return cmd
}
