package decision

import (
	"encoding/json"
	"fmt"
	"log"
	"os"
)

// LoadOptions controls pagination
type LoadOptions struct {
	Offset int // Starting index (0-based)
	Limit  int // Max decisions to return (0 = no limit)
}

// LoadAll loads all decisions from the log
func (l *Logger) LoadAll() ([]Decision, error) {
	log.Printf("[decision] LoadAll: start, path=%s", l.filePath)

	file, err := os.Open(l.filePath)
	if err != nil {
		if os.IsNotExist(err) {
			log.Printf("[decision] LoadAll: file not found")
			return []Decision{}, nil // No decisions yet
		}
		log.Printf("[decision] ERROR: failed to open: path=%s error=%v", l.filePath, err)
		return nil, fmt.Errorf("failed to open decisions file: %w", err)
	}
	defer file.Close()

	var decisions []Decision
	decoder := json.NewDecoder(file)
	parseErrors := 0

	for decoder.More() {
		var decision Decision
		if err := decoder.Decode(&decision); err != nil {
			parseErrors++
			log.Printf("[decision] WARN: parse error #%d: %v", parseErrors, err)
			break // End of file or error
		}
		decisions = append(decisions, decision)
	}

	log.Printf("[decision] LoadAll: success, count=%d parse_errors=%d", len(decisions), parseErrors)
	return decisions, nil
}

// Load loads decisions with pagination
func (l *Logger) Load(opts LoadOptions) ([]Decision, error) {
	log.Printf("[decision] Load: offset=%d limit=%d path=%s", opts.Offset, opts.Limit, l.filePath)

	all, err := l.LoadAll()
	if err != nil {
		return nil, err
	}

	// Apply offset
	if opts.Offset >= len(all) {
		log.Printf("[decision] Load: offset exceeds total, returning empty")
		return []Decision{}, nil
	}

	start := opts.Offset
	end := len(all)

	// Apply limit
	if opts.Limit > 0 && start+opts.Limit < end {
		end = start + opts.Limit
	}

	result := all[start:end]
	log.Printf("[decision] Load: success, returned=%d total=%d", len(result), len(all))
	return result, nil
}
