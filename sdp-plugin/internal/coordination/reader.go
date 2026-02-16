package coordination

import (
	"bufio"
	"encoding/json"
	"fmt"
	"os"
)

// scanEvents reads all events from a file
func scanEvents(path string, fn func(*AgentEvent) error) error {
	f, err := os.Open(path)
	if err != nil {
		if os.IsNotExist(err) {
			return nil
		}
		return err
	}
	defer f.Close()

	sc := bufio.NewScanner(f)
	for sc.Scan() {
		line := sc.Bytes()
		if len(line) == 0 {
			continue
		}
		var ev AgentEvent
		if err := json.Unmarshal(line, &ev); err != nil {
			continue
		}
		if fn != nil {
			if err := fn(&ev); err != nil {
				return err
			}
		}
	}
	return sc.Err()
}

// verifyHashChainWithCallback verifies hash chain and calls onError on failure
func verifyHashChain(path string) error {
	var prevHash string
	lineNum := 0

	return scanEvents(path, func(e *AgentEvent) error {
		lineNum++
		if e.PrevHash != prevHash && lineNum > 1 {
			return fmt.Errorf("hash chain broken at line %d: expected prev_hash %q, got %q",
				lineNum, prevHash, e.PrevHash)
		}
		prevHash = e.Hash
		return nil
	})
}
