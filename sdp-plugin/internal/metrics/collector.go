package metrics

import (
	"bufio"
	"encoding/json"
	"fmt"
	"os"
	"path/filepath"
	"strconv"
)

// Metrics represents collected metrics from evidence events (AC2-AC5).
type Metrics struct {
	CatchRate           float64            `json:"catch_rate"`
	TotalVerifications  int                `json:"total_verifications"`
	FailedVerifications int                `json:"failed_verifications"`
	IterationCount      map[string]int     `json:"iteration_count"`
	ModelPassRate       map[string]float64 `json:"model_pass_rate"`
	TotalApprovals      int                `json:"total_approvals"`
	FailedApprovals     int                `json:"failed_approvals"`
	AcceptanceCatchRate float64            `json:"acceptance_catch_rate"`
}

// workstreamState tracks workstream iteration state.
type workstreamState struct {
	generationCount int
	lastWasGen      bool
}

// modelVerificationStats tracks verification stats per model.
type modelVerificationStats struct {
	Passed   int
	Total    int
	PassRate float64
}

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
	var processedEventIDs map[string]bool
	if c.watermarkPath != "" {
		processedEventIDs = c.loadWatermark()
	}

	// Read and process events
	events, err := c.readEvents()
	if err != nil {
		return nil, fmt.Errorf("read events: %w", err)
	}

	// Track workstream state for iteration counting
	wsState := make(map[string]*workstreamState)
	// Track model verification stats
	modelStats := make(map[string]*modelVerificationStats)

	newProcessedIDs := make(map[string]bool)

	for _, ev := range events {
		// Skip if already processed (incremental mode)
		if processedEventIDs != nil && processedEventIDs[ev.ID] {
			continue
		}
		newProcessedIDs[ev.ID] = true

		switch ev.Type {
		case "verification":
			c.processVerification(ev, metrics, wsState, modelStats)
		case "approval":
			c.processApproval(ev, metrics)
		case "generation":
			c.processGeneration(ev, wsState)
		}
	}

	// Compute model pass rates
	for modelID, stats := range modelStats {
		if stats.Total > 0 {
			modelStats[modelID].PassRate = float64(stats.Passed) / float64(stats.Total)
		}
		metrics.ModelPassRate[modelID] = stats.PassRate
	}

	// Compute catch rate
	if metrics.TotalVerifications > 0 {
		metrics.CatchRate = float64(metrics.FailedVerifications) / float64(metrics.TotalVerifications)
	}

	// Compute acceptance catch rate
	if metrics.TotalApprovals > 0 {
		metrics.AcceptanceCatchRate = float64(metrics.FailedApprovals) / float64(metrics.TotalApprovals)
	}

	// Write output
	if err := c.writeOutput(metrics); err != nil {
		return nil, fmt.Errorf("write output: %w", err)
	}

	// Save watermark
	if c.watermarkPath != "" && len(newProcessedIDs) > 0 {
		if err := c.saveWatermark(newProcessedIDs); err != nil {
			return nil, fmt.Errorf("save watermark: %w", err)
		}
	}

	return metrics, nil
}

func (c *Collector) processVerification(ev evidenceEvent, metrics *Metrics, wsState map[string]*workstreamState, modelStats map[string]*modelVerificationStats) {
	metrics.TotalVerifications++

	passed, ok := ev.Data["passed"].(bool)
	if !ok {
		return
	}

	if !passed {
		metrics.FailedVerifications++
	}

	// Track model pass rate
	modelID := c.wsModel[ev.WSID]
	if modelID != "" {
		stats, exists := modelStats[modelID]
		if !exists {
			stats = &modelVerificationStats{}
			modelStats[modelID] = stats
		}
		stats.Total++
		if passed {
			stats.Passed++
		}
	}

	// Track iteration count
	if _, exists := metrics.IterationCount[ev.WSID]; !exists {
		metrics.IterationCount[ev.WSID] = 0
	}
	metrics.IterationCount[ev.WSID]++
}

func (c *Collector) processApproval(ev evidenceEvent, metrics *Metrics) {
	metrics.TotalApprovals++

	approved, ok := ev.Data["approved"].(bool)
	if !ok {
		return
	}

	if !approved {
		metrics.FailedApprovals++
	}
}

func (c *Collector) processGeneration(ev evidenceEvent, wsState map[string]*workstreamState) {
	state := wsState[ev.WSID]
	if state == nil {
		state = &workstreamState{}
		wsState[ev.WSID] = state
	}
	state.generationCount++
	state.lastWasGen = true

	// Extract model ID from generation data
	modelID, ok := ev.Data["model_id"].(string)
	if ok && modelID != "" {
		c.wsModel[ev.WSID] = modelID
	}
}

func (c *Collector) readEvents() ([]evidenceEvent, error) {
	f, err := os.Open(c.logPath)
	if err != nil {
		if os.IsNotExist(err) {
			return []evidenceEvent{}, nil
		}
		return nil, fmt.Errorf("open log: %w", err)
	}
	defer f.Close()

	var events []evidenceEvent
	sc := bufio.NewScanner(f)
	for sc.Scan() {
		line := sc.Bytes()
		if len(line) == 0 {
			continue
		}
		var ev evidenceEvent
		if err := json.Unmarshal(line, &ev); err != nil {
			continue // Skip invalid lines
		}
		// Convert data to map[string]interface{} for easier access
		if ev.Data == nil {
			ev.Data = make(map[string]interface{})
		}
		events = append(events, ev)
	}

	return events, sc.Err()
}

// evidenceEvent represents an evidence log event.
type evidenceEvent struct {
	ID        string                 `json:"id"`
	Type      string                 `json:"type"`
	Timestamp string                 `json:"timestamp"`
	WSID      string                 `json:"ws_id"`
	Data      map[string]interface{} `json:"data"`
}

func (c *Collector) writeOutput(metrics *Metrics) error {
	// Ensure output directory exists
	dir := filepath.Dir(c.outputPath)
	if err := os.MkdirAll(dir, 0755); err != nil {
		return fmt.Errorf("create output dir: %w", err)
	}

	data, err := json.MarshalIndent(metrics, "", "  ")
	if err != nil {
		return fmt.Errorf("marshal metrics: %w", err)
	}

	return os.WriteFile(c.outputPath, data, 0644)
}

func (c *Collector) loadWatermark() map[string]bool {
	data, err := os.ReadFile(c.watermarkPath)
	if err != nil {
		return make(map[string]bool)
	}

	var ids []string
	if err := json.Unmarshal(data, &ids); err != nil {
		return make(map[string]bool)
	}

	result := make(map[string]bool, len(ids))
	for _, id := range ids {
		result[id] = true
	}
	return result
}

func (c *Collector) saveWatermark(processedIDs map[string]bool) error {
	ids := make([]string, 0, len(processedIDs))
	for id := range processedIDs {
		ids = append(ids, id)
	}

	data, err := json.Marshal(ids)
	if err != nil {
		return fmt.Errorf("marshal watermark: %w", err)
	}

	// Ensure watermark directory exists
	dir := filepath.Dir(c.watermarkPath)
	if err := os.MkdirAll(dir, 0755); err != nil {
		return fmt.Errorf("create watermark dir: %w", err)
	}

	return os.WriteFile(c.watermarkPath, data, 0644)
}

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
