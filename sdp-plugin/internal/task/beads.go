package task

import (
	"bufio"
	"fmt"
	"os"
	"os/exec"
	"path/filepath"
	"strconv"
	"strings"

	beadscli "github.com/fall-out-bug/sdp/internal/beads"
)

// BeadsIntegration handles beads CLI operations for tasks
type BeadsIntegration struct {
	client  *beadscli.Client
	enabled bool
}

// NewBeadsIntegration creates a new beads integration
func NewBeadsIntegration() *BeadsIntegration {
	return newBeadsIntegration(".")
}

func newBeadsIntegration(startDir string) *BeadsIntegration {
	projectRoot := findBeadsProjectRoot(startDir)
	oldWd, err := os.Getwd()
	if err != nil {
		return &BeadsIntegration{}
	}
	if err := os.Chdir(projectRoot); err != nil {
		return &BeadsIntegration{}
	}
	defer func() { _ = os.Chdir(oldWd) }()

	client, err := beadscli.NewClient()
	if err != nil {
		return &BeadsIntegration{}
	}
	return &BeadsIntegration{
		client:  client,
		enabled: detectBeads(projectRoot),
	}
}

// IsEnabled returns whether beads is available
func (b *BeadsIntegration) IsEnabled() bool {
	return b.enabled
}

// CreateBeadsIssue creates a beads issue and returns the ID
func (b *BeadsIntegration) CreateBeadsIssue(t *Task) (string, error) {
	if !b.enabled {
		return "", nil
	}
	beadsID, err := b.client.Create(t.Title, beadscli.CreateOptions{
		Type:     taskTypeToIssueType(t.Type),
		Priority: strconv.Itoa(int(t.Priority)),
	})
	if err != nil {
		return "", err
	}
	if err := b.client.Sync(); err != nil {
		return "", err
	}
	return beadsID, nil
}

// LinkWorkstreamToBeads updates workstream frontmatter with beads_id
func (b *BeadsIntegration) LinkWorkstreamToBeads(wsPath, beadsID string) error {
	if !b.enabled || beadsID == "" {
		return nil
	}

	content, err := os.ReadFile(wsPath)
	if err != nil {
		return fmt.Errorf("failed to read workstream: %w", err)
	}

	contentStr := string(content)
	if strings.Contains(contentStr, "beads_id:") {
		contentStr = updateBeadsIDLine(contentStr, beadsID)
	} else {
		contentStr = insertBeadsID(contentStr, beadsID)
	}

	if err := os.WriteFile(wsPath, []byte(contentStr), 0o644); err != nil {
		return fmt.Errorf("failed to write workstream: %w", err)
	}
	return nil
}

// updateBeadsIDLine replaces existing beads_id value
func updateBeadsIDLine(content, beadsID string) string {
	lines := strings.Split(content, "\n")
	for i, line := range lines {
		if strings.HasPrefix(line, "beads_id:") {
			lines[i] = fmt.Sprintf("beads_id: %s", beadsID)
			break
		}
	}
	return strings.Join(lines, "\n")
}

// insertBeadsID inserts beads_id into frontmatter before closing ---
func insertBeadsID(content, beadsID string) string {
	lines := strings.Split(content, "\n")
	for i := len(lines) - 1; i >= 0; i-- {
		if lines[i] == "---" && i > 0 {
			newLines := make([]string, 0, len(lines)+1)
			newLines = append(newLines, lines[:i]...)
			newLines = append(newLines, fmt.Sprintf("beads_id: %s", beadsID))
			newLines = append(newLines, lines[i:]...)
			return strings.Join(newLines, "\n")
		}
	}
	return content
}

// detectBeads checks if beads CLI is available
func detectBeads(projectRoot string) bool {
	if _, err := exec.LookPath("bd"); err != nil {
		return false
	}
	info, err := os.Stat(filepath.Join(projectRoot, ".beads"))
	return err == nil && info.IsDir()
}

