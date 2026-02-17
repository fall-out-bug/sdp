package sdpinit

import (
	"bufio"
	"fmt"
	"io"
	"os"
	"path/filepath"
	"strconv"
	"strings"
)

// WizardAnswers contains the user's responses to wizard prompts.
type WizardAnswers struct {
	ProjectName string
	ProjectType string
	Skills      []string
	NoEvidence  bool
	SkipBeads   bool
}

// WizardPrompt defines a single prompt in the wizard.
type WizardPrompt struct {
	ID          string
	Question    string
	Default     string
	Options     []string
	AllowCustom bool
}

// Wizard controls the interactive onboarding flow.
type Wizard struct {
	reader    *bufio.Reader
	writer    io.Writer
	preflight *PreflightResult
}

// NewWizard creates a new interactive wizard.
func NewWizard(reader io.Reader, writer io.Writer, preflight *PreflightResult) *Wizard {
	var r *bufio.Reader
	if reader != nil {
		r = bufio.NewReader(reader)
	} else {
		r = bufio.NewReader(os.Stdin)
	}

	var w io.Writer
	if writer != nil {
		w = writer
	} else {
		w = os.Stdout
	}

	return &Wizard{
		reader:    r,
		writer:    w,
		preflight: preflight,
	}
}

// Run executes the interactive wizard and returns the user's answers.
func (w *Wizard) Run() (*WizardAnswers, error) {
	answers := &WizardAnswers{}

	// Display welcome
	w.printWelcome()

	// Show preflight results if available
	if w.preflight != nil {
		w.printPreflightSummary()
	}

	// Prompt for project name
	name, err := w.promptString("Project name", w.getDefaultProjectName())
	if err != nil {
		return nil, fmt.Errorf("reading project name: %w", err)
	}
	answers.ProjectName = name

	// Prompt for project type
	projectType, err := w.promptProjectType()
	if err != nil {
		return nil, fmt.Errorf("reading project type: %w", err)
	}
	answers.ProjectType = projectType

	// Prompt for skills selection
	skills, err := w.promptSkills(projectType)
	if err != nil {
		return nil, fmt.Errorf("reading skills: %w", err)
	}
	answers.Skills = skills

	// Prompt for evidence logging
	noEvidence, err := w.promptBool("Disable evidence logging", false)
	if err != nil {
		return nil, fmt.Errorf("reading evidence preference: %w", err)
	}
	answers.NoEvidence = noEvidence

	// Prompt for Beads integration
	skipBeads, err := w.promptBool("Skip Beads integration", false)
	if err != nil {
		return nil, fmt.Errorf("reading Beads preference: %w", err)
	}
	answers.SkipBeads = skipBeads

	// Show summary and confirm
	w.printSummary(answers)
	confirm, err := w.promptBool("Proceed with these settings", true)
	if err != nil {
		return nil, fmt.Errorf("reading confirmation: %w", err)
	}

	if !confirm {
		return nil, fmt.Errorf("initialization cancelled by user")
	}

	return answers, nil
}

func (w *Wizard) printWelcome() {
	fmt.Fprintln(w.writer, "")
	fmt.Fprintln(w.writer, "Welcome to SDP Onboarding")
	fmt.Fprintln(w.writer, "=========================")
	fmt.Fprintln(w.writer, "This wizard will help you set up SDP for your project.")
	fmt.Fprintln(w.writer, "")
}

func (w *Wizard) printPreflightSummary() {
	fmt.Fprintln(w.writer, "Detected Environment:")
	fmt.Fprintf(w.writer, "  Project type: %s\n", w.preflight.ProjectType)

	if w.preflight.HasGit {
		fmt.Fprintln(w.writer, "  Git: detected")
	} else {
		fmt.Fprintln(w.writer, "  Git: not found (recommended)")
	}

	if w.preflight.HasSDP {
		fmt.Fprintln(w.writer, "  SDP: already initialized")
	}

	if w.preflight.HasClaude {
		fmt.Fprintln(w.writer, "  Claude: already initialized")
	}

	if len(w.preflight.Conflicts) > 0 {
		fmt.Fprintln(w.writer, "  Conflicts:")
		for _, c := range w.preflight.Conflicts {
			fmt.Fprintf(w.writer, "    - %s\n", c)
		}
	}

	if len(w.preflight.Warnings) > 0 {
		fmt.Fprintln(w.writer, "  Warnings:")
		for _, warn := range w.preflight.Warnings {
			fmt.Fprintf(w.writer, "    - %s\n", warn)
		}
	}

	fmt.Fprintln(w.writer, "")
}

func (w *Wizard) getDefaultProjectName() string {
	cwd, err := os.Getwd()
	if err != nil {
		return "my-project"
	}
	return filepath.Base(cwd)
}

