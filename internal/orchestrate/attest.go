package orchestrate

import (
	"bufio"
	"encoding/json"
	"fmt"
	"os"
	"os/exec"
	"path/filepath"
	"regexp"
	"strings"
	"time"

	intoto "github.com/in-toto/in-toto-golang/in_toto"
	"github.com/in-toto/in-toto-golang/in_toto/slsa_provenance/common"

	"github.com/fall-out-bug/sdp/internal/evidenceenv"
)

// GenerateOrchestratorAttestation creates an in-toto attestation from a checkpoint.
// Called by sdp-orchestrate --advance after each phase transition.
// The attestation captures what the orchestrator knows: intent, plan, execution boundary.
// CI auto-attestation later adds verification (test results, lint, coverage).
func GenerateOrchestratorAttestation(projectRoot string, cp *Checkpoint) (evidenceenv.CodingWorkflowStatement, error) {
	branch := cp.Branch
	headSHA, err := gitHeadSHA(projectRoot)
	if err != nil {
		headSHA = "unknown"
	}

	// Extract beads IDs from the workstream mapping for this feature
	beadsIDs := lookupBeadsIDsForFeature(projectRoot, cp.FeatureID)
	issueID := firstBeadsID(beadsIDs)
	if issueID == "" {
		issueID = cp.FeatureID
	}

	// Collect workstream IDs in order
	wsIDs := make([]string, 0, len(cp.Workstreams))
	for _, ws := range cp.Workstreams {
		wsIDs = append(wsIDs, ws.ID)
	}

	// Get changed files since branch diverged from master
	changedFiles := getChangedFilesSinceBranch(projectRoot, "master")

	// Determine scope from workstream files
	scopePrefixes := collectWorkstreamScopePrefixes(projectRoot, wsIDs)
	outOfBoundary := checkOutOfBoundary(changedFiles, scopePrefixes)
	scopeOK := len(outOfBoundary) == 0

	scopeReason := fmt.Sprintf("all %d changed files within declared scope", len(changedFiles))
	if !scopeOK {
		scopeReason = fmt.Sprintf("%d files outside declared scope: %s", len(outOfBoundary), strings.Join(outOfBoundary, ", "))
	}

	subjects := []intoto.Subject{{
		Name:   fmt.Sprintf("branch:%s", branch),
		Digest: common.DigestSet{"sha256": headSHA},
	}}

	predicate := evidenceenv.CodingWorkflowPredicate{
		Intent: evidenceenv.Intent{
			IssueID: issueID,
			Trigger: "sdp-orchestrate",
		},
		Plan: evidenceenv.Plan{
			Workstreams:       wsIDs,
			OrderingRationale: "sequential execution via sdp-orchestrate state machine",
		},
		Execution: evidenceenv.Execution{
			ClaimedIssueIDs: beadsIDs,
			Branch:          branch,
			ChangedFiles:    changedFiles,
		},
		Verification: evidenceenv.Verification{
			// Tests filled by CI auto-attestation; leave empty with a note
			Tests: []evidenceenv.GateResult{{
				Name:   "orchestrator-phase",
				Status: fmt.Sprintf("phase=%s", cp.Phase),
			}},
		},
		Boundary: evidenceenv.Boundary{
			Declared: evidenceenv.DeclaredBoundary{
				AllowedPathPrefixes: scopePrefixes,
			},
			Observed: evidenceenv.ObservedBoundary{
				TouchedPaths:       changedFiles,
				OutOfBoundaryPaths: outOfBoundary,
			},
			Compliance: evidenceenv.BoundaryCompliance{
				OK:     scopeOK,
				Reason: scopeReason,
			},
		},
		Provenance: evidenceenv.Provenance{
			RunID:        fmt.Sprintf("orch-%s-%s", cp.FeatureID, headSHA[:minLen(len(headSHA), 8)]),
			Orchestrator: "sdp-orchestrate",
			Runtime:      "local",
			Phase:        cp.Phase,
			SourceIssueID: issueID,
			CapturedAt:   time.Now().UTC().Format(time.RFC3339),
		},
		Trace: evidenceenv.Trace{
			BeadsIDs: beadsIDs,
			Branch:   branch,
			Commits:  []string{headSHA},
			PRURL:    cp.PRURL,
		},
	}

	if cp.Review != nil && cp.Review.Status == "approved" {
		predicate.Review.SelfReview = []evidenceenv.ReviewItem{{
			Reviewer: "sdp-orchestrate",
			Verdict:  "APPROVED",
			Notes:    fmt.Sprintf("iteration %d", cp.Review.Iteration),
		}}
	}

	return evidenceenv.NewStatement(subjects, predicate), nil
}

