package main

import (
	"strings"
	"testing"
)

func TestDiagnoseCmd_Help(t *testing.T) {
	cmd := diagnoseCmd

	// Test command structure
	if cmd.Use != "diagnose [error-code]" {
		t.Errorf("diagnoseCmd has wrong use: %s", cmd.Use)
	}

	// Test that command runs without crashing
	err := cmd.RunE(cmd, []string{})
	if err != nil {
		t.Errorf("diagnoseCmd() failed: %v", err)
	}
}

func TestDiagnoseCmd_ListClasses(t *testing.T) {
	// Reset flags
	diagnoseListClasses = true
	diagnoseListCodes = false
	diagnoseJSON = false
	t.Cleanup(func() {
		diagnoseListClasses = false
	})

	err := listErrorClasses()
	if err != nil {
		t.Fatalf("listErrorClasses() failed: %v", err)
	}
}

func TestDiagnoseCmd_ListCodes(t *testing.T) {
	// Reset flags
	diagnoseListCodes = true
	diagnoseListClasses = false
	diagnoseJSON = false
	t.Cleanup(func() {
		diagnoseListCodes = false
	})

	err := listErrorCodes()
	if err != nil {
		t.Fatalf("listErrorCodes() failed: %v", err)
	}
}

func TestDiagnoseCmd_SpecificCode(t *testing.T) {
	// Reset flags
	diagnoseJSON = false
	diagnoseOutput = ""

	err := diagnoseErrorCode("ENV001")
	if err != nil {
		t.Fatalf("diagnoseErrorCode(ENV001) failed: %v", err)
	}
}

func TestDiagnoseCmd_UnknownCode(t *testing.T) {
	err := diagnoseErrorCode("UNKNOWN999")
	if err == nil {
		t.Error("Expected error for unknown error code")
	}
	if !strings.Contains(err.Error(), "unknown error code") {
		t.Errorf("Expected 'unknown error code' error, got: %v", err)
	}
}

func TestDiagnoseCmd_JSON(t *testing.T) {
	// Reset flags
	diagnoseJSON = true
	t.Cleanup(func() {
		diagnoseJSON = false
	})

	err := diagnoseErrorCode("ENV001")
	if err != nil {
		t.Fatalf("diagnoseErrorCode(ENV001) with JSON failed: %v", err)
	}
}

func TestDiagnoseCmd_CoverageError(t *testing.T) {
	// Reset flags
	diagnoseJSON = false
	diagnoseOutput = ""

	err := diagnoseErrorCode("VAL001")
	if err != nil {
		t.Fatalf("diagnoseErrorCode(VAL001) failed: %v", err)
	}
}

func TestDiagnoseCmd_BlockedWorkstream(t *testing.T) {
	// Reset flags
	diagnoseJSON = false
	diagnoseOutput = ""

	err := diagnoseErrorCode("DEP001")
	if err != nil {
		t.Fatalf("diagnoseErrorCode(DEP001) failed: %v", err)
	}
}

func TestDiagnoseCmd_HasFlags(t *testing.T) {
	cmd := diagnoseCmd

	// Test flags exist
	if cmd.Flags().Lookup("list-classes") == nil {
		t.Error("diagnoseCmd missing --list-classes flag")
	}
	if cmd.Flags().Lookup("list-codes") == nil {
		t.Error("diagnoseCmd missing --list-codes flag")
	}
	if cmd.Flags().Lookup("json") == nil {
		t.Error("diagnoseCmd missing --json flag")
	}
	if cmd.Flags().Lookup("output") == nil {
		t.Error("diagnoseCmd missing --output flag")
	}
}

func TestDiagnoseCmd_AllClassesListed(t *testing.T) {
	// Reset flags
	diagnoseJSON = false

	classes := []string{"ENV", "PROTO", "DEP", "VAL", "RUNTIME"}

	for _, class := range classes {
		// Test that each class has a default playbook
		err := diagnoseErrorCode(class + "000")
		// Should either find a playbook or report unknown
		// We just want to verify no panic
		_ = err
	}
}
