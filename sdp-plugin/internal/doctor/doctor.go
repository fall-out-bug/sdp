package doctor

import (
	"fmt"
	"os"
	"os/exec"
	"strings"
)

type CheckResult struct {
	Name    string
	Status  string // "ok", "warning", "error"
	Message string
}

func Run() []CheckResult {
	results := []CheckResult{}

	// Check 1: Git
	results = append(results, checkGit())

	// Check 2: Claude Code
	results = append(results, checkClaudeCode())

	// Check 3: Go (for building binary)
	results = append(results, checkGo())

	// Check 4: .claude/ directory
	results = append(results, checkClaudeDir())

	return results
}

func checkGit() CheckResult {
	if _, err := exec.LookPath("git"); err != nil {
		return CheckResult{
			Name:    "Git",
			Status:  "error",
			Message: "Git not found. Install from https://git-scm.com",
		}
	}

	// Get version
	cmd := exec.Command("git", "--version")
	output, _ := cmd.Output()
	version := strings.TrimSpace(string(output))

	return CheckResult{
		Name:    "Git",
		Status:  "ok",
		Message: fmt.Sprintf("Installed (%s)", version),
	}
}

func checkClaudeCode() CheckResult {
	if _, err := exec.LookPath("claude"); err != nil {
		return CheckResult{
			Name:    "Claude Code",
			Status:  "warning",
			Message: "Claude Code CLI not found. Plugin will work in Claude Desktop app.",
		}
	}

	// Get version
	cmd := exec.Command("claude", "--version")
	output, _ := cmd.Output()
	version := strings.TrimSpace(string(output))

	return CheckResult{
		Name:    "Claude Code",
		Status:  "ok",
		Message: fmt.Sprintf("Installed (%s)", version),
	}
}

func checkGo() CheckResult {
	if _, err := exec.LookPath("go"); err != nil {
		return CheckResult{
			Name:    "Go",
			Status:  "warning",
			Message: "Go not found. Required only for building SDP binary.",
		}
	}

	// Get version
	cmd := exec.Command("go", "version")
	output, _ := cmd.Output()
	version := strings.TrimSpace(string(output))

	return CheckResult{
		Name:    "Go",
		Status:  "ok",
		Message: fmt.Sprintf("Installed (%s)", version),
	}
}

func checkClaudeDir() CheckResult {
	if _, err := os.Stat(".claude"); os.IsNotExist(err) {
		return CheckResult{
			Name:    ".claude/ directory",
			Status:  "error",
			Message: "Not found. Run 'sdp init' to initialize",
		}
	}

	// Check if it has expected structure
	dirs := []string{"skills", "agents", "validators"}
	missing := []string{}
	for _, dir := range dirs {
		if _, err := os.Stat(".claude/" + dir); os.IsNotExist(err) {
			missing = append(missing, dir)
		}
	}

	if len(missing) > 0 {
		return CheckResult{
			Name:    ".claude/ directory",
			Status:  "warning",
			Message: fmt.Sprintf("Incomplete (missing: %s)", strings.Join(missing, ", ")),
		}
	}

	return CheckResult{
		Name:    ".claude/ directory",
		Status:  "ok",
		Message: "SDP prompts installed",
	}
}
