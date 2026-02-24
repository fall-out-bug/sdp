package orchestrate

import (
	"context"
	"crypto/sha256"
	"encoding/hex"
	"encoding/json"
	"fmt"
	"os"
	"os/exec"
	"path/filepath"
	"strings"
)

// buildPromptWithContext injects the pre-hydrated context packet into the prompt.
func buildPromptWithContext(dir, basePrompt string) string {
	pkt, err := LoadContextPacket(dir)
	if err != nil || pkt == nil {
		return basePrompt
	}
	return basePrompt + pkt.FormatForPrompt()
}

// ComputePromptHash returns SHA-256 hex of the rendered prompt (captures exactly what was sent to the LLM).
func ComputePromptHash(prompt string) string {
	h := sha256.Sum256([]byte(prompt))
	return hex.EncodeToString(h[:])
}

// ContextSource records an input that entered the agent's context (F026 prompt provenance).
type ContextSource struct {
	Type string `json:"type"`
	Path string `json:"path"`
	Hash string `json:"hash"`
}

// BuildContextSources builds the list of context sources for prompt provenance.
// Paths are relative to projectRoot for portability.
func BuildContextSources(projectRoot, featureID, wsID string, scopeFiles []string) []ContextSource {
	hashFile := func(absPath string) string {
		b, err := os.ReadFile(absPath)
		if err != nil {
			return ""
		}
		h := sha256.Sum256(b)
		return hex.EncodeToString(h[:])
	}
	var out []ContextSource
	wsRel := filepath.Join("docs", "workstreams", "backlog", wsID+".md")
	wsPath := filepath.Join(projectRoot, wsRel)
	if h := hashFile(wsPath); h != "" {
		out = append(out, ContextSource{Type: "workstream_spec", Path: wsRel, Hash: h})
	}
	cpRel := filepath.Join(".sdp", "checkpoints", featureID+".json")
	cpPath := filepath.Join(projectRoot, cpRel)
	if h := hashFile(cpPath); h != "" {
		out = append(out, ContextSource{Type: "checkpoint", Path: cpRel, Hash: h})
	}
	for _, f := range scopeFiles {
		p := filepath.Join(projectRoot, f)
		if h := hashFile(p); h != "" {
			out = append(out, ContextSource{Type: "scope_file", Path: f, Hash: h})
		}
	}
	agentsRel := "AGENTS.md"
	if h := hashFile(filepath.Join(projectRoot, agentsRel)); h != "" {
		out = append(out, ContextSource{Type: "agents_md", Path: agentsRel, Hash: h})
	}
	skillRel := filepath.Join(".cursor", "skills", "build", "SKILL.md")
	if h := hashFile(filepath.Join(projectRoot, skillRel)); h != "" {
		out = append(out, ContextSource{Type: "skill", Path: skillRel, Hash: h})
	}
	ctxPktRel := filepath.Join(".sdp", "context-packet.json")
	if h := hashFile(filepath.Join(projectRoot, ctxPktRel)); h != "" {
		out = append(out, ContextSource{Type: "context_packet", Path: ctxPktRel, Hash: h})
	}
	return out
}

// WritePromptProvenance writes prompt_hash and context_sources to .sdp/prompt-provenance.json.
// Downstream (evidence builder, post-build hook) can merge into the evidence envelope.
// Uses tmp+rename for atomic write.
func WritePromptProvenance(projectRoot string, promptHash string, sources []ContextSource) error {
	sdpDir := filepath.Join(projectRoot, ".sdp")
	if err := os.MkdirAll(sdpDir, 0o755); err != nil {
		return err
	}
	path := filepath.Join(sdpDir, "prompt-provenance.json")
	tmpPath := path + ".tmp"
	body := map[string]any{"prompt_hash": promptHash, "context_sources": sources}
	data, err := json.MarshalIndent(body, "", "  ")
	if err != nil {
		return err
	}
	if err := os.WriteFile(tmpPath, data, 0o644); err != nil {
		return err
	}
	if err := os.Rename(tmpPath, path); err != nil {
		_ = os.Remove(tmpPath)
		return err
	}
	return nil
}

// InvokeOpenCode runs `opencode run --agent orchestrator` with the given prompt.
// Returns the combined stdout+stderr and exit code.
func InvokeOpenCode(ctx context.Context, dir, agent, prompt string) (string, int, error) {
	if agent == "" {
		agent = "orchestrator"
	}
	cmd := exec.CommandContext(ctx, "opencode", "run", "--agent", agent)
	cmd.Dir = dir
	cmd.Stdin = strings.NewReader(prompt)
	out, err := cmd.CombinedOutput()
	if err != nil {
		if exitErr, ok := err.(*exec.ExitError); ok {
			return string(out), exitErr.ExitCode(), nil
		}
		return string(out), -1, fmt.Errorf("opencode run: %w", err)
	}
	return string(out), 0, nil
}

// RunBuildPhase invokes opencode to execute a single @build workstream.
// Computes prompt_hash and context_sources before LLM invocation (F026 prompt provenance).
func RunBuildPhase(ctx context.Context, projectRoot, featureID, wsID string) (commit string, err error) {
	prompt := buildPromptWithContext(projectRoot, fmt.Sprintf("Execute @build %s. Output only code and commit message. After commit, output the commit hash.", wsID))
	promptHash := ComputePromptHash(prompt)
	var scopeFiles []string
	if pkt, err := LoadContextPacket(projectRoot); err == nil && pkt != nil {
		scopeFiles = pkt.ScopeFiles
	}
	sources := BuildContextSources(projectRoot, featureID, wsID, scopeFiles)
	_ = WritePromptProvenance(projectRoot, promptHash, sources)
	out, code, err := InvokeOpenCode(ctx, projectRoot, "implementer", prompt)
	if err != nil {
		return "", err
	}
	if code != 0 {
		return "", fmt.Errorf("opencode build exited %d: %s", code, out)
	}
	// Extract last line as commit hash if it looks like a SHA
	lines := strings.Split(strings.TrimSpace(out), "\n")
	for i := len(lines) - 1; i >= 0; i-- {
		s := strings.TrimSpace(lines[i])
		if len(s) == 40 && isHex(s) {
			return s, nil
		}
	}
	return "", nil
}

// RunReviewPhase invokes opencode to execute @review for a feature.
func RunReviewPhase(ctx context.Context, dir, featureID string) (approved bool, err error) {
	prompt := buildPromptWithContext(dir, fmt.Sprintf("Execute @review %s. Fix P0/P1 findings. Output APPROVED when done.", featureID))
	out, code, err := InvokeOpenCode(ctx, dir, "reviewer", prompt)
	if err != nil {
		return false, err
	}
	approved = code == 0 && strings.Contains(strings.ToUpper(out), "APPROVED")
	return approved, nil
}

func isHex(s string) bool {
	for _, c := range s {
		if (c < '0' || c > '9') && (c < 'a' || c > 'f') && (c < 'A' || c > 'F') {
			return false
		}
	}
	return true
}
