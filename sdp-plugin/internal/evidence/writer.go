package evidence

import (
	"crypto/sha256"
	"encoding/hex"
	"encoding/json"
	"fmt"
	"os"
	"path/filepath"
	"sync"
	"syscall"
)

const genesisHash = "genesis"

// Writer appends events to .sdp/log/events.jsonl with hash chain (AC1–AC4, AC8).
// Safe for concurrent use within a process (mutex) and across processes (flock).
type Writer struct {
	path     string
	mu       sync.Mutex
	lastHash string
}

// NewWriter creates a writer for the given path; creates parent dirs (AC1).
func NewWriter(path string) (*Writer, error) {
	if path == "" {
		return nil, fmt.Errorf("path is empty")
	}
	dir := filepath.Dir(path)
	if err := os.MkdirAll(dir, 0755); err != nil {
		return nil, fmt.Errorf("create log dir: %w", err)
	}
	w := &Writer{path: path, lastHash: genesisHash}
	if b, err := os.ReadFile(path); err == nil && len(b) > 0 {
		lastLine := lastLineBytes(b)
		if len(lastLine) > 0 {
			w.lastHash = hashLine(lastLine)
		}
	}
	return w, nil
}

// Append writes event with prev_hash and fsync (AC2, AC3, AC4).
// Uses flock for inter-process safety: re-reads lastHash under lock
// in case another process appended since our last write.
func (w *Writer) Append(ev *Event) error {
	w.mu.Lock()
	defer w.mu.Unlock()

	lockPath := w.path + ".lock"
	lf, err := os.OpenFile(lockPath, os.O_CREATE|os.O_WRONLY, 0600)
	if err != nil {
		return fmt.Errorf("open lock file: %w", err)
	}
	defer lf.Close()

	if err := syscall.Flock(int(lf.Fd()), syscall.LOCK_EX); err != nil {
		return fmt.Errorf("acquire file lock: %w", err)
	}
	defer func() { _ = syscall.Flock(int(lf.Fd()), syscall.LOCK_UN) }()

	// Re-read last hash under flock — another process may have appended.
	// This ensures prev_hash is always derived from on-disk state (atomic with append).
	if b, err := os.ReadFile(w.path); err == nil && len(b) > 0 {
		if last := lastLineBytes(b); len(last) > 0 {
			w.lastHash = hashLine(last)
		}
	}

	ev.PrevHash = w.lastHash
	data, err := json.Marshal(ev)
	if err != nil {
		return fmt.Errorf("marshal event: %w", err)
	}
	line := append(data, '\n')
	if err := appendToFile(w.path, line); err != nil {
		return err
	}
	w.lastHash = hashLine(data)
	return nil
}

func hashLine(data []byte) string {
	h := sha256.Sum256(data)
	return hex.EncodeToString(h[:])
}

func lastLineBytes(b []byte) []byte {
	if len(b) == 0 {
		return nil
	}
	lastNewline := -1
	for i := len(b) - 1; i >= 0; i-- {
		if b[i] == '\n' {
			lastNewline = i
			break
		}
	}
	if lastNewline < 0 {
		return b
	}
	prevNewline := -1
	for i := 0; i < lastNewline; i++ {
		if b[i] == '\n' {
			prevNewline = i
		}
	}
	return b[prevNewline+1 : lastNewline]
}

func appendToFile(path string, data []byte) error {
	f, err := os.OpenFile(path, os.O_APPEND|os.O_CREATE|os.O_WRONLY, 0600)
	if err != nil {
		return fmt.Errorf("open file: %w", err)
	}
	_, err = f.Write(data)
	if err != nil {
		if cErr := f.Close(); cErr != nil {
			return fmt.Errorf("write: %w (close: %v)", err, cErr)
		}
		return fmt.Errorf("write: %w", err)
	}
	if err := f.Sync(); err != nil {
		if cErr := f.Close(); cErr != nil {
			return fmt.Errorf("fsync: %w (close: %v)", err, cErr)
		}
		return fmt.Errorf("fsync: %w", err)
	}
	return f.Close()
}
