package main

import (
	"testing"
)

func TestHasSuffix(t *testing.T) {
	tests := []struct {
		s    string
		suf  string
		want bool
	}{
		{"test.md", ".md", true},
		{"test.txt", ".md", false},
		{"test.md", "md", true},
		{"test", ".md", false},
		{"", ".md", false},
		{"test.md", "", true},
		{"short", "longer", false},
		{"file.go", ".go", true},
		{"file_test.go", "_test.go", true},
	}

	for _, tt := range tests {
		t.Run(tt.s+"_"+tt.suf, func(t *testing.T) {
			got := hasSuffix(tt.s, tt.suf)
			if got != tt.want {
				t.Errorf("hasSuffix(%q, %q) = %v, want %v", tt.s, tt.suf, got, tt.want)
			}
		})
	}
}

func TestLoadInProgressScopes_EmptyDir(t *testing.T) {
	// Test with a path that doesn't have workstreams directory
	scopes, err := loadInProgressScopes("/nonexistent/path")
	if err != nil {
		t.Errorf("loadInProgressScopes should not return error for nonexistent path: %v", err)
	}
	if scopes != nil {
		t.Errorf("Expected nil scopes for nonexistent path, got %v", scopes)
	}
}

func TestScopeFilesForWS_NoProject(t *testing.T) {
	// This will fail to find project root from temp directory
	// Just verify it doesn't panic
	files := scopeFilesForWS("nonexistent-ws")
	// May return nil or empty slice depending on project root detection
	_ = files
}
