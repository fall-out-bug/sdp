package evidenceenv

import (
	"crypto/sha256"
	"encoding/hex"
	"encoding/json"
	"fmt"
	"os"
	"strings"

	intoto "github.com/in-toto/in-toto-golang/in_toto"
)

const (
	PredicateTypeCodingWorkflow = "https://sdp.dev/attestation/coding-workflow/v1"
	StatementType               = intoto.StatementInTotoV01
)

type CodingWorkflowStatement struct {
	intoto.StatementHeader
	Predicate CodingWorkflowPredicate `json:"predicate"`
}

type CodingWorkflowPredicate struct {
	Intent       Intent       `json:"intent"`
	Plan         Plan         `json:"plan"`
	Execution    Execution    `json:"execution"`
	Verification Verification `json:"verification"`
	Review       Review       `json:"review"`
	RiskNotes    RiskNotes    `json:"risk_notes"`
	Boundary     Boundary     `json:"boundary"`
	Provenance   Provenance   `json:"provenance"`
	Trace        Trace        `json:"trace"`
}

type Intent struct {
	IssueID            string   `json:"issue_id"`
	Trigger            string   `json:"trigger"`
	AcceptanceCriteria []string `json:"acceptance_criteria"`
	RiskClass          string   `json:"risk_class"`
}

type Plan struct {
	Workstreams       []string `json:"workstreams"`
	OrderingRationale string   `json:"ordering_rationale"`
}

type Execution struct {
	ClaimedIssueIDs []string `json:"claimed_issue_ids"`
	Branch          string   `json:"branch"`
	ChangedFiles    []string `json:"changed_files"`
}

type Verification struct {
	Tests    []GateResult `json:"tests"`
	Lint     []GateResult `json:"lint"`
	Coverage *Coverage    `json:"coverage,omitempty"`
}

type GateResult struct {
	Name   string `json:"name"`
	Status string `json:"status"`
}

type Coverage struct {
	Value     float64 `json:"value"`
	Threshold float64 `json:"threshold"`
}

type Review struct {
	SelfReview       []ReviewItem `json:"self_review"`
	AdversarialItems []ReviewItem `json:"adversarial_review"`
}

type ReviewItem struct {
	Reviewer string `json:"reviewer"`
	Verdict  string `json:"verdict"`
	Notes    string `json:"notes,omitempty"`
}

type RiskNotes struct {
	ResidualRisks []string `json:"residual_risks"`
	OutOfScope    []string `json:"out_of_scope"`
}

type Boundary struct {
	Declared   DeclaredBoundary   `json:"declared"`
	Observed   ObservedBoundary   `json:"observed"`
	Compliance BoundaryCompliance `json:"compliance"`
}

type DeclaredBoundary struct {
	AllowedPathPrefixes   []string `json:"allowed_path_prefixes"`
	ControlPathPrefixes   []string `json:"control_path_prefixes"`
	ForbiddenPathPrefixes []string `json:"forbidden_path_prefixes"`
}

type ObservedBoundary struct {
	TouchedPaths       []string `json:"touched_paths"`
	OutOfBoundaryPaths []string `json:"out_of_boundary_paths"`
}

type BoundaryCompliance struct {
	OK     bool   `json:"ok"`
	Reason string `json:"reason"`
}

type Provenance struct {
	RunID          string          `json:"run_id"`
	Orchestrator   string          `json:"orchestrator"`
	Runtime        string          `json:"runtime"`
	Model          string          `json:"model"`
	Phase          string          `json:"phase"`
	Role           string          `json:"role"`
	CapturedAt     string          `json:"captured_at"`
	SourceIssueID  string          `json:"source_issue_id"`
	PromptHash     string          `json:"prompt_hash,omitempty"`
	ContextSources []ContextSource `json:"context_sources,omitempty"`
}

type ContextSource struct {
	Type string `json:"type"`
	Path string `json:"path"`
	Hash string `json:"hash"`
}

