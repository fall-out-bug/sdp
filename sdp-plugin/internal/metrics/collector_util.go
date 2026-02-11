package metrics

import (
	"encoding/json"
	"fmt"
	"os"
	"strconv"
)

// GetLatestWatermark reads the last processed event ID from watermark file.
func GetLatestWatermark(watermarkPath string) (string, error) {
	data, err := os.ReadFile(watermarkPath)
	if err != nil {
		if os.IsNotExist(err) {
			return "", nil
		}
		return "", fmt.Errorf("read watermark: %w", err)
	}

	var ids []string
	if err := json.Unmarshal(data, &ids); err != nil {
		return "", fmt.Errorf("parse watermark: %w", err)
	}

	if len(ids) == 0 {
		return "", nil
	}

	// Return the last ID as the watermark
	return ids[len(ids)-1], nil
}

// ParseIntFromPath extracts an integer from a file path component.
func ParseIntFromPath(s string) int {
	i, err := strconv.Atoi(s)
	if err != nil {
		return 0
	}
	return i
}
