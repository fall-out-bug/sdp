package guard

import (
	"os"
	"path/filepath"
	"strings"
	"testing"

	"github.com/fall-out-bug/sdp/internal/config"
)

// TestCheckMaxFileLOC tests the max file LOC rule
func TestCheckMaxFileLOC(t *testing.T) {
	tests := []struct {
		name        string
		content     string
		maxLines    int
		severity    string
		wantFinding bool
		wantCount   int
	}{
		{
			name:        "Small file - no finding",
			content:     "package main\n\nfunc main() {}\n",
			maxLines:    200,
			severity:    "error",
			wantFinding: false,
		},
		{
			name:        "File exactly at limit - no finding",
			content:     generateLines(200),
			maxLines:    200,
			severity:    "error",
			wantFinding: false,
		},
		{
			name:        "File exceeds limit - error finding",
			content:     generateLines(201),
			maxLines:    200,
			severity:    "error",
			wantFinding: true,
			wantCount:   1,
		},
		{
			name:        "File exceeds limit - warning finding",
			content:     generateLines(201),
			maxLines:    200,
			severity:    "warning",
			wantFinding: true,
			wantCount:   1,
		},
		{
			name:        "Custom limit - file exceeds custom limit",
			content:     generateLines(151),
			maxLines:    150,
			severity:    "error",
			wantFinding: true,
			wantCount:   1,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			configDir := t.TempDir()
			skill := NewSkill(configDir)

			rule := config.GuardRule{
				ID:       "max-file-loc",
				Enabled:  true,
				Severity: tt.severity,
				Config: map[string]interface{}{
					"max_lines": tt.maxLines,
				},
			}

			findings := skill.checkMaxFileLOC("test.go", []byte(tt.content), rule)

			if tt.wantFinding && len(findings) == 0 {
				t.Error("Expected finding, got none")
			}

			if !tt.wantFinding && len(findings) > 0 {
				t.Errorf("Expected no finding, got %d", len(findings))
			}

			if tt.wantCount > 0 && len(findings) != tt.wantCount {
				t.Errorf("Expected %d finding(s), got %d", tt.wantCount, len(findings))
			}

			// Verify severity
			if len(findings) > 0 {
				expectedSeverity := SeverityError
				if tt.severity == "warning" {
					expectedSeverity = SeverityWarning
				}
				if findings[0].Severity != expectedSeverity {
					t.Errorf("Severity = %s, want %s", findings[0].Severity, expectedSeverity)
				}
			}
		})
	}
}

// TestCheckMaxFileLOCDefaultLimit tests default limit (200 lines)
func TestCheckMaxFileLOCDefaultLimit(t *testing.T) {
	configDir := t.TempDir()
	skill := NewSkill(configDir)

	// Rule without max_lines config - should use default 200
	rule := config.GuardRule{
		ID:       "max-file-loc",
		Enabled:  true,
		Severity: "error",
		Config:   map[string]interface{}{},
	}

	// 201 lines should trigger finding
	content := generateLines(201)
	findings := skill.checkMaxFileLOC("test.go", []byte(content), rule)

	if len(findings) != 1 {
		t.Errorf("Expected 1 finding, got %d", len(findings))
	}

	// 200 lines should not trigger finding
	content = generateLines(200)
	findings = skill.checkMaxFileLOC("test.go", []byte(content), rule)

	if len(findings) != 0 {
		t.Errorf("Expected 0 findings, got %d", len(findings))
	}
}

