package metrics

import (
	"encoding/json"
	"fmt"
)

// GenerateJSON creates JSON benchmark report (AC5).
func (r *Reporter) GenerateJSON() (string, error) {
	data, err := r.loadReportData()
	if err != nil {
		return "", fmt.Errorf("load report data: %w", err)
	}

	reportJSON, err := json.MarshalIndent(data, "", "  ")
	if err != nil {
		return "", fmt.Errorf("marshal report: %w", err)
	}

	return string(reportJSON), nil
}
