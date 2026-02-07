package main

import (
	"os"
	"testing"
)

// TestPrdDetectTypeCmd tests the prd detect-type command
func TestPrdDetectTypeCmd(t *testing.T) {
	tests := []struct {
		name         string
		createFile   string
		expectedType string
	}{
		{
			name:         "go project",
			createFile:   "go.mod",
			expectedType: "go",
		},
		{
			name:         "python project",
			createFile:   "pyproject.toml",
			expectedType: "python",
		},
		{
			name:         "java maven project",
			createFile:   "pom.xml",
			expectedType: "java",
		},
		{
			name:         "java gradle project",
			createFile:   "build.gradle",
			expectedType: "java",
		},
		{
			name:         "default library project",
			createFile:   "",
			expectedType: "library",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			tmpDir := t.TempDir()

			// Create test file if specified
			if tt.createFile != "" {
				filePath := tmpDir + "/" + tt.createFile
				if err := os.WriteFile(filePath, []byte("test"), 0644); err != nil {
					t.Fatalf("Failed to create test file: %v", err)
				}
			}

			// Change to temp directory
			originalWd, _ := os.Getwd()
			t.Cleanup(func() { os.Chdir(originalWd) })
			if err := os.Chdir(tmpDir); err != nil {
				t.Fatalf("Failed to chdir: %v", err)
			}

			cmd := prdDetectType()
			err := cmd.RunE(cmd, []string{})

			// Should succeed
			if err != nil {
				t.Errorf("prdDetectType() failed: %v", err)
			}
		})
	}
}

// TestPrdDetectTypeWithPath tests the prd detect-type command with a path argument
func TestPrdDetectTypeWithPath(t *testing.T) {
	tmpDir := t.TempDir()

	// Create a go.mod file
	if err := os.WriteFile(tmpDir+"/go.mod", []byte("module test"), 0644); err != nil {
		t.Fatalf("Failed to create go.mod: %v", err)
	}

	// Change to temp directory
	originalWd, _ := os.Getwd()
	t.Cleanup(func() { os.Chdir(originalWd) })
	if err := os.Chdir(tmpDir); err != nil {
		t.Fatalf("Failed to chdir: %v", err)
	}

	cmd := prdDetectType()
	err := cmd.RunE(cmd, []string{"."})

	// Should succeed
	if err != nil {
		t.Errorf("prdDetectType() with path failed: %v", err)
	}
}
