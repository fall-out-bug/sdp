// Package evidence provides event emission for SDP workflows.
//
// CLI usage: Use EmitSync from CLI entry points (sdp verify, sdp quality, sdp apply)
// so process exit does not drop evidence. Emit() is async and may lose events on exit.
package evidence

import (
	"log/slog"
	"os"
	"path/filepath"
	"strconv"
	"sync"
	"time"

	"github.com/fall-out-bug/sdp/internal/config"
)

var (
	globalWriter     *Writer
	globalWriterOnce sync.Once
	globalWriterErr  error
	globalWriterPath string
)

// ResetGlobalWriter clears the singleton for testing.
func ResetGlobalWriter() {
	globalWriterOnce = sync.Once{}
	globalWriter = nil
	globalWriterErr = nil
	globalWriterPath = ""
}

func getOrCreateWriter() (*Writer, error) {
	globalWriterOnce.Do(func() {
		root, err := config.FindProjectRoot()
		if err != nil {
			globalWriterErr = err
			return
		}
		cfg, err := config.Load(root)
		if err != nil {
			globalWriterErr = err
			return
		}
		if cfg == nil || !cfg.Evidence.Enabled {
			return
		}
		logPath := cfg.Evidence.LogPath
		if logPath == "" {
			logPath = ".sdp/log/events.jsonl"
		}
		globalWriterPath = filepath.Join(root, logPath)
		globalWriter, globalWriterErr = NewWriter(globalWriterPath)
	})
	return globalWriter, globalWriterErr
}

func fillDefaults(ev *Event) {
	if ev.ID == "" {
		ev.ID = "evt-" + strconv.FormatInt(time.Now().UnixNano(), 10)
	}
	if ev.Timestamp == "" {
		ev.Timestamp = time.Now().UTC().Format(time.RFC3339)
	}
}

// Emit appends an event asynchronously in a goroutine. Use only when the
// process will not exit immediately (e.g. long-running services). For CLI
// entry points (verify, quality, oneshot) use EmitSync so process exit does
// not drop evidence before the goroutine completes.
func Emit(ev *Event) {
	if ev == nil {
		return
	}
	ev2 := *ev
	fillDefaults(&ev2)
	go func() {
		if err := emitSync(&ev2); err != nil {
			slog.Error("evidence emission failed",
				"event_id", ev2.ID,
				"event_type", ev2.Type,
				"error", err,
			)
		}
	}()
}

// EmitSync writes the event immediately. Use from CLI entry points (verify,
// quality, oneshot) so process exit does not drop evidence.
func EmitSync(ev *Event) error {
	if ev == nil {
		return nil
	}
	ev2 := *ev
	fillDefaults(&ev2)
	return emitSync(&ev2)
}

// emitSync writes event to the singleton writer.
func emitSync(ev *Event) error {
	if ev.ID == "" || ev.Type == "" || ev.Timestamp == "" {
		slog.Warn("evidence event missing required fields, skipping",
			"id", ev.ID, "type", ev.Type, "timestamp", ev.Timestamp)
		return nil
	}
	w, err := getOrCreateWriter()
	if err != nil {
		return err
	}
	if w == nil {
		return nil
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

// ModelID returns best-effort model identifier from environment (AC5).
func ModelID() string {
	keys := []string{
		"SDP_MODEL_ID",
		"OPENCODE_MODEL",
		"ANTHROPIC_MODEL",
		"OPENAI_MODEL",
		"MODEL",
	}
	for _, key := range keys {
		if s := os.Getenv(key); s != "" {
			return s
		}
	}
	return "unknown"
}
