package sdpinit

import (
	"testing"
)

func TestGetDefaults_Go(t *testing.T) {
	defaults := GetDefaults("go")

	if defaults == nil {
		t.Fatal("GetDefaults returned nil for go")
	}

	// Check skills
	if len(defaults.Skills) == 0 {
		t.Error("Go defaults should have skills")
	}

	// Check specific skills
	expectedSkills := []string{"feature", "idea", "design", "build", "review"}
	for _, skill := range expectedSkills {
		found := false
		for _, s := range defaults.Skills {
			if s == skill {
				found = true
				break
			}
		}
		if !found {
			t.Errorf("Go defaults missing skill: %s", skill)
		}
	}

	// Check commands
	if defaults.TestCommand != "go test ./..." {
		t.Errorf("Go test command wrong: %s", defaults.TestCommand)
	}

	if defaults.BuildCommand != "go build ./..." {
		t.Errorf("Go build command wrong: %s", defaults.BuildCommand)
	}

	if defaults.PackageManager != "go" {
		t.Errorf("Go package manager wrong: %s", defaults.PackageManager)
	}

	// Evidence should be enabled by default
	if !defaults.EvidenceEnabled {
		t.Error("Evidence should be enabled by default for go")
	}
}

func TestGetDefaults_Node(t *testing.T) {
	defaults := GetDefaults("node")

	if defaults == nil {
		t.Fatal("GetDefaults returned nil for node")
	}

	// Check commands
	if defaults.TestCommand != "npm test" {
		t.Errorf("Node test command wrong: %s", defaults.TestCommand)
	}

	if defaults.PackageManager != "npm" {
		t.Errorf("Node package manager wrong: %s", defaults.PackageManager)
	}

	// Evidence should be enabled
	if !defaults.EvidenceEnabled {
		t.Error("Evidence should be enabled by default for node")
	}
}

func TestGetDefaults_Python(t *testing.T) {
	defaults := GetDefaults("python")

	if defaults == nil {
		t.Fatal("GetDefaults returned nil for python")
	}

	// Check commands
	if defaults.TestCommand != "pytest" {
		t.Errorf("Python test command wrong: %s", defaults.TestCommand)
	}

	if defaults.PackageManager != "pip" {
		t.Errorf("Python package manager wrong: %s", defaults.PackageManager)
	}

	// Evidence should be enabled
	if !defaults.EvidenceEnabled {
		t.Error("Evidence should be enabled by default for python")
	}
}

func TestGetDefaults_Mixed(t *testing.T) {
	defaults := GetDefaults("mixed")

	if defaults == nil {
		t.Fatal("GetDefaults returned nil for mixed")
	}

	// Mixed projects should have skills
	if len(defaults.Skills) == 0 {
		t.Error("Mixed defaults should have skills")
	}

	// Mixed projects may not have specific commands
	if defaults.PackageManager != "mixed" {
		t.Errorf("Mixed package manager wrong: %s", defaults.PackageManager)
	}
}

func TestGetDefaults_Unknown(t *testing.T) {
	defaults := GetDefaults("unknown")

	if defaults == nil {
		t.Fatal("GetDefaults returned nil for unknown")
	}

	// Unknown should still have some skills
	if len(defaults.Skills) == 0 {
		t.Error("Unknown defaults should have skills")
	}

	// Unknown should have safe defaults
	if defaults.PackageManager != "unknown" {
		t.Errorf("Unknown package manager wrong: %s", defaults.PackageManager)
	}
}

func TestGetDefaults_EmptyString(t *testing.T) {
	// Empty string should return unknown defaults
	defaults := GetDefaults("")

	if defaults == nil {
		t.Fatal("GetDefaults returned nil for empty string")
	}

	if defaults.PackageManager != "unknown" {
		t.Errorf("Empty string should return unknown defaults, got: %s", defaults.PackageManager)
	}
}

func TestGetDefaults_InvalidType(t *testing.T) {
	// Invalid types should return unknown defaults
	defaults := GetDefaults("invalid-type-xyz")

	if defaults == nil {
		t.Fatal("GetDefaults returned nil for invalid type")
	}

	// Should get safe defaults
	if len(defaults.Skills) == 0 {
		t.Error("Invalid type should still return skills")
	}
}

func TestMergeDefaults_WithUserSkills(t *testing.T) {
	userConfig := &Config{
		Skills: []string{"custom", "skills"},
	}

	merged := MergeDefaults("go", userConfig)

	// Should use user skills
	if len(merged.Skills) != 2 {
		t.Errorf("Expected 2 skills, got %d", len(merged.Skills))
	}

	if merged.Skills[0] != "custom" || merged.Skills[1] != "skills" {
		t.Errorf("User skills not preserved: %v", merged.Skills)
	}
}

func TestMergeDefaults_NoEvidence(t *testing.T) {
	userConfig := &Config{
		NoEvidence: true,
	}

	merged := MergeDefaults("go", userConfig)

	// Should disable evidence
	if merged.EvidenceEnabled {
		t.Error("Evidence should be disabled when user requests it")
	}
}

func TestMergeDefaults_EmptyUserConfig(t *testing.T) {
	userConfig := &Config{}

	merged := MergeDefaults("go", userConfig)

	// Should use all defaults
	if merged.TestCommand != "go test ./..." {
		t.Errorf("Default test command not used: %s", merged.TestCommand)
	}

	if !merged.EvidenceEnabled {
		t.Error("Evidence should be enabled by default")
	}
}

func TestMergeDefaults_NilUserConfig(t *testing.T) {
	// MergeDefaults should handle nil gracefully by using defaults
	merged := MergeDefaults("go", &Config{})

	if merged == nil {
		t.Fatal("MergeDefaults returned nil")
	}

	// Should use defaults
	if len(merged.Skills) == 0 {
		t.Error("Should have default skills")
	}
}
