package sdpinit

import (
	"os"
	"testing"
)

func TestRunPreflight(t *testing.T) {
	tmpDir := t.TempDir()
	oldWd, _ := os.Getwd()
	os.Chdir(tmpDir)
	defer os.Chdir(oldWd)

	result := RunPreflight()

	if result == nil {
		t.Fatal("RunPreflight returned nil")
	}

	// Should have detected something
	if result.ProjectType == "" {
		t.Error("ProjectType should not be empty")
	}
}

func TestDetectProjectType_Unknown(t *testing.T) {
	tmpDir := t.TempDir()
	oldWd, _ := os.Getwd()
	os.Chdir(tmpDir)
	defer os.Chdir(oldWd)

	pt := DetectProjectType()
	if pt != "unknown" {
		t.Errorf("Expected unknown, got %s", pt)
	}
}

func TestDetectProjectType_Go(t *testing.T) {
	tmpDir := t.TempDir()
	oldWd, _ := os.Getwd()
	os.Chdir(tmpDir)
	defer os.Chdir(oldWd)

	os.WriteFile("go.mod", []byte("module test"), 0644)

	pt := DetectProjectType()
	if pt != "go" {
		t.Errorf("Expected go, got %s", pt)
	}
}

func TestDetectProjectType_Node(t *testing.T) {
	tmpDir := t.TempDir()
	oldWd, _ := os.Getwd()
	os.Chdir(tmpDir)
	defer os.Chdir(oldWd)

	os.WriteFile("package.json", []byte("{}"), 0644)

	pt := DetectProjectType()
	if pt != "node" {
		t.Errorf("Expected node, got %s", pt)
	}
}

func TestDetectProjectType_Python(t *testing.T) {
	tmpDir := t.TempDir()
	oldWd, _ := os.Getwd()
	os.Chdir(tmpDir)
	defer os.Chdir(oldWd)

	os.WriteFile("setup.py", []byte("# setup"), 0644)

	pt := DetectProjectType()
	if pt != "python" {
		t.Errorf("Expected python, got %s", pt)
	}
}

func TestDetectProjectType_Mixed(t *testing.T) {
	tmpDir := t.TempDir()
	oldWd, _ := os.Getwd()
	os.Chdir(tmpDir)
	defer os.Chdir(oldWd)

	os.WriteFile("go.mod", []byte("module test"), 0644)
	os.WriteFile("package.json", []byte("{}"), 0644)

	pt := DetectProjectType()
	if pt != "mixed" {
		t.Errorf("Expected mixed, got %s", pt)
	}
}

func TestIsGoProject(t *testing.T) {
	tests := []struct {
		name     string
		setup    func(string)
		expected bool
	}{
		{"no files", func(d string) {}, false},
		{"go.mod", func(d string) { os.WriteFile("go.mod", []byte{}, 0644) }, true},
		{"go file", func(d string) { os.WriteFile("main.go", []byte{}, 0644) }, true},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			tmpDir := t.TempDir()
			oldWd, _ := os.Getwd()
			os.Chdir(tmpDir)
			defer os.Chdir(oldWd)

			tt.setup(tmpDir)
			result := isGoProject()
			if result != tt.expected {
				t.Errorf("Expected %v, got %v", tt.expected, result)
			}
		})
	}
}

func TestIsNodeProject(t *testing.T) {
	tests := []struct {
		name     string
		setup    func()
		expected bool
	}{
		{"no files", func() {}, false},
		{"package.json", func() { os.WriteFile("package.json", []byte{}, 0644) }, true},
		{"ts file", func() { os.WriteFile("index.ts", []byte{}, 0644) }, true},
		{"js file", func() { os.WriteFile("index.js", []byte{}, 0644) }, true},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			tmpDir := t.TempDir()
			oldWd, _ := os.Getwd()
			os.Chdir(tmpDir)
			defer os.Chdir(oldWd)

			tt.setup()
			result := isNodeProject()
			if result != tt.expected {
				t.Errorf("Expected %v, got %v", tt.expected, result)
			}
		})
	}
}