// TestCheckCommentedCode tests the commented code detection rule
func TestCheckCommentedCode(t *testing.T) {
	tests := []struct {
		name        string
		content     string
		fileName    string
		severity    string
		wantFinding bool
		wantCount   int
	}{
		{
			name:        "No commented code",
			content:     "package main\n\nfunc main() {\n\tprintln(\"hello\")\n}\n",
			fileName:    "test.go",
			severity:    "error",
			wantFinding: false,
		},
		{
			name:        "Single comment - no finding",
			content:     "// This is a comment\npackage main\n",
			fileName:    "test.go",
			severity:    "error",
			wantFinding: false,
		},
		{
			name:        "Two consecutive comments - no finding",
			content:     "// Comment 1\n// Comment 2\npackage main\n",
			fileName:    "test.go",
			severity:    "error",
			wantFinding: false,
		},
		{
			name:        "Three consecutive comments - finding",
			content:     "// Looks like code\n// var x = 1\n// var y = 2\npackage main\n",
			fileName:    "test.go",
			severity:    "error",
			wantFinding: true,
			wantCount:   1,
		},
		{
			name:        "TODO comment - ignored",
			content:     "// TODO(WS-063-01): implement this\n// var x = 1\n// var y = 2\npackage main\n",
			fileName:    "test.go",
			severity:    "error",
			wantFinding: false,
		},
		{
			name:        "FIXME comment - ignored",
			content:     "// FIXME: bug here\n// var x = 1\n// var y = 2\npackage main\n",
			fileName:    "test.go",
			severity:    "error",
			wantFinding: false,
		},
		{
			name:        "NOTE comment - ignored",
			content:     "// NOTE: important\n// var x = 1\n// var y = 2\npackage main\n",
			fileName:    "test.go",
			severity:    "error",
			wantFinding: false,
		},
		{
			name:        "Python file with # comments",
			content:     "# x = 12345\n# y = 23456\n# z = 34567\nprint('hello')\n",
			fileName:    "test.py",
			severity:    "warning",
			wantFinding: true,
			wantCount:   1,
		},
		{
			name:        "Short comment ignored",
			content:     "// short\n// var x = 1\n// var y = 2\npackage main\n",
			fileName:    "test.go",
			severity:    "error",
			wantFinding: false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			configDir := t.TempDir()
			skill := NewSkill(configDir)

			rule := config.GuardRule{
				ID:       "no-commented-code",
				Enabled:  true,
				Severity: tt.severity,
				Config:   map[string]interface{}{},
			}

			findings := skill.checkCommentedCode(tt.fileName, []byte(tt.content), rule)

			if tt.wantFinding && len(findings) == 0 {
				t.Errorf("Expected finding, got none for: %s", tt.content)
			}

			if !tt.wantFinding && len(findings) > 0 {
				t.Errorf("Expected no finding, got %d", len(findings))
			}

			if tt.wantCount > 0 && len(findings) != tt.wantCount {
				t.Errorf("Expected %d finding(s), got %d", tt.wantCount, len(findings))
			}

			// Verify severity
			if len(findings) > 0 {
				expectedSeverity := SeverityError
				if tt.severity == "warning" {
					expectedSeverity = SeverityWarning
				}
				if findings[0].Severity != expectedSeverity {
					t.Errorf("Severity = %s, want %s", findings[0].Severity, expectedSeverity)
				}
			}
		})
	}
}

// TestCheckOrphanedTODOs tests the orphaned TODO detection rule
func TestCheckOrphanedTODOs(t *testing.T) {
	tests := []struct {
		name        string
		content     string
		severity    string
		wantFinding bool
		wantCount   int
	}{
		{
			name:        "No TODOs",
			content:     "package main\n\nfunc main() {}\n",
			severity:    "error",
			wantFinding: false,
		},
		{
			name:        "TODO with valid WS ID",
			content:     "// TODO(WS-063-01): implement feature\npackage main\n",
			severity:    "error",
			wantFinding: false,
		},
		{
			name:        "TODO with lowercase ws ID",
			content:     "// TODO(ws-063-01): implement feature\npackage main\n",
			severity:    "error",
			wantFinding: false,
		},
		{
			name:        "TODO without WS ID - error",
			content:     "// TODO: implement this\npackage main\n",
			severity:    "error",
			wantFinding: true,
			wantCount:   1,
		},
		{
			name:        "TODO without WS ID - warning",
			content:     "// TODO: implement this\npackage main\n",
			severity:    "warning",
			wantFinding: true,
			wantCount:   1,
		},
		{
			name:        "Multiple TODOs without WS ID",
			content:     "// TODO: fix this\n// TODO: fix that\npackage main\n",
			severity:    "error",
			wantFinding: true,
			wantCount:   2,
		},
		{
			name:        "TODO with malformed WS ID",
			content:     "// TODO(WS-123): incomplete ID\npackage main\n",
			severity:    "error",
			wantFinding: true,
			wantCount:   1,
		},
		{
			name:        "Mixed - some with WS ID, some without",
			content:     "// TODO(WS-063-01): valid\n// TODO: invalid\npackage main\n",
			severity:    "error",
			wantFinding: true,
			wantCount:   1,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			configDir := t.TempDir()
			skill := NewSkill(configDir)

			rule := config.GuardRule{
				ID:       "no-orphaned-todos",
				Enabled:  true,
				Severity: tt.severity,
				Config:   map[string]interface{}{},
			}

			findings := skill.checkOrphanedTODOs("test.go", []byte(tt.content), rule)

			if tt.wantFinding && len(findings) == 0 {
				t.Errorf("Expected finding, got none for: %s", tt.content)
			}

			if !tt.wantFinding && len(findings) > 0 {
				t.Errorf("Expected no finding, got %d", len(findings))
			}

			if tt.wantCount > 0 && len(findings) != tt.wantCount {
				t.Errorf("Expected %d finding(s), got %d", tt.wantCount, len(findings))
			}

			// Verify severity
			if len(findings) > 0 {
				expectedSeverity := SeverityError
				if tt.severity == "warning" {
					expectedSeverity = SeverityWarning
				}
				if findings[0].Severity != expectedSeverity {
					t.Errorf("Severity = %s, want %s", findings[0].Severity, expectedSeverity)
				}
			}
		})
	}
}

