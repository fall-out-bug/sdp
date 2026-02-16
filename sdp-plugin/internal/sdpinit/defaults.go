package sdpinit

// ProjectDefaults contains safe defaults for a project type.
type ProjectDefaults struct {
	// Skills to enable by default
	Skills []string
	// Validators to configure
	Validators []string
	// Whether evidence logging is enabled
	EvidenceEnabled bool
	// Path to documentation
	DocPath string
	// Test command for the project
	TestCommand string
	// Build command for the project
	BuildCommand string
	// Lint command for the project
	LintCommand string
	// Package manager for the project
	PackageManager string
}

// GetDefaults returns safe defaults for the given project type.
func GetDefaults(projectType string) *ProjectDefaults {
	switch projectType {
	case "go":
		return &ProjectDefaults{
			Skills: []string{
				"feature", "idea", "design", "build",
				"review", "deploy", "debug", "bugfix", "hotfix", "oneshot",
			},
			Validators:      []string{"go-vet", "go-test"},
			EvidenceEnabled: true,
			DocPath:         "docs",
			TestCommand:     "go test ./...",
			BuildCommand:    "go build ./...",
			LintCommand:     "go vet ./...",
			PackageManager:  "go",
		}
	case "node":
		return &ProjectDefaults{
			Skills: []string{
				"feature", "idea", "design", "build",
				"review", "deploy", "debug", "bugfix", "hotfix", "oneshot",
			},
			Validators:      []string{"npm-test", "eslint"},
			EvidenceEnabled: true,
			DocPath:         "docs",
			TestCommand:     "npm test",
			BuildCommand:    "npm run build",
			LintCommand:     "npm run lint",
			PackageManager:  "npm",
		}
	case "python":
		return &ProjectDefaults{
			Skills: []string{
				"feature", "idea", "design", "build",
				"review", "deploy", "debug", "bugfix", "hotfix", "oneshot",
			},
			Validators:      []string{"pytest", "ruff"},
			EvidenceEnabled: true,
			DocPath:         "docs",
			TestCommand:     "pytest",
			BuildCommand:    "python -m build",
			LintCommand:     "ruff check .",
			PackageManager:  "pip",
		}
	case "mixed":
		return &ProjectDefaults{
			Skills: []string{
				"feature", "idea", "design", "build",
				"review", "deploy", "debug", "bugfix", "hotfix", "oneshot",
			},
			Validators:      []string{},
			EvidenceEnabled: true,
			DocPath:         "docs",
			TestCommand:     "",
			BuildCommand:    "",
			LintCommand:     "",
			PackageManager:  "mixed",
		}
	default:
		// Unknown project type - minimal safe defaults
		return &ProjectDefaults{
			Skills: []string{
				"feature", "idea", "design", "build",
				"review", "deploy", "debug", "bugfix", "hotfix", "oneshot",
			},
			Validators:      []string{},
			EvidenceEnabled: true,
			DocPath:         "docs",
			TestCommand:     "",
			BuildCommand:    "",
			LintCommand:     "",
			PackageManager:  "unknown",
		}
	}
}

// MergeDefaults combines user-provided config with defaults.
// User-provided values take precedence.
func MergeDefaults(projectType string, userConfig *Config) *ProjectDefaults {
	defaults := GetDefaults(projectType)

	// If user specified skills, use those instead
	if len(userConfig.Skills) > 0 {
		defaults.Skills = userConfig.Skills
	}

	// If user disabled evidence, respect that
	if userConfig.NoEvidence {
		defaults.EvidenceEnabled = false
	}

	return defaults
}
