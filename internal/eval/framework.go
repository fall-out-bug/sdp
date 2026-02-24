package eval

import (
	"bufio"
	"encoding/json"
	"fmt"
	"os"
	"path/filepath"
	"strings"

	"gopkg.in/yaml.v3"
)

// Case defines a single eval case.
type Case struct {
	Name             string   `yaml:"name"`
	Skill            string   `yaml:"skill"`
	InputTranscript  string   `yaml:"input_transcript"`
	ForbiddenPatterns []string `yaml:"forbidden_patterns"`
	RequiredPatterns []string `yaml:"required_patterns"`
	Verdict          string   `yaml:"verdict"` // PASS or FAIL
}

// Result is the outcome of running one case.
type Result struct {
	Case   string
	Pass   bool
	Reason string
}

// RunCase loads the transcript, extracts agent output, and checks patterns.
// For verdict=PASS: case passes when no forbidden patterns and all required present.
// For verdict=FAIL: case passes when we correctly flag violations (expect transcript to fail).
func RunCase(c *Case, projectRoot string) Result {
	path := filepath.Join(projectRoot, c.InputTranscript)
	data, err := os.ReadFile(path)
	if err != nil {
		return Result{Case: c.Name, Pass: false, Reason: fmt.Sprintf("read transcript: %v", err)}
	}
	output := extractAgentOutput(data)
	hasForbidden := false
	var forbiddenFound []string
	for _, p := range c.ForbiddenPatterns {
		if strings.Contains(output, p) {
			hasForbidden = true
			forbiddenFound = append(forbiddenFound, p)
		}
	}
	missingRequired := false
	var missing []string
	for _, p := range c.RequiredPatterns {
		if !strings.Contains(output, p) {
			missingRequired = true
			missing = append(missing, p)
		}
	}
	rawPass := !hasForbidden && !missingRequired
	var reason string
	if hasForbidden {
		reason = fmt.Sprintf("forbidden patterns found: %s", strings.Join(forbiddenFound, ", "))
	}
	if missingRequired {
		if reason != "" {
			reason += "; "
		}
		reason += fmt.Sprintf("missing required patterns: %s", strings.Join(missing, ", "))
	}
	// verdict FAIL = we expect transcript to violate; "pass" means we correctly caught it
	expectFail := strings.ToUpper(c.Verdict) == "FAIL"
	pass := (expectFail && !rawPass) || (!expectFail && rawPass)
	return Result{Case: c.Name, Pass: pass, Reason: reason}
}

// extractAgentOutput parses JSONL transcript and concatenates assistant message content.
func extractAgentOutput(data []byte) string {
	var sb strings.Builder
	sc := bufio.NewScanner(strings.NewReader(string(data)))
	for sc.Scan() {
		line := strings.TrimSpace(sc.Text())
		if line == "" {
			continue
		}
		var msg struct {
			Role    string `json:"role"`
			Content string `json:"content"`
			Message *struct {
				Content []struct {
					Type string `json:"type"`
					Text string `json:"text"`
				} `json:"content"`
			} `json:"message"`
		}
		if err := json.Unmarshal([]byte(line), &msg); err != nil {
			continue
		}
		if msg.Role != "assistant" {
			continue
		}
		if msg.Content != "" {
			sb.WriteString(msg.Content)
			sb.WriteString("\n")
		}
		if msg.Message != nil {
			for _, c := range msg.Message.Content {
				if c.Type == "text" && c.Text != "" {
					sb.WriteString(c.Text)
					sb.WriteString("\n")
				}
			}
		}
	}
	return sb.String()
}

// LoadCases reads YAML case files from a directory.
func LoadCases(casesDir, skill string) ([]Case, error) {
	pattern := filepath.Join(casesDir, "*.yaml")
	matches, err := filepath.Glob(pattern)
	if err != nil {
		return nil, err
	}
	var cases []Case
	for _, p := range matches {
		data, err := os.ReadFile(p)
		if err != nil {
			return nil, err
		}
		var c Case
		if err := yaml.Unmarshal(data, &c); err != nil {
			return nil, fmt.Errorf("%s: %w", p, err)
		}
		if skill != "" && c.Skill != skill {
			continue
		}
		cases = append(cases, c)
	}
	return cases, nil
}

// Run runs all cases for a skill and returns results.
func Run(projectRoot, casesDir, skill string) ([]Result, error) {
	cases, err := LoadCases(casesDir, skill)
	if err != nil {
		return nil, err
	}
	var results []Result
	for _, c := range cases {
		results = append(results, RunCase(&c, projectRoot))
	}
	return results, nil
}
