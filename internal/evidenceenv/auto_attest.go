package evidenceenv

import (
	"bufio"
	"encoding/json"
	"fmt"
	"os"
	"os/exec"
	"path/filepath"
	"regexp"
	"strconv"
	"strings"
	"time"

	intoto "github.com/in-toto/in-toto-golang/in_toto"
	"github.com/in-toto/in-toto-golang/in_toto/slsa_provenance/common"
)

type AutoAttestOptions struct {
	BaseBranch string
	PRNumber   string
	PRURL      string
	RepoRoot   string
}

// AutoAttest collects facts from CI (git diff, tests, lint, scope) and generates
// an in-toto CodingWorkflowStatement. No agent action required — CI is the observer.
func AutoAttest(opts AutoAttestOptions) (CodingWorkflowStatement, error) {
	changedFiles, err := gitChangedFiles(opts.RepoRoot, opts.BaseBranch)
	if err != nil {
		return CodingWorkflowStatement{}, fmt.Errorf("git changed files: %w", err)
	}

	branch, err := gitCurrentBranch(opts.RepoRoot)
	if err != nil {
		return CodingWorkflowStatement{}, fmt.Errorf("git branch: %w", err)
	}

	headSHA, err := gitHeadSHA(opts.RepoRoot)
	if err != nil {
		return CodingWorkflowStatement{}, fmt.Errorf("git head SHA: %w", err)
	}

	commits, err := gitCommitsSinceBase(opts.RepoRoot, opts.BaseBranch)
	if err != nil {
		commits = []string{headSHA}
	}

	beadsIDs := extractBeadsIDsFromCommits(opts.RepoRoot, opts.BaseBranch)
	issueID := firstOrEmpty(beadsIDs)
	if issueID == "" {
		issueID = fmt.Sprintf("ci-auto-pr%s", opts.PRNumber)
	}

	testResults, coverage := collectTestResults(opts.RepoRoot)
	lintResults := collectLintResults(opts.RepoRoot)

	boundary, boundaryOK := checkScopeCompliance(opts.RepoRoot, changedFiles)

	subjectName := opts.PRURL
	if subjectName == "" {
		subjectName = fmt.Sprintf("PR #%s", opts.PRNumber)
	}

	subjects := []intoto.Subject{{
		Name:   subjectName,
		Digest: common.DigestSet{"sha256": headSHA},
	}}

	predicate := CodingWorkflowPredicate{
		Intent: Intent{
			IssueID: issueID,
			Trigger: "ci-auto-attestation",
		},
		Plan: Plan{
			Workstreams:       extractWorkstreamsFromBranch(branch),
			OrderingRationale: "auto-detected from branch name",
		},
		Execution: Execution{
			ClaimedIssueIDs: beadsIDs,
			Branch:          branch,
			ChangedFiles:    changedFiles,
		},
		Verification: Verification{
			Tests: testResults,
			Lint:  lintResults,
			Coverage: func() *Coverage {
				if coverage >= 0 {
					return &Coverage{Value: coverage, Threshold: 80}
				}
				return nil
			}(),
		},
		Boundary: boundary,
		Provenance: Provenance{
			RunID:        fmt.Sprintf("ci-auto-%s-%s", opts.PRNumber, headSHA[:minLen(len(headSHA), 8)]),
			Orchestrator: "github-actions",
			Runtime:      "ci",
			CapturedAt:   time.Now().UTC().Format(time.RFC3339),
		},
		Trace: Trace{
			BeadsIDs: beadsIDs,
			Branch:   branch,
			Commits:  commits,
			PRURL:    opts.PRURL,
		},
	}
	_ = boundaryOK

	return NewStatement(subjects, predicate), nil
}

func gitChangedFiles(repoRoot, baseBranch string) ([]string, error) {
	if baseBranch == "" {
		baseBranch = "master"
	}
	out, err := runGit(repoRoot, "diff", "--name-only", "origin/"+baseBranch+"...HEAD")
	if err != nil {
		return nil, err
	}
	return splitLines(out), nil
}

func gitCurrentBranch(repoRoot string) (string, error) {
	out, err := runGit(repoRoot, "branch", "--show-current")
	if err != nil {
		return "", err
	}
	return strings.TrimSpace(out), nil
}

func gitHeadSHA(repoRoot string) (string, error) {
	out, err := runGit(repoRoot, "rev-parse", "HEAD")
	if err != nil {
		return "", err
	}
	return strings.TrimSpace(out), nil
}

func gitCommitsSinceBase(repoRoot, baseBranch string) ([]string, error) {
	if baseBranch == "" {
		baseBranch = "master"
	}
	out, err := runGit(repoRoot, "log", "--format=%H", "origin/"+baseBranch+"...HEAD")
	if err != nil {
		return nil, err
	}
	return splitLines(out), nil
}

var beadsIDRe = regexp.MustCompile(`sdp_dev-[a-z0-9]{4}`)

