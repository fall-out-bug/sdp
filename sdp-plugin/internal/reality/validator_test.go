package reality

import (
	"os"
	"path/filepath"
	"testing"
)

func TestValidateOSS_SucceedsOnEmittedArtifacts(t *testing.T) {
	projectRoot := t.TempDir()
	seedProject(t, projectRoot)

	if _, err := EmitOSSWithOptions(projectRoot, Options{Mode: ModeDeep}); err != nil {
		t.Fatalf("EmitOSSWithOptions failed: %v", err)
	}

	validated, err := ValidateOSS(projectRoot)
	if err != nil {
		t.Fatalf("ValidateOSS failed: %v", err)
	}
	if len(validated) != len(requiredOSSValidationArtifacts) {
		t.Fatalf("expected %d validated artifacts, got %d", len(requiredOSSValidationArtifacts), len(validated))
	}
}

func TestValidateOSS_SucceedsOnEmptyRepository(t *testing.T) {
	projectRoot := t.TempDir()

	if _, err := EmitOSSWithOptions(projectRoot, Options{Mode: ModeDeep}); err != nil {
		t.Fatalf("EmitOSSWithOptions failed: %v", err)
	}

	validated, err := ValidateOSS(projectRoot)
	if err != nil {
		t.Fatalf("ValidateOSS failed for empty repository: %v", err)
	}
	if len(validated) != len(requiredOSSValidationArtifacts) {
		t.Fatalf("expected %d validated artifacts, got %d", len(requiredOSSValidationArtifacts), len(validated))
	}
}

func TestValidateOSS_FailsOnInvalidArtifact(t *testing.T) {
	projectRoot := t.TempDir()
	seedProject(t, projectRoot)

	if _, err := EmitOSS(projectRoot); err != nil {
		t.Fatalf("EmitOSS failed: %v", err)
	}
	if err := os.WriteFile(filepath.Join(projectRoot, ".sdp", "reality", "quality-report.json"), []byte("{bad json"), 0o644); err != nil {
		t.Fatalf("corrupt quality report: %v", err)
	}

	issues, err := ValidateOSS(projectRoot)
	if err == nil {
		t.Fatal("expected ValidateOSS to fail")
	}
	if len(issues) == 0 {
		t.Fatal("expected validation issues")
	}
}