// WriteOrchestratorAttestation saves the attestation to .sdp/evidence/FXXX.json.
func WriteOrchestratorAttestation(projectRoot string, cp *Checkpoint) error {
	stmt, err := GenerateOrchestratorAttestation(projectRoot, cp)
	if err != nil {
		return fmt.Errorf("generate attestation: %w", err)
	}

	evidenceDir := filepath.Join(projectRoot, ".sdp", "evidence")
	if err := os.MkdirAll(evidenceDir, 0o755); err != nil {
		return fmt.Errorf("mkdir evidence: %w", err)
	}

	outPath := filepath.Join(evidenceDir, cp.FeatureID+".json")
	return evidenceenv.WriteAttestation(outPath, stmt)
}

var beadsIDRe = regexp.MustCompile(`sdp_dev-[a-z0-9]{4}`)

// lookupBeadsIDsForFeature reads the beads mapping file to find issues for a feature.
func lookupBeadsIDsForFeature(projectRoot, featureID string) []string {
	mappingPath := filepath.Join(projectRoot, ".beads-sdp-mapping.jsonl")
	f, err := os.Open(mappingPath)
	if err != nil {
		return nil
	}
	defer f.Close()

	// Feature ID "F028" → workstream prefix "00-028"
	featureNum := extractFeatureNum(featureID)
	if featureNum == "" {
		return nil
	}

	prefix := fmt.Sprintf("00-%s-", featureNum)
	var ids []string
	scanner := bufio.NewScanner(f)
	for scanner.Scan() {
		var entry struct {
			SDPID   string `json:"sdp_id"`
			BeadsID string `json:"beads_id"`
		}
		if json.Unmarshal(scanner.Bytes(), &entry) == nil {
			if strings.HasPrefix(entry.SDPID, prefix) {
				ids = append(ids, entry.BeadsID)
			}
		}
	}
	return ids
}

var featureNumRe = regexp.MustCompile(`[Ff](\d+)`)

func extractFeatureNum(featureID string) string {
	m := featureNumRe.FindStringSubmatch(featureID)
	if m == nil {
		return ""
	}
	n := m[1]
	// Pad to 3 digits
	for len(n) < 3 {
		n = "0" + n
	}
	return n
}

// collectWorkstreamScopePrefixes reads workstream files and extracts declared scope.
func collectWorkstreamScopePrefixes(projectRoot string, wsIDs []string) []string {
	backlogDir := filepath.Join(projectRoot, "docs", "workstreams", "backlog")
	var prefixes []string
	seen := map[string]bool{}

	for _, wsID := range wsIDs {
		wsPath := filepath.Join(backlogDir, wsID+".md")
		f, err := os.Open(wsPath)
		if err != nil {
			continue
		}

		inScope := false
		scanner := bufio.NewScanner(f)
		for scanner.Scan() {
			line := scanner.Text()
			if strings.HasPrefix(line, "## Scope Files") {
				inScope = true
				continue
			}
			if inScope && strings.HasPrefix(line, "##") {
				break
			}
			if inScope && strings.HasPrefix(line, "- ") {
				path := strings.TrimPrefix(line, "- ")
				path = strings.TrimSpace(strings.Trim(path, "`"))
				if path != "" && !seen[path] {
					seen[path] = true
					prefixes = append(prefixes, path)
				}
			}
		}
		f.Close()
	}
	return prefixes
}

func checkOutOfBoundary(files, prefixes []string) []string {
	if len(prefixes) == 0 {
		return nil
	}
	var out []string
	for _, f := range files {
		if !matchesPrefix(f, prefixes) {
			out = append(out, f)
		}
	}
	return out
}

func matchesPrefix(file string, prefixes []string) bool {
	for _, p := range prefixes {
		if strings.HasPrefix(file, p) || file == p {
			return true
		}
	}
	return false
}

// GetChangedFiles returns changed files vs origin/master (for policy input construction).
func GetChangedFiles(projectRoot string) []string {
	return getChangedFilesSinceBranch(projectRoot, "master")
}

func getChangedFilesSinceBranch(projectRoot, baseBranch string) []string {
	cmd := exec.Command("git", "diff", "--name-only", "origin/"+baseBranch+"...HEAD")
	cmd.Dir = projectRoot
	out, err := cmd.Output()
	if err != nil {
		// Fallback: uncommitted changes
		cmd2 := exec.Command("git", "diff", "--name-only", "HEAD")
		cmd2.Dir = projectRoot
		out2, _ := cmd2.Output()
		return splitLines(string(out2))
	}
	return splitLines(string(out))
}

func gitHeadSHA(projectRoot string) (string, error) {
	cmd := exec.Command("git", "rev-parse", "HEAD")
	cmd.Dir = projectRoot
	out, err := cmd.Output()
	if err != nil {
		return "", err
	}
	return strings.TrimSpace(string(out)), nil
}

func splitLines(s string) []string {
	lines := strings.Split(strings.TrimSpace(s), "\n")
	result := make([]string, 0, len(lines))
	for _, l := range lines {
		l = strings.TrimSpace(l)
		if l != "" {
			result = append(result, l)
		}
	}
	return result
}

func firstBeadsID(ids []string) string {
	if len(ids) > 0 {
		return ids[0]
	}
	return ""
}

func minLen(a, b int) int {
	if a < b {
		return a
	}
	return b
}