func extractBeadsIDsFromCommits(repoRoot, baseBranch string) []string {
	if baseBranch == "" {
		baseBranch = "master"
	}
	out, _ := runGit(repoRoot, "log", "--format=%s %b", "origin/"+baseBranch+"...HEAD")
	seen := map[string]bool{}
	var ids []string
	for _, id := range beadsIDRe.FindAllString(out, -1) {
		if !seen[id] {
			seen[id] = true
			ids = append(ids, id)
		}
	}
	return ids
}

func extractWorkstreamsFromBranch(branch string) []string {
	// Parse workstream IDs from branch names like feature/F031-something or ws/00-031-01
	wsRe := regexp.MustCompile(`00-\d{3}-\d{2}`)
	if matches := wsRe.FindAllString(branch, -1); len(matches) > 0 {
		return matches
	}
	return nil
}

// collectTestResults runs go test with -count=1 -cover and parses JSON output.
func collectTestResults(repoRoot string) ([]GateResult, float64) {
	cmd := exec.Command("go", "test", "./...", "-count=1", "-cover", "-json")
	cmd.Dir = repoRoot
	out, err := cmd.Output()

	passed := 0
	failed := 0
	totalCoverage := 0.0
	coverageCount := 0

	for _, line := range strings.Split(string(out), "\n") {
		if strings.TrimSpace(line) == "" {
			continue
		}
		var evt map[string]any
		if json.Unmarshal([]byte(line), &evt) != nil {
			continue
		}
		action, _ := evt["Action"].(string)
		switch action {
		case "pass":
			if _, hasTest := evt["Test"]; hasTest {
				passed++
			}
		case "fail":
			if _, hasTest := evt["Test"]; hasTest {
				failed++
			}
		}
		// Package-level coverage output appears in "output" lines
		if action == "output" {
			output, _ := evt["Output"].(string)
			if pct := parseCoverageLine(output); pct >= 0 {
				totalCoverage += pct
				coverageCount++
			}
		}
	}

	status := "pass"
	if err != nil || failed > 0 {
		status = "fail"
	}

	avgCoverage := -1.0
	if coverageCount > 0 {
		avgCoverage = totalCoverage / float64(coverageCount)
	}

	return []GateResult{{
		Name:   "go-test",
		Status: fmt.Sprintf("%s (%d passed, %d failed)", status, passed, failed),
	}}, avgCoverage
}

// parseCoverageLine extracts coverage percentage from a line like:
// "ok  	sdp_dev/internal/evidence	2.481s	coverage: 85.3% of statements"
func parseCoverageLine(line string) float64 {
	re := regexp.MustCompile(`coverage:\s+([\d.]+)%`)
	m := re.FindStringSubmatch(line)
	if m == nil {
		return -1
	}
	pct, err := strconv.ParseFloat(m[1], 64)
	if err != nil {
		return -1
	}
	return pct
}

// collectLintResults runs go vet and golangci-lint if available.
func collectLintResults(repoRoot string) []GateResult {
	var results []GateResult

	// Always run go vet
	cmd := exec.Command("go", "vet", "./...")
	cmd.Dir = repoRoot
	vetOut, vetErr := cmd.CombinedOutput()
	vetStatus := "pass"
	if vetErr != nil {
		vetStatus = fmt.Sprintf("fail: %s", strings.TrimSpace(string(vetOut)))
	}
	results = append(results, GateResult{Name: "go-vet", Status: vetStatus})

	// Run golangci-lint if available
	lintPath, err := exec.LookPath("golangci-lint")
	if err == nil {
		lintCmd := exec.Command(lintPath, "run", "--out-format=line-number", "--timeout=120s", "./...")
		lintCmd.Dir = repoRoot
		lintOut, lintErr := lintCmd.CombinedOutput()
		lintStatus := "pass"
		if lintErr != nil {
			lines := countNonEmptyLines(string(lintOut))
			lintStatus = fmt.Sprintf("fail (%d issues)", lines)
		}
		results = append(results, GateResult{Name: "golangci-lint", Status: lintStatus})
	}

	return results
}

// checkScopeCompliance checks changed files against declared workstream scope files.
// Returns a Boundary and whether it's compliant.
func checkScopeCompliance(repoRoot string, changedFiles []string) (Boundary, bool) {
	boundary := Boundary{
		Observed: ObservedBoundary{
			TouchedPaths: changedFiles,
		},
	}

	// Try to find declared scope from workstream files in the backlog
	declaredPrefixes := collectDeclaredScopePrefixes(repoRoot)

	if len(declaredPrefixes) == 0 {
		boundary.Compliance = BoundaryCompliance{
			OK:     true,
			Reason: "no declared scope — auto-attested from CI observation",
		}
		return boundary, true
	}

	boundary.Declared = DeclaredBoundary{AllowedPathPrefixes: declaredPrefixes}

	var outOfBoundary []string
	for _, f := range changedFiles {
		if !matchesAnyPrefix(f, declaredPrefixes) {
			outOfBoundary = append(outOfBoundary, f)
		}
	}

	boundary.Observed.OutOfBoundaryPaths = outOfBoundary

	if len(outOfBoundary) == 0 {
		boundary.Compliance = BoundaryCompliance{
			OK:     true,
			Reason: fmt.Sprintf("all %d changed files within declared scope (%d prefixes)", len(changedFiles), len(declaredPrefixes)),
		}
		return boundary, true
	}

	boundary.Compliance = BoundaryCompliance{
		OK:     false,
		Reason: fmt.Sprintf("%d files outside declared scope: %s", len(outOfBoundary), strings.Join(outOfBoundary, ", ")),
	}
	return boundary, false
}

