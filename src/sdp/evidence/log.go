package evidence

import (
	"bufio"
	"io"
	"os"
	"path/filepath"
	"sync"
)

// EvidenceLog manages the append-only evidence log.
type EvidenceLog struct {
	path string
	mu   sync.Mutex
}

// NewEvidenceLog creates a new evidence log at the specified path.
func NewEvidenceLog(path string) (*EvidenceLog, error) {
	dir := filepath.Dir(path)
	if err := os.MkdirAll(dir, 0755); err != nil {
		return nil, err
	}
	return &EvidenceLog{path: path}, nil
}

// Append adds an event to the log.
func (l *EvidenceLog) Append(event BaseEvent) error {
	l.mu.Lock()
	defer l.mu.Unlock()
	line, err := event.ToJSONL()
	if err != nil {
		return err
	}
	f, err := os.OpenFile(l.path, os.O_APPEND|os.O_CREATE|os.O_WRONLY, 0644)
	if err != nil {
		return err
	}
	defer f.Close()
	_, err = f.WriteString(line)
	return err
}

// AppendEvent creates and appends a new event with automatic hash chain.
func (l *EvidenceLog) AppendEvent(eventType EventType, extra func(*BaseEvent)) (BaseEvent, error) {
	l.mu.Lock()
	defer l.mu.Unlock()
	lastHash, err := l.getLastHashLocked()
	if err != nil && !os.IsNotExist(err) {
		return BaseEvent{}, err
	}
	event := NewBaseEvent(eventType, lastHash)
	if extra != nil {
		extra(&event)
		event.CalculateHash()
	}
	f, err := os.OpenFile(l.path, os.O_APPEND|os.O_CREATE|os.O_WRONLY, 0644)
	if err != nil {
		return BaseEvent{}, err
	}
	defer f.Close()
	line, err := event.ToJSONL()
	if err != nil {
		return BaseEvent{}, err
	}
	_, err = f.WriteString(line)
	return event, err
}

// ReadAll returns all events from the log.
func (l *EvidenceLog) ReadAll() ([]BaseEvent, error) {
	l.mu.Lock()
	defer l.mu.Unlock()
	return l.readAllLocked()
}

// GetLastHash returns the hash of the last event.
func (l *EvidenceLog) GetLastHash() (string, error) {
	l.mu.Lock()
	defer l.mu.Unlock()
	hash, err := l.getLastHashLocked()
	if os.IsNotExist(err) {
		return "", nil
	}
	return hash, err
}

// Truncate removes events from the end of the log.
func (l *EvidenceLog) Truncate(keepCount int) error {
	l.mu.Lock()
	defer l.mu.Unlock()
	events, err := l.readAllLocked()
	if err != nil || keepCount >= len(events) {
		return err
	}
	f, err := os.Create(l.path)
	if err != nil {
		return err
	}
	defer f.Close()
	for _, event := range events[:keepCount] {
		if line, err := event.ToJSONL(); err == nil {
			f.WriteString(line)
		}
	}
	return nil
}

// Exists returns true if the log file exists.
func (l *EvidenceLog) Exists() bool {
	_, err := os.Stat(l.path)
	return err == nil
}

// Path returns the log file path.
func (l *EvidenceLog) Path() string {
	return l.path
}

// copyToInternal copies all events to a writer.
func (l *EvidenceLog) copyToInternal(w io.Writer) error {
	l.mu.Lock()
	defer l.mu.Unlock()
	f, err := os.Open(l.path)
	if err != nil {
		if os.IsNotExist(err) {
			return nil
		}
		return err
	}
	defer f.Close()
	_, err = io.Copy(w, f)
	return err
}

func (l *EvidenceLog) getLastHashLocked() (string, error) {
	f, err := os.Open(l.path)
	if err != nil {
		return "", err
	}
	defer f.Close()
	var lastHash string
	scanner := bufio.NewScanner(f)
	for scanner.Scan() {
		if line := scanner.Text(); line != "" {
			if e, err := ParseJSONL(line); err == nil {
				lastHash = e.Hash
			}
		}
	}
	return lastHash, scanner.Err()
}

func (l *EvidenceLog) readAllLocked() ([]BaseEvent, error) {
	f, err := os.Open(l.path)
	if err != nil {
		if os.IsNotExist(err) {
			return []BaseEvent{}, nil
		}
		return nil, err
	}
	defer f.Close()
	var events []BaseEvent
	scanner := bufio.NewScanner(f)
	for scanner.Scan() {
		if line := scanner.Text(); line != "" {
			if e, err := ParseJSONL(line); err == nil {
				events = append(events, e)
			}
		}
	}
	return events, scanner.Err()
}
