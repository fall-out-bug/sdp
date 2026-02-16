package beads

import (
	"encoding/json"
	"fmt"
	"os"
	"time"
)

// UpdateMapping updates the mapping file with a new entry
func (c *Client) UpdateMapping(wsID, beadsID string) error {
	entries, err := c.readMapping()
	if err != nil {
		return err
	}

	// Check if entry already exists
	found := false
	for i, entry := range entries {
		if entry.SdpID == wsID {
			// Update existing entry
			entries[i].BeadsID = beadsID
			entries[i].UpdatedAt = time.Now().Format(time.RFC3339)
			found = true
			break
		}
	}

	if !found {
		// Add new entry
		entry := mappingEntry{
			SdpID:     wsID,
			BeadsID:   beadsID,
			UpdatedAt: time.Now().Format(time.RFC3339),
		}
		entries = append(entries, entry)
	}

	// Write back to file
	return c.writeMapping(entries)
}

// writeMapping writes entries to the mapping file
func (c *Client) writeMapping(entries []mappingEntry) error {
	file, err := os.Create(c.mappingPath)
	if err != nil {
		return fmt.Errorf("failed to create mapping file: %w", err)
	}
	defer func() {
		if cerr := file.Close(); cerr != nil {
			fmt.Fprintf(os.Stderr, "warning: failed to close mapping file: %v\n", cerr)
		}
	}()

	for _, entry := range entries {
		data, err := json.Marshal(entry)
		if err != nil {
			return fmt.Errorf("failed to marshal entry: %w", err)
		}

		if _, err := file.Write(append(data, '\n')); err != nil {
			return fmt.Errorf("failed to write entry: %w", err)
		}
	}

	return nil
}