// TestApplyGuardRules tests the applyGuardRules function
func TestApplyGuardRules(t *testing.T) {
	tests := []struct {
		name            string
		files           []string
		rules           *config.GuardRules
		wantMinFindings int
	}{
		{
			name:  "No files - no findings",
			files: []string{},
			rules: &config.GuardRules{
				Version: 1,
				Rules: []config.GuardRule{
					{
						ID:       "max-file-loc",
						Enabled:  true,
						Severity: "error",
						Config:   map[string]interface{}{"max_lines": 200},
					},
				},
			},
			wantMinFindings: 0,
		},
		{
			name:  "File exceeds max LOC - finding",
			files: []string{"large.go"},
			rules: &config.GuardRules{
				Version: 1,
				Rules: []config.GuardRule{
					{
						ID:       "max-file-loc",
						Enabled:  true,
						Severity: "error",
						Config:   map[string]interface{}{"max_lines": 10},
					},
				},
			},
			wantMinFindings: 1,
		},
		{
			name:  "Disabled rule - no findings",
			files: []string{"large.go"},
			rules: &config.GuardRules{
				Version: 1,
				Rules: []config.GuardRule{
					{
						ID:       "max-file-loc",
						Enabled:  false,
						Severity: "error",
						Config:   map[string]interface{}{"max_lines": 10},
					},
				},
			},
			wantMinFindings: 0,
		},
		{
			name:  "Commented code detection",
			files: []string{"commented.go"},
			rules: &config.GuardRules{
				Version: 1,
				Rules: []config.GuardRule{
					{
						ID:       "no-commented-code",
						Enabled:  true,
						Severity: "warning",
					},
				},
			},
			wantMinFindings: 1,
		},
		{
			name:  "Orphaned TODO detection",
			files: []string{"todo.go"},
			rules: &config.GuardRules{
				Version: 1,
				Rules: []config.GuardRule{
					{
						ID:       "no-orphaned-todos",
						Enabled:  true,
						Severity: "error",
					},
				},
			},
			wantMinFindings: 1,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			configDir := t.TempDir()
			skill := NewSkill(configDir)

			// Create test files and build absolute paths list
			absFiles := []string{}
			for _, fileName := range tt.files {
				var content string
				switch fileName {
				case "large.go":
					content = generateLines(201)
				case "commented.go":
					content = "// var x = 123456\n// var y = 234567\n// var z = 345678\npackage main\n"
				case "todo.go":
					content = "// TODO: implement this\npackage main\n"
				default:
					content = "package main\n"
				}
				absPath := filepath.Join(configDir, fileName)
				if err := os.WriteFile(absPath, []byte(content), 0644); err != nil {
					t.Fatalf("Failed to create test file: %v", err)
				}
				absFiles = append(absFiles, absPath)
			}

			findings := skill.applyGuardRules(absFiles, tt.rules)

			if len(findings) < tt.wantMinFindings {
				t.Errorf("Expected at least %d finding(s), got %d", tt.wantMinFindings, len(findings))
			}
		})
	}
}

// TestApplyGuardRulesSkipsCoverageAndComplexity tests that coverage and complexity rules are skipped
func TestApplyGuardRulesSkipsCoverageAndComplexity(t *testing.T) {
	configDir := t.TempDir()
	skill := NewSkill(configDir)

	rules := &config.GuardRules{
		Version: 1,
		Rules: []config.GuardRule{
			{
				ID:       "coverage-threshold",
				Enabled:  true,
				Severity: "error",
			},
			{
				ID:       "max-cyclomatic-complexity",
				Enabled:  true,
				Severity: "error",
			},
		},
	}

	// Create test file
	content := "package main\n\nfunc main() {}\n"
	absPath := filepath.Join(configDir, "test.go")
	if err := os.WriteFile(absPath, []byte(content), 0644); err != nil {
		t.Fatalf("Failed to create test file: %v", err)
	}

	findings := skill.applyGuardRules([]string{absPath}, rules)

	// Coverage and complexity rules should be skipped (no findings)
	if len(findings) != 0 {
		t.Errorf("Coverage and complexity rules should be skipped, got %d findings", len(findings))
	}
}

// Helper function to generate a file with a specific number of lines
func generateLines(n int) string {
	var lines []string
	for i := 0; i < n; i++ {
		lines = append(lines, "// Line "+string(rune('A'+(i%26))))
	}
	return strings.Join(lines, "\n")
}
