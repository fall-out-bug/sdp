package telemetry

import "fmt"

// Analyzer analyzes telemetry data to generate insights
type Analyzer struct {
	filePath string
}

// NewAnalyzer creates a new telemetry analyzer
func NewAnalyzer(filePath string) (*Analyzer, error) {
	if filePath == "" {
		return nil, fmt.Errorf("file path cannot be empty")
	}

	return &Analyzer{
		filePath: filePath,
	}, nil
}
