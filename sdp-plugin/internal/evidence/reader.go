package evidence

import (
	"bufio"
	"encoding/json"
	"fmt"
	"os"
)

// Reader reads and validates the evidence log hash chain (AC5, AC6).
type Reader struct {
	path string
}

// NewReader creates a reader for the given path.
func NewReader(path string) *Reader {
	return &Reader{path: path}
}

// Verify validates hash chain integrity; returns first broken link or nil (AC6).
func (r *Reader) Verify() error {
	f, err := os.Open(r.path)
	if err != nil {
		if os.IsNotExist(err) {
			return nil
		}
		return fmt.Errorf("open: %w", err)
	}
	defer f.Close()
	sc := bufio.NewScanner(f)
	var prevHash string = genesisHash
	lineNum := 0
	for sc.Scan() {
		lineNum++
		line := sc.Bytes()
		if len(line) == 0 {
			continue
		}
		// Parse prev_hash from line (minimal: we need only prev_hash to validate chain)
		ev := struct {
			PrevHash string `json:"prev_hash"`
		}{}
		if err := json.Unmarshal(line, &ev); err != nil {
			return fmt.Errorf("line %d: invalid json: %w", lineNum, err)
		}
		if ev.PrevHash != prevHash {
			return fmt.Errorf("line %d: chain broken (prev_hash %q != expected %q)", lineNum, ev.PrevHash, prevHash)
		}
		prevHash = hashLine(line)
	}
	return sc.Err()
}

// ReadAll reads all events from the log file (for trace/show).
func (r *Reader) ReadAll() ([]Event, error) {
	f, err := os.Open(r.path)
	if err != nil {
		if os.IsNotExist(err) {
			return nil, nil
		}
		return nil, fmt.Errorf("open: %w", err)
	}
	defer f.Close()
	var out []Event
	sc := bufio.NewScanner(f)
	for sc.Scan() {
		line := sc.Bytes()
		if len(line) == 0 {
			continue
		}
		var ev Event
		if err := json.Unmarshal(line, &ev); err != nil {
			continue
		}
		out = append(out, ev)
	}
	return out, sc.Err()
}
