package evidence

import (
	"os"
	"path/filepath"
	"strconv"
	"time"

	"github.com/fall-out-bug/sdp/internal/config"
)

// Emit appends an event to the evidence log (AC6, AC7). Non-blocking; errors are ignored.
func Emit(ev *Event) {
	if ev == nil {
		return
	}
	ev2 := *ev
	if ev2.ID == "" {
		ev2.ID = "evt-" + strconv.FormatInt(time.Now().UnixNano(), 10)
	}
	if ev2.Timestamp == "" {
		ev2.Timestamp = time.Now().UTC().Format(time.RFC3339)
	}
	go func() {
		_ = emitSync(&ev2)
	}()
}

// emitSync writes event to log; returns error (caller may ignore).
func emitSync(ev *Event) error {
	root, err := config.FindProjectRoot()
	if err != nil {
		return err
	}
	cfg, err := config.Load(root)
	if err != nil {
		return err
	}
	if cfg == nil || !cfg.Evidence.Enabled {
		return nil
	}
	logPath := cfg.Evidence.LogPath
	if logPath == "" {
		logPath = ".sdp/log/events.jsonl"
	}
	path := filepath.Join(root, logPath)
	w, err := NewWriter(path)
	if err != nil {
		return err
	}
	return w.Append(ev)
}

// Enabled returns whether evidence emission is enabled (AC8).
func Enabled() bool {
	root, err := config.FindProjectRoot()
	if err != nil {
		return false
	}
	cfg, err := config.Load(root)
	if err != nil || cfg == nil {
		return true
	}
	return cfg.Evidence.Enabled
}

// ModelID returns SDP_MODEL_ID or ANTHROPIC_MODEL or "unknown" (AC5).
func ModelID() string {
	if s := os.Getenv("SDP_MODEL_ID"); s != "" {
		return s
	}
	if s := os.Getenv("ANTHROPIC_MODEL"); s != "" {
		return s
	}
	return "unknown"
}

// PlanEvent builds a plan event (AC1).
func PlanEvent(wsID string, scopeFiles []string) *Event {
	return &Event{
		Type: "plan",
		WSID: wsID,
		Data: map[string]interface{}{
			"scope_files": scopeFiles,
			"action":      "activate",
		},
	}
}

// GenerationEvent builds a generation event (AC2).
func GenerationEvent(wsID string, filesChanged []string) *Event {
	return &Event{
		Type: "generation",
		WSID: wsID,
		Data: map[string]interface{}{
			"model_id":      ModelID(),
			"model_version": "",
			"prompt_hash":   "",
			"files_changed": filesChanged,
		},
	}
}

// VerificationEvent builds a verification event (AC3, AC4).
func VerificationEvent(wsID string, passed bool, gateName string, coverage float64) *Event {
	return &Event{
		Type: "verification",
		WSID: wsID,
		Data: map[string]interface{}{
			"passed":    passed,
			"gate_name": gateName,
			"coverage":  coverage,
		},
	}
}