func TestIsPythonProject(t *testing.T) {
	tests := []struct {
		name     string
		setup    func()
		expected bool
	}{
		{"no files", func() {}, false},
		{"setup.py", func() { os.WriteFile("setup.py", []byte{}, 0644) }, true},
		{"pyproject.toml", func() { os.WriteFile("pyproject.toml", []byte{}, 0644) }, true},
		{"requirements.txt", func() { os.WriteFile("requirements.txt", []byte{}, 0644) }, true},
		{"py file", func() { os.WriteFile("main.py", []byte{}, 0644) }, true},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			tmpDir := t.TempDir()
			oldWd, _ := os.Getwd()
			os.Chdir(tmpDir)
			defer os.Chdir(oldWd)

			tt.setup()
			result := isPythonProject()
			if result != tt.expected {
				t.Errorf("Expected %v, got %v", tt.expected, result)
			}
		})
	}
}

func TestDirExists(t *testing.T) {
	tmpDir := t.TempDir()
	oldWd, _ := os.Getwd()
	os.Chdir(tmpDir)
	defer os.Chdir(oldWd)

	os.Mkdir("testdir", 0755)
	os.WriteFile("testfile", []byte{}, 0644)

	if !dirExists("testdir") {
		t.Error("testdir should exist")
	}
	if dirExists("testfile") {
		t.Error("testfile is not a directory")
	}
	if dirExists("nonexistent") {
		t.Error("nonexistent should not exist")
	}
}

func TestCheckConflicts(t *testing.T) {
	tmpDir := t.TempDir()
	oldWd, _ := os.Getwd()
	os.Chdir(tmpDir)
	defer os.Chdir(oldWd)

	// No conflicts
	conflicts := checkConflicts()
	if len(conflicts) != 0 {
		t.Errorf("Expected no conflicts, got %v", conflicts)
	}

	// Create conflict
	os.MkdirAll(".claude", 0755)
	os.WriteFile(".claude/settings.json", []byte{}, 0644)

	conflicts = checkConflicts()
	if len(conflicts) != 1 {
		t.Errorf("Expected 1 conflict, got %v", conflicts)
	}
}

func TestPreflightResult_HasSDP(t *testing.T) {
	tmpDir := t.TempDir()
	oldWd, _ := os.Getwd()
	os.Chdir(tmpDir)
	defer os.Chdir(oldWd)

	result := RunPreflight()
	if result.HasSDP {
		t.Error("Should not have .sdp")
	}

	os.Mkdir(".sdp", 0755)
	result = RunPreflight()
	if !result.HasSDP {
		t.Error("Should have .sdp")
	}
}

func TestPreflightResult_HasGit(t *testing.T) {
	tmpDir := t.TempDir()
	oldWd, _ := os.Getwd()
	os.Chdir(tmpDir)
	defer os.Chdir(oldWd)

	result := RunPreflight()
	if result.HasGit {
		t.Error("Should not have .git")
	}

	os.Mkdir(".git", 0755)
	result = RunPreflight()
	if !result.HasGit {
		t.Error("Should have .git")
	}
}

func TestPreflightResult_Warnings(t *testing.T) {
	tmpDir := t.TempDir()
	oldWd, _ := os.Getwd()
	os.Chdir(tmpDir)
	defer os.Chdir(oldWd)

	// No git = warning
	result := RunPreflight()
	if len(result.Warnings) == 0 {
		t.Error("Expected warning about no git")
	}

	// With git = no warning
	os.Mkdir(".git", 0755)
	result = RunPreflight()
	hasGitWarning := false
	for _, w := range result.Warnings {
		if w == "Not a git repository - version control recommended" {
			hasGitWarning = true
		}
	}
	if hasGitWarning {
		t.Error("Should not have git warning when .git exists")
	}
}