type Trace struct {
	BeadsIDs []string `json:"beads_ids"`
	Branch   string   `json:"branch"`
	Commits  []string `json:"commits"`
	PRURL    string   `json:"pr_url"`
}

func NewStatement(subjects []intoto.Subject, predicate CodingWorkflowPredicate) CodingWorkflowStatement {
	return CodingWorkflowStatement{
		StatementHeader: intoto.StatementHeader{
			Type:          StatementType,
			PredicateType: PredicateTypeCodingWorkflow,
			Subject:       subjects,
		},
		Predicate: predicate,
	}
}

func WriteAttestation(path string, stmt CodingWorkflowStatement) error {
	b, err := json.MarshalIndent(stmt, "", "  ")
	if err != nil {
		return fmt.Errorf("marshal attestation: %w", err)
	}
	b = append(b, '\n')
	return os.WriteFile(path, b, 0o644)
}

func ReadAttestation(path string) (CodingWorkflowStatement, error) {
	b, err := os.ReadFile(path)
	if err != nil {
		return CodingWorkflowStatement{}, err
	}
	var stmt CodingWorkflowStatement
	if err := json.Unmarshal(b, &stmt); err != nil {
		return CodingWorkflowStatement{}, fmt.Errorf("parse attestation: %w", err)
	}
	return stmt, nil
}

func ValidateAttestation(stmt CodingWorkflowStatement, requirePRURL bool) Result {
	if stmt.Type != StatementType {
		return Result{OK: false, Reason: fmt.Sprintf("invalid statement type: %s (expected %s)", stmt.Type, StatementType)}
	}
	if stmt.PredicateType != PredicateTypeCodingWorkflow {
		return Result{OK: false, Reason: fmt.Sprintf("invalid predicate type: %s (expected %s)", stmt.PredicateType, PredicateTypeCodingWorkflow)}
	}
	if len(stmt.Subject) == 0 {
		return Result{OK: false, Reason: "no subjects in statement"}
	}

	p := stmt.Predicate

	if strings.TrimSpace(p.Intent.IssueID) == "" {
		return Result{OK: false, Reason: "missing intent.issue_id"}
	}
	if !p.Boundary.Compliance.OK && p.Boundary.Compliance.Reason == "" {
		return Result{OK: false, Reason: "boundary compliance failed with no reason"}
	}
	if strings.TrimSpace(p.Provenance.RunID) == "" {
		return Result{OK: false, Reason: "missing provenance.run_id"}
	}
	if strings.TrimSpace(p.Provenance.CapturedAt) == "" {
		return Result{OK: false, Reason: "missing provenance.captured_at"}
	}

	if p.Provenance.PromptHash != "" && !isSHA256Hex(p.Provenance.PromptHash) {
		return Result{OK: false, Reason: "invalid provenance.prompt_hash: not SHA-256 hex"}
	}
	for _, cs := range p.Provenance.ContextSources {
		if cs.Type == "" || cs.Path == "" || cs.Hash == "" {
			return Result{OK: false, Reason: "context_source missing type, path, or hash"}
		}
		if !isSHA256Hex(cs.Hash) {
			return Result{OK: false, Reason: fmt.Sprintf("context_source hash not SHA-256 hex: %s", cs.Path)}
		}
	}

	if requirePRURL && strings.TrimSpace(p.Trace.PRURL) == "" {
		return Result{OK: false, Reason: "missing trace.pr_url"}
	}

	return Result{OK: true, Reason: "ok"}
}

func ValidateAttestationFile(path string, requirePRURL bool) (Result, error) {
	stmt, err := ReadAttestation(path)
	if err != nil {
		return Result{}, err
	}
	return ValidateAttestation(stmt, requirePRURL), nil
}

func isSHA256Hex(s string) bool {
	if len(s) != 64 {
		return false
	}
	_, err := hex.DecodeString(s)
	return err == nil
}

func DigestOfBytes(b []byte) string {
	h := sha256.Sum256(b)
	return hex.EncodeToString(h[:])
}