// collectDeclaredScopePrefixes reads active workstream files and extracts scope paths.
func collectDeclaredScopePrefixes(repoRoot string) []string {
	backlogDir := filepath.Join(repoRoot, "docs", "workstreams", "backlog")
	entries, err := os.ReadDir(backlogDir)
	if err != nil {
		return nil
	}

	var prefixes []string
	seen := map[string]bool{}

	for _, e := range entries {
		if e.IsDir() || !strings.HasSuffix(e.Name(), ".md") {
			continue
		}
		f, err := os.Open(filepath.Join(backlogDir, e.Name()))
		if err != nil {
			continue
		}
		defer f.Close() //nolint:gocritic // defer in loop is acceptable here
		inScopeSection := false
		scanner := bufio.NewScanner(f)
		for scanner.Scan() {
			line := scanner.Text()
			if strings.HasPrefix(line, "## Scope Files") {
				inScopeSection = true
				continue
			}
			if inScopeSection && strings.HasPrefix(line, "##") {
				break
			}
			if inScopeSection && strings.HasPrefix(line, "- ") {
				path := strings.TrimPrefix(line, "- ")
				path = strings.TrimSpace(strings.Trim(path, "`"))
				if path != "" && !seen[path] {
					seen[path] = true
					prefixes = append(prefixes, path)
				}
			}
		}
	}
	return prefixes
}

func matchesAnyPrefix(file string, prefixes []string) bool {
	for _, p := range prefixes {
		if strings.HasPrefix(file, p) || file == p {
			return true
		}
	}
	return false
}

func countNonEmptyLines(s string) int {
	count := 0
	for _, line := range strings.Split(s, "\n") {
		if strings.TrimSpace(line) != "" {
			count++
		}
	}
	return count
}

func runGit(dir string, args ...string) (string, error) {
	cmd := exec.Command("git", args...)
	cmd.Dir = dir
	out, err := cmd.Output()
	if err != nil {
		return "", fmt.Errorf("git %s: %w", strings.Join(args, " "), err)
	}
	return string(out), nil
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

func firstOrEmpty(s []string) string {
	if len(s) > 0 {
		return s[0]
	}
	return ""
}

func minLen(a, b int) int {
	if a < b {
		return a
	}
	return b
}

// WriteAutoAttestationReport writes a human-readable summary JSON alongside the attestation.
func WriteAutoAttestationReport(outputPath string, stmt CodingWorkflowStatement) error {
	allTestsPass := true
	for _, t := range stmt.Predicate.Verification.Tests {
		if strings.HasPrefix(t.Status, "fail") {
			allTestsPass = false
		}
	}
	allLintPass := true
	for _, l := range stmt.Predicate.Verification.Lint {
		if strings.HasPrefix(l.Status, "fail") {
			allLintPass = false
		}
	}

	report := map[string]any{
		"type":             "ci-auto-attestation",
		"generated_at":     stmt.Predicate.Provenance.CapturedAt,
		"attestation_id":   stmt.Predicate.Provenance.RunID,
		"branch":           stmt.Predicate.Trace.Branch,
		"head_commit":      firstOrEmpty(stmt.Predicate.Trace.Commits),
		"beads_ids":        stmt.Predicate.Trace.BeadsIDs,
		"changed_files":    len(stmt.Predicate.Execution.ChangedFiles),
		"test_results":     stmt.Predicate.Verification.Tests,
		"all_tests_pass":   allTestsPass,
		"lint_results":     stmt.Predicate.Verification.Lint,
		"all_lint_pass":    allLintPass,
		"scope_compliance": stmt.Predicate.Boundary.Compliance,
		"out_of_scope":     stmt.Predicate.Boundary.Observed.OutOfBoundaryPaths,
	}
	if stmt.Predicate.Verification.Coverage != nil {
		report["coverage_pct"] = stmt.Predicate.Verification.Coverage.Value
		report["coverage_threshold"] = stmt.Predicate.Verification.Coverage.Threshold
		report["coverage_ok"] = stmt.Predicate.Verification.Coverage.Value >= stmt.Predicate.Verification.Coverage.Threshold
	}

	b, err := json.MarshalIndent(report, "", "  ")
	if err != nil {
		return err
	}
	b = append(b, '\n')
	return os.WriteFile(outputPath, b, 0o644)
}