func (w *Wizard) promptString(label string, def string) (string, error) {
	if def != "" {
		fmt.Fprintf(w.writer, "%s [%s]: ", label, def)
	} else {
		fmt.Fprintf(w.writer, "%s: ", label)
	}

	input, err := w.reader.ReadString('\n')
	if err != nil && err != io.EOF {
		return "", err
	}

	input = strings.TrimSpace(input)
	if input == "" {
		return def, nil
	}
	return input, nil
}

func (w *Wizard) promptProjectType() (string, error) {
	options := []string{"go", "node", "python", "mixed", "unknown"}
	defaultIdx := 0

	// Use detected type as default
	if w.preflight != nil {
		for i, opt := range options {
			if opt == w.preflight.ProjectType {
				defaultIdx = i
				break
			}
		}
	}

	fmt.Fprintln(w.writer, "Project type:")
	for i, opt := range options {
		marker := " "
		if i == defaultIdx {
			marker = "*"
		}
		fmt.Fprintf(w.writer, "  %d) %s %s\n", i+1, marker, opt)
	}
	fmt.Fprintln(w.writer, "  (* = detected)")

	choice, err := w.promptInt("Select", defaultIdx+1)
	if err != nil {
		return "", err
	}

	if choice < 1 || choice > len(options) {
		return options[defaultIdx], nil
	}

	return options[choice-1], nil
}

func (w *Wizard) promptSkills(projectType string) ([]string, error) {
	defaults := GetDefaults(projectType)

	fmt.Fprintln(w.writer, "Skills to enable (comma-separated numbers, or 'all'):")
	allSkills := []string{
		"feature", "idea", "design", "build",
		"review", "deploy", "debug", "bugfix", "hotfix", "oneshot",
	}

	for i, skill := range allSkills {
		marker := " "
		for _, def := range defaults.Skills {
			if def == skill {
				marker = "*"
				break
			}
		}
		fmt.Fprintf(w.writer, "  %d) %s %s\n", i+1, marker, skill)
	}
	fmt.Fprintln(w.writer, "  (* = recommended)")

	input, err := w.promptString("Select skills", "")
	if err != nil {
		return defaults.Skills, nil
	}

	input = strings.TrimSpace(strings.ToLower(input))
	if input == "" || input == "all" {
		return defaults.Skills, nil
	}

	// Parse comma-separated numbers
	selected := []string{}
	parts := strings.Split(input, ",")
	for _, part := range parts {
		part = strings.TrimSpace(part)
		if idx, err := strconv.Atoi(part); err == nil {
			if idx >= 1 && idx <= len(allSkills) {
				selected = append(selected, allSkills[idx-1])
			}
		}
	}

	if len(selected) == 0 {
		return defaults.Skills, nil
	}

	return selected, nil
}

func (w *Wizard) promptBool(label string, def bool) (bool, error) {
	defStr := "y/N"
	if def {
		defStr = "Y/n"
	}

	fmt.Fprintf(w.writer, "%s [%s]: ", label, defStr)

	input, err := w.reader.ReadString('\n')
	if err != nil && err != io.EOF {
		return def, err
	}

	input = strings.TrimSpace(strings.ToLower(input))
	if input == "" {
		return def, nil
	}

	return input == "y" || input == "yes", nil
}

func (w *Wizard) promptInt(label string, def int) (int, error) {
	fmt.Fprintf(w.writer, "%s [%d]: ", label, def)

	input, err := w.reader.ReadString('\n')
	if err != nil && err != io.EOF {
		return def, err
	}

	input = strings.TrimSpace(input)
	if input == "" {
		return def, nil
	}

	val, err := strconv.Atoi(input)
	if err != nil {
		return def, nil
	}

	return val, nil
}

func (w *Wizard) printSummary(answers *WizardAnswers) {
	fmt.Fprintln(w.writer, "")
	fmt.Fprintln(w.writer, "Configuration Summary:")
	fmt.Fprintln(w.writer, "----------------------")
	fmt.Fprintf(w.writer, "  Project name: %s\n", answers.ProjectName)
	fmt.Fprintf(w.writer, "  Project type: %s\n", answers.ProjectType)
	fmt.Fprintf(w.writer, "  Skills: %v\n", answers.Skills)
	fmt.Fprintf(w.writer, "  Evidence logging: %v\n", !answers.NoEvidence)
	fmt.Fprintf(w.writer, "  Beads integration: %v\n", !answers.SkipBeads)
	fmt.Fprintln(w.writer, "")
}

// RunWizard is a convenience function to run the wizard with default I/O.
func RunWizard(preflight *PreflightResult) (*WizardAnswers, error) {
	wizard := NewWizard(nil, nil, preflight)
	return wizard.Run()
}
