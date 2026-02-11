package metrics

import (
	"fmt"
)

// Collector reads evidence events and computes metrics (AC1).
type Collector struct {
	logPath       string
	outputPath    string
	watermarkPath string
	wsModel       map[string]string
}

// NewCollector creates a metrics collector for given log path.
func NewCollector(logPath, outputPath string) *Collector {
	return &Collector{
		logPath:    logPath,
		outputPath: outputPath,
		wsModel:    make(map[string]string),
	}
}

// SetWatermarkPath sets the path for storing incremental collection watermark (AC7).
func (c *Collector) SetWatermarkPath(path string) {
	c.watermarkPath = path
}

// Collect reads evidence log and computes metrics (AC1-AC7).
func (c *Collector) Collect() (*Metrics, error) {
	metrics := &Metrics{
		IterationCount: make(map[string]int),
		ModelPassRate:  make(map[string]float64),
	}

	// Load watermark for incremental processing
	processedEventIDs := c.loadWatermarkIfExists()

	// Read and process events
	events, err := c.readEvents()
	if err != nil {
		return nil, fmt.Errorf("read events: %w", err)
	}

	_, modelStats, newProcessedIDs := c.processEvents(events, processedEventIDs, metrics)

	c.computeModelPassRates(modelStats, metrics)
	c.computeCatchRates(metrics)

	if err := c.writeOutput(metrics); err != nil {
		return nil, fmt.Errorf("write output: %w", err)
	}

	if err := c.saveWatermarkIfNeeded(newProcessedIDs); err != nil {
		return nil, fmt.Errorf("save watermark: %w", err)
	}

	return metrics, nil
}

// loadWatermarkIfExists loads watermark if watermark path is set.
func (c *Collector) loadWatermarkIfExists() map[string]bool {
	if c.watermarkPath == "" {
		return nil
	}
	return c.loadWatermark()
}

// saveWatermarkIfNeeded saves watermark if watermark path is set and there are new events.
func (c *Collector) saveWatermarkIfNeeded(newProcessedIDs map[string]bool) error {
	if c.watermarkPath == "" || len(newProcessedIDs) == 0 {
		return nil
	}
	return c.saveWatermark(newProcessedIDs)
}