// generateBeadsID creates a beads-style ID from title
func generateBeadsID(title string) string {
	words := strings.Fields(strings.ToLower(title))
	if len(words) == 0 {
		return "sdp-unknown"
	}
	slug := words[0]
	if len(words) > 1 {
		slug += "-" + words[1]
	}
	if len(slug) > 10 {
		slug = slug[:10]
	}
	return "sdp-" + slug
}

// CreateWorkstreamWithBeads creates a workstream and optionally links to beads
func (c *Creator) CreateWorkstreamWithBeads(task *Task) (*Workstream, error) {
	ws, err := c.CreateWorkstream(task)
	if err != nil {
		return nil, err
	}

	beads := newBeadsIntegration(c.projectRoot())
	if beads.IsEnabled() && task.BeadsID == "" {
		beadsID, err := beads.CreateBeadsIssue(task)
		if err != nil {
			fmt.Fprintf(os.Stderr, "warning: failed to create beads issue: %v\n", err)
		} else if beadsID != "" {
			if err := beads.LinkWorkstreamToBeads(ws.Path, beadsID); err != nil {
				fmt.Fprintf(os.Stderr, "warning: failed to link beads: %v\n", err)
			}
			if err := beads.UpdateMapping(ws.WSID, beadsID); err != nil {
				fmt.Fprintf(os.Stderr, "warning: failed to update beads mapping: %v\n", err)
			}
			ws.BeadsID = beadsID
		}
	}
	return ws, nil
}

func (c *Creator) projectRoot() string {
	return findBeadsProjectRoot(c.config.WorkstreamDir)
}

func findBeadsProjectRoot(startDir string) string {
	if startDir == "" {
		startDir = "."
	}
	absStart, err := filepath.Abs(startDir)
	if err != nil {
		absStart = startDir
	}
	current := absStart
	for {
		for _, marker := range []string{".beads", ".sdp", ".git"} {
			if _, err := os.Stat(filepath.Join(current, marker)); err == nil {
				return current
			}
		}
		parent := filepath.Dir(current)
		if parent == current {
			break
		}
		current = parent
	}
	if info, err := os.Stat(absStart); err == nil && !info.IsDir() {
		return filepath.Dir(absStart)
	}
	return absStart
}

func (b *BeadsIntegration) UpdateMapping(wsID, beadsID string) error {
	if !b.enabled || b.client == nil || beadsID == "" {
		return nil
	}
	return b.client.UpdateMapping(wsID, beadsID)
}

func taskTypeToIssueType(t Type) string {
	switch t {
	case TypeBug:
		return "bug"
	case TypeHotfix:
		return "bug"
	default:
		return "task"
	}
}

// ReadBeadsMapping reads the beads mapping file
func ReadBeadsMapping(path string) (map[string]string, error) {
	file, err := os.Open(path)
	if err != nil {
		if os.IsNotExist(err) {
			return make(map[string]string), nil
		}
		return nil, err
	}
	defer func() { _ = file.Close() }() //nolint:errcheck // cleanup

	mapping := make(map[string]string)
	scanner := bufio.NewScanner(file)
	for scanner.Scan() {
		line := scanner.Text()
		if line == "" {
			continue
		}
		if strings.Contains(line, "sdp_id") && strings.Contains(line, "beads_id") {
			sdpID := extractJSONValue(line, "sdp_id")
			beadsID := extractJSONValue(line, "beads_id")
			if sdpID != "" && beadsID != "" {
				mapping[sdpID] = beadsID
				mapping[beadsID] = sdpID
			}
		}
	}
	return mapping, scanner.Err()
}

// extractJSONValue extracts a value from a simple JSON line
func extractJSONValue(line, key string) string {
	search := `"` + key + `":"`
	start := strings.Index(line, search)
	if start == -1 {
		return ""
	}
	start += len(search)
	end := strings.Index(line[start:], `"`)
	if end == -1 {
		return ""
	}
	return line[start : start+end]
}
